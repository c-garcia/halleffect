package poller

import (
	"github.com/c-garcia/halleffect/internal/pkg/concourse"
	concourseMocks "github.com/c-garcia/halleffect/internal/pkg/concourse/mocks"
	"github.com/c-garcia/halleffect/internal/pkg/publisher"
	metricsMocks "github.com/c-garcia/halleffect/internal/pkg/publisher/mocks"
	timingMocks "github.com/c-garcia/halleffect/internal/pkg/timing/mocks"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const ConcourseHost = "concourse"

var SamplingTime = time.Now().Unix()

func Test_ExportsMetrics_PublishesAllFinishedBuilds(t *testing.T) {
	b1 := concourse.Build{StartTime: 100, EndTime: 120, PipelineName: "p1", JobName: "j1", Status: "finished", TeamName: "main"}
	b2 := concourse.Build{StartTime: 121, EndTime: 200, PipelineName: "p1", JobName: "j2", Status: "failed", TeamName: "not-main"}
	builds := []concourse.Build{b1, b2}
	m1 := publisher.JobDurationMetric{
		Concourse: ConcourseHost, Timestamp: 100, EndTime: 120, PipelineName: "p1", JobName: "j1", Status: "finished",
		TeamName: "main",
	}
	m2 := publisher.JobDurationMetric{
		Concourse: ConcourseHost, Timestamp: 121, EndTime: 200, PipelineName: "p1", JobName: "j2", Status: "failed",
		TeamName: "not-main",
	}
	ctrl := gomock.NewController(t)
	mockConcourse := concourseMocks.NewMockAPI(ctrl)
	mockConcourse.EXPECT().Name().Return(ConcourseHost)
	mockConcourse.EXPECT().FindLastBuilds().Return(builds, nil)
	mockConcourse.EXPECT().SupportsJobsEndpoint().Return(true, nil)
	noStatuses := make([]concourse.JobStatus, 0)
	mockConcourse.EXPECT().FindJobStatuses().Return(noStatuses, nil)
	mockExporter := metricsMocks.NewMockMetricsPublisher(ctrl)
	mockExporter.EXPECT().PublishJobDuration(m1).Return(nil)
	mockExporter.EXPECT().PublishJobDuration(m2).Return(nil)
	mockClock := timingMocks.NewMockClock(ctrl)
	mockClock.EXPECT().UnixTime().Return(SamplingTime)

	sut := New(mockConcourse, mockExporter, mockClock)
	err := sut.ExportJobDurationMetrics()

	assert.NoError(t, err)
	ctrl.Finish()
}

func Test_ExportsMetrics_PublishesOnlyFinishedBuilds(t *testing.T) {
	b1 := concourse.Build{StartTime: 100, EndTime: 0, PipelineName: "p1", JobName: "j1", Status: "started", TeamName: "main"}
	b2 := concourse.Build{StartTime: 121, EndTime: 200, PipelineName: "p1", JobName: "j2", Status: "failed", TeamName: "not-main"}
	builds := []concourse.Build{b1, b2}
	m2 := publisher.JobDurationMetric{
		Concourse: ConcourseHost, Timestamp: 121, EndTime: 200, PipelineName: "p1", JobName: "j2", Status: "failed",
		TeamName: "not-main",
	}
	ctrl := gomock.NewController(t)
	mockConcourse := concourseMocks.NewMockAPI(ctrl)
	mockConcourse.EXPECT().FindLastBuilds().Return(builds, nil)
	mockConcourse.EXPECT().Name().Return(ConcourseHost)
	mockConcourse.EXPECT().SupportsJobsEndpoint().Return(true, nil)
	noStatuses := make([]concourse.JobStatus, 0)
	mockConcourse.EXPECT().FindJobStatuses().Return(noStatuses, nil)
	mockExporter := metricsMocks.NewMockMetricsPublisher(ctrl)
	mockExporter.EXPECT().PublishJobDuration(m2).Return(nil)
	mockClock := timingMocks.NewMockClock(ctrl)
	mockClock.EXPECT().UnixTime().Return(SamplingTime)

	sut := New(mockConcourse, mockExporter, mockClock)
	err := sut.ExportJobDurationMetrics()

	assert.NoError(t, err)
	ctrl.Finish()
}

func Test_ExportsMetrics_DoesNotExportStatusesIfEndpointDoesNotSupportIt(t *testing.T) {
	b1 := concourse.Build{StartTime: 100, EndTime: 0, PipelineName: "p1", JobName: "j1", Status: "started", TeamName: "main"}
	b2 := concourse.Build{StartTime: 121, EndTime: 200, PipelineName: "p1", JobName: "j2", Status: "failed", TeamName: "not-main"}
	builds := []concourse.Build{b1, b2}
	m2 := publisher.JobDurationMetric{
		Concourse: ConcourseHost, Timestamp: 121, EndTime: 200, PipelineName: "p1", JobName: "j2", Status: "failed",
		TeamName: "not-main",
	}
	ctrl := gomock.NewController(t)
	mockConcourse := concourseMocks.NewMockAPI(ctrl)
	mockConcourse.EXPECT().Name().Return(ConcourseHost)
	mockConcourse.EXPECT().FindLastBuilds().Return(builds, nil)
	mockExporter := metricsMocks.NewMockMetricsPublisher(ctrl)
	mockExporter.EXPECT().PublishJobDuration(m2).Return(nil)
	mockConcourse.EXPECT().SupportsJobsEndpoint().Return(false, nil)
	mockClock := timingMocks.NewMockClock(ctrl)
	mockClock.EXPECT().UnixTime().Return(SamplingTime)

	sut := New(mockConcourse, mockExporter, mockClock)
	err := sut.ExportJobDurationMetrics()

	assert.NoError(t, err)
	ctrl.Finish()
}

func Test_ExportsMetrics_DoesNotExportStatusesIfEndpointTestingFails(t *testing.T) {
	b1 := concourse.Build{StartTime: 100, EndTime: 0, PipelineName: "p1", JobName: "j1", Status: "started", TeamName: "main"}
	b2 := concourse.Build{StartTime: 121, EndTime: 200, PipelineName: "p1", JobName: "j2", Status: "failed", TeamName: "not-main"}
	builds := []concourse.Build{b1, b2}
	m2 := publisher.JobDurationMetric{
		Concourse: ConcourseHost, Timestamp: 121, EndTime: 200, PipelineName: "p1", JobName: "j2", Status: "failed",
		TeamName: "not-main",
	}
	ctrl := gomock.NewController(t)
	mockConcourse := concourseMocks.NewMockAPI(ctrl)
	mockConcourse.EXPECT().Name().Return(ConcourseHost)
	mockConcourse.EXPECT().FindLastBuilds().Return(builds, nil)
	mockExporter := metricsMocks.NewMockMetricsPublisher(ctrl)
	mockExporter.EXPECT().PublishJobDuration(m2).Return(nil)
	mockConcourse.EXPECT().SupportsJobsEndpoint().Return(false, assert.AnError)
	mockClock := timingMocks.NewMockClock(ctrl)
	mockClock.EXPECT().UnixTime().Return(SamplingTime)

	sut := New(mockConcourse, mockExporter, mockClock)
	err := sut.ExportJobDurationMetrics()

	assert.Error(t, err)
	ctrl.Finish()
}

func Test_ExportsMetrics_PropagatesConcourseErrors(t *testing.T) {
	b1 := concourse.Build{StartTime: 100, EndTime: 0, PipelineName: "p1", JobName: "j1", Status: "started", TeamName: "main"}
	b2 := concourse.Build{StartTime: 121, EndTime: 200, PipelineName: "p1", JobName: "j2", Status: "failed", TeamName: "not-main"}
	builds := []concourse.Build{b1, b2}
	ctrl := gomock.NewController(t)
	mockConcourse := concourseMocks.NewMockAPI(ctrl)
	mockConcourse.EXPECT().Name().Return(ConcourseHost)
	mockConcourse.EXPECT().FindLastBuilds().Return(builds, assert.AnError)
	mockExporter := metricsMocks.NewMockMetricsPublisher(ctrl)
	mockClock := timingMocks.NewMockClock(ctrl)
	mockClock.EXPECT().UnixTime().Return(SamplingTime)

	sut := New(mockConcourse, mockExporter, mockClock)
	err := sut.ExportJobDurationMetrics()

	assert.Error(t, err)
	assert.Equal(t, assert.AnError, errors.Cause(err))
	ctrl.Finish()
}

func Test_ExportsMetrics_AbortsAtFirstPublishingError(t *testing.T) {
	b1 := concourse.Build{StartTime: 100, EndTime: 120, PipelineName: "p1", JobName: "j1", Status: "finished", TeamName: "main"}
	b2 := concourse.Build{StartTime: 121, EndTime: 200, PipelineName: "p1", JobName: "j2", Status: "failed", TeamName: "not-main"}
	builds := []concourse.Build{b1, b2}
	m1 := publisher.JobDurationMetric{
		Concourse: ConcourseHost, Timestamp: 100, EndTime: 120, PipelineName: "p1", JobName: "j1", Status: "finished",
		TeamName: "main",
	}
	ctrl := gomock.NewController(t)
	mockConcourse := concourseMocks.NewMockAPI(ctrl)
	mockConcourse.EXPECT().Name().Return(ConcourseHost)
	mockConcourse.EXPECT().FindLastBuilds().Return(builds, nil)
	mockExporter := metricsMocks.NewMockMetricsPublisher(ctrl)
	mockExporter.EXPECT().PublishJobDuration(m1).Return(assert.AnError)

	mockClock := timingMocks.NewMockClock(ctrl)
	mockClock.EXPECT().UnixTime().Return(SamplingTime)

	sut := New(mockConcourse, mockExporter, mockClock)
	err := sut.ExportJobDurationMetrics()

	assert.Error(t, err)
	assert.Equal(t, assert.AnError, errors.Cause(err))
	ctrl.Finish()
}
