package main

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	rabbitmq "github.com/sejamuchhal/taskhub/notification/events"

	"github.com/sejamuchhal/taskhub/notification/common"
	"github.com/sejamuchhal/taskhub/notification/server"
)

func main() {

	conf := common.LoadConfig()
	logger := conf.Logger
	logger.Info("Configuration loaded successfully")

	emailSender := worker.NewEmailSender(conf)

	var conn *amqp.Connection
	var err error

	backoff := time.Millisecond * 500
	maxBackoff := time.Second * 60

	logger.Info("Starting RabbitMQ connection attempts")

	for {
		conn, err = amqp.Dial(conf.RabbitMQURL)
		if err == nil {
			logger.Info("Successfully connected to RabbitMQ")
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
	worker := worker.NewWorker(emailSender, rmq)
	rmq.MsgHandler = worker.NotificationHandler

	// Start consuming messages
	logger.Info("Starting message consumption from RabbitMQ")
	rmq.Consume()
}
