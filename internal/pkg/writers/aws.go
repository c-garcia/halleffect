package writers

//go:generate mockgen -source=aws.go -destination=mocks/aws.go -package=mocks

import "github.com/aws/aws-sdk-go/service/cloudwatch"

type AWSCloudWatchMetricWriter interface {
	PutMetricData(in *cloudwatch.PutMetricDataInput) (*cloudwatch.PutMetricDataOutput, error)
}
