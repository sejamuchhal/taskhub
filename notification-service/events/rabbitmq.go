package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type RabbitMQBroker struct {
	QueueName  string
	Connection *amqp.Connection
	MsgHandler func(queue string, msg amqp.Delivery, err error)
	Logger     *logrus.Entry
}

func (rmq *RabbitMQBroker) OnError(err error, msg string) {
	if err != nil {
		rmq.Logger.WithError(err).Error(msg)
		rmq.MsgHandler(rmq.QueueName, amqp.Delivery{}, err)
	}
}

func (rmq *RabbitMQBroker) Consume() {
	channel, err := rmq.Connection.Channel()
	rmq.OnError(err, "Failed to open a channel")
	defer channel.Close()

	q, err := channel.QueueDeclare(
		rmq.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	rmq.OnError(err, "Failed to declare a queue")

	msgs, err := channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	rmq.OnError(err, "Failed to register a consumer")

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			rmq.MsgHandler(rmq.QueueName, d, nil)
		}
	}()
	<-forever

}
