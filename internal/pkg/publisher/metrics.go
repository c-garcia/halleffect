package publisher

import (
	"github.com/c-garcia/halleffect/internal/pkg/concourse"
	"github.com/pkg/errors"
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

type JobStatusMetric struct {
	Concourse    string
	TeamName     string
	PipelineName string
	JobName      string
	Status       string
	SamplingTime int64
}

func FromConcourseJobStatus(name string, samplingTime int64, s concourse.JobStatus) JobStatusMetric {
	acceptableMetrics := map[string]string{
		"succeeded": "up",
		"failed":    "down",
		"errored":   "down",
		"aborted":   "down",
	}
	status, ok := acceptableMetrics[s.Status]
	if !ok {
		panic(errors.Errorf("Unacceptable status '%s'", s.Status))
	}
	return JobStatusMetric{
		Concourse:    name,
		TeamName:     s.TeamName,
		PipelineName: s.PipelineName,
		JobName:      s.JobName,
		Status:       status,
		SamplingTime: samplingTime,
	}
}
