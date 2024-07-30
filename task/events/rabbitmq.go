package events

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type RabbitMQBrokerInterface interface {
	Publish(message []byte) error
}

type RabbitMQBroker struct {
	QueueName  string
	Connection *amqp.Connection
	Logger     *logrus.Entry
}

func (rmq *RabbitMQBroker) Publish(message []byte) error {
	channel, err := rmq.Connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	_, err = channel.QueueDeclare(
		rmq.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = channel.Publish(
		"",
		rmq.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
	return err
}