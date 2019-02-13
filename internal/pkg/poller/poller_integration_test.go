//+build integration

package poller_test

import (
	"fmt"
	"github.com/c-garcia/halleffect/internal/pkg/concourse"
	"github.com/c-garcia/halleffect/internal/pkg/doubles"
	"github.com/c-garcia/halleffect/internal/pkg/poller"
	"github.com/c-garcia/halleffect/internal/pkg/publisher"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const ImposterPort = 3333
const ConcourseName = "test"

var SamplingTime = time.Now().Unix()

func TestExportMetrics_ExportsJobDurations(t *testing.T) {
	successfulBuild := concourse.BuildDTO{
		Id:           1,
		TeamName:     "main",
		Name:         "1",
		Status:       "succeeded",
		JobName:      "j1",
		APIURL:       "/api/v1/builds/1",
		PipelineName: "p1",
		StartTime:    1000,
		EndTime:      1100,
	}
	failedBuild := concourse.BuildDTO{
		Id:           2,
		TeamName:     "main",
		Name:         "2",
		Status:       "failed",
		JobName:      "j2",
		APIURL:       "/api/v1/builds/2",
		PipelineName: "p2",
		StartTime:    1200,
		EndTime:      1300,
	}
	successfulBuildMetric := publisher.JobDurationMetric{
		Timestamp:    1000,
		EndTime:      1100,
		PipelineName: "p1",
		JobName:      "j1",
		Status:       "succeeded",
		Concourse:    ConcourseName,
		TeamName:     "main",
	}

	failedBuildMetric := publisher.JobDurationMetric{
		Timestamp:    1200,
		EndTime:      1300,
		PipelineName: "p2",
		JobName:      "j2",
		Status:       "failed",
		Concourse:    ConcourseName,
		TeamName:     "main",
	}
	builds := []concourse.BuildDTO{successfulBuild, failedBuild}
	expectedMetrics := []publisher.JobDurationMetric{successfulBuildMetric, failedBuildMetric}
	doubles.GivenACouncourseServerWithBuilds(ImposterPort, builds)
	defer doubles.ShutdownConcourseServer(ImposterPort)
	metricsPublisher := publisher.NewInMemory()
	sut := poller.New(
		concourse.New(ConcourseName, doubles.ImposterURL(ImposterPort)),
		metricsPublisher,
		doubles.NewStoppedClock(SamplingTime),
	)

	err := sut.ExportJobDurationMetrics()

	assert.NoError(t, err)
	for _, metric := range expectedMetrics {
		assert.True(t, metricsPublisher.JobDurationHasBeenPublished(metric), "Metric should have been exported: %v", metric)
	}
}

func TestExportMetrics_WhenConcourseFails(t *testing.T) {
	doubles.GivenAFailingCouncourseServer(ImposterPort)
	defer doubles.ShutdownConcourseServer(ImposterPort)
	metricsPublisher := publisher.NewInMemory()
	sut := poller.New(
		concourse.New("test", doubles.ImposterURL(ImposterPort)),
		metricsPublisher,
		doubles.NewStoppedClock(SamplingTime),
	)

	err := sut.ExportJobDurationMetrics()

	assert.Error(t, err)
	assert.Equal(t, 0, metricsPublisher.NumberOfPublishedJobDurationMetrics())
}

func TestExportMetrics_ExportsJobStatus(t *testing.T) {
	const SamplingTime = 1000

	jobNoBuilds := concourse.JobDTO{
		Id:            1,
		Name:          "job-1",
		PipelineName:  "p-1",
		TeamName:      "main",
		FinishedBuild: nil,
	}
	jobSucceeded := concourse.JobDTO{
		Id:           2,
		Name:         "job-2",
		PipelineName: "p-2",
		TeamName:     "main",
		FinishedBuild: &concourse.BuildDTO{
			Id:           2,
			TeamName:     "main",
			Name:         "2",
			Status:       "succeeded",
			JobName:      "job-2",
			APIURL:       "/ap1/v1/builds/2",
			PipelineName: "p-2",
			StartTime:    1000,
			EndTime:      1010,
		},
	}
	jobFailed := concourse.JobDTO{
		Id:           3,
		Name:         "job-3",
		PipelineName: "p-3",
		TeamName:     "main",
		FinishedBuild: &concourse.BuildDTO{
			Id:           3,
			TeamName:     "main",
			Name:         "3",
			Status:       "failed",
			JobName:      "job-3",
			APIURL:       "/ap1/v1/builds/3",
			PipelineName: "p-3",
			StartTime:    1020,
			EndTime:      1021,
		},
	}
	jobSucceededStatus := publisher.JobStatusMetric{
		Concourse:    "test",
		TeamName:     "main",
		PipelineName: "p-2",
		JobName:      "job-2",
		Status:       "up",
		SamplingTime: SamplingTime,
	}
	jobFailedStatus := publisher.JobStatusMetric{
		Concourse:    "test",
		TeamName:     "main",
		PipelineName: "p-3",
		JobName:      "job-3",
		Status:       "down",
		SamplingTime: SamplingTime,
	}

	doubles.GivenACouncourseServerWithJobs(ImposterPort, []concourse.JobDTO{jobNoBuilds, jobSucceeded, jobFailed})
	defer doubles.ShutdownConcourseServer(ImposterPort)
	metricsPublisher := publisher.NewInMemory()
	sut := poller.New(
		concourse.New("test", doubles.ImposterURL(ImposterPort)),
		metricsPublisher,
		doubles.NewStoppedClock(SamplingTime),
	)

	err := sut.ExportJobDurationMetrics()

	fmt.Printf("In-Memory Publisher content: %s\n", metricsPublisher.String())
	assert.NoError(t, err)
	assert.True(t, metricsPublisher.JobStatusHasBeenPublished(jobSucceededStatus))
	assert.True(t, metricsPublisher.JobStatusHasBeenPublished(jobFailedStatus))
	assert.Equal(t, 2, metricsPublisher.NumberOfPublishedJobStatusMetrics())
}
