// Code generated by MockGen. DO NOT EDIT.
// Source: log.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockLogger is a mock of Logger interface
type MockLogger struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerMockRecorder
}

// MockLoggerMockRecorder is the mock recorder for MockLogger
type MockLoggerMockRecorder struct {
	mock *MockLogger
}

// NewMockLogger creates a new mock instance
func NewMockLogger(ctrl *gomock.Controller) *MockLogger {
	mock := &MockLogger{ctrl: ctrl}
	mock.recorder = &MockLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLogger) EXPECT() *MockLoggerMockRecorder {
	return m.recorder
}

// Println mocks base method
func (m *MockLogger) Println(v ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range v {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Println", varargs...)
}

// Println indicates an expected call of Println
func (mr *MockLoggerMockRecorder) Println(v ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Println", reflect.TypeOf((*MockLogger)(nil).Println), v...)
}

// Printf mocks base method
func (m *MockLogger) Printf(format string, v ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{format}
	for _, a := range v {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Printf", varargs...)
}

// Printf indicates an expected call of Printf
func (mr *MockLoggerMockRecorder) Printf(format interface{}, v ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{format}, v...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Printf", reflect.TypeOf((*MockLogger)(nil).Printf), varargs...)
}
