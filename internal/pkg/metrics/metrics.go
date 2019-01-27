package metrics

import (
	"github.com/c-garcia/halleffect/internal/pkg/concourse"
)

type Metric struct {
	Timestamp int
	EndTime int
}

func (m Metric) Duration() int {
	return m.EndTime - m.Timestamp
}

func FromConcourseBuild(b concourse.Build) Metric{
	return Metric{
		Timestamp: b.StartTime,
		EndTime: b.EndTime,
	}
}

