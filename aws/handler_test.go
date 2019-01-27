package aws

import (
	"github.com/c-garcia/halleffect/internal/pkg/exporter/mocks"
	handlerMocks "github.com/c-garcia/halleffect/aws/mocks"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"context"
)

func Test_NewLambdaHandler_CallsService(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockService(ctrl)
	mockService.EXPECT().ExportMetrics().Return(nil)
	mockLogger := handlerMocks.NewMockLogger(ctrl)
	mockLogger.EXPECT().Println("Metrics export done")

	sut := NewLambdaHandler(mockService, mockLogger)

	res, err := sut(context.Background(), ExportMetricsLambdaEvent{})

	assert.Zero(t, "", res)
	assert.NoError(t, err)
	ctrl.Finish()
}

func Test_NewLambdaHandler_PropagatesError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockService(ctrl)
	mockService.EXPECT().ExportMetrics().Return(assert.AnError)
	mockLogger := handlerMocks.NewMockLogger(ctrl)
	mockLogger.EXPECT().Printf("%+v", gomock.Any())

	sut := NewLambdaHandler(mockService, mockLogger)

	_, err := sut(context.Background(), ExportMetricsLambdaEvent{})

	assert.Error(t, err)
	assert.Equal(t, errors.Cause(err), assert.AnError)
	ctrl.Finish()
}
