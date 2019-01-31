//+build integration

package concourse_test

import (
	"fmt"
	"github.com/c-garcia/halleffect/internal/pkg/concourse"
	"github.com/durmaze/gobank"
	"github.com/durmaze/gobank/responses"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
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
		gobank.Stub().Responses(
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
		Id: 160, StartTime: 1548573140, EndTime: 0, PipelineName: "p2", JobName: "build-node", Status: "started",
	})
	assert.Equal(t, builds[1], concourse.Build{
		Id: 159, StartTime: 1548573115, EndTime: 1548573122, PipelineName: "p1", JobName: "show-time", Status: "succeeded",
	})
	assert.Equal(t, builds[2], concourse.Build{
		Id: 158, StartTime: 1548573055, EndTime: 1548573063, PipelineName: "p1", JobName: "show-time", Status: "failed",
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
