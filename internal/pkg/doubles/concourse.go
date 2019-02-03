package doubles

import (
	"fmt"
	"github.com/c-garcia/halleffect/internal/pkg/publisher"
	"github.com/durmaze/gobank"
	"github.com/durmaze/gobank/predicates"
	"github.com/durmaze/gobank/responses"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

func MounteBankURL() string {
	urlStr := "http://localhost:2525"
	if os.Getenv("MB_URL") != "" {
		urlStr = os.Getenv("MB_URL")
	}
	_, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}
	return urlStr
}

func ImposterURL(port int) string {
	portRe := regexp.MustCompile(`:\d{1,5}$`)
	mb, _ := url.Parse(MounteBankURL())
	if portRe.MatchString(mb.Host) {
		mb.Host = portRe.ReplaceAllString(mb.Host, fmt.Sprintf(":%d", port))
	} else {
		mb.Host = fmt.Sprintf("%s:%d", mb.Host, port)
	}
	return mb.String()
}

func GivenAConcourseServer(name string, port int) []publisher.JobDurationMetric {
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
	mbClient := gobank.NewClient(MounteBankURL())
	imposter := gobank.NewImposterBuilder().Protocol("http").Port(port).Stubs(
		gobank.Stub().
			Predicates(
				predicates.Equals().Method(http.MethodGet).Build(),
				predicates.Equals().Path("/api/v1/builds").Build(),
			).Responses(
			responses.Is().
				StatusCode(200).
				Header("Content-type", "application/json").
				Body(buildsJSON).
				Build(),
		).Build(),
	).Build()
	if _, err := mbClient.CreateImposter(imposter); err != nil {
		panic(errors.Wrap(err, "Error setting up mountebank"))
	}
	return []publisher.JobDurationMetric{
		publisher.JobDurationMetric{
			Timestamp: 1548573115, EndTime: 1548573122, PipelineName: "p1", JobName: "show-time", Status: "succeeded", TeamName: "main", Concourse: name,
		},
		publisher.JobDurationMetric{
			Timestamp: 1548573055, EndTime: 1548573063, PipelineName: "p1", JobName: "show-time", Status: "failed", TeamName: "main", Concourse: name,
		},
	}
}

func ShutdownConcourseServer(port int) {
	mbClient := gobank.NewClient(MounteBankURL())
	if _, err := mbClient.DeleteImposter(port); err != nil {
		panic(err)
	}
}
