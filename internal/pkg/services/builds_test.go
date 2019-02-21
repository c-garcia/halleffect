//+build !integration, !service

package services

import (
	"github.com/c-garcia/halleffect/internal/pkg/concourse"
	mocks2 "github.com/c-garcia/halleffect/internal/pkg/concourse/mocks"
	"github.com/c-garcia/halleffect/internal/pkg/metrics"
	mocks3 "github.com/c-garcia/halleffect/internal/pkg/store/mocks"
	"github.com/c-garcia/halleffect/internal/pkg/timing/mocks"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLastSuccessfulBuildDurationImpl_SaveMetrics(t *testing.T) {
	const Team = "t1"
	const Pipeline = "p1"
	const Job = "j1"
	const Concourse = "test-concourse"
	var StartTime = time.Now().Add(-24 * time.Hour)
	now := time.Now().Truncate(time.Second)

	failedBuild := concourse.Build{
		Id:           101,
		StartTime:    int(StartTime.Add(2 * time.Minute).Unix()),
		EndTime:      int(StartTime.Add(3 * time.Minute).Unix()),
		PipelineName: Pipeline,
		JobName:      Job,
		Status:       concourse.StatusFailed,
		TeamName:     Team,
	}
	successfulBuild := concourse.Build{
		Id:           100,
		StartTime:    int(StartTime.Unix()),
		EndTime:      int(StartTime.Add(1 * time.Minute).Unix()),
		PipelineName: Pipeline,
		JobName:      Job,
		Status:       concourse.StatusSucceeded,
		TeamName:     Team,
	}
	metricForSuccessful := metrics.JobLastSuccessfulDuration{
		Timestamp: now,
		Concourse: Concourse,
		Team:      Team,
		Pipeline:  Pipeline,
		Job:       Job,
		Duration:  successfulBuild.Duration(),
	}
	ctrl := gomock.NewController(t)
	mockClock := mocks.NewMockClock(ctrl)
	mockClock.EXPECT().UnixTime().Return(now.Unix())
	mockConcourse := mocks2.NewMockAPI(ctrl)
	mockConcourse.EXPECT().FindLastBuilds().Return([]concourse.Build{failedBuild, successfulBuild}, nil)
	mockStore := mocks3.NewMockJobLastSuccessfulDuration(ctrl)
	mockStore.EXPECT().Store(metricForSuccessful)
	sut := LastSuccessfulBuildDurationImpl{
		Name:         Concourse,
		Clock:        mockClock,
		Concourse:    mockConcourse,
		MetricsStore: mockStore,
	}
	err := sut.SaveMetrics()
	assert.NoError(t, err)
}

func TestLastSuccessfulBuildDurationImpl_SaveMetrics_PropagatesFindErrors(t *testing.T) {
	const Concourse = "test-concourse"
	now := time.Now().Truncate(time.Second)
	ctrl := gomock.NewController(t)
	mockClock := mocks.NewMockClock(ctrl)
	mockClock.EXPECT().UnixTime().Return(now.Unix())
	mockConcourse := mocks2.NewMockAPI(ctrl)
	mockConcourse.EXPECT().FindLastBuilds().Return(make([]concourse.Build, 0), assert.AnError)
	mockStore := mocks3.NewMockJobLastSuccessfulDuration(ctrl)
	sut := LastSuccessfulBuildDurationImpl{
		Name:         Concourse,
		Clock:        mockClock,
		Concourse:    mockConcourse,
		MetricsStore: mockStore,
	}
	err := sut.SaveMetrics()
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, errors.Cause(err))
}

func TestLastSuccessfulBuildDurationImpl_SaveMetrics_AbortsOnFirstStoreError(t *testing.T) {
	const Team = "t1"
	const Pipeline = "p1"
	const Job1 = "j1"
	const Job2 = "j2"
	const Concourse = "test-concourse"
	var StartTime = time.Now().Add(-24 * time.Hour)
	now := time.Now().Truncate(time.Second)

	successful1 := concourse.Build{
		Id:           101,
		StartTime:    int(StartTime.Add(2 * time.Minute).Unix()),
		EndTime:      int(StartTime.Add(3 * time.Minute).Unix()),
		PipelineName: Pipeline,
		JobName:      Job1,
		Status:       concourse.StatusSucceeded,
		TeamName:     Team,
	}
	successful2 := concourse.Build{
		Id:           100,
		StartTime:    int(StartTime.Unix()),
		EndTime:      int(StartTime.Add(1 * time.Minute).Unix()),
		PipelineName: Pipeline,
		JobName:      Job2,
		Status:       concourse.StatusSucceeded,
		TeamName:     Team,
	}
	ctrl := gomock.NewController(t)
	mockClock := mocks.NewMockClock(ctrl)
	mockClock.EXPECT().UnixTime().Return(now.Unix())
	mockConcourse := mocks2.NewMockAPI(ctrl)
	mockConcourse.EXPECT().FindLastBuilds().Return([]concourse.Build{successful1, successful2}, nil)
	mockStore := mocks3.NewMockJobLastSuccessfulDuration(ctrl)
	mockStore.EXPECT().Store(gomock.Any()).Return(assert.AnError)
	sut := LastSuccessfulBuildDurationImpl{
		Name:         Concourse,
		Clock:        mockClock,
		Concourse:    mockConcourse,
		MetricsStore: mockStore,
	}
	err := sut.SaveMetrics()
	assert.Error(t, err)
	assert.Equal(t, assert.AnError, errors.Cause(err))
}
