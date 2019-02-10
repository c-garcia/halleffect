//+build integration

package concourse_test

import (
	"fmt"
	"github.com/c-garcia/halleffect/internal/pkg/concourse"
	"github.com/c-garcia/halleffect/internal/pkg/doubles"
	"github.com/durmaze/gobank"
	"github.com/durmaze/gobank/predicates"
	"github.com/durmaze/gobank/responses"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

const (
	PORT = 3500
)

var (
	imposter gobank.ImposterElement
)

func givenACouncourseServer(port int, jsonText string) {
	imposter = gobank.NewImposterBuilder().Protocol("http").Port(port).Stubs(
		gobank.Stub().
			Predicates(
				predicates.Equals().Method(http.MethodGet).Build(),
				predicates.Equals().Path("/api/v1/builds").Build(),
			).Responses(
			responses.Is().
				StatusCode(200).
				Header("Content-type", "application/json").
				Body(jsonText).
				Build(),
		).Build(),
	).Build()
	if _, err := mbClient.CreateImposter(imposter); err != nil {
		panic(errors.Wrap(err, "Error setting up mountebank"))
	}
}

func shutdownServer() {
	mbClient.DeleteImposter(PORT)
}

func TestAPI_FindLastBuilds_RetrievesLastBuilds(t *testing.T) {
	defer shutdownServer()
	buildsJSON := `[
{
    "id": 160,
    "team_name": "main",
    "name": "14",
    "status": "started",
    "job_name": "build-node",
    "api_url": "/api/v1/buildsJSON/160",
    "pipeline_name": "p2",
    "start_time": 1548573140
},
{
    "id": 159,
    "team_name": "main",
    "name": "136",
    "status": "succeeded",
    "job_name": "show-time",
    "api_url": "/api/v1/buildsJSON/159",
    "pipeline_name": "p1",
    "start_time": 1548573115,
    "end_time": 1548573122
},
{
    "id": 158,
    "team_name": "main",
    "name": "135",
    "status": "failed",
    "job_name": "show-time",
    "api_url": "/api/v1/buildsJSON/158",
    "pipeline_name": "p1",
    "start_time": 1548573055,
    "end_time": 1548573063
}
]
`
	givenACouncourseServer(PORT, buildsJSON)
	sut := concourse.New("test", fmt.Sprintf("http://%s:%d/", getMBHost(), PORT))
	builds, err := sut.FindLastBuilds()

	assert.NoError(t, err)
	assert.Len(t, builds, 3)
	assert.Equal(t, builds[0], concourse.Build{
		Id: 160, StartTime: 1548573140, EndTime: 0, PipelineName: "p2", JobName: "build-node", Status: "started", TeamName: "main",
	})
	assert.Equal(t, builds[1], concourse.Build{
		Id: 159, StartTime: 1548573115, EndTime: 1548573122, PipelineName: "p1", JobName: "show-time", Status: "succeeded", TeamName: "main",
	})
	assert.Equal(t, builds[2], concourse.Build{
		Id: 158, StartTime: 1548573055, EndTime: 1548573063, PipelineName: "p1", JobName: "show-time", Status: "failed", TeamName: "main",
	})
}

func givenAFailingCouncourseServer(port int) {
	imposter = gobank.NewImposterBuilder().Protocol("http").Port(port).Stubs(
		gobank.Stub().Responses(
			responses.Is().
				StatusCode(500).
				Header("Content-type", "application/json").
				Body("").
				Build(),
		).Build(),
	).Build()
	if _, err := mbClient.CreateImposter(imposter); err != nil {
		panic(errors.Wrap(err, "Error setting up mountebank"))
	}
}

func TestAPI_FindLastBuilds_PropagatesServerErrors(t *testing.T) {
	defer shutdownServer()

	givenAFailingCouncourseServer(PORT)
	sut := concourse.New("test", fmt.Sprintf("http://%s:%d/", getMBHost(), PORT))
	_, err := sut.FindLastBuilds()

	assert.Error(t, err)
}

func givenNoConcourseServer() {
	return
}

func TestAPI_FindLastBuilds_PropagatesConnectionErrors(t *testing.T) {
	defer shutdownServer()

	givenNoConcourseServer()
	sut := concourse.New("test", fmt.Sprintf("http://%s:%d/", getMBHost(), PORT))
	_, err := sut.FindLastBuilds()

	assert.Error(t, err)
}

func TestApi_FindJobStatuses(t *testing.T) {
	numberOfJobsWithFinishedBuilds := func(jobs []concourse.JobDTO) int {
		res := 0
		for _, job := range jobs {
			if job.FinishedBuild != nil {
				res++
			}
		}
		return res
	}
	successFulJobJSON := concourse.JobDTO{
		Id:           1,
		Name:         "job-1",
		PipelineName: "p1",
		TeamName:     "main",
		FinishedBuild: &concourse.BuildDTO{
			Id:           99,
			TeamName:     "main",
			Name:         "99",
			Status:       "succeeded",
			JobName:      "job-1",
			APIURL:       "/api/v1/jobs/99",
			PipelineName: "p1",
			StartTime:    1000,
			EndTime:      1100,
		},
	}

	successFulJob := concourse.JobStatus{
		Id:           1,
		TeamName:     "main",
		JobName:      "job-1",
		PipelineName: "p1",
		Status:       "succeeded",
	}

	failedJobJSON := concourse.JobDTO{
		Id:           2,
		Name:         "job-2",
		PipelineName: "p2",
		TeamName:     "main",
		FinishedBuild: &concourse.BuildDTO{
			Id:           98,
			TeamName:     "main",
			Name:         "98",
			Status:       "failed",
			JobName:      "job-2",
			APIURL:       "/api/v1/jobs/98",
			PipelineName: "p2",
			StartTime:    1101,
			EndTime:      1200,
		},
	}

	failedJob := concourse.JobStatus{
		Id:           2,
		TeamName:     "main",
		JobName:      "job-2",
		PipelineName: "p2",
		Status:       "failed",
	}

	neverRunJob := concourse.JobDTO{
		Id:            3,
		Name:          "job-3",
		PipelineName:  "p2",
		TeamName:      "main",
		FinishedBuild: nil,
	}

	jobsJSON := []concourse.JobDTO{successFulJobJSON, failedJobJSON, neverRunJob}
	doubles.GivenACouncourseServerWithJobs(PORT, jobsJSON)
	defer doubles.ShutdownConcourseServer(PORT)
	sut := concourse.New("test", doubles.ImposterURL(PORT))
	metrics, err := sut.FindJobStatuses()
	assert.NoError(t, err, "There should be no error")
	assert.Contains(t, metrics, successFulJob, "The successful job belongs to the returned ones")
	assert.Contains(t, metrics, failedJob, "The failed job belongs to the returned ones")
	assert.Len(t, metrics, numberOfJobsWithFinishedBuilds(jobsJSON), "There should be only jobs with finished builds")
}
