//+build !integration, !service

package services

import (
	"github.com/c-garcia/halleffect/internal/pkg/concourse"
	"github.com/c-garcia/halleffect/internal/pkg/metrics"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBuildToDurationMetric(t *testing.T) {
	const Team = "team"
	const Job = "j1"
	const Pipeline = "p1"
	const Concourse = "concourse"
	var samplingTime = time.Now()
	var startTime = samplingTime.Add(-24 * time.Hour)
	var endTime = startTime.Add(1 * time.Minute)
	build := concourse.Build{
		Id:           0,
		StartTime:    int(startTime.Unix()),
		EndTime:      int(endTime.Unix()),
		PipelineName: Pipeline,
		JobName:      Job,
		Status:       concourse.StatusSucceeded,
		TeamName:     Team,
	}
	expectedMetric := metrics.JobLastSuccessfulDuration{
		Timestamp: samplingTime,
		Concourse: Concourse,
		Team:      Team,
		Pipeline:  Pipeline,
		Job:       Job,
		Duration:  endTime.Sub(startTime),
	}
	assert.Equal(t, expectedMetric, BuildToDurationMetric(Concourse, samplingTime, build))
}
