package metrics

import (
	"github.com/c-garcia/halleffect/internal/pkg/concourse"
)

type Metric struct {
	Timestamp    int
	EndTime      int
	PipelineName string
	JobName      string
	Status       string
	Concourse    string
}

func (m Metric) Duration() int {
	return m.EndTime - m.Timestamp
}

func FromConcourseBuild(name string, b concourse.Build) Metric {
	return Metric{
		Concourse:    name,
		PipelineName: b.PipelineName,
		JobName:      b.JobName,
		Status:       b.Status,
		Timestamp:    b.StartTime,
		EndTime:      b.EndTime,
	}
}
