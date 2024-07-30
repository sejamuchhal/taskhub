package worker

import (
	"testing"

	event "github.com/sejamuchhal/taskhub/notification/pb"
)

func TestPrepareEmailContent(t *testing.T) {
	type args struct {
		message event.TaskUpdateEvent
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name: "Basic transformation",
			args: args{
				message: event.TaskUpdateEvent{
					Title:  "Test Task",
					Status: "created",
					Email:  "harry@hogwarts.edu",
				},
			},
			want:  "Test Task: created",
			want1: "Task Title: Test Task\nStatus: created",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := prepareEmailContent(tt.args.message)
			if got != tt.want {
				t.Errorf("prepareEmailContent() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("prepareEmailContent() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
