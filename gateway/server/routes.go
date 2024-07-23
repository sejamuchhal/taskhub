package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up all the required routes
func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	r.Use(RecordRequestLatency)
	r.GET("/health", s.Health)
	r.GET("/metrics", prometheusHandler())

	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/signup", s.SignupUser)
		authRoutes.POST("/login", s.LoginUser)

	}

	taskRoutes := r.Group("/tasks")
	{
		taskRoutes.Use(Authenticate(s))
		taskRoutes.POST("", s.CreateTask)
		taskRoutes.GET("/:id", s.GetTask)
		taskRoutes.GET("", s.ListTasks)
		taskRoutes.DELETE("/:id", s.DeleteTask)
		taskRoutes.PUT("/:id", s.UpdateTask)
		taskRoutes.PUT("/:id/complete", s.CompleteTask)
	}

	return r
}
