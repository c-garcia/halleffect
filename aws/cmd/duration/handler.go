package main

import (
	"context"
	"github.com/c-garcia/halleffect/internal/pkg/microlog"
	"github.com/c-garcia/halleffect/internal/pkg/services"
	"github.com/pkg/errors"
)

type GetDurationEvent struct{}

type DurationMetricsEventHandler func(c context.Context, event GetDurationEvent) (string, error)

func NewLambdaHandler(svc services.Metrics, logger microlog.Logger) DurationMetricsEventHandler {
	return func(ctx context.Context, event GetDurationEvent) (string, error) {
		if err := svc.SaveMetrics(); err != nil {
			newErr := errors.Wrap(err, "Service error")
			logger.Printf("%+v", newErr)
			return "", newErr
		}
		logger.Println("Metrics export done")
		return "", nil
	}
}
