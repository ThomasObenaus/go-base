// Code generated by MockGen. DO NOT EDIT.
// Source: signal/signal.go

// Package signal is a generated GoMock package.
package signal

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockListener is a mock of Listener interface.
type MockListener struct {
	ctrl     *gomock.Controller
	recorder *MockListenerMockRecorder
}

// MockListenerMockRecorder is the mock recorder for MockListener.
type MockListenerMockRecorder struct {
	mock *MockListener
}

// NewMockListener creates a new mock instance.
func NewMockListener(ctrl *gomock.Controller) *MockListener {
	mock := &MockListener{ctrl: ctrl}
	mock.recorder = &MockListenerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockListener) EXPECT() *MockListenerMockRecorder {
	return m.recorder
}

// ShutdownSignalReceived mocks base method.
func (m *MockListener) ShutdownSignalReceived() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShutdownSignalReceived")
}

// ShutdownSignalReceived indicates an expected call of ShutdownSignalReceived.
func (mr *MockListenerMockRecorder) ShutdownSignalReceived() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShutdownSignalReceived", reflect.TypeOf((*MockListener)(nil).ShutdownSignalReceived))
}
