package worker

import (
	"github.com/sejamuchhal/taskhub/notification/rabbitmq"
)

type Worker struct {
	EmailSender    EmailSenderInterface
	RabbitMQBroker *rabbitmq.RabbitMQBroker
}

func NewWorker(emailSender *EmailSender, rmq *rabbitmq.RabbitMQBroker) *Worker {
	return &Worker{
		EmailSender:    emailSender,
		RabbitMQBroker: rmq,
	}
}
