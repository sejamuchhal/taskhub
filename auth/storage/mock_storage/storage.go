// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/sejamuchhal/taskhub/auth/storage (interfaces: StorageInterface)
//
// Generated by this command:
//
//	mockgen --build_flags=--mod=mod --destination=./auth/storage/mock_storage/storage.go github.com/sejamuchhal/taskhub/auth/storage StorageInterface
//

// Package mock_storage is a generated GoMock package.
package mock_storage

import (
	reflect "reflect"

	storage "github.com/sejamuchhal/taskhub/auth/storage"
	gomock "go.uber.org/mock/gomock"
)

// MockStorageInterface is a mock of StorageInterface interface.
type MockStorageInterface struct {
	ctrl     *gomock.Controller
	recorder *MockStorageInterfaceMockRecorder
}

// MockStorageInterfaceMockRecorder is the mock recorder for MockStorageInterface.
type MockStorageInterfaceMockRecorder struct {
	mock *MockStorageInterface
}

// NewMockStorageInterface creates a new mock instance.
func NewMockStorageInterface(ctrl *gomock.Controller) *MockStorageInterface {
	mock := &MockStorageInterface{ctrl: ctrl}
	mock.recorder = &MockStorageInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorageInterface) EXPECT() *MockStorageInterfaceMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockStorageInterface) CreateUser(arg0 *storage.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockStorageInterfaceMockRecorder) CreateUser(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockStorageInterface)(nil).CreateUser), arg0)
}

// GetSessionByID mocks base method.
func (m *MockStorageInterface) GetSessionByID(arg0 string) (*storage.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSessionByID", arg0)
	ret0, _ := ret[0].(*storage.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSessionByID indicates an expected call of GetSessionByID.
func (mr *MockStorageInterfaceMockRecorder) GetSessionByID(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSessionByID", reflect.TypeOf((*MockStorageInterface)(nil).GetSessionByID), arg0)
}

// GetUserByEmail mocks base method.
func (m *MockStorageInterface) GetUserByEmail(arg0 string) (*storage.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", arg0)
	ret0, _ := ret[0].(*storage.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockStorageInterfaceMockRecorder) GetUserByEmail(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockStorageInterface)(nil).GetUserByEmail), arg0)
}

// UpdateSession mocks base method.
func (m *MockStorageInterface) UpdateSession(arg0 *storage.Session) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSession", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSession indicates an expected call of UpdateSession.
func (mr *MockStorageInterfaceMockRecorder) UpdateSession(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSession", reflect.TypeOf((*MockStorageInterface)(nil).UpdateSession), arg0)
}
