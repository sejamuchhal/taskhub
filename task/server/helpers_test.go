package server

import (
	"reflect"
	"testing"
	"time"

	task_pb "github.com/sejamuchhal/taskhub/task/pb/task"
	"github.com/sejamuchhal/taskhub/task/storage"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestTransformTask(t *testing.T) {
	type args struct {
		st *storage.Task
	}
	tests := []struct {
		name string
		args args
		want *task_pb.Task
	}{
		{
			name: "Basic transformation",
			args: args{
				st: &storage.Task{
					ID:          "123",
					Title:       "Test Task",
					Description: "This is a test task",
					Status:      "completed",
					DueDate:     time.Date(2024, 7, 28, 0, 0, 0, 0, time.UTC),
					CreatedAt:   time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC),
				},
			},
			want: &task_pb.Task{
				Id:          "123",
				Title:       "Test task",
				Description: "This is a test task",
				Status:      "completed",
				DueDate: timestamppb.New(time.Date(2024, 7, 28, 0, 0, 0, 0, time.UTC).In(time.Local)),
				CreatedAt: timestamppb.New(time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC).In(time.Local)),
				UpdatedAt: timestamppb.New(time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC).In(time.Local)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TransformTask(tt.args.st); !reflect.DeepEqual(got.Id, tt.want.Id) {
				t.Errorf("TransformTask() = %v, want %v", got, tt.want)
			}
		})
	}
}
