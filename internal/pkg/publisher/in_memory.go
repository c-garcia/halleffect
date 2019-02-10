package publisher

import (
	"fmt"
	"strings"
)

type InMemoryImpl struct {
	JobDurationMetrics map[JobDurationMetric]bool
	JobStatusMetrics   map[JobStatusMetric]bool
	publications       int
	statusPublications int
}

func (p *InMemoryImpl) PublishJobStatus(m JobStatusMetric) error {
	p.JobStatusMetrics[m] = true
	p.statusPublications++
	return nil
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
	builder.WriteString("Durations\n")
	for k := range p.JobDurationMetrics {
		builder.WriteString(fmt.Sprintf("* %v\n", k))
	}
	builder.WriteString("\nStatuses\n")
	for k := range p.JobStatusMetrics {
		builder.WriteString(fmt.Sprintf("* %v\n", k))
	}
	return builder.String()
}

func (p *InMemoryImpl) JobStatusHasBeenPublished(metric JobStatusMetric) bool {
	return p.JobStatusMetrics[metric]
}

func (p *InMemoryImpl) NumberOfPublishedJobStatusMetrics() int {
	return p.statusPublications
}

func NewInMemory() *InMemoryImpl {
	return &InMemoryImpl{
		JobDurationMetrics: make(map[JobDurationMetric]bool),
		JobStatusMetrics:   make(map[JobStatusMetric]bool),
	}

}
