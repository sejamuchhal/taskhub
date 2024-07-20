package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"

	"github.com/sejamuchhal/task-management/user-service/internal/database"
	"github.com/sejamuchhal/task-management/user-service/internal/token"
)

type Server struct {
	port         int
	db           *database.Storage
	tokenHandler token.TokenHandler
	logger       *logrus.Entry
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:         port,
		db:           database.New(),
		tokenHandler: token.NewTokenHandler(os.Getenv("JWT_SECRET")),
		logger:       utils.Logger,
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
