// Code generated by MockGen. DO NOT EDIT.
// Source: shutdown/interfaces.go

// Package shutdown is a generated GoMock package.
package shutdown

import (
	reflect "reflect"

	stop "github.com/ThomasObenaus/go-base/shutdown/stop"
	gomock "github.com/golang/mock/gomock"
)

// MockstopIF is a mock of stopIF interface.
type MockstopIF struct {
	ctrl     *gomock.Controller
	recorder *MockstopIFMockRecorder
}

// MockstopIFMockRecorder is the mock recorder for MockstopIF.
type MockstopIFMockRecorder struct {
	mock *MockstopIF
}

// NewMockstopIF creates a new mock instance.
func NewMockstopIF(ctrl *gomock.Controller) *MockstopIF {
	mock := &MockstopIF{ctrl: ctrl}
	mock.recorder = &MockstopIFMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockstopIF) EXPECT() *MockstopIFMockRecorder {
	return m.recorder
}

// AddToBack mocks base method.
func (m *MockstopIF) AddToBack(stoppable1 stop.Stoppable) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToBack", stoppable1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddToBack indicates an expected call of AddToBack.
func (mr *MockstopIFMockRecorder) AddToBack(stoppable1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToBack", reflect.TypeOf((*MockstopIF)(nil).AddToBack), stoppable1)
}

// AddToFront mocks base method.
func (m *MockstopIF) AddToFront(stoppable stop.Stoppable) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToFront", stoppable)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddToFront indicates an expected call of AddToFront.
func (mr *MockstopIFMockRecorder) AddToFront(stoppable interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToFront", reflect.TypeOf((*MockstopIF)(nil).AddToFront), stoppable)
}

// StopAllInOrder mocks base method.
func (m *MockstopIF) StopAllInOrder(listener stop.Listener) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StopAllInOrder", listener)
}

// StopAllInOrder indicates an expected call of StopAllInOrder.
func (mr *MockstopIFMockRecorder) StopAllInOrder(listener interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopAllInOrder", reflect.TypeOf((*MockstopIF)(nil).StopAllInOrder), listener)
}

// MocksignalHandlerIF is a mock of signalHandlerIF interface.
type MocksignalHandlerIF struct {
	ctrl     *gomock.Controller
	recorder *MocksignalHandlerIFMockRecorder
}

// MocksignalHandlerIFMockRecorder is the mock recorder for MocksignalHandlerIF.
type MocksignalHandlerIFMockRecorder struct {
	mock *MocksignalHandlerIF
}

// NewMocksignalHandlerIF creates a new mock instance.
func NewMocksignalHandlerIF(ctrl *gomock.Controller) *MocksignalHandlerIF {
	mock := &MocksignalHandlerIF{ctrl: ctrl}
	mock.recorder = &MocksignalHandlerIFMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocksignalHandlerIF) EXPECT() *MocksignalHandlerIFMockRecorder {
	return m.recorder
}

// WaitForSignal mocks base method.
func (m *MocksignalHandlerIF) WaitForSignal() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WaitForSignal")
}

// WaitForSignal indicates an expected call of WaitForSignal.
func (mr *MocksignalHandlerIFMockRecorder) WaitForSignal() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitForSignal", reflect.TypeOf((*MocksignalHandlerIF)(nil).WaitForSignal))
}

// MocklogIF is a mock of logIF interface.
type MocklogIF struct {
	ctrl     *gomock.Controller
	recorder *MocklogIFMockRecorder
}

// MocklogIFMockRecorder is the mock recorder for MocklogIF.
type MocklogIFMockRecorder struct {
	mock *MocklogIF
}

// NewMocklogIF creates a new mock instance.
func NewMocklogIF(ctrl *gomock.Controller) *MocklogIF {
	mock := &MocklogIF{ctrl: ctrl}
	mock.recorder = &MocklogIFMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocklogIF) EXPECT() *MocklogIFMockRecorder {
	return m.recorder
}

// LogCanNotAddService mocks base method.
func (m *MocklogIF) LogCanNotAddService(serviceName string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "LogCanNotAddService", serviceName)
}

// LogCanNotAddService indicates an expected call of LogCanNotAddService.
func (mr *MocklogIFMockRecorder) LogCanNotAddService(serviceName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogCanNotAddService", reflect.TypeOf((*MocklogIF)(nil).LogCanNotAddService), serviceName)
}

// ServiceWasStopped mocks base method.
func (m *MocklogIF) ServiceWasStopped(name string, err ...error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{name}
	for _, a := range err {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "ServiceWasStopped", varargs...)
}

// ServiceWasStopped indicates an expected call of ServiceWasStopped.
func (mr *MocklogIFMockRecorder) ServiceWasStopped(name interface{}, err ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{name}, err...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServiceWasStopped", reflect.TypeOf((*MocklogIF)(nil).ServiceWasStopped), varargs...)
}

// ServiceWillBeStopped mocks base method.
func (m *MocklogIF) ServiceWillBeStopped(name string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ServiceWillBeStopped", name)
}

// ServiceWillBeStopped indicates an expected call of ServiceWillBeStopped.
func (mr *MocklogIFMockRecorder) ServiceWillBeStopped(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServiceWillBeStopped", reflect.TypeOf((*MocklogIF)(nil).ServiceWillBeStopped), name)
}

// ShutdownSignalReceived mocks base method.
func (m *MocklogIF) ShutdownSignalReceived() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShutdownSignalReceived")
}

// ShutdownSignalReceived indicates an expected call of ShutdownSignalReceived.
func (mr *MocklogIFMockRecorder) ShutdownSignalReceived() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShutdownSignalReceived", reflect.TypeOf((*MocklogIF)(nil).ShutdownSignalReceived))
}

// MockhealthIF is a mock of healthIF interface.
type MockhealthIF struct {
	ctrl     *gomock.Controller
	recorder *MockhealthIFMockRecorder
}

// MockhealthIFMockRecorder is the mock recorder for MockhealthIF.
type MockhealthIFMockRecorder struct {
	mock *MockhealthIF
}

// NewMockhealthIF creates a new mock instance.
func NewMockhealthIF(ctrl *gomock.Controller) *MockhealthIF {
	mock := &MockhealthIF{ctrl: ctrl}
	mock.recorder = &MockhealthIFMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockhealthIF) EXPECT() *MockhealthIFMockRecorder {
	return m.recorder
}

// IsHealthy mocks base method.
func (m *MockhealthIF) IsHealthy() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsHealthy")
	ret0, _ := ret[0].(error)
	return ret0
}

// IsHealthy indicates an expected call of IsHealthy.
func (mr *MockhealthIFMockRecorder) IsHealthy() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsHealthy", reflect.TypeOf((*MockhealthIF)(nil).IsHealthy))
}

// ShutdownSignalReceived mocks base method.
func (m *MockhealthIF) ShutdownSignalReceived() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ShutdownSignalReceived")
}

// ShutdownSignalReceived indicates an expected call of ShutdownSignalReceived.
func (mr *MockhealthIFMockRecorder) ShutdownSignalReceived() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShutdownSignalReceived", reflect.TypeOf((*MockhealthIF)(nil).ShutdownSignalReceived))
}

// String mocks base method.
func (m *MockhealthIF) String() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String.
func (mr *MockhealthIFMockRecorder) String() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockhealthIF)(nil).String))
}
