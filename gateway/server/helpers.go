package server

import (
	"time"

	pb "github.com/sejamuchhal/taskhub/gateway/client/task"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TransformTask(task *pb.Task) *TaskDetails {
	// Convert protobuf Timestamp to a formatted string
	formatTimestamp := func(ts *timestamppb.Timestamp) string {
		if ts == nil {
			return ""
		}
		t := ts.AsTime()
		return t.Format(time.RFC3339)
	}

	return &TaskDetails{
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		DueDate:     formatTimestamp(task.DueDate),
		CreatedAt:   formatTimestamp(task.CreatedAt),
		UpdatedAt:   formatTimestamp(task.UpdatedAt),
	}
}
