package server

import (

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sejamuchhal/taskhub/task-service/common"
	rabbitmq "github.com/sejamuchhal/taskhub/task-service/events"
	pb "github.com/sejamuchhal/taskhub/task-service/pb/task"
	"github.com/sejamuchhal/taskhub/task-service/storage"
	"github.com/sirupsen/logrus"
)

// Server implements the TaskServiceServer interface
type Server struct {
	pb.UnimplementedTaskServiceServer
	Publisher *rabbitmq.RabbitMQBroker
	Storage   *storage.Storage
	Logger    *logrus.Entry
}

func NewServer(cfg *common.Config) (*Server, error) {
	logger := common.Logger

	conn, err := amqp.Dial(cfg.RMQUrl)
	if err != nil {
		logger.WithError(err).Error("Failed to connect to rabbitmq")
		return nil, err
	}

	rmq := &rabbitmq.RabbitMQBroker{
		Connection: conn,
		Logger:     logger,
		QueueName:  cfg.RMQQueue,
	}

	server := &Server{
		Publisher: rmq,
		Storage:   storage.New(),
		Logger:    logger,
	}
	return server, nil
}
