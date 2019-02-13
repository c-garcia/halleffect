package aws

import (
	"context"
	handlerMocks "github.com/c-garcia/halleffect/aws/mocks"
	"github.com/c-garcia/halleffect/internal/pkg/poller/mocks"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewLambdaHandler_CallsService(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockService(ctrl)
	mockService.EXPECT().ExportJobDurationMetrics().Return(nil)
	mockLogger := handlerMocks.NewMockLogger(ctrl)
	mockLogger.EXPECT().Println("Metrics export done")

	sut := NewLambdaHandler(mockService, mockLogger)

	res, err := sut(context.Background(), PublishConcourseMetricsLambdaEvent{})

	assert.Zero(t, "", res)
	assert.NoError(t, err)
	ctrl.Finish()
}

func Test_NewLambdaHandler_PropagatesError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockService(ctrl)
	mockService.EXPECT().ExportJobDurationMetrics().Return(assert.AnError)
	mockLogger := handlerMocks.NewMockLogger(ctrl)
	mockLogger.EXPECT().Printf("%+v", gomock.Any())

	sut := NewLambdaHandler(mockService, mockLogger)

	_, err := sut(context.Background(), PublishConcourseMetricsLambdaEvent{})

	assert.Error(t, err)
	assert.Equal(t, errors.Cause(err), assert.AnError)
	ctrl.Finish()
}
