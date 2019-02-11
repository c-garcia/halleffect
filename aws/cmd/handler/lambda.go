package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/c-garcia/halleffect/aws"
	"github.com/c-garcia/halleffect/internal/pkg/concourse"
	"github.com/c-garcia/halleffect/internal/pkg/poller"
	"github.com/c-garcia/halleffect/internal/pkg/publisher"
	"github.com/c-garcia/halleffect/internal/pkg/timing"
	"github.com/c-garcia/halleffect/internal/pkg/writers"
	"github.com/pkg/errors"
	"log"
	"os"
)
import "github.com/aws/aws-lambda-go/lambda"

var (
	handler aws.MetricsHandler
)

type SystemParams struct {
	ConcourseName string
	ConcourseURL  string
}

type SystemConfiguration struct {
	ConcourseAPI     concourse.API
	CloudwatchAPI    writers.AWSCloudWatchMetricWriter
	MetricsPublisher publisher.MetricsPublisher
	Clock            timing.Clock
	Poller           poller.Service
	Logger           aws.Logger
}

const Namespace = "Concourse/Jobs"

func configureSystem(pars SystemParams) SystemConfiguration {
	concourseName := pars.ConcourseName
	concourseURL := pars.ConcourseURL
	concourseAPI := concourse.New(concourseName, concourseURL)
	cloudwatchAPI := cloudwatch.New(session.Must(session.NewSession()))
	metricsPublisher := publisher.NewAWS(Namespace, cloudwatchAPI)
	clock := timing.NewSystemClock()
	thePoller := poller.New(concourseAPI, metricsPublisher, clock)
	logger := log.New(os.Stderr, "hall-effect", log.LstdFlags)

	return SystemConfiguration{
		ConcourseAPI:     concourseAPI,
		CloudwatchAPI:    cloudwatchAPI,
		MetricsPublisher: metricsPublisher,
		Clock:            clock,
		Poller:           thePoller,
		Logger:           logger,
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
	handler = aws.NewLambdaHandler(cfg.Poller, cfg.Logger)
	lambda.Start(handler)
}
