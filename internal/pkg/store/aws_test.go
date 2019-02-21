//+build !integration, !service

package store

import (
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/c-garcia/halleffect/internal/pkg/metrics"
	"github.com/c-garcia/halleffect/internal/pkg/writers/mocks"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestJobLastSuccessfulDurationAWS_Store(t *testing.T) {
	now := time.Now()
	const Concourse = "concourse-test"
	const Team = "main"
	const Pipeline = "p1"
	const Job = "j1"
	var Duration = 120 * time.Second
	metric := metrics.JobLastSuccessfulDuration{
		Timestamp: now,
		Concourse: Concourse,
		Team:      Team,
		Pipeline:  Pipeline,
		Job:       Job,
		Duration:  Duration,
	}
	pipelineDimension := "pipeline"
	jobNameDimension := "job"
	concourseDimension := "concourse"
	teamDimension := "team"
	datum := &cloudwatch.MetricDatum{}
	datum.
		SetTimestamp(now).
		SetMetricName("LastSuccessfulDuration").
		SetDimensions([]*cloudwatch.Dimension{
			{Name: &concourseDimension, Value: &metric.Concourse},
			{Name: &teamDimension, Value: &metric.Team},
			{Name: &pipelineDimension, Value: &metric.Pipeline},
			{Name: &jobNameDimension, Value: &metric.Job},
		}).
		SetUnit("Seconds")

	dataInput := &cloudwatch.PutMetricDataInput{}
	dataInput.
		SetNamespace("Concourse/Jobs").
		SetMetricData([]*cloudwatch.MetricDatum{datum})

	testCases := []struct {
		Desc     string
		Duration time.Duration
	}{{"120 seconds", 120 * time.Second}}

	for _, testCase := range testCases {
		t.Run(testCase.Desc, func(t *testing.T) {
			metric.Duration = testCase.Duration
			datum.SetValue(testCase.Duration.Seconds())

			ctrl := gomock.NewController(t)
			mockWriter := mocks.NewMockAWSCloudWatchMetricWriter(ctrl)
			mockWriter.EXPECT().PutMetricData(dataInput).Return(&cloudwatch.PutMetricDataOutput{}, nil)

			sut := NewJobLastSuccessfulDurationAWS("Concourse/Jobs", mockWriter)
			err := sut.Store(metric)

			assert.NoError(t, err)
			ctrl.Finish()
		})
	}
}

func TestJobLastSuccessfulDurationAWS_Store_PropagatesErrors(t *testing.T) {
	now := time.Now()
	const Concourse = "concourse-test"
	const Team = "main"
	const Pipeline = "p1"
	const Job = "j1"
	var Duration = 120 * time.Second
	metric := metrics.JobLastSuccessfulDuration{
		Timestamp: now,
		Concourse: Concourse,
		Team:      Team,
		Pipeline:  Pipeline,
		Job:       Job,
		Duration:  Duration,
	}
	pipelineDimension := "pipeline"
	jobNameDimension := "job"
	concourseDimension := "concourse"
	teamDimension := "team"
	datum := &cloudwatch.MetricDatum{}
	datum.
		SetTimestamp(now).
		SetMetricName("LastSuccessfulDuration").
		SetDimensions([]*cloudwatch.Dimension{
			{Name: &concourseDimension, Value: &metric.Concourse},
			{Name: &teamDimension, Value: &metric.Team},
			{Name: &pipelineDimension, Value: &metric.Pipeline},
			{Name: &jobNameDimension, Value: &metric.Job},
		}).
		SetUnit("Seconds")

	dataInput := &cloudwatch.PutMetricDataInput{}
	dataInput.
		SetNamespace("Concourse/Jobs").
		SetMetricData([]*cloudwatch.MetricDatum{datum})

	testCases := []struct {
		Desc     string
		Duration time.Duration
	}{{"120 seconds", 120 * time.Second}}

	for _, testCase := range testCases {
		t.Run(testCase.Desc, func(t *testing.T) {
			metric.Duration = testCase.Duration
			datum.SetValue(testCase.Duration.Seconds())

			ctrl := gomock.NewController(t)
			mockWriter := mocks.NewMockAWSCloudWatchMetricWriter(ctrl)
			mockWriter.EXPECT().PutMetricData(dataInput).Return(nil, assert.AnError)

			sut := NewJobLastSuccessfulDurationAWS("Concourse/Jobs", mockWriter)
			err := sut.Store(metric)

			assert.Error(t, err)
			assert.Equal(t, assert.AnError, errors.Cause(err))
			ctrl.Finish()
		})
	}
}
