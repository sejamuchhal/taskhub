package server

import (
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	pb "github.com/sejamuchhal/taskhub/gateway/pb/task"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestTransformTask(t *testing.T) {
	type args struct {
		task *pb.Task
	}
	tests := []struct {
		name string
		args args
		want *TaskDetails
	}{
		{
			name: "Basic transformation",
			args: args{
				task: &pb.Task{
					Id:          "123",
					Title:       "Test Task",
					Description: "This is a test task",
					Status:      "created",
					DueDate:     timestamppb.New(time.Date(2024, 7, 28, 0, 0, 0, 0, time.UTC)),
					CreatedAt:   timestamppb.New(time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)),
					UpdatedAt:   timestamppb.New(time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC)),
				},
			},
			want: &TaskDetails{
				Title:       "Test Task",
				Description: "This is a test task",
				Status:      "created",
				DueDate:     "2024-07-28T00:00:00Z",
				CreatedAt:   "2024-07-01T00:00:00Z",
				UpdatedAt:   "2024-07-15T00:00:00Z",
			},
		},
		{
			name: "Partial task",
			args: args{
				task: &pb.Task{
					Id:          "456",
					Title:       "Partial Task",
					Description: "Partially filled task",
					Status:      "In Progress",
					DueDate:     nil,
					CreatedAt:   timestamppb.New(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)),
					UpdatedAt:   nil,
				},
			},
			want: &TaskDetails{
				Title:       "Partial Task",
				Description: "Partially filled task",
				Status:      "In Progress",
				DueDate:     "",
				CreatedAt:   "2024-06-01T00:00:00Z",
				UpdatedAt:   "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TransformTask(tt.args.task); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransformTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hasPermission(t *testing.T) {
	type args struct {
		c               *gin.Context
		accessibleRoles []string
	}

	validCtx := &gin.Context{}
	validCtx.Set("role", "user")

	invalidCtx := &gin.Context{}
	invalidCtx.Set("role", "user")
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid role",
			args: args{
				c:               validCtx,
				accessibleRoles: []string{"user", "admin"},
			},
			want: true,
		},
		{
			name: "Invalid role",
			args: args{
				c:               invalidCtx,
				accessibleRoles: []string{"admin"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasPermission(tt.args.c, tt.args.accessibleRoles); got != tt.want {
				t.Errorf("hasPermission() = %v, want %v", got, tt.want)
			}
		})
	}
}
