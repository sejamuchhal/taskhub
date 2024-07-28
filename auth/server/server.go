package server

import (
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/sejamuchhal/taskhub/auth/common"
	auth "github.com/sejamuchhal/taskhub/auth/pb"
	"github.com/sejamuchhal/taskhub/auth/storage"
	"github.com/sejamuchhal/taskhub/auth/util"
	"github.com/sirupsen/logrus"
)

// Server implements the TaskServiceServer interface
type Server struct {
	auth.UnsafeAuthServiceServer
	TokenHandler util.TokenHandler
	Storage      *storage.Storage
	Logger       *logrus.Entry
	Config       *common.Config
}

func NewServer(cfg *common.Config) (*Server, error) {
	logger := common.Logger
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	if rdb == nil {
		logger.Error("Unable to connect to redis")
		return nil, errors.New("unable to connect to redis")
	}

	server := &Server{
		Storage:      storage.New(),
		TokenHandler: util.NewTokenHandler(cfg.JWTSecret, rdb),
		Logger:       logger,
		Config:       cfg,
	}
	return server, nil
}
