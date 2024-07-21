package server

import (
	"time"

	task_pb "github.com/sejamuchhal/taskhub/task-service/pb/task"
	"github.com/sejamuchhal/taskhub/task-service/storage"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TransformTask(st *storage.Task) *task_pb.Task {
	return &task_pb.Task{
		Id:          st.ID,
		Title:       st.Title,
		Description: st.Description,
		Status:      st.Status,
		UserId:      st.UserID,
		DueDate:     timestamppb.New(st.DueDate),
		CreatedAt:   timestamppb.New(st.CreatedAt.In(time.Local)),
		UpdatedAt:   timestamppb.New(st.UpdatedAt.In(time.Local)),
	}
}
