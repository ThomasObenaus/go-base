// Code generated by MockGen. DO NOT EDIT.
// Source: config/provider.go

// Package mock_config is a generated GoMock package.
package mock_config

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockProvider is a mock of Provider interface
type MockProvider struct {
	ctrl     *gomock.Controller
	recorder *MockProviderMockRecorder
}

// MockProviderMockRecorder is the mock recorder for MockProvider
type MockProviderMockRecorder struct {
	mock *MockProvider
}

// NewMockProvider creates a new mock instance
func NewMockProvider(ctrl *gomock.Controller) *MockProvider {
	mock := &MockProvider{ctrl: ctrl}
	mock.recorder = &MockProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockProvider) EXPECT() *MockProviderMockRecorder {
	return m.recorder
}

// ReadConfig mocks base method
func (m *MockProvider) ReadConfig(args []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadConfig", args)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReadConfig indicates an expected call of ReadConfig
func (mr *MockProviderMockRecorder) ReadConfig(args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadConfig", reflect.TypeOf((*MockProvider)(nil).ReadConfig), args)
}

// Get mocks base method
func (m *MockProvider) Get(key string) interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", key)
	ret0, _ := ret[0].(interface{})
	return ret0
}

// Get indicates an expected call of Get
func (mr *MockProviderMockRecorder) Get(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockProvider)(nil).Get), key)
}

// GetString mocks base method
func (m *MockProvider) GetString(key string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetString", key)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetString indicates an expected call of GetString
func (mr *MockProviderMockRecorder) GetString(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetString", reflect.TypeOf((*MockProvider)(nil).GetString), key)
}

// GetBool mocks base method
func (m *MockProvider) GetBool(key string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBool", key)
	ret0, _ := ret[0].(bool)
	return ret0
}

// GetBool indicates an expected call of GetBool
func (mr *MockProviderMockRecorder) GetBool(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBool", reflect.TypeOf((*MockProvider)(nil).GetBool), key)
}

// GetInt mocks base method
func (m *MockProvider) GetInt(key string) int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInt", key)
	ret0, _ := ret[0].(int)
	return ret0
}

// GetInt indicates an expected call of GetInt
func (mr *MockProviderMockRecorder) GetInt(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInt", reflect.TypeOf((*MockProvider)(nil).GetInt), key)
}

// GetInt32 mocks base method
func (m *MockProvider) GetInt32(key string) int32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInt32", key)
	ret0, _ := ret[0].(int32)
	return ret0
}

// GetInt32 indicates an expected call of GetInt32
func (mr *MockProviderMockRecorder) GetInt32(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInt32", reflect.TypeOf((*MockProvider)(nil).GetInt32), key)
}

// GetInt64 mocks base method
func (m *MockProvider) GetInt64(key string) int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInt64", key)
	ret0, _ := ret[0].(int64)
	return ret0
}

// GetInt64 indicates an expected call of GetInt64
func (mr *MockProviderMockRecorder) GetInt64(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInt64", reflect.TypeOf((*MockProvider)(nil).GetInt64), key)
}

// GetUint mocks base method
func (m *MockProvider) GetUint(key string) uint {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUint", key)
	ret0, _ := ret[0].(uint)
	return ret0
}

// GetUint indicates an expected call of GetUint
func (mr *MockProviderMockRecorder) GetUint(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUint", reflect.TypeOf((*MockProvider)(nil).GetUint), key)
}

// GetUint32 mocks base method
func (m *MockProvider) GetUint32(key string) uint32 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUint32", key)
	ret0, _ := ret[0].(uint32)
	return ret0
}

// GetUint32 indicates an expected call of GetUint32
func (mr *MockProviderMockRecorder) GetUint32(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUint32", reflect.TypeOf((*MockProvider)(nil).GetUint32), key)
}

// GetUint64 mocks base method
func (m *MockProvider) GetUint64(key string) uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUint64", key)
	ret0, _ := ret[0].(uint64)
	return ret0
}

// GetUint64 indicates an expected call of GetUint64
func (mr *MockProviderMockRecorder) GetUint64(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUint64", reflect.TypeOf((*MockProvider)(nil).GetUint64), key)
}

// GetFloat64 mocks base method
func (m *MockProvider) GetFloat64(key string) float64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFloat64", key)
	ret0, _ := ret[0].(float64)
	return ret0
}

// GetFloat64 indicates an expected call of GetFloat64
func (mr *MockProviderMockRecorder) GetFloat64(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFloat64", reflect.TypeOf((*MockProvider)(nil).GetFloat64), key)
}

// GetTime mocks base method
func (m *MockProvider) GetTime(key string) time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTime", key)
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// GetTime indicates an expected call of GetTime
func (mr *MockProviderMockRecorder) GetTime(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTime", reflect.TypeOf((*MockProvider)(nil).GetTime), key)
}

// GetDuration mocks base method
func (m *MockProvider) GetDuration(key string) time.Duration {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDuration", key)
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

// GetDuration indicates an expected call of GetDuration
func (mr *MockProviderMockRecorder) GetDuration(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDuration", reflect.TypeOf((*MockProvider)(nil).GetDuration), key)
}

// GetIntSlice mocks base method
func (m *MockProvider) GetIntSlice(key string) []int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIntSlice", key)
	ret0, _ := ret[0].([]int)
	return ret0
}

// GetIntSlice indicates an expected call of GetIntSlice
func (mr *MockProviderMockRecorder) GetIntSlice(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIntSlice", reflect.TypeOf((*MockProvider)(nil).GetIntSlice), key)
}

// GetStringSlice mocks base method
func (m *MockProvider) GetStringSlice(key string) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStringSlice", key)
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetStringSlice indicates an expected call of GetStringSlice
func (mr *MockProviderMockRecorder) GetStringSlice(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStringSlice", reflect.TypeOf((*MockProvider)(nil).GetStringSlice), key)
}

// GetStringMap mocks base method
func (m *MockProvider) GetStringMap(key string) map[string]interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStringMap", key)
	ret0, _ := ret[0].(map[string]interface{})
	return ret0
}

// GetStringMap indicates an expected call of GetStringMap
func (mr *MockProviderMockRecorder) GetStringMap(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStringMap", reflect.TypeOf((*MockProvider)(nil).GetStringMap), key)
}

// GetStringMapString mocks base method
func (m *MockProvider) GetStringMapString(key string) map[string]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStringMapString", key)
	ret0, _ := ret[0].(map[string]string)
	return ret0
}

// GetStringMapString indicates an expected call of GetStringMapString
func (mr *MockProviderMockRecorder) GetStringMapString(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStringMapString", reflect.TypeOf((*MockProvider)(nil).GetStringMapString), key)
}

// GetStringMapStringSlice mocks base method
func (m *MockProvider) GetStringMapStringSlice(key string) map[string][]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStringMapStringSlice", key)
	ret0, _ := ret[0].(map[string][]string)
	return ret0
}

// GetStringMapStringSlice indicates an expected call of GetStringMapStringSlice
func (mr *MockProviderMockRecorder) GetStringMapStringSlice(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStringMapStringSlice", reflect.TypeOf((*MockProvider)(nil).GetStringMapStringSlice), key)
}

// GetSizeInBytes mocks base method
func (m *MockProvider) GetSizeInBytes(key string) uint {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSizeInBytes", key)
	ret0, _ := ret[0].(uint)
	return ret0
}

// GetSizeInBytes indicates an expected call of GetSizeInBytes
func (mr *MockProviderMockRecorder) GetSizeInBytes(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSizeInBytes", reflect.TypeOf((*MockProvider)(nil).GetSizeInBytes), key)
}

// IsSet mocks base method
func (m *MockProvider) IsSet(key string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsSet", key)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsSet indicates an expected call of IsSet
func (mr *MockProviderMockRecorder) IsSet(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsSet", reflect.TypeOf((*MockProvider)(nil).IsSet), key)
}

// String mocks base method
func (m *MockProvider) String() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String
func (mr *MockProviderMockRecorder) String() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockProvider)(nil).String))
}