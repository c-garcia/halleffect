//+build integration

package poller_test

import (
	"github.com/c-garcia/halleffect/internal/pkg/concourse"
	"github.com/c-garcia/halleffect/internal/pkg/doubles"
	"github.com/c-garcia/halleffect/internal/pkg/poller"
	"github.com/c-garcia/halleffect/internal/pkg/publisher"
	"github.com/stretchr/testify/assert"
	"testing"
)

const IMPOSTER_PORT = 3333

func TestExportMetrics(t *testing.T) {
	expectedMetrics := doubles.GivenAConcourseServer("test", IMPOSTER_PORT)
	defer doubles.ShutdownConcourseServer(IMPOSTER_PORT)
	metricsPublisher := publisher.NewInMemory()
	sut := poller.New(
		concourse.New("test", doubles.ImposterURL(IMPOSTER_PORT)),
		metricsPublisher,
	)

	err := sut.ExportMetrics()

	assert.NoError(t, err)
	for _, metric := range expectedMetrics {
		assert.True(t, metricsPublisher.JobDurationHasBeenPublished(metric), "Metric should have been exported: %v", metric)
	}
}

func TestExportMetrics_WhenConcourseFails(t *testing.T) {
	doubles.GivenAFailingCouncourseServer("test", IMPOSTER_PORT)
	defer doubles.ShutdownConcourseServer(IMPOSTER_PORT)
	metricsPublisher := publisher.NewInMemory()
	sut := poller.New(
		concourse.New("test", doubles.ImposterURL(IMPOSTER_PORT)),
		metricsPublisher,
	)

	err := sut.ExportMetrics()

	assert.Error(t, err)
	assert.Equal(t, 0, metricsPublisher.NumberOfPublishedJobDurationMetrics())
}
