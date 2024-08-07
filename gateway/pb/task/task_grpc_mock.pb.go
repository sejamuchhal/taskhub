// Code generated by protoc-gen-go-grpc-mock. DO NOT EDIT.
// source: task.proto

package task

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockTaskServiceClient is a mock of TaskServiceClient interface.
type MockTaskServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockTaskServiceClientMockRecorder
}

// MockTaskServiceClientMockRecorder is the mock recorder for MockTaskServiceClient.
type MockTaskServiceClientMockRecorder struct {
	mock *MockTaskServiceClient
}

// NewMockTaskServiceClient creates a new mock instance.
func NewMockTaskServiceClient(ctrl *gomock.Controller) *MockTaskServiceClient {
	mock := &MockTaskServiceClient{ctrl: ctrl}
	mock.recorder = &MockTaskServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskServiceClient) EXPECT() *MockTaskServiceClientMockRecorder {
	return m.recorder
}

// CreateTask mocks base method.
func (m *MockTaskServiceClient) CreateTask(ctx context.Context, in *CreateTaskRequest, opts ...grpc.CallOption) (*CreateTaskResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateTask", varargs...)
	ret0, _ := ret[0].(*CreateTaskResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTask indicates an expected call of CreateTask.
func (mr *MockTaskServiceClientMockRecorder) CreateTask(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTask", reflect.TypeOf((*MockTaskServiceClient)(nil).CreateTask), varargs...)
}

// DeleteTask mocks base method.
func (m *MockTaskServiceClient) DeleteTask(ctx context.Context, in *DeleteTaskRequest, opts ...grpc.CallOption) (*DeleteTaskResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteTask", varargs...)
	ret0, _ := ret[0].(*DeleteTaskResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteTask indicates an expected call of DeleteTask.
func (mr *MockTaskServiceClientMockRecorder) DeleteTask(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTask", reflect.TypeOf((*MockTaskServiceClient)(nil).DeleteTask), varargs...)
}

// GetTask mocks base method.
func (m *MockTaskServiceClient) GetTask(ctx context.Context, in *GetTaskRequest, opts ...grpc.CallOption) (*GetTaskResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetTask", varargs...)
	ret0, _ := ret[0].(*GetTaskResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTask indicates an expected call of GetTask.
func (mr *MockTaskServiceClientMockRecorder) GetTask(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTask", reflect.TypeOf((*MockTaskServiceClient)(nil).GetTask), varargs...)
}

// ListTasks mocks base method.
func (m *MockTaskServiceClient) ListTasks(ctx context.Context, in *ListTasksRequest, opts ...grpc.CallOption) (*ListTasksResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListTasks", varargs...)
	ret0, _ := ret[0].(*ListTasksResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTasks indicates an expected call of ListTasks.
func (mr *MockTaskServiceClientMockRecorder) ListTasks(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTasks", reflect.TypeOf((*MockTaskServiceClient)(nil).ListTasks), varargs...)
}

// UpdateTask mocks base method.
func (m *MockTaskServiceClient) UpdateTask(ctx context.Context, in *UpdateTaskRequest, opts ...grpc.CallOption) (*UpdateTaskResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateTask", varargs...)
	ret0, _ := ret[0].(*UpdateTaskResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTask indicates an expected call of UpdateTask.
func (mr *MockTaskServiceClientMockRecorder) UpdateTask(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTask", reflect.TypeOf((*MockTaskServiceClient)(nil).UpdateTask), varargs...)
}

// MockTaskServiceServer is a mock of TaskServiceServer interface.
type MockTaskServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockTaskServiceServerMockRecorder
}

// MockTaskServiceServerMockRecorder is the mock recorder for MockTaskServiceServer.
type MockTaskServiceServerMockRecorder struct {
	mock *MockTaskServiceServer
}

// NewMockTaskServiceServer creates a new mock instance.
func NewMockTaskServiceServer(ctrl *gomock.Controller) *MockTaskServiceServer {
	mock := &MockTaskServiceServer{ctrl: ctrl}
	mock.recorder = &MockTaskServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskServiceServer) EXPECT() *MockTaskServiceServerMockRecorder {
	return m.recorder
}

// CreateTask mocks base method.
func (m *MockTaskServiceServer) CreateTask(ctx context.Context, in *CreateTaskRequest) (*CreateTaskResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTask", ctx, in)
	ret0, _ := ret[0].(*CreateTaskResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTask indicates an expected call of CreateTask.
func (mr *MockTaskServiceServerMockRecorder) CreateTask(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTask", reflect.TypeOf((*MockTaskServiceServer)(nil).CreateTask), ctx, in)
}

// DeleteTask mocks base method.
func (m *MockTaskServiceServer) DeleteTask(ctx context.Context, in *DeleteTaskRequest) (*DeleteTaskResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTask", ctx, in)
	ret0, _ := ret[0].(*DeleteTaskResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteTask indicates an expected call of DeleteTask.
func (mr *MockTaskServiceServerMockRecorder) DeleteTask(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTask", reflect.TypeOf((*MockTaskServiceServer)(nil).DeleteTask), ctx, in)
}

// GetTask mocks base method.
func (m *MockTaskServiceServer) GetTask(ctx context.Context, in *GetTaskRequest) (*GetTaskResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTask", ctx, in)
	ret0, _ := ret[0].(*GetTaskResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTask indicates an expected call of GetTask.
func (mr *MockTaskServiceServerMockRecorder) GetTask(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTask", reflect.TypeOf((*MockTaskServiceServer)(nil).GetTask), ctx, in)
}

// ListTasks mocks base method.
func (m *MockTaskServiceServer) ListTasks(ctx context.Context, in *ListTasksRequest) (*ListTasksResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTasks", ctx, in)
	ret0, _ := ret[0].(*ListTasksResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTasks indicates an expected call of ListTasks.
func (mr *MockTaskServiceServerMockRecorder) ListTasks(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTasks", reflect.TypeOf((*MockTaskServiceServer)(nil).ListTasks), ctx, in)
}

// UpdateTask mocks base method.
func (m *MockTaskServiceServer) UpdateTask(ctx context.Context, in *UpdateTaskRequest) (*UpdateTaskResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTask", ctx, in)
	ret0, _ := ret[0].(*UpdateTaskResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTask indicates an expected call of UpdateTask.
func (mr *MockTaskServiceServerMockRecorder) UpdateTask(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTask", reflect.TypeOf((*MockTaskServiceServer)(nil).UpdateTask), ctx, in)
}
