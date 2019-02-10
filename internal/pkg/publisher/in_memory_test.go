package publisher

import (
	"testing"

	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
)

func TestInMemory_Publish(t *testing.T) {
	var m1 JobDurationMetric
	err := faker.FakeData(&m1)
	if err != nil {
		panic(err)
	}

	sut := NewInMemory()
	err = sut.PublishJobDuration(m1)

	assert.NoError(t, err)
	assert.True(t, sut.JobDurationHasBeenPublished(m1))
}

func TestInMemoryImpl_NumberOfPublishedJobDurationMetrics(t *testing.T) {
	var m1 JobDurationMetric
	err := faker.FakeData(&m1)
	if err != nil {
		panic(err)
	}

	sut := NewInMemory()
	assert.Equal(t, 0, sut.NumberOfPublishedJobDurationMetrics())

	err = sut.PublishJobDuration(m1)
	assert.NoError(t, err)
	assert.Equal(t, 1, sut.NumberOfPublishedJobDurationMetrics())

	err = sut.PublishJobDuration(m1)
	assert.NoError(t, err)
	assert.Equal(t, 2, sut.NumberOfPublishedJobDurationMetrics())
}

func TestInMemoryImpl_PublishJobStatus(t *testing.T) {
	status := JobStatusMetric{
		Concourse:    "concourse",
		TeamName:     "main",
		PipelineName: "p12",
		JobName:      "j32",
		Status:       "up",
		SamplingTime: 1050,
	}
	statusUnPublished := JobStatusMetric{
		Concourse:    "concourse",
		TeamName:     "main",
		PipelineName: "p13",
		JobName:      "j33",
		Status:       "down",
		SamplingTime: 1090,
	}
	sut := NewInMemory()
	err := sut.PublishJobStatus(status)
	assert.NoError(t, err)
	assert.True(t, sut.JobStatusHasBeenPublished(status))
	assert.False(t, sut.JobStatusHasBeenPublished(statusUnPublished))
	assert.Equal(t, 1, sut.NumberOfPublishedJobStatusMetrics())
	assert.Equal(t, 0, sut.NumberOfPublishedJobDurationMetrics())
}
