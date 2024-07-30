package server_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/sejamuchhal/taskhub/auth/common"
	pb "github.com/sejamuchhal/taskhub/auth/pb"
	"github.com/sejamuchhal/taskhub/auth/server"
	"github.com/sejamuchhal/taskhub/auth/storage"
	"github.com/sejamuchhal/taskhub/auth/storage/mock_storage"
	"github.com/sejamuchhal/taskhub/auth/util"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type ServerTestSuite struct {
	suite.Suite
	Server      *server.Server
	MockCtrl    *gomock.Controller
	MockStorage *mock_storage.MockStorageInterface
	MockRedis   *redismock.ClientMock
}

func (suite *ServerTestSuite) SetupSuite() {
	logger := logrus.NewEntry(logrus.New())
	config := &common.Config{
		JWTSecret:            "test-secret",
		AccessTokenDuration:  time.Hour,
		RefreshTokenDuration: 24 * time.Hour,
	}
	db, mock := redismock.NewClientMock()

	suite.MockCtrl = gomock.NewController(suite.T())
	suite.MockStorage = mock_storage.NewMockStorageInterface(suite.MockCtrl)
	suite.MockRedis = &mock
	suite.Server = &server.Server{
		Storage:      suite.MockStorage,
		TokenHandler: util.NewTokenHandler(config.JWTSecret, db),
		Logger:       logger,
		Config:       config,
	}
}

func (suite *ServerTestSuite) TearDownSuite() {
	suite.MockCtrl.Finish()

}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

// req := &pb.SignupRequest{
// 	Name:     "Harry Potter",
// 	Email:    "harry@hogwarts.edu",
// 	Password: "password",
// }

func (suite *ServerTestSuite) TestSignup_Success() {
	suite.MockStorage.EXPECT().GetUserByEmail("harry@hogwarts.edu").Return(nil, gorm.ErrRecordNotFound)
	suite.MockStorage.EXPECT().CreateUser(gomock.Any()).Return(nil)

	req := &pb.SignupRequest{
		Name:     "Harry Potter",
		Email:    "harry@hogwarts.edu",
		Password: "password",
	}
	log.Printf("req: %v", req)
	resp, err := suite.Server.Signup(context.Background(), req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), "User signup successful", resp.Message)
}

func (suite *ServerTestSuite) TestSignup_UserAlreadyExistError() {
	user := &storage.User{
		ID:    "1",
		Name:  "Harry Potter",
		Email: "harry@hogwarts.edu",
	}
	suite.MockStorage.EXPECT().GetUserByEmail("harry@hogwarts.edu").Return(user, nil)

	req := &pb.SignupRequest{
		Name:     "Harry Potter",
		Email:    "harry@hogwarts.edu",
		Password: "password",
	}
	log.Printf("req: %v", req)
	resp, err := suite.Server.Signup(context.Background(), req)

	assert.Equal(suite.T(), codes.AlreadyExists, status.Code(err))
	assert.Nil(suite.T(), resp)
}

func (suite *ServerTestSuite) TestLogin_Success() {
	hashedPassword, _ := util.HashPassword("password")
	user := &storage.User{
		ID:       "1",
		Name:     "Harry Potter",
		Email:    "harry@hogwarts.edu",
		Password: hashedPassword,
	}

	suite.MockStorage.EXPECT().GetUserByEmail("harry@hogwarts.edu").Return(user, nil)
	suite.MockStorage.EXPECT().CreateSession(gomock.Any()).Return(nil)

	req := &pb.LoginRequest{
		Email:    "harry@hogwarts.edu",
		Password: "password",
	}
	resp, err := suite.Server.Login(context.Background(), req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), "harry@hogwarts.edu", resp.User.Email)
}

func (suite *ServerTestSuite) TestLogin_InvalidPassword() {
	hashedPassword, _ := util.HashPassword("password")
	user := &storage.User{
		ID:       "1",
		Name:     "Harry Potter",
		Email:    "harry@hogwarts.edu",
		Password: hashedPassword,
	}

	suite.MockStorage.EXPECT().GetUserByEmail("harry@hogwarts.edu").Return(user, nil)

	req := &pb.LoginRequest{
		Email:    "harry@hogwarts.edu",
		Password: "wrongpassword",
	}
	resp, err := suite.Server.Login(context.Background(), req)

	assert.Nil(suite.T(), resp)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), codes.Unauthenticated, status.Code(err))
}

func (suite *ServerTestSuite) TestValidate_Success() {
	user := &storage.User{
		ID:    "1",
		Name:  "Harry Potter",
		Email: "harry@hogwarts.edu",
	}
	tokenString, claims, err := suite.Server.TokenHandler.CreateToken(user, time.Hour, "access")
	assert.NoError(suite.T(), err)
	mock := *suite.MockRedis
	mock.ExpectGet(tokenString).SetVal("")
	req := &pb.ValidateRequest{
		Token: tokenString,
	}

	resp, err := suite.Server.Validate(context.Background(), req)

	assert.NotNil(suite.T(), resp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), claims.UserID, resp.UserId)
	assert.Equal(suite.T(), claims.Role, resp.Role)
}

func (suite *ServerTestSuite) TestValidate_InvalidTokenError() {
	token := "invalid-token"
	mock := *suite.MockRedis
	mock.ExpectGet(token).SetVal("")
	req := &pb.ValidateRequest{
		Token: token,
	}

	resp, err := suite.Server.Validate(context.Background(), req)
	log.Printf("resp: %v", resp)
	assert.Nil(suite.T(), resp)
	log.Printf("err: %v", err)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), codes.Unauthenticated, status.Code(err))
}

func (suite *ServerTestSuite) TestRenewAccessToken_Success() {
	user := &storage.User{
		ID:    "1",
		Name:  "Harry Potter",
		Email: "harry@hogwarts.edu",
	}

	refreshToken, refreshClaims, err := suite.Server.TokenHandler.CreateToken(user, 24*time.Hour, "refresh")
	assert.NoError(suite.T(), err)
	mock := *suite.MockRedis
	mock.ExpectGet(refreshToken).SetVal("")
	req := &pb.RenewAccessTokenRequest{
		RefreshToken: refreshToken,
	}

	session := &storage.Session{
		IsBlocked:    false,
		Email:        "harry@hogwarts.edu",
		RefreshToken: refreshToken,
		ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
	}

	suite.MockStorage.EXPECT().GetSessionByID(refreshClaims.RegisteredClaims.ID).Return(session, nil)

	resp, err := suite.Server.RenewAccessToken(context.Background(), req)

	assert.NotNil(suite.T(), resp.AccessToken)
	assert.NoError(suite.T(), err)
}

func (suite *ServerTestSuite) TestRenewAccessToken_SessionBlocked() {
	expectedErr := "rpc error: code = Unauthenticated desc = Session blocked"
	user := &storage.User{
		ID:    "1",
		Name:  "Harry Potter",
		Email: "harry@hogwarts.edu",
	}

	refreshToken, refreshClaims, err := suite.Server.TokenHandler.CreateToken(user, 24*time.Hour, "refresh")
	assert.NoError(suite.T(), err)
	mock := *suite.MockRedis
	mock.ExpectGet(refreshToken).SetVal("")
	req := &pb.RenewAccessTokenRequest{
		RefreshToken: refreshToken,
	}

	session := &storage.Session{
		IsBlocked:    true,
		Email:        "harry@hogwarts.edu",
		RefreshToken: refreshToken,
		ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
	}

	suite.MockStorage.EXPECT().GetSessionByID(refreshClaims.RegisteredClaims.ID).Return(session, nil)

	resp, err := suite.Server.RenewAccessToken(context.Background(), req)

	assert.Nil(suite.T(), resp)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), codes.Unauthenticated, status.Code(err))
	assert.EqualError(suite.T(), err, expectedErr)
}

func (suite *ServerTestSuite) TestLogout_Success() {
	mock := *suite.MockRedis

	user := &storage.User{
		ID:    "1",
		Name:  "Harry Potter",
		Email: "harry@hogwarts.edu",
	}
	accessToken, _, err := suite.Server.TokenHandler.CreateToken(user, time.Hour, "access")
	assert.NoError(suite.T(), err)

	mock.ExpectGet(accessToken).SetVal("")
	mock.ExpectSet(accessToken, "blacklisted", time.Hour).SetVal("OK")

	refreshToken, refreshClaims, err := suite.Server.TokenHandler.CreateToken(user, 24*time.Hour, "refresh")
	assert.NoError(suite.T(), err)

	mock.ExpectGet(refreshToken).SetVal("")

	session := &storage.Session{
		ID: refreshClaims.RegisteredClaims.ID,
		IsBlocked:    false,
		Email:        "harry@hogwarts.edu",
		RefreshToken: refreshToken,
		ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
	}

	suite.MockStorage.EXPECT().GetSessionByID(refreshClaims.RegisteredClaims.ID).Return(session, nil)
	suite.MockStorage.EXPECT().BlockSessionByID(refreshClaims.RegisteredClaims.ID).Return(nil)

	req := &pb.LogoutRequest{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	resp, err := suite.Server.Logout(context.Background(), req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
}
