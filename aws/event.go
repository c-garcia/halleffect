package aws

import "context"

type ExportMetricsLambdaEvent struct{}

type MetricsHandler func (ctx context.Context, event ExportMetricsLambdaEvent) (string, error)


