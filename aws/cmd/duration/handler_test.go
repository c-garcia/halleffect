// +build !integration, !service

package main

import (
	"context"
	"github.com/c-garcia/halleffect/internal/pkg/microlog/mocks"
	mocks2 "github.com/c-garcia/halleffect/internal/pkg/services/mocks"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewLambdaHandler_CallsService(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks2.NewMockMetrics(ctrl)
	mockService.EXPECT().SaveMetrics().Return(nil)
	mockLogger := mocks.NewMockLogger(ctrl)
	mockLogger.EXPECT().Println("Metrics export done")

	sut := NewLambdaHandler(mockService, mockLogger)

	res, err := sut(context.Background(), GetDurationEvent{})

	assert.Zero(t, "", res)
	assert.NoError(t, err)
	ctrl.Finish()
}

func Test_NewLambdaHandler_PropagatesError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks2.NewMockMetrics(ctrl)
	mockService.EXPECT().SaveMetrics().Return(assert.AnError)
	mockLogger := mocks.NewMockLogger(ctrl)
	mockLogger.EXPECT().Printf("%+v", gomock.Any())

	sut := NewLambdaHandler(mockService, mockLogger)

	_, err := sut(context.Background(), GetDurationEvent{})

	assert.Error(t, err)
	assert.Equal(t, errors.Cause(err), assert.AnError)
	ctrl.Finish()
}
