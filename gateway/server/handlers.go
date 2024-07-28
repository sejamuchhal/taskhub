package server

import (
	"context"
	"net/http"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/gin-gonic/gin"
	"github.com/golodash/galidator"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sejamuchhal/taskhub/gateway/pb/auth"
	"github.com/sejamuchhal/taskhub/gateway/pb/task"
)

var g = galidator.New()

func (s *Server) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "It's healthy"})
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (s *Server) SignupUser(c *gin.Context) {
	logger := s.Logger.WithField("method", "SignupUser")
	logger.Debug("Incoming request")
	customizer := g.Validator(SignupUserRequest{})
	var req SignupUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithError(err).Error("Error parsing signup request payload")
		c.JSON(http.StatusBadRequest, gin.H{"message": customizer.DecryptErrors(err)})
		return
	}

	res, err := s.AuthClient.Signup(context.Background(), &auth.SignupRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.AlreadyExists:
				logger.WithError(err).Error("User already exists")
				c.JSON(http.StatusConflict, gin.H{"message": "User already exists"})
				return
			default:
				logger.WithError(err).Error("Internal server error")
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
				return
			}
		}
		logger.WithError(err).Error("Unknown error")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
		return
	}

	logger.Info("User signup successful")
	c.JSON(http.StatusOK, gin.H{"message": res.Message})
}

func (s *Server) LoginUser(c *gin.Context) {
	logger := s.Logger.WithField("method", "LoginUser")
	logger.Debug("Incoming request")

	var req LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithError(err).Error("Error parsing login request payload")
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	res, err := s.AuthClient.Login(context.Background(), &auth.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound, codes.Unauthenticated:
				logger.WithError(err).Error("Invalid email or password")
				c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email or password"})
				return
			default:
				logger.WithError(err).Error("Internal server error")
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
				return
			}
		}
		logger.WithError(err).Error("Unknown error")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
		return
	}

	logger.Info("User login successful")
	c.JSON(http.StatusOK, LoginUserResponse{
		AccessToken:           res.AccessToken,
		AccessTokenExpiresAt:  res.AccessTokenExpiresAt.AsTime(),
		RefreshToken:          res.RefreshToken,
		RefreshTokenExpiresAt: res.RefreshTokenExpiresAt.AsTime(),
		SessionID:             res.SessionId,

		User: UserDetail{
			Name:  res.User.Name,
			Email: res.User.Email,
		},
	})
}

func (s *Server) RenewAccessToken(c *gin.Context) {
	logger := s.Logger.WithField("method", "RenewAccessToken")
	logger.Debug("Incoming request")

	refresh := c.GetHeader("Refresh")

	if refresh == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "refresh token header is missing"})
		c.Abort()
		return
	}

	res, err := s.AuthClient.RenewAccessToken(context.Background(), &auth.RenewAccessTokenRequest{
		RefreshToken: refresh,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound, codes.Unauthenticated:
				logger.WithError(err).Error("Invalid refresh token")
				c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid refresh token"})
				return
			case codes.InvalidArgument:
				logger.WithError(err).Error("Invalid request")
				c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
				return
			default:
				logger.WithError(err).Error("Internal server error")
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
				return
			}
		}
		logger.WithError(err).Error("Unknown error")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
		return
	}

	logger.Info("Token renewed")
	c.JSON(http.StatusOK, RenewAccessTokenResponse{
		AccessToken:          res.AccessToken,
		AccessTokenExpiresAt: res.AccessTokenExpiresAt.AsTime(),
	})
}

func (s *Server) Logout(c *gin.Context) {
	logger := s.Logger.WithField("method", "Logout")
	logger.Debug("Incoming request")

	refresh := c.GetHeader("Refresh")

	if refresh == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "refresh token header is missing"})
		c.Abort()
		return
	}

	access := c.GetHeader("Access")

	if access == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "access token header is missing"})
		c.Abort()
		return
	}

	_, err := s.AuthClient.Logout(context.Background(), &auth.LogoutRequest{
		RefreshToken: refresh,
		AccessToken:  access,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				logger.WithError(err).Error("Invalid access or refresh token")
				c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid access or refresh token"})
				return
			default:
				logger.WithError(err).Error("Internal server error")
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
				return
			}
		}
		logger.WithError(err).Error("Unknown error")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unknown error"})
		return
	}

	logger.Debug("User logout successful")
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func (s *Server) CreateTask(c *gin.Context) {
	logger := s.Logger.WithField("method", "CreateTask")
	logger.Debug("Incoming request")
	if !hasPermission(c, []string{"user", "admin"}) {
		c.JSON(http.StatusForbidden, gin.H{"message": "Insufficient permissions"})
		return
	}

	var taskReq CreateTaskRequest
	if err := c.ShouldBindJSON(&taskReq); err != nil {
		logger.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var dueDate *timestamppb.Timestamp
	if taskReq.DueDateTime != "" {
		const layout = "January 2, 2006 3:04 PM MST"
		dueDateTime, err := time.Parse(layout, taskReq.DueDateTime)
		if err != nil {
			logger.WithError(err).Error("Invalid date time format")
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid date time format, please use: January 2, 2006 3:04 PM MST"})
			return
		}
		dueDate = timestamppb.New(dueDateTime)
	}

	md, ok := getGRPCMetadataFromGin(c, logger)
	if !ok {
		return
	}
	ctxWithMetadata := metadata.NewOutgoingContext(context.Background(), md)

	logger.Debug("Sending CreateTask gRPC request")
	resp, err := s.TaskClient.CreateTask(ctxWithMetadata, &task.CreateTaskRequest{
		Task: &task.Task{
			Title:       taskReq.Title,
			Description: taskReq.Description,
			DueDate:     dueDate,
		},
	})
	if err != nil {
		logger.WithError(err).Error("Failed to create task via gRPC")
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	logger.WithField("response", resp).Debug("Received response from gRPC")
	c.JSON(http.StatusOK, gin.H{"task_id": resp.Id})
}

func (s *Server) GetTask(c *gin.Context) {
	logger := s.Logger.WithField("method", "GetTask")
	logger.Debug("Incoming request")
	if !hasPermission(c, []string{"user", "admin"}) {
		c.JSON(http.StatusForbidden, gin.H{"message": "Insufficient permissions"})
		return
	}

	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Task ID is required"})
		return
	}

	// Permission check
	if !hasPermission(c, []string{"user", "admin"}) {
		logger.Warn("Permission denied")
		c.JSON(http.StatusForbidden, gin.H{"message": "Permission denied"})
		return
	}

	md, ok := getGRPCMetadataFromGin(c, logger)
	if !ok {
		return
	}
	ctxWithMetadata := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := s.TaskClient.GetTask(ctxWithMetadata, &task.GetTaskRequest{Id: taskID})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				logger.WithError(err).Error("Task not found")
				c.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
				return
			default:
				logger.WithError(err).Error("Failed to get task via gRPC")
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}

		}
	}

	taskDetails := TransformTask(resp.Task)
	c.JSON(http.StatusOK, taskDetails)
}

func (s *Server) ListTasks(c *gin.Context) {
	logger := s.Logger.WithField("method", "ListTasks")
	logger.Debug("Incoming request")
	if !hasPermission(c, []string{"user", "admin"}) {
		c.JSON(http.StatusForbidden, gin.H{"message": "Insufficient permissions"})
		return
	}

	var req ListTasksRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.WithError(err).Error("Failed to bind query parameters")
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	md, ok := getGRPCMetadataFromGin(c, logger)
	if !ok {
		return
	}
	ctxWithMetadata := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := s.TaskClient.ListTasks(ctxWithMetadata, &task.ListTasksRequest{
		Limit:   int32(req.Limit),
		Offset:  int32(req.Offset),
		Pending: req.Pending,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				logger.WithError(err).Error("No tasks found")
				c.JSON(http.StatusNoContent, gin.H{"message": "No tasks found"})
				return
			default:
				logger.WithError(err).Error("Failed to list task.")
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to list tasks. Please try again"})
				return
			}
		}
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
	logger := s.Logger.WithField("method", "DeleteTask")
	logger.Debug("Incoming request")
	if !hasPermission(c, []string{"user", "admin"}) {
		c.JSON(http.StatusForbidden, gin.H{"message": "Insufficient permissions"})
		return
	}

	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Task ID is required"})
		return
	}
	md, ok := getGRPCMetadataFromGin(c, logger)
	if !ok {
		return
	}
	ctxWithMetadata := metadata.NewOutgoingContext(context.Background(), md)

	_, err := s.TaskClient.DeleteTask(ctxWithMetadata, &task.DeleteTaskRequest{Id: taskID})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				logger.WithError(err).Error("Task not found")
				c.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
				return
			default:
				logger.WithError(err).Error("Failed to delete task.")
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete task. Please try again"})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func (s *Server) UpdateTask(c *gin.Context) {
	logger := s.Logger.WithField("method", "UpdateTask")
	logger.Debug("Incoming request")
	if !hasPermission(c, []string{"user", "admin"}) {
		c.JSON(http.StatusForbidden, gin.H{"message": "Insufficient permissions"})
		return
	}

	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Task ID is required"})
		return
	}

	var taskReq UpdateTaskRequest
	if err := c.ShouldBindJSON(&taskReq); err != nil {
		logger.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var dueDate *timestamppb.Timestamp
	if taskReq.DueDateTime != "" {
		const layout = "January 2, 2006 3:04 PM MST"
		dueDateTime, err := time.Parse(layout, taskReq.DueDateTime)
		if err != nil {
			logger.WithError(err).Error("Invalid date time format")
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid date time format, please use: January 2, 2006 3:04 PM MST"})
			return
		}
		dueDate = timestamppb.New(dueDateTime)
	}
	md, ok := getGRPCMetadataFromGin(c, logger)
	if !ok {
		return
	}
	ctxWithMetadata := metadata.NewOutgoingContext(context.Background(), md)

	_, err := s.TaskClient.UpdateTask(ctxWithMetadata, &task.UpdateTaskRequest{
		Task: &task.Task{
			Id:          taskID,
			Title:       taskReq.Title,
			Description: taskReq.Description,
			DueDate:     dueDate,
		},
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				logger.WithError(err).Error("Task not found")
				c.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
				return
			default:
				logger.WithError(err).Error("Failed to update task.")
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update task. Please try again"})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully", "task_id": taskID})
}

func (s *Server) CompleteTask(c *gin.Context) {
	logger := s.Logger.WithField("method", "CompleteTask")
	logger.Debug("Incoming request")

	if !hasPermission(c, []string{"user", "admin"}) {
		c.JSON(http.StatusForbidden, gin.H{"message": "Insufficient permissions"})
		return
	}

	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Task ID is required"})
		return
	}

	md, ok := getGRPCMetadataFromGin(c, logger)
	if !ok {
		return
	}
	ctxWithMetadata := metadata.NewOutgoingContext(context.Background(), md)

	getTaskResp, err := s.TaskClient.GetTask(ctxWithMetadata, &task.GetTaskRequest{Id: taskID})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				logger.WithError(err).Error("Task not found")
				c.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
				return
			default:
				logger.WithError(err).Error("Failed to mark task as complete. Please try again")
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
		}
	}

	getTaskResp.Task.Status = "completed"

	_, err = s.TaskClient.UpdateTask(ctxWithMetadata, &task.UpdateTaskRequest{
		Task: getTaskResp.Task,
	})
	if err != nil {
		logger.WithError(err).Error("Failed to mark task as complete. Please try again")
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task marked as completed", "task_id": getTaskResp.Task.Id})
}
