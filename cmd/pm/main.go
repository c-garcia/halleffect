package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"os"
	"time"
)

type Options struct {
	Namespace      string
	Metric         string
	Units          string
	DimensionName  string
	DimensionValue string
	Value          float64
}

func parseArgs() (*Options, error) {
	res := Options{}
	flag.StringVar(&res.Namespace, "namespace", "", "namespace")
	flag.StringVar(&res.Metric, "metric", "", "metric name")
	flag.StringVar(&res.Units, "units", "", "metric units")
	flag.StringVar(&res.DimensionName, "dname", "", "dimension name")
	flag.StringVar(&res.DimensionValue, "dvalue", "", "dimension value")
	flag.Float64Var(&res.Value, "value", 0.0, "metric value")
	flag.Parse()
	return &res, nil
}

func main() {

	args, err := parseArgs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sess := session.Must(session.NewSession())
	cw := cloudwatch.New(sess)
	dim := &cloudwatch.Dimension{}
	dim.SetName(args.DimensionName).SetValue(args.DimensionValue)
	datum := &cloudwatch.MetricDatum{}
	datum.SetTimestamp(time.Now()).
		SetMetricName(args.Metric).
		SetValue(args.Value).
		SetDimensions([]*cloudwatch.Dimension{dim})
	in := &cloudwatch.PutMetricDataInput{}
	in.SetNamespace(args.Namespace).SetMetricData([]*cloudwatch.MetricDatum{datum})

	out, err := cw.PutMetricData(in)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Printf("%v\n", out)
}
