package server

import (
	"github.com/sejamuchhal/taskhub/task-service/common"
	"github.com/sejamuchhal/taskhub/task-service/events"
	pb "github.com/sejamuchhal/taskhub/task-service/proto"
	"github.com/sejamuchhal/taskhub/task-service/storage"
	"github.com/sirupsen/logrus"
)

// Server implements the TaskServiceServer interface
type Server struct {
	pb.UnimplementedTaskServiceServer
	RMQ     *events.RabbitMQ
	Storage *storage.Storage
	Logger  *logrus.Entry
}

func NewServer(cfg *common.Config) (*Server, error) {
	logger := common.Logger
	rmq, err := events.NewRabbitMQ(cfg.RMQQueue, cfg.RMQUser, cfg.RMQPassword)
	if err != nil {
		logger.WithError(err).Error("Failed to connect to RabbitMQ")
		return nil, err
	}

	server := &Server{
		RMQ:     rmq,
		Storage: storage.New(),
		Logger:  logger,
	}
	return server, nil
}
