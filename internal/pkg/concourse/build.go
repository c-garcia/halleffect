package concourse

import "time"

type Build struct {
	Id           int
	StartTime    time.Time
	EndTime      time.Time
	PipelineName string
	JobName      string
	Status       string
	TeamName     string
}

const (
	StatusSucceeded = "succeeded"
	StatusStarted   = "started"
	StatusFailed    = "failed"
	StatusErrored   = "errored"
)

func (b Build) Finished() bool {
	return !b.EndTime.IsZero()
}

func (b Build) Duration() time.Duration {
	return b.EndTime.Sub(b.StartTime)
}

func (b Build) Succeeded() bool {
	return b.Status == StatusSucceeded
}
