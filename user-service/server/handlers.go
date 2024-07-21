package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/sejamuchhal/taskhub/user-service/client/task"
	"github.com/sejamuchhal/taskhub/user-service/database"
)

func (s *Server) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "It's healthy"})
}

func (s *Server) SignupUser(c *gin.Context) {
	s.logger.Info("Incoming signup request")

	var req SignupUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.logger.WithError(err).Error("Error parsing signup request payload")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		s.logger.WithError(err).Error("Error hashing password")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = s.db.CreateUser(&database.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			s.logger.WithError(err).Error("User already exists")
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}
		s.logger.WithError(err).Error("Error creating user in the database")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	s.logger.Info("User signup successful")
	c.JSON(http.StatusOK, gin.H{"message": "User signup successful"})
}

func (s *Server) LoginUser(c *gin.Context) {
	s.logger.Info("Incoming login request")

	var req LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.logger.WithError(err).Error("Error parsing login request payload")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := s.db.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.WithError(err).Warn("User not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid email or password"})
			return
		}
		s.logger.WithError(err).Error("Error fetching user from the database")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := CheckPasswordHash(req.Password, user.Password); err != nil {
		s.logger.WithError(err).Warn("Invalid password")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	expiry := time.Now().Add(24 * time.Hour)
	token, err := s.tokenHandler.CreateToken(user.ID, expiry)
	if err != nil {
		s.logger.WithError(err).Error("Error creating access token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := LoginUserResponse{
		AccessToken:          token,
		AccessTokenExpiresAt: expiry,
		User: userDetail{
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}

	s.logger.Info("User login successful")
	c.JSON(http.StatusOK, res)
}

func (s *Server) CreateTask(c *gin.Context) {
	s.logger.Info("Incoming create task request")

	var taskReq CreateTaskRequest
	if err := c.ShouldBindJSON(&taskReq); err != nil {
		s.logger.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	layout := "January 2, 2006 3:04 PM MST"
	dueDateTime, err := time.Parse(layout, taskReq.DueDateTime)
	if err != nil {
		s.logger.WithError(err).Error("Invalid date time format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date time format, please use: January 2, 2006 3:04 PM MST"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		s.logger.Error("User ID not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		s.logger.Error("Invalid user ID type")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	s.logger.Debug("Sending CreateTask gRPC request")
	resp, err := s.taskClient.CreateTask(context.Background(), &task.CreateTaskRequest{
		Task: &task.Task{
			Title:       taskReq.Title,
			Description: taskReq.Description,
			UserId:      userIDStr,
			DueDate:     timestamppb.New(dueDateTime),
		},
	})
	if err != nil {
		s.logger.WithError(err).Error("Failed to create task via gRPC")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	s.logger.WithField("response", resp).Debug("Received response from gRPC")
	c.JSON(http.StatusOK, gin.H{"task_id": resp.Id})
}

func (s *Server) GetTask(c *gin.Context) {
	s.logger.Info("Incoming get task request")

	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	resp, err := s.taskClient.GetTask(context.Background(), &task.GetTaskRequest{Id: taskID})
	if err != nil {
		s.logger.WithError(err).Error("Failed to get task via gRPC")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	taskDetails := TransformTask(resp.Task)
	c.JSON(http.StatusOK, taskDetails)
}

func (s *Server) ListTasks(c *gin.Context) {
	s.logger.Info("Incoming list tasks request")

	var req ListTasksRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		s.logger.WithError(err).Error("Failed to bind query parameters")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		s.logger.Error("User ID not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		s.logger.Error("Invalid user ID type")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	resp, err := s.taskClient.ListTasks(context.Background(), &task.ListTasksRequest{
		Limit:   int32(req.Limit),
		Offset:  int32(req.Offset),
		UserId:  userIDStr,
		Pending: req.Pending,
	})
	if err != nil {
		s.logger.WithError(err).Error("Failed to list tasks via gRPC")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(resp.GetTasks()) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No tasks found"})
		return
	}

	taskDetails := make([]*TaskDetails, len(resp.GetTasks()))
	for i, t := range resp.GetTasks() {
		taskDetails[i] = TransformTask(t)
	}

	c.JSON(http.StatusOK, ListTasksResponse{
		Count: int(resp.TotalCount),
		Tasks: taskDetails,
	})
}

func (s *Server) DeleteTask(c *gin.Context) {
	s.logger.Info("Incoming delete task request")

	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	_, err := s.taskClient.DeleteTask(context.Background(), &task.DeleteTaskRequest{Id: taskID})
	if err != nil {
		s.logger.WithError(err).Error("Failed to delete task via gRPC")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func (s *Server) UpdateTask(c *gin.Context) {
	s.logger.Info("Incoming update task request")

	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	var taskReq UpdateTaskRequest
	if err := c.ShouldBindJSON(&taskReq); err != nil {
		s.logger.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	layout := "January 2, 2006 3:04 PM MST"
	dueDateTime, err := time.Parse(layout, taskReq.DueDateTime)
	if err != nil {
		s.logger.WithError(err).Error("Invalid date time format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date time format, please use: January 2, 2006 3:04 PM MST"})
		return
	}

	_, err = s.taskClient.UpdateTask(context.Background(), &task.UpdateTaskRequest{
		Task: &task.Task{
			Id:          taskID,
			Title:       taskReq.Title,
			Description: taskReq.Description,
			DueDate:     timestamppb.New(dueDateTime),
		},
	})
	if err != nil {
		s.logger.WithError(err).Error("Failed to update task via gRPC")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully", "task_id": taskID})
}

func (s *Server) CompleteTask(c *gin.Context) {
	s.logger.Info("Incoming complete task request")

	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	getTaskResp, err := s.taskClient.GetTask(context.Background(), &task.GetTaskRequest{Id: taskID})
	if err != nil {
		s.logger.WithError(err).Error("Failed to get task via gRPC")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	getTaskResp.Task.Status = "completed"

	_, err = s.taskClient.UpdateTask(context.Background(), &task.UpdateTaskRequest{
		Task: getTaskResp.Task,
	})
	if err != nil {
		s.logger.WithError(err).Error("Failed to update task via gRPC")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task marked as completed", "task_id": getTaskResp.Task.Id})
}
