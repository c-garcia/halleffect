package services

import (
	"github.com/c-garcia/halleffect/internal/pkg/concourse"
	"github.com/c-garcia/halleffect/internal/pkg/metrics"
	"time"
)

func BuildToDurationMetric(concourse string, when time.Time, b concourse.Build) metrics.JobLastSuccessfulDuration {
	return metrics.JobLastSuccessfulDuration{
		Timestamp: when,
		Concourse: concourse,
		Team:      b.TeamName,
		Pipeline:  b.PipelineName,
		Job:       b.JobName,
		Duration:  b.EndTime.Sub(b.StartTime),
	}
}
