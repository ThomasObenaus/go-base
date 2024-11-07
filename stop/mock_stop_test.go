// Code generated by MockGen. DO NOT EDIT.
// Source: stop/interfaces.go

// Package stop is a generated GoMock package.
package stop

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockStoppable is a mock of Stoppable interface.
type MockStoppable struct {
	ctrl     *gomock.Controller
	recorder *MockStoppableMockRecorder
}

// MockStoppableMockRecorder is the mock recorder for MockStoppable.
type MockStoppableMockRecorder struct {
	mock *MockStoppable
}

// NewMockStoppable creates a new mock instance.
func NewMockStoppable(ctrl *gomock.Controller) *MockStoppable {
	mock := &MockStoppable{ctrl: ctrl}
	mock.recorder = &MockStoppableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStoppable) EXPECT() *MockStoppableMockRecorder {
	return m.recorder
}

// Stop mocks base method.
func (m *MockStoppable) Stop() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockStoppableMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockStoppable)(nil).Stop))
}

// String mocks base method.
func (m *MockStoppable) String() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String.
func (mr *MockStoppableMockRecorder) String() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockStoppable)(nil).String))
}