//+build !integration, !service

package concourse

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	sut := New("concourse", "http://concourse.example.com")
	assert.NotNil(t, sut)
}

func TestNew_PanicsIfBadURL(t *testing.T) {
	var testCases = []struct {
		URI     string
		Message string
	}{
		{URI: "hk>123", Message: "Malformed URL"},
		{URI: "", Message: "Empty URL"},
		{URI: "a", Message: "Relative URL is not valid"},
		{URI: "http://", Message: "No host present is not valid"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Message, func(st *testing.T) {
			assert.Panics(st, func() { New("concourse", testCase.URI) })
		})
	}
}

func TestJobDTO_Unmarshal(t *testing.T) {
	var jsonNoBuilds = `{
        "id": 5,
        "name": "more-show-time",
        "pipeline_name": "p3",
        "team_name": "main",
        "next_build": null,
        "finished_build": null,
        "inputs": [
            {
                "name": "once-every-5m",
                "resource": "once-every-5m",
                "trigger": true
            }
        ],
        "outputs": [],
        "groups": []
    }`
	var jsonFinishedBuild = `{
        "id": 3,
        "name": "integration-tests",
        "pipeline_name": "p2",
        "team_name": "main",
        "next_build": null,
        "finished_build": {
            "id": 130,
            "team_name": "main",
            "name": "7",
            "status": "succeeded",
            "job_name": "integration-tests",
            "api_url": "/api/v1/builds/130",
            "pipeline_name": "p2",
            "start_time": 1549136078,
            "end_time": 1549136080
        },
        "transition_build": {
            "id": 3,
            "team_name": "main",
            "name": "1",
            "status": "succeeded",
            "job_name": "integration-tests",
            "api_url": "/api/v1/builds/3",
            "pipeline_name": "p2",
            "start_time": 1548970000,
            "end_time": 1548970013
        },
        "inputs": [
            {
                "name": "micro-node",
                "resource": "micro-node",
                "passed": [
                    "build-node"
                ],
                "trigger": true
            }
        ],
        "outputs": [],
        "groups": []
    }`
	var testCases = []struct {
		Desc string
		Json string
		DTO  JobDTO
	}{
		{"No builds job", jsonNoBuilds, JobDTO{Id: 5, Name: "more-show-time", PipelineName: "p3", TeamName: "main", FinishedBuild: nil}},
		{"Job with finished build", jsonFinishedBuild,
			JobDTO{Id: 3, Name: "integration-tests",
				PipelineName: "p2", TeamName: "main", FinishedBuild: &BuildDTO{
					Id: 130, TeamName: "main", Name: "7", Status: "succeeded", JobName: "integration-tests", APIURL: "/api/v1/builds/130", PipelineName: "p2",
					StartTime: 1549136078, EndTime: 1549136080,
				},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Desc, func(t *testing.T) {
			var res JobDTO
			err := json.Unmarshal([]byte(testCase.Json), &res)
			assert.NoError(t, err)
			assert.Equal(t, testCase.DTO, res)
		})
	}
}
