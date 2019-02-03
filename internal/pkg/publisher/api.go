package publisher

//go:generate mockgen -source=api.go -destination=mocks/api.go -package=mocks

import (
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/c-garcia/halleffect/internal/pkg/writers"
	"time"
)

type MetricsPublisher interface {
	PublishJobDuration(m JobDurationMetric) error
}

type AWSImpl struct {
	Namespace string
	Writer    writers.AWSCloudWatchMetricWriter
}

func cloudwatchDimension(n string, v string) *cloudwatch.Dimension {
	return &cloudwatch.Dimension{Name: &n, Value: &v}
}

func cloudwatchDimensions(d ...*cloudwatch.Dimension) []*cloudwatch.Dimension {
	return d
}

func (s *AWSImpl) PublishJobDuration(m JobDurationMetric) error {

	datum := &cloudwatch.MetricDatum{}
	datum.
		SetTimestamp(time.Unix(int64(m.Timestamp), 0)).
		SetMetricName("Duration").
		SetUnit("Seconds").
		SetDimensions(cloudwatchDimensions(
			cloudwatchDimension("pipeline", m.PipelineName),
			cloudwatchDimension("job_name", m.JobName),
			cloudwatchDimension("status", m.Status),
			cloudwatchDimension("concourse", m.Concourse),
			cloudwatchDimension("team", m.TeamName),
		)).
		SetValue(float64(m.Duration()))
	data := []*cloudwatch.MetricDatum{datum}
	in := &cloudwatch.PutMetricDataInput{}
	in.
		SetNamespace(s.Namespace).
		SetMetricData(data)
	_, err := s.Writer.PutMetricData(in)
	return err
}

func NewAWS(n string, w writers.AWSCloudWatchMetricWriter) *AWSImpl {
	return &AWSImpl{
		Namespace: n,
		Writer:    w,
	}
}
