package publisher

import (
	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
	"testing"
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
