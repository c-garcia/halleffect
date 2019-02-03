package publisher

import (
	"github.com/c-garcia/halleffect/internal/pkg/concourse"
)

type JobDurationMetric struct {
	Timestamp    int
	EndTime      int
	PipelineName string
	JobName      string
	Status       string
	Concourse    string
	TeamName     string
}

func (m JobDurationMetric) Duration() int {
	return m.EndTime - m.Timestamp
}

func FromConcourseBuild(name string, b concourse.Build) JobDurationMetric {
	return JobDurationMetric{
		Concourse:    name,
		PipelineName: b.PipelineName,
		JobName:      b.JobName,
		Status:       b.Status,
		Timestamp:    b.StartTime,
		EndTime:      b.EndTime,
		TeamName:     b.TeamName,
	}
}
