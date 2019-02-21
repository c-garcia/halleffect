package store

import "github.com/c-garcia/halleffect/internal/pkg/metrics"

//go:generate mockgen -source=contract.go -destination=mocks/mock_contract.go -package=mocks

type JobLastSuccessfulDuration interface {
	Store(m metrics.JobLastSuccessfulDuration) error
}
