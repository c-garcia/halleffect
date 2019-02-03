package concourse

type Build struct {
	Id           int
	StartTime    int
	EndTime      int
	PipelineName string
	JobName      string
	Status       string
	TeamName     string
}

func (b Build) Finished() bool {
	return b.EndTime == 0
}
