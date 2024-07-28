package server

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	event_pb "github.com/sejamuchhal/taskhub/task/pb/event"
	task_pb "github.com/sejamuchhal/taskhub/task/pb/task"
	"github.com/sejamuchhal/taskhub/task/storage"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// Createtask creates a task with given details
func (s *Server) CreateTask(ctx context.Context, req *task_pb.CreateTaskRequest) (*task_pb.CreateTaskResponse, error) {
	logger := s.Logger.WithFields(logrus.Fields{
		"req":    req,
		"method": "CreateTask",
	})

	logger.Info("Received CreateTask request")

	email, err := ExtractUserEmail(ctx)
	if err != nil {
		return nil, err
	}

	userID, err := ExtractUserID(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetTask().GetTitle() == "" {
		s.Logger.Error("CreateTask failed: Title is required")
		return nil, status.Error(codes.InvalidArgument, "Title is required")
	}

	taskID := uuid.NewString()

	storageTask := &storage.Task{
		ID:          taskID,
		UserID:      userID,
		Title:       req.GetTask().GetTitle(),
		Description: req.GetTask().GetDescription(),
		DueDate:     req.GetTask().GetDueDate().AsTime(),
	}

	err = s.Storage.CreateTask(storageTask)
	if err != nil {
		s.Logger.WithError(err).Error("Could not insert Task into the database")
		return nil, status.Errorf(codes.Internal, "Could not insert Task into the database: %v", err)
	}

	s.Logger.WithFields(logrus.Fields{
		"task_id": storageTask.ID,
		"title":   storageTask.Title,
	}).Info("Task successfully inserted into the database")

	event := event_pb.TaskUpdateEvent{
		Status: "created",
		Title:  storageTask.Title,
		Email:  email,
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to marshal task update event to JSON")
		return &task_pb.CreateTaskResponse{}, status.Errorf(codes.Internal, "Failed to marshal event: %v", err)
	}

	err = s.Publisher.Publish(eventJSON)
	if err != nil {
		s.Logger.WithError(err).Error("failed to send task update event")
	} else {
		s.Logger.Info("Task update event sent to queue")
	}

	s.Logger.Info("Task created successfully")
	return &task_pb.CreateTaskResponse{Id: taskID}, nil
}

func (s *Server) GetTask(ctx context.Context, req *task_pb.GetTaskRequest) (*task_pb.GetTaskResponse, error) {
	logger := s.Logger.WithFields(logrus.Fields{
		"req":    req,
		"method": "GetTask",
	})

	logger.Info("Received GetTask request")
	taskID := req.GetId()
	if taskID == "" {
		errMsg := "GetTask failed: Task ID is required"
		logger.Error(errMsg)
		return nil, status.Error(codes.InvalidArgument, errMsg)
	}

	task, err := s.Storage.GetTaskByID(taskID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			errMsg := "Task not found"
			logger.WithError(err).Error(errMsg)
			return nil, status.Errorf(codes.NotFound, "%s: %v", errMsg, err)
		}
		errMsg := "Failed to retrieve Task from the database"
		logger.WithError(err).Error(errMsg)
		return nil, status.Errorf(codes.Internal, "%s: %v", errMsg, err)
	}

	res := &task_pb.GetTaskResponse{
		Task: TransformTask(task),
	}
	logger.Info("GetTask request processed successfully")
	return res, nil
}

func (s *Server) ListTasks(ctx context.Context, req *task_pb.ListTasksRequest) (*task_pb.ListTasksResponse, error) {
	logger := s.Logger.WithFields(logrus.Fields{
		"req":    req,
		"method": "ListTasks",
	})

	logger.Info("Received ListTasks request")

	// Set default values if not provided
	if req.Offset == 0 {
		req.Offset = 0
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	userID, err := ExtractUserID(ctx)
	if err != nil {
		return nil, err
	}

	tasks, count, err := s.Storage.ListTasksWithCount(userID, int(req.Limit), int(req.Offset))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			errMsg := "No tasks found for the given user ID"
			logger.WithError(err).Error(errMsg)
			return nil, status.Errorf(codes.NotFound, "%s: %v", errMsg, err)
		}
		errMsg := "Failed to list tasks from the database"
		logger.WithError(err).Error(errMsg)
		return nil, status.Errorf(codes.Internal, "%s: %v", errMsg, err)
	}

	pbTasks := make([]*task_pb.Task, len(tasks))
	for i, t := range tasks {
		pbTasks[i] = TransformTask(t)
	}

	logger.WithFields(logrus.Fields{
		"task_count":  len(pbTasks),
		"total_count": count,
	}).Info("Successfully retrieved tasks")

	return &task_pb.ListTasksResponse{
		Tasks:      pbTasks,
		TotalCount: count,
	}, nil
}

func (s *Server) DeleteTask(ctx context.Context, req *task_pb.DeleteTaskRequest) (*task_pb.DeleteTaskResponse, error) {
	logger := s.Logger.WithFields(logrus.Fields{
		"req":    req,
		"method": "DeleteTask",
	})

	logger.Info("Received DeleteTask request")
	email, err := ExtractUserEmail(ctx)
	if err != nil {
		return nil, err
	}

	task, err := s.Storage.GetTaskByID(req.GetId())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			errMsg := "Task not found"
			logger.WithError(err).Error(errMsg)
			return nil, status.Errorf(codes.NotFound, "%s: %v", errMsg, err)
		}
		errMsg := "Failed to retrieve Task from the database"
		logger.WithError(err).Error(errMsg)
		return nil, status.Errorf(codes.Internal, "%s: %v", errMsg, err)
	}

	err = s.Storage.DeleteTask(req.GetId())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			s.Logger.WithError(err).Error("Task not found")
			return nil, status.Errorf(codes.NotFound, "Task not found: %v", err)
		}
		s.Logger.WithError(err).Error("Could not delete Task from the database")
		return nil, status.Errorf(codes.Internal, "Could not delete Task from the database: %v", err)
	}

	s.Logger.WithFields(logrus.Fields{
		"task_id": req.GetId(),
	}).Info("Task successfully deleted from the database")

	event := event_pb.TaskUpdateEvent{
		Status: "deleted",
		Title:  task.Title,
		Email:  email,
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to marshal task update event to JSON")
		return &task_pb.DeleteTaskResponse{}, status.Errorf(codes.Internal, "Failed to marshal event: %v", err)
	}

	err = s.Publisher.Publish(eventJSON)
	if err != nil {
		s.Logger.WithError(err).Error("failed to send task update event")
	} else {
		s.Logger.Info("Task update event sent to queue")
	}

	s.Logger.Info("Task deleted successfully")
	return &task_pb.DeleteTaskResponse{}, nil
}

func (s *Server) UpdateTask(ctx context.Context, req *task_pb.UpdateTaskRequest) (*task_pb.UpdateTaskResponse, error) {
	logger := s.Logger.WithFields(logrus.Fields{
		"req":    req,
		"method": "UpdateTask",
	})

	logger.Info("Received UpdateTask request")

	email, err := ExtractUserEmail(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetTask().GetTitle() == "" {
		s.Logger.Error("UpdateTask failed: Title is required")
		return nil, status.Error(codes.InvalidArgument, "Title is required")
	}

	task, err := s.Storage.GetTaskByID(req.GetTask().GetId())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			errMsg := "Task not found"
			logger.WithError(err).Error(errMsg)
			return nil, status.Errorf(codes.NotFound, "%s: %v", errMsg, err)
		}
		errMsg := "Failed to retrieve Task from the database"
		logger.WithError(err).Error(errMsg)
		return nil, status.Errorf(codes.Internal, "%s: %v", errMsg, err)
	}

	task.Title = req.Task.Title
	task.Description = req.Task.Description
	task.DueDate = req.Task.DueDate.AsTime()

	err = s.Storage.UpdateTask(task)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			s.Logger.WithError(err).Error("Task not found")
			return nil, status.Errorf(codes.NotFound, "Task not found: %v", err)
		}
		s.Logger.WithError(err).Error("Could not update Task in the database")
		return nil, status.Errorf(codes.Internal, "Could not update Task in the database: %v", err)
	}

	s.Logger.WithFields(logrus.Fields{
		"task_id": task.ID,
		"title":   task.Title,
	}).Info("Task successfully updated in the database")

	event := event_pb.TaskUpdateEvent{
		Status: "updated",
		Title:  task.Title,
		Email:  email,
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to marshal task update event to JSON")
		return &task_pb.UpdateTaskResponse{}, status.Errorf(codes.Internal, "Failed to marshal event: %v", err)
	}

	err = s.Publisher.Publish(eventJSON)
	if err != nil {
		s.Logger.WithError(err).Error("failed to send task update event")
	} else {
		s.Logger.Info("Task update event sent to queue")
	}

	s.Logger.Info("Task updated successfully")
	return &task_pb.UpdateTaskResponse{}, nil
}
