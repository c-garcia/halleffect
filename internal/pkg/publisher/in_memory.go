package publisher

import (
	"fmt"
	"strings"
)

type InMemoryImpl struct {
	JobDurationMetrics map[JobDurationMetric]bool
	publications       int
}

func (p *InMemoryImpl) PublishJobDuration(m JobDurationMetric) error {
	p.publications++
	p.JobDurationMetrics[m] = true
	return nil
}

func (p *InMemoryImpl) JobDurationHasBeenPublished(m JobDurationMetric) bool {
	return p.JobDurationMetrics[m]
}

func (p *InMemoryImpl) NumberOfPublishedJobDurationMetrics() int {
	return p.publications
}

func (p *InMemoryImpl) String() string {
	builder := &strings.Builder{}
	builder.WriteString("In memory publisher\n")
	for k := range p.JobDurationMetrics {
		builder.WriteString(fmt.Sprintf("* %v\n", k))
	}
	return builder.String()
}

func NewInMemory() *InMemoryImpl {
	return &InMemoryImpl{
		JobDurationMetrics: make(map[JobDurationMetric]bool),
	}

}