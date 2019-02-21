package services

//go:generate mockgen -source=metrics.go -destination=mocks/metrics.go -package=mocks

type Metrics interface {
	SaveMetrics() error
}
