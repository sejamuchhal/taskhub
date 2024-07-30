package server_test

import (
	"context"
	"errors"
	"testing"
	"time"

	mock_events "github.com/sejamuchhal/taskhub/task/events/mock_rabbitmq"
	// event_pb "github.com/sejamuchhal/taskhub/task/pb/event"
	task_pb "github.com/sejamuchhal/taskhub/task/pb/task"
	"github.com/sejamuchhal/taskhub/task/server"
	"github.com/sejamuchhal/taskhub/task/storage/mock_storage"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ServerTestSuite struct {
	suite.Suite
	Server       *server.Server
	MockCtrl     *gomock.Controller
	MockStorage  *mock_storage.MockStorageInterface
	MockRabbitMQ *mock_events.MockRabbitMQBrokerInterface
}

func (s *ServerTestSuite) SetupSuite() {
	logger := logrus.NewEntry(logrus.New())
	s.MockCtrl = gomock.NewController(s.T())
	s.MockStorage = mock_storage.NewMockStorageInterface(s.MockCtrl)
	s.MockRabbitMQ = mock_events.NewMockRabbitMQBrokerInterface(s.MockCtrl)

	s.Server = &server.Server{
		Publisher: s.MockRabbitMQ,
		Storage:   s.MockStorage,
		Logger:    logger,
	}
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}


func (s *ServerTestSuite) TearDownSuite() {
	s.MockCtrl.Finish()

}

func createTestContext() context.Context {
	md := metadata.New(map[string]string{
		"email":   "user@example.com",
		"user_id": "1234",
	})
	return metadata.NewIncomingContext(context.Background(), md)
}

func (s *ServerTestSuite) TestCreateTask_Success() {
	ctx := createTestContext()

	task := &task_pb.Task{
		Title:       "Test Task",
		Description: "This is a test task",
		DueDate:     timestamppb.New(time.Now().Add(24 * time.Hour)),
	}

	taskRequest := &task_pb.CreateTaskRequest{
		Task: task,
	}

	s.Server.Logger.WithField("task_request", taskRequest).Info("Create task request")

	s.MockStorage.EXPECT().CreateTask(gomock.Any()).Return(nil)
	s.MockRabbitMQ.EXPECT().Publish(gomock.Any()).Return(nil)

	resp, err := s.Server.CreateTask(ctx, taskRequest)
	s.NoError(err)
	s.NotNil(resp)
	s.NotEmpty(resp.Id)
}

func (s *ServerTestSuite) TestCreateTask_MissingTitle() {
	ctx := createTestContext()

	taskRequest := &task_pb.CreateTaskRequest{
		Task: &task_pb.Task{
			Description: "This is a test task",
			DueDate:     timestamppb.New(time.Now().Add(24 * time.Hour)),
		},
	}

	resp, err := s.Server.CreateTask(ctx, taskRequest)
	s.Nil(resp)
	s.Error(err)
	s.Equal(codes.InvalidArgument, status.Code(err))
	s.Equal("Title is required", status.Convert(err).Message())
}

func (s *ServerTestSuite) TestCreateTask_DatabaseError() {
	ctx := createTestContext()

	taskRequest := &task_pb.CreateTaskRequest{
		Task: &task_pb.Task{
			Title:       "Test Task",
			Description: "This is a test task",
			DueDate:     timestamppb.New(time.Now().Add(24 * time.Hour)),
		},
	}

	s.MockStorage.EXPECT().CreateTask(gomock.Any()).Return(errors.New("database error"))

	resp, err := s.Server.CreateTask(ctx, taskRequest)
	s.Nil(resp)
	s.Error(err)
	s.Equal(codes.Internal, status.Code(err))
	s.Contains(status.Convert(err).Message(), "Could not insert Task into the database")
}

func (s *ServerTestSuite) TestCreateTask_EventPublishError() {
	ctx := createTestContext()

	taskRequest := &task_pb.CreateTaskRequest{
		Task: &task_pb.Task{
			Title:       "Test Task",
			Description: "This is a test task",
			DueDate:     timestamppb.New(time.Now().Add(24 * time.Hour)),
		},
	}

	s.MockStorage.EXPECT().CreateTask(gomock.Any()).Return(nil)
	s.MockRabbitMQ.EXPECT().Publish(gomock.Any()).Return(errors.New("publish error"))

	resp, err := s.Server.CreateTask(ctx, taskRequest)
	s.NoError(err)
	s.NotNil(resp)
	s.NotEmpty(resp.Id)
}

