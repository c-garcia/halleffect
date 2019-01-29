package concourse

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

//go:generate mockgen -source=api.go -destination=mocks/api.go -package=mocks

type API interface {
	Name() string
	FindLastBuilds() ([]Build, error)
}

type ApiImpl struct {
	Concourse string
	URI       string
}

type BuildDTO struct {
	Id           int    `json:"id"`
	TeamName     string `json:"team_name"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	JobName      string `json:"job_name"`
	APIURL       string `json:"api_url"`
	PipelineName string `json:"pipeline_name"`
	StartTime    int    `json:"start_time"`
	EndTime      int    `json:"end_time"`
}

func dtoToBuild(b BuildDTO) Build {
	return Build{
		Id:           b.Id,
		StartTime:    b.StartTime,
		EndTime:      b.EndTime,
		PipelineName: b.PipelineName,
		JobName:      b.JobName,
		Status:       b.Status,
	}
}

func (s *ApiImpl) FindLastBuilds() ([]Build, error) {

	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, s.URI, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-type", "application/json")
	if err != nil {
		return nil, errors.Wrap(err, "Impossible to get build list")
	}
	resp, err := client.Do(req)

	if err != nil {
		return nil, errors.Wrap(err, "Impossible to get build list")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Impossible to get build list")
	}

	var jsonBuilds []BuildDTO

	err = json.NewDecoder(resp.Body).Decode(&jsonBuilds)
	if err != nil {
		return nil, errors.Wrap(err, "Impossible to get build list")
	}
	res := make([]Build, len(jsonBuilds))
	for i, item := range jsonBuilds {
		res[i] = dtoToBuild(item)
	}
	return res, nil
}

func (s *ApiImpl) Name() string {
	return s.Concourse
}

func New(name string, uri string) *ApiImpl {
	return &ApiImpl{Concourse: name, URI: uri}
}
