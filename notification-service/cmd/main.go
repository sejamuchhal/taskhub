package main

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	rabbitmq "github.com/sejamuchhal/taskhub/notification-service/events"

	"github.com/sejamuchhal/taskhub/notification-service/common"
	"github.com/sejamuchhal/taskhub/notification-service/worker"
)

func main() {
	conf := common.LoadConfig()
	emailSender := worker.NewEmailSender(conf.MailersendAPIKey)
	var conn *amqp.Connection
	var err error

	logger := conf.Logger
	backoff := time.Millisecond * 500
	maxBackoff := time.Second * 60

	for {
		conn, err = amqp.Dial(conf.RabbitMQURL)
		if err == nil {
			break
		}

		logger.WithError(err).Error("Error while connecting to RabbitMQ, retrying...")
		time.Sleep(backoff)

		backoff *= 2
		if backoff > maxBackoff {
			backoff = maxBackoff
		}
	}

	rmq := &rabbitmq.RabbitMQBroker{
		Connection: conn,
		Logger:     logger,
		QueueName:  conf.TaskQueue,
	}
	worker.NewWorker(emailSender, rmq)
}
