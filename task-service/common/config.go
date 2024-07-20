package common

import (
	"fmt"
	"os"
)

type Config struct {
	RMQUser     string
	RMQPassword string
	RMQQueue    string
	GRPCAddress string
}

func LoadConfig() (*Config, error) {
	config := &Config{
		RMQUser:     getEnv("RABBITMQ_USER", ""),
		RMQPassword: getEnv("RABBITMQ_PASSWORD", ""),
		RMQQueue:    getEnv("RABBITMQ_QUEUE", ""),
		GRPCAddress: getEnv("GRPC_ADDRESS", ":8005"),
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
