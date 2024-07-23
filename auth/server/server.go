package server

import (
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
}

func NewServer(cfg *common.Config) (*Server, error) {
	logger := common.Logger

	server := &Server{
		Storage:      storage.New(),
		TokenHandler: util.NewTokenHandler(cfg.JWTSecret),
		Logger:       logger,
	}
	return server, nil
}
