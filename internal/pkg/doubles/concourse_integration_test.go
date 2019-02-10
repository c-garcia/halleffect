//+build integration

package doubles_test

import (
	"encoding/json"
	"github.com/bxcodec/faker"
	"github.com/c-garcia/halleffect/internal/pkg/concourse"
	"github.com/c-garcia/halleffect/internal/pkg/doubles"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

const IMPOSTER_PORT = 5555

func TestGivenACouncourseServerWithBuilds(t *testing.T) {
	var builds []concourse.BuildDTO
	err := faker.FakeData(&builds)
	assert.NoError(t, err, "Unexpected error faking data")
	doubles.GivenACouncourseServerWithBuilds(IMPOSTER_PORT, builds)
	defer doubles.ShutdownConcourseServer(IMPOSTER_PORT)
	returnedJobs := getJobs(t, IMPOSTER_PORT)
	assert.Empty(t, returnedJobs, "Should return an empty job list")
	returnedBuilds := getBuilds(t, IMPOSTER_PORT)
	assert.Equal(t, builds, returnedBuilds, "Should return the passwed builds")
}

func TestGivenACouncourseServerWithBuilds_Returns404ForNonAPI(t *testing.T) {
	var builds []concourse.BuildDTO
	err := faker.FakeData(&builds)
	assert.NoError(t, err, "Unexpected error faking data")
	doubles.GivenACouncourseServerWithBuilds(IMPOSTER_PORT, builds)
	defer doubles.ShutdownConcourseServer(IMPOSTER_PORT)
	assert.Equal(t, http.StatusNotFound, statusCodeFor(t, "/nothing-relevant", IMPOSTER_PORT))
}

func TestGivenACouncourseServerWithJobs(t *testing.T) {
	var jobs []concourse.JobDTO
	err := faker.FakeData(&jobs)
	assert.NoError(t, err, "Unexpected error faking data")
	doubles.GivenACouncourseServerWithJobs(IMPOSTER_PORT, jobs)
	defer doubles.ShutdownConcourseServer(IMPOSTER_PORT)
	returnedJobs := getJobs(t, IMPOSTER_PORT)
	assert.Equal(t, jobs, returnedJobs, "Should return passed jobs")
	returnedBuilds := getBuilds(t, IMPOSTER_PORT)
	assert.Empty(t, returnedBuilds, "Should return an empty build list")
}

func TestGivenACouncourseServerWithJobs_Returns404ForNonAPI(t *testing.T) {
	var jobs []concourse.JobDTO
	err := faker.FakeData(&jobs)
	assert.NoError(t, err, "Unexpected error faking data")
	doubles.GivenACouncourseServerWithJobs(IMPOSTER_PORT, jobs)
	defer doubles.ShutdownConcourseServer(IMPOSTER_PORT)
	assert.Equal(t, http.StatusNotFound, statusCodeFor(t, "/nothing-relevant", IMPOSTER_PORT))
}

func getJobs(t *testing.T, port int) []concourse.JobDTO {
	imposter := doubles.ImposterURL(port)
	rootUrl, _ := url.Parse(imposter)
	rootUrl.Path = "/api/v1/jobs"
	req, err := http.NewRequest(http.MethodGet, rootUrl.String(), nil)
	req.Header.Set("Accept", "application/json")
	assert.NoError(t, err)
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err, "Unexpected error contacting imposter endpoint")
	var returnedJobs []concourse.JobDTO
	err = json.NewDecoder(resp.Body).Decode(&returnedJobs)
	assert.NoError(t, err, "There should be no error parsing JSON")
	return returnedJobs
}

func getBuilds(t *testing.T, port int) []concourse.BuildDTO {
	imposter := doubles.ImposterURL(port)
	rootUrl, _ := url.Parse(imposter)
	rootUrl.Path = "/api/v1/builds"
	req, err := http.NewRequest(http.MethodGet, rootUrl.String(), nil)
	req.Header.Set("Accept", "application/json")
	assert.NoError(t, err)
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err, "Unexpected error contacting imposter endpoint")
	var returnedBuilds []concourse.BuildDTO
	err = json.NewDecoder(resp.Body).Decode(&returnedBuilds)
	assert.NoError(t, err, "There should be no error parsing JSON")
	return returnedBuilds
}

func statusCodeFor(t *testing.T, path string, port int) int {
	imposter := doubles.ImposterURL(port)
	rootUrl, _ := url.Parse(imposter)
	rootUrl.Path = path
	req, err := http.NewRequest(http.MethodGet, rootUrl.String(), nil)
	req.Header.Set("Accept", "application/json")
	assert.NoError(t, err)
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err, "Unexpected error contacting imposter endpoint")
	return resp.StatusCode
}
