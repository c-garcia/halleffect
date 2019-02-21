package services

import (
	"github.com/c-garcia/halleffect/internal/pkg/concourse"
	"github.com/c-garcia/halleffect/internal/pkg/store"
	"github.com/c-garcia/halleffect/internal/pkg/timing"
	"github.com/pkg/errors"
	"time"
)

type LastSuccessfulBuildDurationImpl struct {
	Name         string
	Clock        timing.Clock
	Concourse    concourse.API
	MetricsStore store.JobLastSuccessfulDuration
}

func (s *LastSuccessfulBuildDurationImpl) SaveMetrics() error {

	type JobId struct {
		Team     string
		Pipeline string
		Job      string
	}
	latest := make(map[JobId]concourse.Build)
	builds, err := s.Concourse.FindLastBuilds()
	if err != nil {
		return errors.Wrap(err, "Error retrieving builds")
	}
	shouldKeepBuild := func(alreadyExisting bool, newBuild, latestSofar concourse.Build) bool {
		return !alreadyExisting || (alreadyExisting && (newBuild.EndTime > latestSofar.EndTime))
	}
	for _, build := range builds {
		if build.Succeeded() {
			id := JobId{build.TeamName, build.PipelineName, build.JobName}
			if latestSoFar, ok := latest[id]; shouldKeepBuild(ok, build, latestSoFar) {
				latest[id] = build
			}
		}
	}
	now := s.Clock.UnixTime()
	for _, build := range latest {
		m := BuildToDurationMetric(s.Name, time.Unix(now, 0), build)
		err = s.MetricsStore.Store(m)
		if err != nil {
			return errors.Wrap(err, "Error publishing metric")
		}
	}
	return nil
}

func NewLastSuccessfulBuildDurationImpl(concourseName string, concourse concourse.API, store store.JobLastSuccessfulDuration, clock timing.Clock) *LastSuccessfulBuildDurationImpl {
	return &LastSuccessfulBuildDurationImpl{
		Name:         concourseName,
		Clock:        clock,
		Concourse:    concourse,
		MetricsStore: store,
	}
}
