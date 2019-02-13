// Code generated by MockGen. DO NOT EDIT.
// Source: poller.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockService is a mock of Service interface
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// ExportJobDurationMetrics mocks base method
func (m *MockService) ExportJobDurationMetrics() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExportJobDurationMetrics")
	ret0, _ := ret[0].(error)
	return ret0
}

// ExportJobDurationMetrics indicates an expected call of ExportJobDurationMetrics
func (mr *MockServiceMockRecorder) ExportJobDurationMetrics() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExportJobDurationMetrics", reflect.TypeOf((*MockService)(nil).ExportJobDurationMetrics))
}
