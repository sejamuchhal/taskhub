package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/sejamuchhal/taskhub/gateway/pb/auth"
	"github.com/sejamuchhal/taskhub/gateway/pb/task"
	"github.com/sejamuchhal/taskhub/gateway/server"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (suite *ServerTestSuite) TestSignupUser_Success() {
	reqBody := `{"name":"Harry Potter","email":"harry@hogwarts.edu","password":"password"}`
	req := httptest.NewRequest("POST", "/auth/signup", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.mockAuth.EXPECT().Signup(gomock.Any(), &auth.SignupRequest{
		Name:     "Harry Potter",
		Email:    "harry@hogwarts.edu",
		Password: "password",
	}).Return(&auth.SignupResponse{Message: "Signup successful"}, nil)

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.JSONEq(suite.T(), `{"message":"Signup successful"}`, w.Body.String())
}

func (suite *ServerTestSuite) TestSignupUser_InvalidEmail() {
	reqBody := `{"name": "Invalid Email", "email": "invalid-email", "password": "password123"}`
	req := httptest.NewRequest("POST", "/auth/signup", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	assert.JSONEq(suite.T(), `{"message":{"email":"not a valid email address"}}`, w.Body.String())
}

func (suite *ServerTestSuite) TestSignupUser_UserAlreadyExists() {
	reqBody := `{"name": "Harry Potter", "email":"harry@hogwarts.edu","password":"password"}`
	req := httptest.NewRequest("POST", "/auth/signup", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.mockAuth.EXPECT().
		Signup(gomock.Any(), &auth.SignupRequest{
			Name:     "Harry Potter",
			Email:    "harry@hogwarts.edu",
			Password: "password",
		}).
		Return(nil, status.Error(codes.AlreadyExists, "User already exists"))

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusConflict, w.Code)
	assert.JSONEq(suite.T(), `{"message":"User already exists"}`, w.Body.String())
}

func (suite *ServerTestSuite) TestSignupUser_InternalServerError() {
	reqBody := `{"name": "Harry Potter", "email":"harry@hogwarts.edu","password":"password"}`
	req := httptest.NewRequest("POST", "/auth/signup", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.mockAuth.EXPECT().
		Signup(gomock.Any(), &auth.SignupRequest{
			Name:     "Harry Potter",
			Email:    "harry@hogwarts.edu",
			Password: "password",
		}).
		Return(nil, status.Error(codes.Internal, "internal server error"))

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	assert.JSONEq(suite.T(), `{"message":"Internal server error"}`, w.Body.String())
}

func (suite *ServerTestSuite) TestSignupUser_UnknownError() {
	reqBody := `{"name": "Harry Potter", "email":"harry@hogwarts.edu","password":"password"}`
	req := httptest.NewRequest("POST", "/auth/signup", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.mockAuth.EXPECT().
		Signup(gomock.Any(), &auth.SignupRequest{
			Name:     "Harry Potter",
			Email:    "harry@hogwarts.edu",
			Password: "password",
		}).
		Return(nil, errors.New("Unable to process request"))

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	assert.JSONEq(suite.T(), `{"message": "Unknown error"}`, w.Body.String())
}

func (suite *ServerTestSuite) TestLoginUser_Success() {
	reqBody := `{"email":"harry@hogwarts.edu","password":"password"}`
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	expireTime := time.Now().UTC().Add(time.Hour)
	suite.mockAuth.EXPECT().Login(gomock.Any(), &auth.LoginRequest{
		Email:    "harry@hogwarts.edu",
		Password: "password",
	}).Return(&auth.LoginResponse{
		AccessToken:           "access_token",
		AccessTokenExpiresAt:  timestamppb.New(expireTime),
		RefreshToken:          "refresh_token",
		RefreshTokenExpiresAt: timestamppb.New(expireTime),
		SessionId:             "session_id",
		User: &auth.UserDetail{
			Name:  "Harry Potter",
			Email: "harry@hogwarts.edu",
		},
	}, nil)

	suite.router.ServeHTTP(w, req)

	expectedResponse := `{
		"session_id": "session_id",
		"access_token": "access_token",
		"access_token_expires_at": "` + expireTime.Format(time.RFC3339Nano) + `",
		"refresh_token": "refresh_token",
		"refresh_token_expires_at": "` + expireTime.Format(time.RFC3339Nano) + `",
		"user": {
			"name": "Harry Potter",
			"email": "harry@hogwarts.edu"
		}
	}`

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.JSONEq(suite.T(), expectedResponse, w.Body.String())
}

func (suite *ServerTestSuite) TestLoginUser_InvalidCredentials() {
	reqBody := `{"email":"jerry@hogwarts.edu","password":"wrong-password"}`
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.mockAuth.EXPECT().Login(gomock.Any(), &auth.LoginRequest{
		Email:    "jerry@hogwarts.edu",
		Password: "wrong-password",
	}).Return(nil, status.Errorf(codes.Unauthenticated, "Invalid email or password"))

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	assert.JSONEq(suite.T(), `{"message": "Invalid email or password"}`, w.Body.String())
}

func (suite *ServerTestSuite) TestRenewAccessToken_Success() {
	req := httptest.NewRequest("POST", "/auth/renew", nil)
	req.Header.Set("Refresh", "refresh_token")
	w := httptest.NewRecorder()

	expireTime := time.Now().UTC().Add(time.Hour)
	suite.mockAuth.EXPECT().RenewAccessToken(gomock.Any(), &auth.RenewAccessTokenRequest{
		RefreshToken: "refresh_token",
	}).Return(&auth.RenewAccessTokenResponse{
		AccessToken:          "new_access_token",
		AccessTokenExpiresAt: timestamppb.New(expireTime),
	}, nil)

	suite.router.ServeHTTP(w, req)

	expectedResponse := `{
		"access_token": "new_access_token",
		"access_token_expires_at": "` + expireTime.Format(time.RFC3339Nano) + `"
	}`

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.JSONEq(suite.T(), expectedResponse, w.Body.String())
}

func (suite *ServerTestSuite) TestRenewAccessToken_InternalError() {
	req := httptest.NewRequest("POST", "/auth/renew", nil)
	req.Header.Set("Refresh", "refresh_token")
	w := httptest.NewRecorder()

	suite.mockAuth.EXPECT().RenewAccessToken(gomock.Any(), &auth.RenewAccessTokenRequest{
		RefreshToken: "refresh_token",
	}).Return(nil, status.Error(codes.Internal, "Error fethcing session"))

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	assert.JSONEq(suite.T(), `{"message": "Internal server error"}`, w.Body.String())
}

func (suite *ServerTestSuite) TestRenewAccessToken_SessionIsBlocked() {
	req := httptest.NewRequest("POST", "/auth/renew", nil)
	req.Header.Set("Refresh", "refresh_token")
	w := httptest.NewRecorder()

	suite.mockAuth.EXPECT().RenewAccessToken(gomock.Any(), &auth.RenewAccessTokenRequest{
		RefreshToken: "refresh_token",
	}).Return(nil, status.Error(codes.Unauthenticated, "Session revoked"))

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.JSONEq(suite.T(), `{"message": "Invalid refresh token"}`, w.Body.String())
}

func (suite *ServerTestSuite) TestLogout_Success() {
	req := httptest.NewRequest("POST", "/auth/logout", nil)
	req.Header.Set("Refresh", "refresh_token")
	req.Header.Set("Access", "access_token")
	w := httptest.NewRecorder()

	suite.mockAuth.EXPECT().Logout(gomock.Any(), &auth.LogoutRequest{
		RefreshToken: "refresh_token",
		AccessToken:  "access_token",
	}).Return(&auth.LogoutResponse{}, nil)

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.JSONEq(suite.T(), `{"message":"Logout successful"}`, w.Body.String())
}

func (suite *ServerTestSuite) TestLogout_WithoutRefreshToken() {
	req := httptest.NewRequest("POST", "/auth/logout", nil)
	req.Header.Set("Access", "access_token")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnprocessableEntity, w.Code)
	assert.JSONEq(suite.T(), `{"message": "refresh token header is missing"}`, w.Body.String())
}

func (suite *ServerTestSuite) TestLogout_WithoutAccessToken() {
	req := httptest.NewRequest("POST", "/auth/logout", nil)
	req.Header.Set("Refresh", "refresh_token")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnprocessableEntity, w.Code)
	assert.JSONEq(suite.T(), `{"message": "access token header is missing"}`, w.Body.String())
}

func (suite *ServerTestSuite) TestCreateTask_Success() {
	reqBody := `{"title":"New Task","description":"This is a new task","due_date_time":"July 24, 2024 3:04 PM IST"}`
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBufferString(reqBody))
	req.Header.Set("Access", "access_token")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	layout := "January 2, 2006 3:04 PM MST"
	dueDateTime, _ := time.Parse(layout, "July 24, 2024 3:04 PM IST")

	suite.mockAuth.EXPECT().Validate(gomock.Any(), &auth.ValidateRequest{Token: "access_token"}).Return(&auth.ValidateResponse{
		UserId: "test-user-id",
		Email:  "harry@hogwarts.edu",
		Role:   "user",
	}, nil)

	task1 := &task.Task{
		Title:       "New Task",
		Description: "This is a new task",
		DueDate:     timestamppb.New(dueDateTime),
	}

	suite.mockTask.EXPECT().CreateTask(gomock.Any(), &task.CreateTaskRequest{
		Task: task1,
	}).Return(&task.CreateTaskResponse{Id: "12345"}, nil)

	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.JSONEq(suite.T(), `{"task_id": "12345" }`, w.Body.String())
}

func (suite *ServerTestSuite) TestCreateTask_InvalidDateFormat() {
	reqBody := `{"title":"New Task","description":"This is a new task","due_date_time":"30-07-2024 16:00"}`
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBufferString(reqBody))
	req.Header.Set("Access", "access_token")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.mockAuth.EXPECT().Validate(gomock.Any(), &auth.ValidateRequest{Token: "access_token"}).Return(&auth.ValidateResponse{
		UserId: "test-user-id",
		Email:  "harry@hogwarts.edu",
		Role:   "user",
	}, nil)

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	assert.JSONEq(suite.T(), `{"message":"Invalid date time format, please use: January 2, 2006 3:04 PM MST"}`, w.Body.String())
}

func (suite *ServerTestSuite) TestGetTask_Success() {
	req := httptest.NewRequest("GET", "/tasks/12345", nil)
	req.Header.Set("Access", "access_token")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	dueDateTime, _ := time.Parse("January 2, 2006 3:04 PM MST", "July 30, 2024 4:00 PM UTC")

	suite.mockAuth.EXPECT().Validate(gomock.Any(), &auth.ValidateRequest{Token: "access_token"}).Return(&auth.ValidateResponse{
		UserId: "test-user-id",
		Email:  "harry@hogwarts.edu",
		Role:   "user",
	}, nil)

	task1 := &task.Task{
		Id:          "12345",
		Title:       "New Task",
		Description: "This is a new task",
		DueDate:     timestamppb.New(dueDateTime),
		Status:      "pending",
	}

	suite.mockTask.EXPECT().GetTask(gomock.Any(), &task.GetTaskRequest{Id: "12345"}).Return(&task.GetTaskResponse{
		Task: task1,
	}, nil)

	suite.router.ServeHTTP(w, req)

	task2 := server.TransformTask(task1)
	b, _ := json.Marshal(task2)
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.JSONEq(suite.T(), string(b), w.Body.String())
}

func (suite *ServerTestSuite) TestGetTask_NotFound() {
	req := httptest.NewRequest("GET", "/tasks/12345", nil)
	req.Header.Set("Access", "access_token")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.mockAuth.EXPECT().Validate(gomock.Any(), &auth.ValidateRequest{Token: "access_token"}).Return(&auth.ValidateResponse{
		UserId: "test-user-id",
		Email:  "harry@hogwarts.edu",
		Role:   "user",
	}, nil)

	suite.mockTask.EXPECT().GetTask(gomock.Any(), &task.GetTaskRequest{Id: "12345"}).Return(nil, status.Error(codes.NotFound, "Task not found"))

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
	assert.JSONEq(suite.T(), `{"message":"Task not found"}`, w.Body.String())
}

func (suite *ServerTestSuite) TestListTasks_Success() {
	req := httptest.NewRequest("GET", "/tasks?limit=10&offset=0&pending=true", nil)
	req.Header.Set("Access", "access_token")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	dueDateTime, _ := time.Parse("January 2, 2006 3:04 PM MST", "July 30, 2024 4:00 PM UTC")

	suite.mockAuth.EXPECT().Validate(gomock.Any(), &auth.ValidateRequest{Token: "access_token"}).Return(&auth.ValidateResponse{
		UserId: "test-user-id",
		Email:  "harry@hogwarts.edu",
		Role:   "user",
	}, nil)

	tasks := []*task.Task{
		{
			Id:          "12345",
			Title:       "New Task",
			Description: "This is a new task",
			DueDate:     timestamppb.New(dueDateTime),
			Status:      "pending",
		},
	}

	suite.mockTask.EXPECT().ListTasks(gomock.Any(), &task.ListTasksRequest{
		Limit:   10,
		Offset:  0,
		Pending: true,
	}).Return(&task.ListTasksResponse{
		Tasks:      tasks,
		TotalCount: int64(len(tasks)),
	}, nil)

	suite.router.ServeHTTP(w, req)

	taskDetails := make([]*server.TaskDetails, len(tasks))
	for i, t := range tasks {
		taskDetails[i] = server.TransformTask(t)
	}

	response :=server.ListTasksResponse{
		Count: int(len(tasks)),
		Tasks: taskDetails,
	}
	b, _ := json.Marshal(response)
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.JSONEq(suite.T(), string(b), w.Body.String())
}

func (suite *ServerTestSuite) TestListTasks_NoTasksFound() {
	req := httptest.NewRequest("GET", "/tasks?limit=10&offset=0&pending=true", nil)
	req.Header.Set("Access", "access_token")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.mockAuth.EXPECT().Validate(gomock.Any(), &auth.ValidateRequest{Token: "access_token"}).Return(&auth.ValidateResponse{
		UserId: "test-user-id",
		Email:  "harry@hogwarts.edu",
		Role:   "user",
	}, nil)

	suite.mockTask.EXPECT().ListTasks(gomock.Any(), &task.ListTasksRequest{
		Limit:   10,
		Offset:  0,
		Pending: true,
	}).Return(&task.ListTasksResponse{
		Tasks:      []*task.Task{},
		TotalCount: 0,
	}, nil)

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.JSONEq(suite.T(), `{"message":"No tasks found"}`, w.Body.String())
}

func (suite *ServerTestSuite) TestListTasks_InternalServerError() {
	req := httptest.NewRequest("GET", "/tasks?limit=10&offset=0&pending=true", nil)
	req.Header.Set("Access", "access_token")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.mockAuth.EXPECT().Validate(gomock.Any(), &auth.ValidateRequest{Token: "access_token"}).Return(&auth.ValidateResponse{
		UserId: "test-user-id",
		Email:  "harry@hogwarts.edu",
		Role:   "user",
	}, nil)

	suite.mockTask.EXPECT().ListTasks(gomock.Any(), &task.ListTasksRequest{
		Limit:   10,
		Offset:  0,
		Pending: true,
	}).Return(nil, status.Error(codes.Internal, "Internal server error"))

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	assert.JSONEq(suite.T(), `{"message":"Failed to list tasks. Please try again"}`, w.Body.String())
}


func (suite *ServerTestSuite) TestUpdateTask_Success() {
	reqBody := `{"title":"New Task","description":"This is a new task"}`

	req := httptest.NewRequest("PUT", "/tasks/12345", bytes.NewBufferString(reqBody))
	req.Header.Set("Access", "access_token")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.mockAuth.EXPECT().Validate(gomock.Any(), &auth.ValidateRequest{Token: "access_token"}).Return(&auth.ValidateResponse{
		UserId: "test-user-id",
		Email:  "harry@hogwarts.edu",
		Role:   "user",
	}, nil)

	task1 := &task.Task{
		Id: "12345",
		Title:       "New Task",
		Description: "This is a new task",
	}

	md := metadata.New(map[string]string{
		"user_id": "test-user-id",
		"email":   "harry@hogwarts.edu",
	})

	ctxWithMetadata := metadata.NewOutgoingContext(context.Background(), md)

	suite.mockTask.EXPECT().UpdateTask(ctxWithMetadata, &task.UpdateTaskRequest{Task: task1}).Return(&task.UpdateTaskResponse{}, nil)

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.JSONEq(suite.T(), `{"message":"Task updated successfully", "task_id": "12345"}`, w.Body.String())
}

func (suite *ServerTestSuite) TestUpdateTask_NotFound() {
	reqBody := `{"title":"New Task","description":"This is a new task"}`

	req := httptest.NewRequest("PUT", "/tasks/12345", bytes.NewBufferString(reqBody))
	req.Header.Set("Access", "access_token")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.mockAuth.EXPECT().Validate(gomock.Any(), &auth.ValidateRequest{Token: "access_token"}).Return(&auth.ValidateResponse{
		UserId: "test-user-id",
		Email:  "harry@hogwarts.edu",
		Role:   "user",
	}, nil)

	task1 := &task.Task{
		Id: "12345",
		Title:       "New Task",
		Description: "This is a new task",
	}

	md := metadata.New(map[string]string{
		"user_id": "test-user-id",
		"email":   "harry@hogwarts.edu",
	})

	ctxWithMetadata := metadata.NewOutgoingContext(context.Background(), md)

	suite.mockTask.EXPECT().UpdateTask(ctxWithMetadata, &task.UpdateTaskRequest{Task: task1}).Return(nil, status.Error(codes.NotFound, "Task not found"))

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
	assert.JSONEq(suite.T(), `{"message":"Task not found"}`, w.Body.String())
}

func (suite *ServerTestSuite) TestDeleteTask_Success() {

	req := httptest.NewRequest("DELETE", "/tasks/12345", nil)
	req.Header.Set("Access", "access_token")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.mockAuth.EXPECT().Validate(gomock.Any(), &auth.ValidateRequest{Token: "access_token"}).Return(&auth.ValidateResponse{
		UserId: "test-user-id",
		Email:  "harry@hogwarts.edu",
		Role:   "user",
	}, nil)

	md := metadata.New(map[string]string{
		"user_id": "test-user-id",
		"email":   "harry@hogwarts.edu",
	})

	ctxWithMetadata := metadata.NewOutgoingContext(context.Background(), md)

	suite.mockTask.EXPECT().DeleteTask(ctxWithMetadata, &task.DeleteTaskRequest{Id: "12345"}).Return(&task.DeleteTaskResponse{}, nil)

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.JSONEq(suite.T(), `{"message":"Task deleted successfully"}`, w.Body.String())
}

func (suite *ServerTestSuite) TestDeleteTask_NotFound() {

	req := httptest.NewRequest("DELETE", "/tasks/12345", nil)
	req.Header.Set("Access", "access_token")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.mockAuth.EXPECT().Validate(gomock.Any(), &auth.ValidateRequest{Token: "access_token"}).Return(&auth.ValidateResponse{
		UserId: "test-user-id",
		Email:  "harry@hogwarts.edu",
		Role:   "user",
	}, nil)

	md := metadata.New(map[string]string{
		"user_id": "test-user-id",
		"email":   "harry@hogwarts.edu",
	})

	ctxWithMetadata := metadata.NewOutgoingContext(context.Background(), md)

	suite.mockTask.EXPECT().DeleteTask(ctxWithMetadata, &task.DeleteTaskRequest{Id: "12345"}).Return(nil, status.Error(codes.NotFound, "Task not found"))

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
	assert.JSONEq(suite.T(), `{"message":"Task not found"}`, w.Body.String())
}

func (suite *ServerTestSuite) TestCompleteTask_Success() {

	req := httptest.NewRequest("PUT", "/tasks/12345/complete", nil)
	req.Header.Set("Access", "access_token")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.mockAuth.EXPECT().Validate(gomock.Any(), &auth.ValidateRequest{Token: "access_token"}).Return(&auth.ValidateResponse{
		UserId: "test-user-id",
		Email:  "harry@hogwarts.edu",
		Role:   "user",
	}, nil)

	md := metadata.New(map[string]string{
		"user_id": "test-user-id",
		"email":   "harry@hogwarts.edu",
	})

	ctxWithMetadata := metadata.NewOutgoingContext(context.Background(), md)
	dueDateTime, _ := time.Parse("January 2, 2006 3:04 PM MST", "July 30, 2024 4:00 PM UTC")

	task1 := &task.Task{
		Id:          "12345",
		Title:       "New Task",
		Description: "This is a new task",
		DueDate:     timestamppb.New(dueDateTime),
		Status:      "pending",
	}

	suite.mockTask.EXPECT().GetTask(ctxWithMetadata, &task.GetTaskRequest{Id: "12345"}).Return(&task.GetTaskResponse{
		Task: task1,
	}, nil)

	task1.Status = "completed"
	suite.mockTask.EXPECT().UpdateTask(ctxWithMetadata, &task.UpdateTaskRequest{Task: task1}).Return(&task.UpdateTaskResponse{}, nil)


	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.JSONEq(suite.T(), `{"message":"Task marked as completed", "task_id": "12345"}`, w.Body.String())
}


func (suite *ServerTestSuite) TestCompleteTask_NotFound() {

	req := httptest.NewRequest("PUT", "/tasks/12345/complete", nil)
	req.Header.Set("Access", "access_token")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.mockAuth.EXPECT().Validate(gomock.Any(), &auth.ValidateRequest{Token: "access_token"}).Return(&auth.ValidateResponse{
		UserId: "test-user-id",
		Email:  "harry@hogwarts.edu",
		Role:   "user",
	}, nil)

	md := metadata.New(map[string]string{
		"user_id": "test-user-id",
		"email":   "harry@hogwarts.edu",
	})

	ctxWithMetadata := metadata.NewOutgoingContext(context.Background(), md)

	suite.mockTask.EXPECT().GetTask(ctxWithMetadata, &task.GetTaskRequest{Id: "12345"}).Return(nil, status.Error(codes.NotFound, "Task not found"))

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
	assert.JSONEq(suite.T(), `{"message":"Task not found"}`, w.Body.String())
}

func (suite *ServerTestSuite) TestPrometheusHandler_Success() {
	req, _ := http.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *ServerTestSuite) TestHealthCheck_Success() {
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.JSONEq(suite.T(), `{"message":"It's healthy"}`, w.Body.String())
}
