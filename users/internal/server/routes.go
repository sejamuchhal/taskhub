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
	r.POST("/tasks", Authenticate(s), s.Validate)

	return r
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
