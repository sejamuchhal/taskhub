package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sejamuchhal/taskhub/auth/storage"
)

type TokenHandler struct {
	secret string
}

type CustomClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewTokenHandler(secret string) TokenHandler {
	return TokenHandler{
		secret: secret,
	}
}

func (handler *TokenHandler) CreateToken(user *storage.User, expiry time.Time) (string, error) {
	claims := CustomClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
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
