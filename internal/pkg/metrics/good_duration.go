package metrics

import "time"

type JobLastSuccessfulDuration struct {
	Timestamp time.Time
	Concourse string
	Team      string
	Pipeline  string
	Job       string
	Duration  time.Duration
}
