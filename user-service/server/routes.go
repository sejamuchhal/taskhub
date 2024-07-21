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
	r.POST("/tasks", Authenticate(s), s.CreateTask)
	r.GET("/tasks/:id", Authenticate(s), s.GetTask)
	r.GET("/tasks", Authenticate(s), s.ListTasks)
	r.DELETE("/tasks/:id", Authenticate(s), s.DeleteTask)
	r.PUT("/tasks/:id", Authenticate(s), s.UpdateTask)
	r.PUT("/tasks/:id/complete", Authenticate(s), s.CompleteTask)

	return r
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
