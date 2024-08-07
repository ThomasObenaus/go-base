// Code generated by MockGen. DO NOT EDIT.
// Source: shutdown/stopable.go

// Package shutdown is a generated GoMock package.
package shutdown

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockStopable is a mock of Stopable interface.
type MockStopable struct {
	ctrl     *gomock.Controller
	recorder *MockStopableMockRecorder
}

// MockStopableMockRecorder is the mock recorder for MockStopable.
type MockStopableMockRecorder struct {
	mock *MockStopable
}

// NewMockStopable creates a new mock instance.
func NewMockStopable(ctrl *gomock.Controller) *MockStopable {
	mock := &MockStopable{ctrl: ctrl}
	mock.recorder = &MockStopableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStopable) EXPECT() *MockStopableMockRecorder {
	return m.recorder
}

// Stop mocks base method.
func (m *MockStopable) Stop() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockStopableMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockStopable)(nil).Stop))
}

// String mocks base method.
func (m *MockStopable) String() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String.
func (mr *MockStopableMockRecorder) String() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockStopable)(nil).String))
}
