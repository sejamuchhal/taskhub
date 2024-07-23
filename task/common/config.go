package common

import (
	"fmt"
	"os"
)

type Config struct {
	RMQQueue    string
	RMQUrl      string
	GRPCAddress string
}

func LoadConfig() (*Config, error) {
	port := getEnv("TASK_PORT", "8080")
	grpcAddress := fmt.Sprintf("0.0.0.0:%v", port)
	config := &Config{
		RMQQueue:    getEnv("RABBITMQ_QUEUE", ""),
		RMQUrl:      getEnv("RABBITMQ_URL", "amqp://rmq:rmq@rabbit1:5672/"),
		GRPCAddress: grpcAddress,
	}

	if config.RMQUrl == "" || config.RMQQueue == "" {
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
