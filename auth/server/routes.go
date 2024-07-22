package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/health", s.Health)
	r.POST("/signup", s.SignupUser)
	r.POST("/login", s.LoginUser)

	return r
}
