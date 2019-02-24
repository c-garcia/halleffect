//+build !integration, !service

package doubles

import (
	"github.com/c-garcia/halleffect/internal/pkg/metrics"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestJobLastSuccessfulDurationInMemory_JobLastSuccessfulDurationHasBeenPublished(t *testing.T) {
	now := time.Now()
	publishedMetric := metrics.JobLastSuccessfulDuration{
		Timestamp: now,
		Concourse: "concourse",
		Team:      "team",
		Pipeline:  "pipeline",
		Job:       "job",
		Duration:  120 * time.Second,
	}
	unpublishedMetric := metrics.JobLastSuccessfulDuration{
		Timestamp: now,
		Concourse: "concourse2",
		Team:      "team2",
		Pipeline:  "pipeline2",
		Job:       "job2",
		Duration:  120 * time.Second,
	}
	sut := NewJobLastSuccessfulDurationInMemory()
	err := sut.Store(publishedMetric)
	assert.NoError(t, err)
	assert.True(t, sut.JobLastSuccessfulDurationHasBeenPublished(publishedMetric))
	assert.False(t, sut.JobLastSuccessfulDurationHasBeenPublished(unpublishedMetric))
}

func TestJobLastSuccessfulDurationInMemory_Size(t *testing.T) {
	now := time.Now()
	publishedMetric := metrics.JobLastSuccessfulDuration{
		Timestamp: now,
		Concourse: "concourse",
		Team:      "team",
		Pipeline:  "pipeline",
		Job:       "job",
		Duration:  120 * time.Second,
	}
	sut := NewJobLastSuccessfulDurationInMemory()
	err := sut.Store(publishedMetric)
	assert.NoError(t, err)
	assert.Equal(t, 1, sut.Size())
}
