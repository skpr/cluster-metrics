package types

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
)

// CloudwatchInterface provides and interface for Cloudwatch.
type CloudwatchInterface interface {
	PutMetricData(ctx context.Context, params *cloudwatch.PutMetricDataInput, optFns ...func(options *cloudwatch.Options)) (*cloudwatch.PutMetricDataOutput, error)
}
