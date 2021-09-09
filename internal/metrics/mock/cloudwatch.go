package mock

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"

	"github.com/skpr/cluster-metrics/internal/metrics/types"
)

// Cloudwatch provides a mock cloudwatch.
type Cloudwatch struct {
	types.CloudwatchInterface
}

// NewCloudwatch creates a new mock.
func NewCloudwatch() types.CloudwatchInterface {
	return &Cloudwatch{}
}

// PutMetricData implements the interface.
func (c *Cloudwatch) PutMetricData(ctx context.Context, params *cloudwatch.PutMetricDataInput, optFns ...func(options *cloudwatch.Options)) (*cloudwatch.PutMetricDataOutput, error) {
	return &cloudwatch.PutMetricDataOutput{}, nil
}
