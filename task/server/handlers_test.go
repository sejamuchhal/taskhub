package server_test

import (
	"testing"

	mock_events "github.com/sejamuchhal/taskhub/task/events/mock_rabbitmq"
	"github.com/sejamuchhal/taskhub/task/server"
	"github.com/sejamuchhal/taskhub/task/storage/mock_storage"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ServerTestSuite struct {
	suite.Suite
	Server       *server.Server
	MockCtrl     *gomock.Controller
	MockStorage  *mock_storage.MockStorageInterface
	MockRabbitMQ *mock_events.MockRabbitMQBrokerInterface
}

func (s ServerTestSuite) ServerSuite() {
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

func (suite *ServerTestSuite) TearDownSuite() {
	suite.MockCtrl.Finish()

}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
