package worker

import (
	rabbitmq "github.com/sejamuchhal/taskhub/notification/events"
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
