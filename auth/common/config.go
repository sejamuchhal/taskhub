package common

import (
	"fmt"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

type Config struct {
	HTTPPort  int64
	JWTSecret string
	Logger    *logrus.Entry
}

func LoadConfig() (*Config, error) {
	httpPort, err := strconv.ParseInt(getEnv("AUTH_PORT", "3000"), 10, 32)
	if err != nil {
		return nil, fmt.Errorf("missing required port environment variables")
	}

	config := &Config{
		JWTSecret: getEnv("JWT_SECRET", "test"),
		HTTPPort:  httpPort,
		Logger:    Logger,
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
