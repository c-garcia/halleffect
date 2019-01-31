package aws

import "context"

type PublishConcourseMetricsLambdaEvent struct{}

type MetricsHandler func(ctx context.Context, event PublishConcourseMetricsLambdaEvent) (string, error)
