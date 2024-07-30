package server

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sejamuchhal/taskhub/task/common"
	"github.com/sejamuchhal/taskhub/task/events"
	pb "github.com/sejamuchhal/taskhub/task/pb/task"
	"github.com/sejamuchhal/taskhub/task/storage"
	"github.com/sirupsen/logrus"
)

// Server implements the TaskServiceServer interface
type Server struct {
	pb.UnimplementedTaskServiceServer
	Publisher events.RabbitMQBrokerInterface
	Storage   storage.StorageInterface
	Logger    *logrus.Entry
}

func NewServer(cfg *common.Config) (*Server, error) {
	logger := common.Logger

	conn, err := amqp.Dial(cfg.RMQUrl)
	if err != nil {
		logger.WithError(err).Error("Failed to connect to rabbitmq")
		return nil, err
	}

	rmq := &events.RabbitMQBroker{
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
