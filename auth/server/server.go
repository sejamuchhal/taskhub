package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/sejamuchhal/taskhub/auth/common"
	"github.com/sejamuchhal/taskhub/auth/database"
	"github.com/sejamuchhal/taskhub/auth/util"
)

type Server struct {
	port         int64
	db           *database.Storage
	tokenHandler util.TokenHandler
	logger       *logrus.Entry
}

func NewServer(config *common.Config) *http.Server {
	logger := common.Logger

	newServer := &Server{
		port:         config.HTTPPort,
		db:           database.New(),
		tokenHandler: util.NewTokenHandler(config.JWTSecret),
		logger:       logger,
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
