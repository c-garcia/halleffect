package exporter

//go:generate mockgen -source=service.go -destination=mocks/mock_doer.go -package=mocks

import (
	"github.com/c-garcia/halleffect/internal/pkg/concourse"
	"github.com/c-garcia/halleffect/internal/pkg/metrics"
	"github.com/pkg/errors"
)

type Service interface {
	ExportMetrics() error
}

type ServiceImpl struct {
	Concourse concourse.API
	Exporter  metrics.Exporter
}

func New(concourse concourse.API, exporter metrics.Exporter) *ServiceImpl {
	return &ServiceImpl{
		Concourse: concourse,
		Exporter:  exporter,
	}
}

func (s *ServiceImpl) ExportMetrics() error {
	var builds []concourse.Build
	var err error
	if builds, err = s.Concourse.FindLastBuilds(); err != nil {
		return errors.Wrap(err, "Error retrieving builds")
	}
	concourseName := s.Concourse.Name()
	for _, build := range builds {
		if !build.Finished() {
			if err := s.Exporter.Publish(metrics.FromConcourseBuild(concourseName, build)); err != nil {
				return errors.Wrap(err, "Error publishing metrics")
			}
		}
	}
	return nil
}
