package worker

import (
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	rabbitmq "github.com/sejamuchhal/taskhub/notification-service/events"
	"github.com/sejamuchhal/taskhub/notification-service/common"
)

type Worker struct {
	EmailSender    *EmailSender
	RabbitMQBroker *rabbitmq.RabbitMQBroker
}

func NewWorker(emailSender *EmailSender, rmq *rabbitmq.RabbitMQBroker) *Worker {
	return &Worker{
		EmailSender:    emailSender,
		RabbitMQBroker: rmq,
	}
}

func Run() {
	conf := common.LoadConfig()
	fmt.Printf("Starting background worker with config: %v\n", conf)

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

	logger.Info("Successfully connected to RabbitMQ")

	rmq := &rabbitmq.RabbitMQBroker{
		QueueName:  conf.TaskQueue,
		Connection: conn,
		Logger:     conf.Logger,
	}

	emailSender := NewEmailSender(conf.MailersendAPIKey)

	worker := NewWorker(emailSender, rmq)
	rmq.MsgHandler = worker.ReminderHandler

	rmq.Consume()

}