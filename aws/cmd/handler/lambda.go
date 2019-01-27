package handler

import (
	"github.com/c-garcia/halleffect/aws"
	"github.com/c-garcia/halleffect/internal/pkg/exporter"
	"log"
	"os"
)
import "github.com/aws/aws-lambda-go/lambda"

var (
	handler aws.MetricsHandler
)

func main() {
	handler = aws.NewLambdaHandler(exporter.New(), log.New(os.Stderr, "concourse-metrics", log.LstdFlags))
	lambda.Start(handler)
}
