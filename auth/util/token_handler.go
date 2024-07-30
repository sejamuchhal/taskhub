package util

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/sejamuchhal/taskhub/auth/storage"
)

type TokenHandler struct {
	secret      string
	redisClient *redis.Client
}

func NewTokenHandler(secret string, redisClient *redis.Client) TokenHandler {
	return TokenHandler{
		secret:      secret,
		redisClient: redisClient,
	}
}

func (handler *TokenHandler) CreateToken(user *storage.User, duration time.Duration, tokenType string) (string, *UserClaims, error) {
	claims, err := handler.NewUserClaims(user, duration, tokenType)
	if err != nil {
		return "", nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(handler.secret))
	if err != nil {
		return "", nil, err
	}

	return tokenString, claims, nil
}

func (handler *TokenHandler) VerifyToken(tokenString string, expectedType string) (*UserClaims, error) {
	// Check if token is blacklisted
	isBlacklisted, err := handler.redisClient.Get(context.Background(), tokenString).Result()
	if err == nil && isBlacklisted == "blacklisted" {
		return nil, errors.New("token is blacklisted")
	}

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(handler.secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.TokenType != expectedType {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}

// Add token to blacklist in Redis
func (handler *TokenHandler) BlacklistToken(tokenString string, expiry time.Duration) error {
	ctx := context.Background()
	err := handler.redisClient.Set(ctx, tokenString, "blacklisted", expiry).Err()
	if err != nil {
		return err
	}
	return nil
}