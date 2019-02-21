package concourse

import "time"

type Build struct {
	Id           int
	StartTime    int
	EndTime      int
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
	return b.EndTime == 0
}

func (b Build) Duration() time.Duration {
	return time.Unix(int64(b.EndTime), 0).Sub(time.Unix(int64(b.StartTime), 0))
}

func (b Build) Succeeded() bool {
	return b.Status == StatusSucceeded
}
