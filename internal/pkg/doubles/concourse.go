package doubles

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/c-garcia/halleffect/internal/pkg/concourse"
	"github.com/durmaze/gobank"
	"github.com/durmaze/gobank/predicates"
	"github.com/durmaze/gobank/responses"
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

func GivenAFailingCouncourseServer(port int) {
	mbClient := gobank.NewClient(MounteBankURL())
	imposter := gobank.NewImposterBuilder().Port(port).Protocol("http").
		Stubs(
			gobank.Stub().Responses(
				responses.Is().StatusCode(http.StatusInternalServerError).
					Header("Content-type", "application/json").
					Body("").Build()).Build(),
		).Build()
	_, err := mbClient.CreateImposter(imposter)
	if err != nil {
		panic(err)
	}
}

func ShutdownConcourseServer(port int) {
	mbClient := gobank.NewClient(MounteBankURL())
	if _, err := mbClient.DeleteImposter(port); err != nil {
		panic(err)
	}
}

func GivenACouncourseServerWithJobs(port int, jobs []concourse.JobDTO) {
	buff := bytes.Buffer{}
	if err := json.NewEncoder(&buff).Encode(jobs); err != nil {
		panic(err)
	}
	imposter := gobank.NewImposterBuilder().
		Protocol("http").
		Port(port).
		Stubs(
			gobank.Stub().
				Predicates(
					predicates.Equals().Path("/api/v1/jobs").Build(),
				).
				Responses(
					responses.Is().Body(buff.String()).Build(),
				).Build(),
			gobank.Stub().
				Predicates(
					predicates.Equals().Path("/api/v1/builds").Build(),
				).
				Responses(
					responses.Is().Body("[]").Build(),
				).Build(),
			gobank.Stub().
				Predicates().
				Responses(
					responses.Is().StatusCode(http.StatusNotFound).Body("Not found").Build(),
				).Build(),
		).
		Build()
	client := gobank.NewClient(MounteBankURL())
	if _, err := client.CreateImposter(imposter); err != nil {
		panic(err)
	}
}

func GivenACouncourseServerWithBuilds(port int, builds []concourse.BuildDTO) {
	buff := bytes.Buffer{}
	if err := json.NewEncoder(&buff).Encode(builds); err != nil {
		panic(err)
	}
	imposter := gobank.NewImposterBuilder().
		Protocol("http").
		Port(port).
		Stubs(
			gobank.Stub().
				Predicates(
					predicates.Equals().Path("/api/v1/builds").Build(),
				).
				Responses(
					responses.Is().Body(buff.String()).Build(),
				).Build(),
			gobank.Stub().
				Predicates(
					predicates.Equals().Path("/api/v1/jobs").Build(),
				).
				Responses(
					responses.Is().Body("[]").Build(),
				).Build(),
			gobank.Stub().
				Predicates().
				Responses(
					responses.Is().StatusCode(http.StatusNotFound).Body("Not found").Build(),
				).Build(),
		).
		Build()
	client := gobank.NewClient(MounteBankURL())
	if _, err := client.CreateImposter(imposter); err != nil {
		panic(err)
	}
}

func GivenAConcourseServerNotSupportingJobs(port int) {
	imposter := gobank.NewImposterBuilder().
		Protocol("http").
		Port(port).
		Stubs(
			gobank.Stub().
				Predicates(
					predicates.Equals().Path("/api/v1/jobs").Build(),
				).
				Responses(
					responses.Is().StatusCode(http.StatusNotFound).Body("Not found").Build(),
				).Build(),
			gobank.Stub().
				Predicates().
				Responses(
					responses.Is().StatusCode(http.StatusNotFound).Body("Not found").Build(),
				).Build(),
		).
		Build()
	client := gobank.NewClient(MounteBankURL())
	if _, err := client.CreateImposter(imposter); err != nil {
		panic(err)
	}
}

func GivenAConcourseServerSupportingJobs(port int) {
	imposter := gobank.NewImposterBuilder().
		Protocol("http").
		Port(port).
		Stubs(
			gobank.Stub().
				Predicates(
					predicates.Equals().Path("/api/v1/jobs").Build(),
				).
				Responses(
					responses.Is().Body("[]").Build(),
				).Build(),
			gobank.Stub().
				Predicates().
				Responses(
					responses.Is().StatusCode(http.StatusNotFound).Body("Not found").Build(),
				).Build(),
		).
		Build()
	client := gobank.NewClient(MounteBankURL())
	if _, err := client.CreateImposter(imposter); err != nil {
		panic(err)
	}

}

func GivenNoConcourseServer(){
	return
}
