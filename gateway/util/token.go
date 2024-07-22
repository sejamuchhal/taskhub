package util

import (
	"github.com/golang-jwt/jwt/v5"
)

type TokenHandler struct {
	secret string
}

type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func NewTokenHandler(secret string) TokenHandler {
	return TokenHandler{
		secret: secret,
	}
}

func (handler *TokenHandler) VerifyToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(handler.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
