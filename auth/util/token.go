package util

import (
	"time"

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

func (handler *TokenHandler) CreateToken(user_id string, expiry time.Time) (string, error) {
	claims := CustomClaims{
		UserID: user_id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiry),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(handler.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}