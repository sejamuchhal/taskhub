// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/sejamuchhal/taskhub/notification/rabbitmq (interfaces: RabbitMQBrokerInterface)

// Package mock_rabbitmq is a generated GoMock package.
package mock_rabbitmq

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockRabbitMQBrokerInterface is a mock of RabbitMQBrokerInterface interface.
type MockRabbitMQBrokerInterface struct {
	ctrl     *gomock.Controller
	recorder *MockRabbitMQBrokerInterfaceMockRecorder
}

// MockRabbitMQBrokerInterfaceMockRecorder is the mock recorder for MockRabbitMQBrokerInterface.
type MockRabbitMQBrokerInterfaceMockRecorder struct {
	mock *MockRabbitMQBrokerInterface
}

// NewMockRabbitMQBrokerInterface creates a new mock instance.
func NewMockRabbitMQBrokerInterface(ctrl *gomock.Controller) *MockRabbitMQBrokerInterface {
	mock := &MockRabbitMQBrokerInterface{ctrl: ctrl}
	mock.recorder = &MockRabbitMQBrokerInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRabbitMQBrokerInterface) EXPECT() *MockRabbitMQBrokerInterfaceMockRecorder {
	return m.recorder
}

// Consume mocks base method.
func (m *MockRabbitMQBrokerInterface) Consume() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Consume")
}

// Consume indicates an expected call of Consume.
func (mr *MockRabbitMQBrokerInterfaceMockRecorder) Consume() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Consume", reflect.TypeOf((*MockRabbitMQBrokerInterface)(nil).Consume))
}

// OnError mocks base method.
func (m *MockRabbitMQBrokerInterface) OnError(arg0 error, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnError", arg0, arg1)
}

// OnError indicates an expected call of OnError.
func (mr *MockRabbitMQBrokerInterfaceMockRecorder) OnError(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnError", reflect.TypeOf((*MockRabbitMQBrokerInterface)(nil).OnError), arg0, arg1)
}
