package server

import (
	"context"
	"time"

	task_pb "github.com/sejamuchhal/taskhub/task/pb/task"
	"github.com/sejamuchhal/taskhub/task/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TransformTask(st *storage.Task) *task_pb.Task {
	return &task_pb.Task{
		Id:          st.ID,
		Title:       st.Title,
		Description: st.Description,
		Status:      st.Status,
		DueDate:     timestamppb.New(st.DueDate),
		CreatedAt:   timestamppb.New(st.CreatedAt.In(time.Local)),
		UpdatedAt:   timestamppb.New(st.UpdatedAt.In(time.Local)),
	}
}

func ExtractUserEmail(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.InvalidArgument, "No metadata present")
	}

	emails, ok := md["email"]
	if !ok || len(emails) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "Invalid metadata")
	}
	return emails[0], nil
}

func ExtractUserID(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.InvalidArgument, "No metadata present")
	}

	userIds, ok := md["user_id"]
	if !ok || len(userIds) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "Invalid metadata")
	}
	return userIds[0], nil
}
