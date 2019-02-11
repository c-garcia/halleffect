package publisher

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatch"
	writersMocks "github.com/c-garcia/halleffect/internal/pkg/writers/mocks"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestAWSImpl_PublishJobDuration(t *testing.T) {
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
	teamDimension := "team"
	datum := &cloudwatch.MetricDatum{}
	datum.
		SetTimestamp(time.Unix(int64(metric.Timestamp), 0)).
		SetMetricName("Duration").
		SetDimensions([]*cloudwatch.Dimension{
			{Name: &pipelineDimension, Value: &metric.PipelineName},
			{Name: &jobNameDimension, Value: &metric.JobName},
			{Name: &jobStatusDimension, Value: &metric.Status},
			{Name: &concourseDimension, Value: &metric.Concourse},
			{Name: &teamDimension, Value: &metric.TeamName},
		}).
		SetUnit("Seconds").
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

func TestAWSImpl_PublishJobStatus(t *testing.T) {
	metric := JobStatusMetric{
		Concourse:    "concourse-test",
		TeamName:     "main",
		PipelineName: "p1",
		JobName:      "j1",
		Status:       "up",
		SamplingTime: 1100,
	}
	pipelineDimension := "pipeline"
	jobNameDimension := "job_name"
	concourseDimension := "concourse"
	teamDimension := "team"
	datum := &cloudwatch.MetricDatum{}
	datum.
		SetTimestamp(time.Unix(int64(metric.SamplingTime), 0)).
		SetMetricName("Status").
		SetDimensions([]*cloudwatch.Dimension{
			{Name: &pipelineDimension, Value: &metric.PipelineName},
			{Name: &jobNameDimension, Value: &metric.JobName},
			{Name: &concourseDimension, Value: &metric.Concourse},
			{Name: &teamDimension, Value: &metric.TeamName},
		}).
		SetUnit("Percent")

	dataInput := &cloudwatch.PutMetricDataInput{}
	dataInput.
		SetNamespace("Concourse/Jobs").
		SetMetricData([]*cloudwatch.MetricDatum{datum})

	testCases := []struct {
		Desc   string
		Status string
		Value  float64
	}{{"Up->100%", "up", 100.0}, {"Down->0%", "down", 0.0}}

	for _, testCase := range testCases {
		t.Run(testCase.Desc, func(t *testing.T) {
			metric.Status = testCase.Status
			datum.SetValue(testCase.Value)

			ctrl := gomock.NewController(t)
			mockWriter := writersMocks.NewMockAWSCloudWatchMetricWriter(ctrl)
			mockWriter.EXPECT().PutMetricData(dataInput).Return(&cloudwatch.PutMetricDataOutput{}, nil)

			sut := NewAWS("Concourse/Jobs", mockWriter)
			err := sut.PublishJobStatus(metric)

			assert.NoError(t, err)
			ctrl.Finish()
		})
	}

}
