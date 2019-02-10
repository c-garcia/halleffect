package concourse

type JobStatus struct {
	Id           int
	TeamName     string
	JobName      string
	PipelineName string
	Status       string
}
