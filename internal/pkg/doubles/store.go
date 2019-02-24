package doubles

import (
	"fmt"
	"github.com/c-garcia/halleffect/internal/pkg/metrics"
	"strings"
)

type JobLastSuccessfulDurationInMemory struct {
	metrics map[metrics.JobLastSuccessfulDuration]bool
}

func (s *JobLastSuccessfulDurationInMemory) Store(m metrics.JobLastSuccessfulDuration) error {
	s.metrics[m] = true
	return nil
}

func (s *JobLastSuccessfulDurationInMemory) JobLastSuccessfulDurationHasBeenPublished(m metrics.JobLastSuccessfulDuration) bool {
	return s.metrics[m]
}

func (s *JobLastSuccessfulDurationInMemory) Size() int {
	return len(s.metrics)
}

func (s *JobLastSuccessfulDurationInMemory) String() string {
	buff := strings.Builder{}
	buff.WriteString("*** In memory Store ***\n")
	for k, _ := range s.metrics {
		buff.WriteString(fmt.Sprintf("%v\n", k))
	}
	return buff.String()
}

func NewJobLastSuccessfulDurationInMemory() *JobLastSuccessfulDurationInMemory {
	return &JobLastSuccessfulDurationInMemory{
		metrics: make(map[metrics.JobLastSuccessfulDuration]bool),
	}
}
