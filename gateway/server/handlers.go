package server

import (
	"context"
	"net/http"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/gin-gonic/gin"
	"github.com/sejamuchhal/taskhub/gateway/pb/auth"
	"github.com/sejamuchhal/taskhub/gateway/pb/task"
)

func (s *Server) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "It's healthy"})
}

func (s *Server) SignupUser(c *gin.Context) {
    logger := s.Logger.WithField("method", "SignupUser")
    logger.Debug("Incoming request")

    var req SignupUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        logger.WithError(err).Error("Error parsing signup request payload")
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
                c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
                return
            case codes.InvalidArgument:
                logger.WithError(err).Error("Invalid request")
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
                return
            default:
                logger.WithError(err).Error("Internal server error")
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
                return
            }
        }
        logger.WithError(err).Error("Unknown error")
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unknown error"})
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
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
                return
            case codes.InvalidArgument:
                logger.WithError(err).Error("Invalid request")
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
                return
            default:
                logger.WithError(err).Error("Internal server error")
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
                return
            }
        }
        logger.WithError(err).Error("Unknown error")
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unknown error"})
        return
    }

    logger.Info("User login successful")
    c.JSON(http.StatusOK, LoginUserResponse{
        AccessToken: res.Token,
        User: UserDetail{
            Name:  res.User.Name,
            Email: res.User.Email,
        },
    })
}

func (s *Server) CreateTask(c *gin.Context) {
	logger := s.Logger.WithField("method", "CreateTask")
	logger.Debug("Incoming request")

	var taskReq CreateTaskRequest
	if err := c.ShouldBindJSON(&taskReq); err != nil {
		logger.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	layout := "January 2, 2006 3:04 PM MST"
	dueDateTime, err := time.Parse(layout, taskReq.DueDateTime)
	if err != nil {
		logger.WithError(err).Error("Invalid date time format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date time format, please use: January 2, 2006 3:04 PM MST"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		logger.Error("User ID not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		logger.Error("Invalid user ID type")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	logger.Debug("Sending CreateTask gRPC request")
	resp, err := s.TaskClient.CreateTask(context.Background(), &task.CreateTaskRequest{
		Task: &task.Task{
			Title:       taskReq.Title,
			Description: taskReq.Description,
			UserId:      userIDStr,
			DueDate:     timestamppb.New(dueDateTime),
		},
	})
	if err != nil {
		logger.WithError(err).Error("Failed to create task via gRPC")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.WithField("response", resp).Debug("Received response from gRPC")
	c.JSON(http.StatusOK, gin.H{"task_id": resp.Id})
}

func (s *Server) GetTask(c *gin.Context) {
	logger := s.Logger.WithField("method", "GetTask")
	logger.Debug("Incoming request")

	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	resp, err := s.TaskClient.GetTask(context.Background(), &task.GetTaskRequest{Id: taskID})
	if err != nil {
		logger.WithError(err).Error("Failed to get task via gRPC")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	taskDetails := TransformTask(resp.Task)
	c.JSON(http.StatusOK, taskDetails)
}

func (s *Server) ListTasks(c *gin.Context) {
	logger := s.Logger.WithField("method", "ListTasks")
	logger.Debug("Incoming request")

	var req ListTasksRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.WithError(err).Error("Failed to bind query parameters")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		logger.Error("User ID not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		logger.Error("Invalid user ID type")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	resp, err := s.TaskClient.ListTasks(context.Background(), &task.ListTasksRequest{
		Limit:   int32(req.Limit),
		Offset:  int32(req.Offset),
		UserId:  userIDStr,
		Pending: req.Pending,
	})
	if err != nil {
		logger.WithError(err).Error("Failed to list tasks via gRPC")
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
	logger := s.Logger.WithField("method", "DeleteTask")
	logger.Debug("Incoming request")

	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	_, err := s.TaskClient.DeleteTask(context.Background(), &task.DeleteTaskRequest{Id: taskID})
	if err != nil {
		logger.WithError(err).Error("Failed to delete task via gRPC")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func (s *Server) UpdateTask(c *gin.Context) {
	logger := s.Logger.WithField("method", "UpdateTask")
	logger.Debug("Incoming request")

	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	var taskReq UpdateTaskRequest
	if err := c.ShouldBindJSON(&taskReq); err != nil {
		logger.WithError(err).Error("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	layout := "January 2, 2006 3:04 PM MST"
	dueDateTime, err := time.Parse(layout, taskReq.DueDateTime)
	if err != nil {
		logger.WithError(err).Error("Invalid date time format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date time format, please use: January 2, 2006 3:04 PM MST"})
		return
	}

	_, err = s.TaskClient.UpdateTask(context.Background(), &task.UpdateTaskRequest{
		Task: &task.Task{
			Id:          taskID,
			Title:       taskReq.Title,
			Description: taskReq.Description,
			DueDate:     timestamppb.New(dueDateTime),
		},
	})
	if err != nil {
		logger.WithError(err).Error("Failed to update task via gRPC")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully", "task_id": taskID})
}

func (s *Server) CompleteTask(c *gin.Context) {
	logger := s.Logger.WithField("method", "CompleteTask")
	logger.Debug("Incoming request")

	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	getTaskResp, err := s.TaskClient.GetTask(context.Background(), &task.GetTaskRequest{Id: taskID})
	if err != nil {
		logger.WithError(err).Error("Failed to get task via gRPC")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	getTaskResp.Task.Status = "completed"

	_, err = s.TaskClient.UpdateTask(context.Background(), &task.UpdateTaskRequest{
		Task: getTaskResp.Task,
	})
	if err != nil {
		logger.WithError(err).Error("Failed to update task via gRPC")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task marked as completed", "task_id": getTaskResp.Task.Id})
}
