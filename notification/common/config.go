package common

import (
	"os"

	"github.com/sirupsen/logrus"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Logger                *logrus.Entry
	MailersendAPIKey      string
	MailersendSenderName  string
	MailersendSenderEmail string
	RabbitMQURL           string
	TaskQueue             string
}

func LoadConfig() *Config {

	taskQueuequeue := GetEnvDefault("TASK_QUEUE", "task_queue")
	mailersendAPIKey := os.Getenv("MAILERSEND_API_KEY")
	senderName := os.Getenv("MAILERSEND_SENDER_NAME")
	senderEmail := os.Getenv("MAILERSEND_SENDER_EMAIL")
	rabbitMQURL := GetEnvDefault("RABBITMQ_URL", "amqp://rmq:rmq@rabbit1:5672/")

	config := &Config{
		Logger:                Logger,
		MailersendAPIKey:      mailersendAPIKey,
		MailersendSenderName:  senderName,
		MailersendSenderEmail: senderEmail,
		RabbitMQURL:           rabbitMQURL,
		TaskQueue:             taskQueuequeue,
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
