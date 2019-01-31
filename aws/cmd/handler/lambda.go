package handler

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/c-garcia/halleffect/aws"
	"github.com/c-garcia/halleffect/internal/pkg/concourse"
	"github.com/c-garcia/halleffect/internal/pkg/poller"
	"github.com/c-garcia/halleffect/internal/pkg/publisher"
	"github.com/pkg/errors"
	"log"
	"os"
)
import "github.com/aws/aws-lambda-go/lambda"

var (
	handler aws.MetricsHandler
)

func main() {
	concourseName := os.Getenv("CONCOURSE_NAME")
	concourseURL := os.Getenv("CONCOURSE_URL")
	if concourseName == "" || concourseURL == "" {
		panic(errors.New("CONCOURSE_NAME or CONCOURSE_URL environment variables not set"))
	}
	concourseAPI := concourse.New(concourseName, concourseURL)
	cloudwatchAPI := cloudwatch.New(session.Must(session.NewSession()))
	metricsExporter := publisher.NewAWS("Concourse/Jobs", cloudwatchAPI)
	handler = aws.NewLambdaHandler(poller.New(concourseAPI, metricsExporter), log.New(os.Stderr, "concourse-metrics", log.LstdFlags))
	lambda.Start(handler)
}
