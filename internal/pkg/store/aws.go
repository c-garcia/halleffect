package store

import (
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/c-garcia/halleffect/internal/pkg/metrics"
	"github.com/c-garcia/halleffect/internal/pkg/writers"
	"github.com/pkg/errors"
)

type JobLastSuccessfulDurationAWS struct {
	Namespace string
	Writer    writers.AWSCloudWatchMetricWriter
}

func cloudwatchDimension(n string, v string) *cloudwatch.Dimension {
	return &cloudwatch.Dimension{Name: &n, Value: &v}
}

func cloudwatchDimensions(d ...*cloudwatch.Dimension) []*cloudwatch.Dimension {
	return d
}

func (s *JobLastSuccessfulDurationAWS) Store(m metrics.JobLastSuccessfulDuration) error {
	datum := &cloudwatch.MetricDatum{}
	datum.
		SetTimestamp(m.Timestamp).
		SetMetricName("LastSuccessfulDuration").
		SetDimensions(
			cloudwatchDimensions(
				cloudwatchDimension("concourse", m.Concourse),
				cloudwatchDimension("team", m.Team),
				cloudwatchDimension("pipeline", m.Pipeline),
				cloudwatchDimension("job", m.Job),
			)).
		SetUnit("Seconds").
		SetValue(m.Duration.Seconds())

	dataInput := &cloudwatch.PutMetricDataInput{}
	dataInput.
		SetNamespace("Concourse/Jobs").
		SetMetricData([]*cloudwatch.MetricDatum{datum})
	_, err := s.Writer.PutMetricData(dataInput)
	return errors.Wrap(err, "Error publishing duration metric")
}

func NewJobLastSuccessfulDurationAWS(namespace string, w writers.AWSCloudWatchMetricWriter) *JobLastSuccessfulDurationAWS {
	return &JobLastSuccessfulDurationAWS{
		Namespace: namespace,
		Writer:    w,
	}
}
