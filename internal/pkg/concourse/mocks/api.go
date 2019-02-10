// Code generated by MockGen. DO NOT EDIT.
// Source: api.go

// Package mocks is a generated GoMock package.
package mocks

import (
	concourse "github.com/c-garcia/halleffect/internal/pkg/concourse"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockAPI is a mock of API interface
type MockAPI struct {
	ctrl     *gomock.Controller
	recorder *MockAPIMockRecorder
}

// MockAPIMockRecorder is the mock recorder for MockAPI
type MockAPIMockRecorder struct {
	mock *MockAPI
}

// NewMockAPI creates a new mock instance
func NewMockAPI(ctrl *gomock.Controller) *MockAPI {
	mock := &MockAPI{ctrl: ctrl}
	mock.recorder = &MockAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAPI) EXPECT() *MockAPIMockRecorder {
	return m.recorder
}

// Name mocks base method
func (m *MockAPI) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name
func (mr *MockAPIMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockAPI)(nil).Name))
}

// FindLastBuilds mocks base method
func (m *MockAPI) FindLastBuilds() ([]concourse.Build, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindLastBuilds")
	ret0, _ := ret[0].([]concourse.Build)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindLastBuilds indicates an expected call of FindLastBuilds
func (mr *MockAPIMockRecorder) FindLastBuilds() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindLastBuilds", reflect.TypeOf((*MockAPI)(nil).FindLastBuilds))
}

// FindJobStatuses mocks base method
func (m *MockAPI) FindJobStatuses() ([]concourse.JobStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindJobStatuses")
	ret0, _ := ret[0].([]concourse.JobStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindJobStatuses indicates an expected call of FindJobStatuses
func (mr *MockAPIMockRecorder) FindJobStatuses() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindJobStatuses", reflect.TypeOf((*MockAPI)(nil).FindJobStatuses))
}
