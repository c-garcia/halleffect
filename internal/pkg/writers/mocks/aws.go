// Code generated by MockGen. DO NOT EDIT.
// Source: aws.go

// Package mocks is a generated GoMock package.
package mocks

import (
	cloudwatch "github.com/aws/aws-sdk-go/service/cloudwatch"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockAWSCloudWatchMetricWriter is a mock of AWSCloudWatchMetricWriter interface
type MockAWSCloudWatchMetricWriter struct {
	ctrl     *gomock.Controller
	recorder *MockAWSCloudWatchMetricWriterMockRecorder
}

// MockAWSCloudWatchMetricWriterMockRecorder is the mock recorder for MockAWSCloudWatchMetricWriter
type MockAWSCloudWatchMetricWriterMockRecorder struct {
	mock *MockAWSCloudWatchMetricWriter
}

// NewMockAWSCloudWatchMetricWriter creates a new mock instance
func NewMockAWSCloudWatchMetricWriter(ctrl *gomock.Controller) *MockAWSCloudWatchMetricWriter {
	mock := &MockAWSCloudWatchMetricWriter{ctrl: ctrl}
	mock.recorder = &MockAWSCloudWatchMetricWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAWSCloudWatchMetricWriter) EXPECT() *MockAWSCloudWatchMetricWriterMockRecorder {
	return m.recorder
}

// PutMetricData mocks base method
func (m *MockAWSCloudWatchMetricWriter) PutMetricData(in *cloudwatch.PutMetricDataInput) (*cloudwatch.PutMetricDataOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutMetricData", in)
	ret0, _ := ret[0].(*cloudwatch.PutMetricDataOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PutMetricData indicates an expected call of PutMetricData
func (mr *MockAWSCloudWatchMetricWriterMockRecorder) PutMetricData(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutMetricData", reflect.TypeOf((*MockAWSCloudWatchMetricWriter)(nil).PutMetricData), in)
}