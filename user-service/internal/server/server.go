package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"

	"github.com/sejamuchhal/taskhub/user-service/common"
	"github.com/sejamuchhal/taskhub/user-service/internal"
	"github.com/sejamuchhal/taskhub/user-service/internal/database"
)

type Server struct {
	port         int
	db           *database.Storage
	tokenHandler internal.TokenHandler
	logger       *logrus.Entry
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:         port,
		db:           database.New(),
		tokenHandler: internal.NewTokenHandler(os.Getenv("JWT_SECRET")),
		logger:       common.Logger,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
