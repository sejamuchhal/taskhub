package server

import (
	"context"

	pb "github.com/sejamuchhal/taskhub/task-service/proto"
	"github.com/sejamuchhal/taskhub/task-service/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Createtask creates a task with given details
func (s *Server) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	if req.GetItem().GetTitle() == "" {
		s.Logger.Warn("CreateTask failed: Title is required")
		return nil, status.Error(codes.InvalidArgument, "Description is required")
	}

	task := &storage.Task{
		Title:       req.GetItem().GetTitle(),
		Description: req.GetItem().GetDescription(),
	}

	err := s.Storage.CreateTask(task)
	if err != nil {
		s.Logger.WithError(err).Error("Could not insert item into the database")
		return nil, status.Errorf(codes.Internal, "Could not insert item into the database: %v", err)
	}

	s.Logger.Info("Task created successfully")
	return &pb.CreateTaskResponse{}, nil
}

func (s *Server) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "API is not implemented yet")
}

func (s *Server) ListTasks(ctx context.Context, req *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "API is not implemented yet")
}

func (s *Server) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "API is not implemented yet")
}

func (s *Server) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "API is not implemented yet")
}
