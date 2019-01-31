package publisher

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
	metric := JobDurationMetric{
		Timestamp:    100,
		EndTime:      120,
		PipelineName: "pipeline1",
		JobName:      "job",
		Status:       "finished",
		Concourse:    "http://localhost:8080",
	}
	pipelineDimension := "pipeline"
	jobNameDimension := "job_name"
	jobStatusDimension := "status"
	concourseDimension := "concourse"
	datum := &cloudwatch.MetricDatum{}
	datum.
		SetTimestamp(time.Unix(int64(metric.Timestamp), 0)).
		SetMetricName("Duration").
		SetDimensions([]*cloudwatch.Dimension{
			&cloudwatch.Dimension{Name: &pipelineDimension, Value: &metric.PipelineName},
			&cloudwatch.Dimension{Name: &jobNameDimension, Value: &metric.JobName},
			&cloudwatch.Dimension{Name: &jobStatusDimension, Value: &metric.Status},
			&cloudwatch.Dimension{Name: &concourseDimension, Value: &metric.Concourse},
		}).
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
	err := sut.PublishJobDuration(metric)

	assert.NoError(t, err)

	ctrl.Finish()

}
func TestAWSImpl_Publish_PropagatesErrors(t *testing.T) {
	metric := JobDurationMetric{
		Timestamp: 100,
		EndTime:   120,
	}
	ctrl := gomock.NewController(t)
	mockWriter := writersMocks.NewMockAWSCloudWatchMetricWriter(ctrl)
	mockWriter.EXPECT().PutMetricData(gomock.Any()).Return(&cloudwatch.PutMetricDataOutput{}, assert.AnError)

	sut := NewAWS("Concourse/Jobs", mockWriter)
	err := sut.PublishJobDuration(metric)

	assert.Error(t, err)
	assert.Equal(t, assert.AnError, errors.Cause(err))
	ctrl.Finish()
}
