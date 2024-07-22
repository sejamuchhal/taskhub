package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/sejamuchhal/taskhub/gateway/client/task"
	"github.com/sejamuchhal/taskhub/gateway/common"
	"github.com/sejamuchhal/taskhub/gateway/util"
)

type Server struct {
	port         int64
	tokenHandler util.TokenHandler
	logger       *logrus.Entry
	taskClient   task.TaskServiceClient
}

func NewServer(config *common.Config) *http.Server {
	logger := common.Logger
	logger.WithField("config", config).Info("Got config")

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	logger.Debug("Connecting to task grpc server")
	taskConn, err := grpc.Dial("task:5000", opts...)
	if err != nil {
		logger.WithError(err).Error("Error connecting to task grpc server")
		panic(fmt.Sprintf("cannot connect to task grpc server: %v", err))
	} else {
		logger.Info("Connected to task grpc server successfully")
	}

	taskClient := task.NewTaskServiceClient(taskConn)

	newServer := &Server{
		port:         config.HTTPPort,
		tokenHandler: util.NewTokenHandler(os.Getenv("JWT_SECRET")),
		logger:       logger,
		taskClient:   taskClient,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
