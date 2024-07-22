package common

import (
	"fmt"
	"os"
)

type Config struct {
	RMQUser     string
	RMQPassword string
	RMQQueue    string
	RMQPort     string
	RMQUrl      string
	GRPCAddress string
}

func LoadConfig() (*Config, error) {
	config := &Config{
		RMQUser:     getEnv("RABBITMQ_USER", ""),
		RMQPassword: getEnv("RABBITMQ_PASSWORD", ""),
		RMQQueue:    getEnv("RABBITMQ_QUEUE", ""),
		RMQPort:     getEnv("RABBITMQ_PORT", "5672"),
		RMQUrl:      getEnv("RABBITMQ_URL", "amqp://rmq:rmq@rabbit1:5672/"),
		GRPCAddress: getEnv("GRPC_ADDRESS", "0.0.0.0:5000"),
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
