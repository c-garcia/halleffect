package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/c-garcia/halleffect/internal/pkg/concourse"
	"github.com/c-garcia/halleffect/internal/pkg/microlog"
	"github.com/c-garcia/halleffect/internal/pkg/services"
	"github.com/c-garcia/halleffect/internal/pkg/store"
	"github.com/c-garcia/halleffect/internal/pkg/timing"
	"github.com/c-garcia/halleffect/internal/pkg/writers"
	"github.com/pkg/errors"
	"log"
	"os"
)
import "github.com/aws/aws-lambda-go/lambda"

var (
	handler DurationMetricsEventHandler
)

type SystemParams struct {
	ConcourseName string
	ConcourseURL  string
}

type SystemConfiguration struct {
	ConcourseAPI   concourse.API
	CloudwatchAPI  writers.AWSCloudWatchMetricWriter
	Store          store.JobLastSuccessfulDuration
	Clock          timing.Clock
	MetricsService services.Metrics
	Logger         microlog.Logger
}

const Namespace = "Concourse/Jobs"

func configureSystem(pars SystemParams) SystemConfiguration {
	concourseName := pars.ConcourseName
	concourseURL := pars.ConcourseURL
	concourseAPI := concourse.New(concourseName, concourseURL)
	cloudwatchAPI := cloudwatch.New(session.Must(session.NewSession()))
	metricsStore := store.NewJobLastSuccessfulDurationAWS(Namespace, cloudwatchAPI)
	clock := timing.NewSystemClock()
	thePoller := services.NewLastSuccessfulBuildDurationImpl(concourseName, concourseAPI, metricsStore, clock)
	logger := log.New(os.Stderr, "hall-effect", log.LstdFlags)

	return SystemConfiguration{
		ConcourseAPI:   concourseAPI,
		CloudwatchAPI:  cloudwatchAPI,
		Store:          metricsStore,
		Clock:          clock,
		MetricsService: thePoller,
		Logger:         logger,
	}
}

func main() {
	concourseName := os.Getenv("CONCOURSE_NAME")
	concourseURL := os.Getenv("CONCOURSE_URL")
	if concourseName == "" || concourseURL == "" {
		panic(errors.New("CONCOURSE_NAME or CONCOURSE_URL environment variables not set"))
	}
	params := SystemParams{
		ConcourseName: concourseName,
		ConcourseURL:  concourseURL,
	}
	cfg := configureSystem(params)
	handler = NewLambdaHandler(cfg.MetricsService, cfg.Logger)
	lambda.Start(handler)
}
