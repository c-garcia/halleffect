package poller

//go:generate mockgen -source=poller.go -destination=mocks/mock_doer.go -package=mocks

import (
	"github.com/c-garcia/halleffect/internal/pkg/concourse"
	"github.com/c-garcia/halleffect/internal/pkg/publisher"
	"github.com/c-garcia/halleffect/internal/pkg/timing"
	"github.com/pkg/errors"
)

type Service interface {
	ExportMetrics() error
}

type ServiceImpl struct {
	Concourse concourse.API
	Exporter  publisher.MetricsPublisher
	Clock     timing.Clock
}

func (s *ServiceImpl) ExportMetrics() error {
	var builds []concourse.Build
	var err error
	concourseName := s.Concourse.Name()
	samplingTime := s.Clock.UnixTime()

	if builds, err = s.Concourse.FindLastBuilds(); err != nil {
		return errors.Wrap(err, "Error retrieving builds")
	}
	for _, build := range builds {
		if !build.Finished() {
			if err := s.Exporter.PublishJobDuration(publisher.FromConcourseBuild(concourseName, build)); err != nil {
				return errors.Wrap(err, "Error publishing metrics")
			}
		}
	}

	supportsJobs, err := s.Concourse.SupportsJobsEndpoint()
	if err != nil {
		return errors.Wrap(err, "Error determining if Jobs API is supported")
	}
	if !supportsJobs {
		return nil
	}

	var jobs []concourse.JobStatus
	if jobs, err = s.Concourse.FindJobStatuses(); err != nil {
		return errors.Wrap(err, "Error retrieving jobs")
	}
	for _, job := range jobs {
		if err := s.Exporter.PublishJobStatus(publisher.FromConcourseJobStatus(concourseName, samplingTime, job)); err != nil {
			return errors.Wrap(err, "Error publishing status metrics")
		}
	}
	return nil
}

func New(concourse concourse.API, exporter publisher.MetricsPublisher, clock timing.Clock) *ServiceImpl {
	return &ServiceImpl{
		Concourse: concourse,
		Exporter:  exporter,
		Clock:     clock,
	}
}
