package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	auth "github.com/sejamuchhal/taskhub/gateway/pb/auth"
	task "github.com/sejamuchhal/taskhub/gateway/pb/task"
	"github.com/sejamuchhal/taskhub/gateway/server"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	gomock "go.uber.org/mock/gomock"
)

type ServerTestSuite struct {
	suite.Suite
	mockCtrl    *gomock.Controller
	mockAuth    *auth.MockAuthServiceClient
	mockTask    *task.MockTaskServiceClient
	server      *server.Server
	ginRecorder *httptest.ResponseRecorder
	router      http.Handler
}

func (suite *ServerTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockAuth = auth.NewMockAuthServiceClient(suite.mockCtrl)
	suite.mockTask = task.NewMockTaskServiceClient(suite.mockCtrl)

	suite.server = &server.Server{
		AuthClient: suite.mockAuth,
		TaskClient: suite.mockTask,
		Logger:     logrus.NewEntry(logrus.New()),
	}

	gin.SetMode(gin.TestMode)
	suite.ginRecorder = httptest.NewRecorder()
	suite.router = suite.server.RegisterRoutes()
}

func (suite *ServerTestSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}