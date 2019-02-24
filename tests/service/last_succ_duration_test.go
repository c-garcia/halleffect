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

func metricForBuild(at time.Time, from string, b concourse.BuildDTO) metrics.JobLastSuccessfulDuration {
	return metrics.JobLastSuccessfulDuration{
		Timestamp: at,
		Concourse: from,
		Team:      b.TeamName,
		Pipeline:  b.PipelineName,
		Job:       b.JobName,
		Duration:  time.Unix(int64(b.EndTime), 0).Sub(time.Unix(int64(b.StartTime), 0)),
	}
}

func TestPublishesADurationForEachJob(t *testing.T) {
	beginningOfTime := time.Now().Add(-24 * time.Hour).Truncate(time.Second)
	lastLobDuration := time.Duration(10 * time.Minute)
	lastLobTime := beginningOfTime
	lastJobEndTime := lastLobTime.Add(lastLobDuration)
	firstJobDuration := 2 * time.Minute
	firstJobTime := lastLobTime.Add(1 * time.Hour)
	firstJobEndTime := firstJobTime.Add(firstJobDuration)
	samplingTime := beginningOfTime
	const Concourse = "test-concourse"
	const Team = "main"
	const Pipeline1 = "p1"
	const Pipeline2 = "p2"
	const Job1 = "job1"
	const Job2 = "job1"
	firstDTO := concourse.BuildDTO{
		Id:           100,
		TeamName:     Team,
		Name:         "118",
		Status:       "succeeded",
		JobName:      Job1,
		APIURL:       "/api/v1/builds/101",
		PipelineName: Pipeline1,
		StartTime:    int(firstJobTime.Unix()),
		EndTime:      int(firstJobEndTime.Unix()),
	}
	lastDTO := concourse.BuildDTO{
		Id:           101,
		TeamName:     Team,
		Name:         "123",
		Status:       "succeeded",
		JobName:      Job2,
		APIURL:       "/api/v1/builds/100",
		PipelineName: Pipeline2,
		StartTime:    int(lastLobTime.Unix()),
		EndTime:      int(lastJobEndTime.Unix()),
	}
	builds := []concourse.BuildDTO{
		lastDTO,
		firstDTO,
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
	assert.True(t, store.JobLastSuccessfulDurationHasBeenPublished(metricForBuild(samplingTime, Concourse, firstDTO)))
	assert.True(t, store.JobLastSuccessfulDurationHasBeenPublished(metricForBuild(samplingTime, Concourse, lastDTO)))
}

func TestPublishesOnlyLastBuildDurationForEachJob(t *testing.T) {
	beginningOfTime := time.Now().Add(-24 * time.Hour).Truncate(time.Second)
	lastLobDuration := time.Duration(10 * time.Minute)
	lastLobTime := beginningOfTime
	lastJobEndTime := lastLobTime.Add(lastLobDuration)
	firstJobDuration := 2 * time.Minute
	firstJobTime := lastLobTime.Add(1 * time.Hour)
	firstJobEndTime := firstJobTime.Add(firstJobDuration)
	samplingTime := beginningOfTime
	const Concourse = "test-concourse"
	const Team = "main"
	const Pipeline = "p1"
	const Job = "job1"
	firstDTO := concourse.BuildDTO{
		Id:           100,
		TeamName:     Team,
		Name:         "118",
		Status:       "succeeded",
		JobName:      Job,
		APIURL:       "/api/v1/builds/101",
		PipelineName: Pipeline,
		StartTime:    int(firstJobTime.Unix()),
		EndTime:      int(firstJobEndTime.Unix()),
	}
	lastDTO := concourse.BuildDTO{
		Id:           101,
		TeamName:     Team,
		Name:         "123",
		Status:       "succeeded",
		JobName:      Job,
		APIURL:       "/api/v1/builds/100",
		PipelineName: Pipeline,
		StartTime:    int(lastLobTime.Unix()),
		EndTime:      int(lastJobEndTime.Unix()),
	}
	builds := []concourse.BuildDTO{
		lastDTO,
		firstDTO,
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
	assert.True(t, store.JobLastSuccessfulDurationHasBeenPublished(metricForBuild(samplingTime, Concourse, firstDTO)))
	assert.False(t, store.JobLastSuccessfulDurationHasBeenPublished(metricForBuild(samplingTime, Concourse, lastDTO)))
}

func TestSkipsNonSuccessfulJobs(t *testing.T) {
	const Team = "main"
	const Pipeline = "p1"
	const Job = "job1"
	const Concourse = "test-concourse"
	failures := []string{concourse.StatusErrored, concourse.StatusFailed, concourse.StatusAborted, concourse.StatusStarted}
	for _, test := range failures {
		t.Run(test, func(t *testing.T) {
			beginningOfTime := time.Now().Add(-24 * time.Hour).Truncate(time.Second)
			samplingTime := beginningOfTime
			failedJobTime := beginningOfTime
			failedJobDuration := 2 * time.Minute
			failedJobTimeEnd := failedJobTime.Add(failedJobDuration)
			failedDTO := concourse.BuildDTO{
				Id:           101,
				TeamName:     Team,
				Name:         "123",
				Status:       test,
				JobName:      Job,
				APIURL:       "/api/v1/builds/101",
				PipelineName: Pipeline,
				StartTime:    int(failedJobTime.Unix()),
				EndTime:      int(failedJobTimeEnd.Unix()),
			}
			builds := []concourse.BuildDTO{
				failedDTO,
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
			assert.Equal(t, 0, store.Size())
		})
	}
}
