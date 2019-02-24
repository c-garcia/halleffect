// +build service
// +build !integration

package service_test

import (
	"github.com/c-garcia/halleffect/internal/pkg/concourse"
	"github.com/c-garcia/halleffect/internal/pkg/doubles"
	"github.com/c-garcia/halleffect/internal/pkg/metrics"
	"github.com/c-garcia/halleffect/internal/pkg/services"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const ConcoursePort = 3000

func TestPublishesLastSuccessfulBuildDurationForEachJob(t *testing.T) {
	beginningOfTime := time.Now().Truncate(time.Second)
	jobLength := 180 * time.Second
	samplingTime := beginningOfTime
	const Team = "main"
	const Pipeline = "p1"
	const Job = "job1"
	build1DTO := concourse.BuildDTO{
		Id:           100,
		TeamName:     Team,
		Name:         "123",
		Status:       "succeeded",
		JobName:      Job,
		APIURL:       "/api/v1/builds/100",
		PipelineName: Pipeline,
		StartTime:    int(beginningOfTime.Unix()),
		EndTime:      int(beginningOfTime.Add(jobLength).Unix()),
	}
	const Concourse = "test-concourse"
	build1Metric := metrics.JobLastSuccessfulDuration{
		Timestamp: samplingTime,
		Concourse: Concourse,
		Team:      Team,
		Pipeline:  Pipeline,
		Job:       Job,
		Duration:  jobLength,
	}
	builds := []concourse.BuildDTO{
		build1DTO,
	}
	doubles.GivenACouncourseServerWithBuilds(ConcoursePort, builds)
	defer doubles.ShutdownConcourseServer(ConcoursePort)
	store := doubles.NewJobLastSuccessfulDurationInMemory()
	sut := services.LastSuccessfulBuildDurationImpl{
		Name:         Concourse,
		Clock:        doubles.NewStoppedClock(samplingTime.Unix()),
		Concourse:    concourse.New(Concourse, doubles.ImposterURL(ConcoursePort)),
		MetricsStore: store,
	}
	err := sut.SaveMetrics()
	assert.NoError(t, err)
	assert.True(t, store.JobLastSuccessfulDurationHasBeenPublished(build1Metric))
}
