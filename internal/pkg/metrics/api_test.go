package metrics

import (
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	writersMocks "github.com/c-garcia/halleffect/internal/pkg/writers/mocks"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAWSImpl_Publish(t *testing.T) {
	metric := Metric{
		Timestamp: 100,
		EndTime:   120,
	}
	datum := &cloudwatch.MetricDatum{}
	datum.
		SetTimestamp(time.Unix(int64(metric.Timestamp),0)).
		SetMetricName("Duration").
		SetUnit("seconds").
		SetValue(float64(metric.Duration()))

	dataInput := &cloudwatch.PutMetricDataInput{}
	dataInput.
		SetNamespace("Concourse/Jobs").
		SetMetricData([]*cloudwatch.MetricDatum{datum})
	ctrl := gomock.NewController(t)
	mockWriter := writersMocks.NewMockAWSCloudWatchMetricWriter(ctrl)
	mockWriter.EXPECT().PutMetricData(dataInput).Return(&cloudwatch.PutMetricDataOutput{}, nil)

	sut := NewAWS("Concourse/Jobs", mockWriter)
	err := sut.Publish(metric)

	assert.NoError(t, err)

	ctrl.Finish()

}
func TestAWSImpl_Publish_PropagatesErrors(t *testing.T) {
	metric := Metric{
		Timestamp: 100,
		EndTime:   120,
	}
	ctrl := gomock.NewController(t)
	mockWriter := writersMocks.NewMockAWSCloudWatchMetricWriter(ctrl)
	mockWriter.EXPECT().PutMetricData(gomock.Any()).Return(&cloudwatch.PutMetricDataOutput{}, assert.AnError)

	sut := NewAWS("Concourse/Jobs", mockWriter)
	err := sut.Publish(metric)

	assert.Error(t, err)
	assert.Equal(t, assert.AnError, errors.Cause(err))
	ctrl.Finish()
}
