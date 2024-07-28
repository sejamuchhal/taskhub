package common

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Config struct {
	GRPCAddress          string
	JWTSecret            string
	Logger               *logrus.Entry
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

func LoadConfig() (*Config, error) {
	port := getEnv("AUTH_PORT", "4040")

	grpcAddress := fmt.Sprintf("0.0.0.0:%v", port)
	config := &Config{
		GRPCAddress:          grpcAddress,
		JWTSecret:            getEnv("JWT_SECRET", "taskhub"),
		Logger:               Logger,
		AccessTokenDuration:  20 * time.Minute,
		RefreshTokenDuration: 24 * time.Hour,
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
