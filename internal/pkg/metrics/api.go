package metrics

//go:generate mockgen -source=api.go -destination=mocks/api.go -package=mocks

import (
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/c-garcia/halleffect/internal/pkg/writers"
	"time"
)

type Exporter interface {
	Publish(m Metric) error
}

type AWSImpl struct {
	Namespace string
	Writer    writers.AWSCloudWatchMetricWriter
}

func (s *AWSImpl) Publish(m Metric) error {
	datum := &cloudwatch.MetricDatum{}
	datum.
		SetTimestamp(time.Unix(int64(m.Timestamp), 0)).
		SetMetricName("Duration").
		SetUnit("seconds").
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
