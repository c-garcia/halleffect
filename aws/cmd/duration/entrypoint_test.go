// +build !integration, !service

package main

import (
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/c-garcia/halleffect/internal/pkg/concourse"
	"github.com/c-garcia/halleffect/internal/pkg/services"
	"github.com/c-garcia/halleffect/internal/pkg/store"
	"github.com/c-garcia/halleffect/internal/pkg/timing"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestConfig(t *testing.T) {
	params := SystemParams{
		ConcourseName: "test",
		ConcourseURL:  "http://example.com/",
	}
	cfg := configureSystem(params)

	_, ok := cfg.ConcourseAPI.(*concourse.ApiImpl)
	assert.True(t, ok)
	_, ok = cfg.CloudwatchAPI.(*cloudwatch.CloudWatch)
	assert.True(t, ok)
	_, ok = cfg.Store.(*store.JobLastSuccessfulDurationAWS)
	assert.True(t, ok)
	_, ok = cfg.Clock.(*timing.SystemClock)
	assert.True(t, ok)
	_, ok = cfg.MetricsService.(*services.LastSuccessfulBuildDurationImpl)
	assert.True(t, ok)
	_, ok = cfg.Logger.(*log.Logger)
	assert.True(t, ok)
}
