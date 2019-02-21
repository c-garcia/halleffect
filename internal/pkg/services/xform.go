package services

import (
	"github.com/c-garcia/halleffect/internal/pkg/concourse"
	"github.com/c-garcia/halleffect/internal/pkg/metrics"
	"time"
)

func BuildToDurationMetric(concourse string, when time.Time, b concourse.Build) metrics.JobLastSuccessfulDuration {
	unixDateDiff := func(a, b int) time.Duration {
		return time.Unix(int64(a), 0).Sub(time.Unix(int64(b), 0))
	}
	return metrics.JobLastSuccessfulDuration{
		Timestamp: when,
		Concourse: concourse,
		Team:      b.TeamName,
		Pipeline:  b.PipelineName,
		Job:       b.JobName,
		Duration:  unixDateDiff(b.EndTime, b.StartTime),
	}
}
