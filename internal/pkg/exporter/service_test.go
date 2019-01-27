package exporter

import (
	"github.com/c-garcia/halleffect/internal/pkg/concourse"
	concourseMocks "github.com/c-garcia/halleffect/internal/pkg/concourse/mocks"
	"github.com/c-garcia/halleffect/internal/pkg/metrics"
	metricsMocks "github.com/c-garcia/halleffect/internal/pkg/metrics/mocks"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ExportsMetrics_PublishesAllFinishedBuilds(t *testing.T) {
	b1 := concourse.Build{StartTime: 100, EndTime: 120}
	b2 := concourse.Build{StartTime: 121, EndTime: 200}
	builds := []concourse.Build{b1, b2}
	m1 := metrics.Metric{Timestamp: 100, EndTime: 120}
	m2 := metrics.Metric{Timestamp: 121, EndTime: 200}
	ctrl := gomock.NewController(t)
	mockConcourse := concourseMocks.NewMockAPI(ctrl)
	mockConcourse.EXPECT().FindLastBuilds().Return(builds, nil)
	mockExporter := metricsMocks.NewMockExporter(ctrl)
	mockExporter.EXPECT().Publish(m1).Return(nil)
	mockExporter.EXPECT().Publish(m2).Return(nil)

	sut := New(mockConcourse, mockExporter)
	err := sut.ExportMetrics()

	assert.NoError(t, err)
	ctrl.Finish()
}

func Test_ExportsMetrics_PublishesOnlyFinishedBuilds(t *testing.T) {
	b1 := concourse.Build{StartTime: 100, EndTime: 0}
	b2 := concourse.Build{StartTime: 121, EndTime: 200}
	builds := []concourse.Build{b1, b2}
	m2 := metrics.Metric{Timestamp: 121, EndTime: 200}
	ctrl := gomock.NewController(t)
	mockConcourse := concourseMocks.NewMockAPI(ctrl)
	mockConcourse.EXPECT().FindLastBuilds().Return(builds, nil)
	mockExporter := metricsMocks.NewMockExporter(ctrl)
	mockExporter.EXPECT().Publish(m2).Return(nil)

	sut := New(mockConcourse, mockExporter)
	err := sut.ExportMetrics()

	assert.NoError(t, err)
	ctrl.Finish()
}

func Test_ExportsMetrics_PropagatesConcourseErrors(t *testing.T) {
	b1 := concourse.Build{StartTime: 100, EndTime: 120}
	b2 := concourse.Build{StartTime: 121, EndTime: 200}
	builds := []concourse.Build{b1, b2}
	ctrl := gomock.NewController(t)
	mockConcourse := concourseMocks.NewMockAPI(ctrl)
	mockConcourse.EXPECT().FindLastBuilds().Return(builds, assert.AnError)
	mockExporter := metricsMocks.NewMockExporter(ctrl)

	sut := New(mockConcourse, mockExporter)
	err := sut.ExportMetrics()

	assert.Error(t, err)
	assert.Equal(t, assert.AnError, errors.Cause(err))
	ctrl.Finish()
}

func Test_ExportsMetrics_AbortsAtFirstPublishingError(t *testing.T) {
	b1 := concourse.Build{StartTime: 100, EndTime: 120}
	b2 := concourse.Build{StartTime: 121, EndTime: 200}
	builds := []concourse.Build{b1, b2}
	m1 := metrics.Metric{Timestamp: 100, EndTime: 120}
	ctrl := gomock.NewController(t)
	mockConcourse := concourseMocks.NewMockAPI(ctrl)
	mockConcourse.EXPECT().FindLastBuilds().Return(builds, nil)
	mockExporter := metricsMocks.NewMockExporter(ctrl)
	mockExporter.EXPECT().Publish(m1).Return(assert.AnError)

	sut := New(mockConcourse, mockExporter)
	err := sut.ExportMetrics()

	assert.Error(t, err)
	assert.Equal(t, assert.AnError, errors.Cause(err))
	ctrl.Finish()
}
