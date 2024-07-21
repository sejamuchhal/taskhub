package common

import (
	"fmt"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

type Config struct {
	RMQUser     string
	RMQPassword string
	RMQQueue    string
	RMQPort     string
	HTTPPort    int64
	GRPCAddress string
	Logger      *logrus.Entry
}

func LoadConfig() (*Config, error) {
	httpPort, err := strconv.ParseInt(getEnv("HTTPPort", "3000"), 10, 32)
	if err != nil {
		return nil, fmt.Errorf("missing required port environment variables")
	}

	config := &Config{
		RMQUser:     getEnv("RABBITMQ_USER", "rmq"),
		RMQPassword: getEnv("RABBITMQ_PASSWORD", "rmq"),
		RMQQueue:    getEnv("RABBITMQ_QUEUE", "task_queue"),
		RMQPort:     getEnv("RABBITMQ_PORT", "5672"),
		GRPCAddress: getEnv("GRPC_ADDRESS", "0.0.0.0:8080"),
		HTTPPort:    httpPort,
		Logger:      Logger,
	}

	if config.RMQUser == "" || config.RMQPassword == "" || config.RMQQueue == "" {
		return nil, fmt.Errorf("missing required environment variables")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
