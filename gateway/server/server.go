package server

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/sejamuchhal/taskhub/gateway/pb/auth"
	"github.com/sejamuchhal/taskhub/gateway/pb/task"
	"github.com/sejamuchhal/taskhub/gateway/common"
)

type Server struct {
	Logger     *logrus.Entry
	AuthClient auth.AuthServiceClient
	TaskClient task.TaskServiceClient
}

func NewServer(config *common.Config) (*http.Server, error) {
	logger := common.Logger
	logger.WithField("config", config).Info("Got config")

	authClient, err := NewAuthClient(config)
	if err != nil {
		return nil, err
	}

	taskClient, err := NewTaskClient(config)
	if err != nil {
		return nil, err
	}

	newServer := &Server{
		Logger:     logger,
		AuthClient: authClient,
		TaskClient: taskClient,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.HTTPPort),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server, nil
}

func NewAuthClient(config *common.Config) (auth.AuthServiceClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	config.Logger.Debug("Connecting to auth grpc server")
	authConn, err := grpc.Dial(config.AuthServiceUrl, grpc.WithInsecure())
	if err != nil {
		config.Logger.WithError(err).Error("Error connecting to auth grpc server")
		return nil, fmt.Errorf("cannot connect to auth grpc server: %v", err)
	} else {
		config.Logger.Info("Connected to auth grpc server successfully")
	}

	authClient := auth.NewAuthServiceClient(authConn)
	return authClient, nil
}

func NewTaskClient(config *common.Config) (task.TaskServiceClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	config.Logger.Debug("Connecting to task grpc server")
	taskConn, err := grpc.Dial(config.TaskServiceUrl, grpc.WithInsecure())
	if err != nil {
		config.Logger.WithError(err).Error("Error connecting to task grpc server")
		return nil, fmt.Errorf("cannot connect to task grpc server: %v", err)
	} else {
		config.Logger.Info("Connected to task grpc server successfully")
	}

	taskClient := task.NewTaskServiceClient(taskConn)
	return taskClient, nil
}
