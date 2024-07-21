package common

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Logger           *logrus.Entry
	MailersendAPIKey string
	RabbitMQURL      string
	TaskQueue        string
}

func LoadConfig() *Config {

	taskQueuequeue := GetEnvDefault("TASK_QUEUE", "task_queue")
	mailersendAPIKey := os.Getenv("MAILERSEND_API_KEY")
	rabbitMQURL := os.Getenv("RABBITMQ_URL")

	config := &Config{
		Logger:           Logger,
		MailersendAPIKey: mailersendAPIKey,
		RabbitMQURL:      rabbitMQURL,
		TaskQueue:        taskQueuequeue,
	}
	return config
}

func GetEnvDefault(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		val = defaultValue
	}
	return val
}
