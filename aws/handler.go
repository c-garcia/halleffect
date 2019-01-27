package aws

import (
	"context"
	"github.com/c-garcia/halleffect/internal/pkg/exporter"
	"github.com/pkg/errors"
)

func NewLambdaHandler(exporter exporter.Service, logger Logger) MetricsHandler {
	return func(ctx context.Context, event ExportMetricsLambdaEvent) (s string, e error) {
		if err := exporter.ExportMetrics(); err != nil {
			newErr := errors.Wrap(err, "Service error")
			logger.Printf("%+v", newErr)
			return "", newErr
		}
		logger.Println("Metrics export done")
		return "", nil
	}
}
