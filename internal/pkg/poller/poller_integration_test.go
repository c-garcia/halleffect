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

func TestExportMetrics(t *testing.T) {
	const IMPOSTER_PORT = 3333
	expectedMetrics := doubles.GivenAConcourseServer("test", IMPOSTER_PORT)
	defer doubles.ShutdownConcourseServer(IMPOSTER_PORT)
	metricsPublisher := publisher.NewInMemory()
	sut := poller.New(
		concourse.New("test", doubles.ImposterURL(IMPOSTER_PORT)),
		metricsPublisher,
	)

	err := sut.ExportMetrics()

	assert.NoError(t, err)
	for _, metric := range expectedMetrics{
		assert.True(t, metricsPublisher.JobDurationHasBeenPubished(metric), "Metric should have been exported: %v", metric)
	}
}
