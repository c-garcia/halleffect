package publisher

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/bxcodec/faker"
)

func TestInMemory_Publish(t *testing.T){
	var m1 JobDurationMetric
	err := faker.FakeData(&m1)
	if err != nil {
		panic(err)
	}

	sut := NewInMemory()
	err = sut.PublishJobDuration(m1)

	assert.NoError(t, err)
	assert.True(t, sut.JobDurationHasBeenPubished(m1))
}


