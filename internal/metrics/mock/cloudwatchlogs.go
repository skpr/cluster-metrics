package mock

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"

	"github.com/skpr/cluster-metrics/internal/metrics/types"
)

// CloudwatchLogs is the mock cloudwatch logs client.
type CloudwatchLogs struct {
	types.CloudwatchLogsInterface
}

// NewCloudwatchLogs creates a new mock cloudwatch logs client.
func NewCloudwatchLogs() *CloudwatchLogs {
	return &CloudwatchLogs{}
}

// CreateLogGroup implements the interface.
func (l CloudwatchLogs) CreateLogGroup(ctx context.Context, params *cloudwatchlogs.CreateLogGroupInput, optFns ...func(options *cloudwatchlogs.Options)) (*cloudwatchlogs.CreateLogGroupOutput, error) {
	return nil, nil
}

// CreateLogStream implements the interface.
func (l CloudwatchLogs) CreateLogStream(ctx context.Context, params *cloudwatchlogs.CreateLogStreamInput, optFns ...func(options *cloudwatchlogs.Options)) (*cloudwatchlogs.CreateLogStreamOutput, error) {
	return nil, nil
}

// PutLogEvents implements the interface.
func (l CloudwatchLogs) PutLogEvents(ctx context.Context, params *cloudwatchlogs.PutLogEventsInput, optFns ...func(options *cloudwatchlogs.Options)) (*cloudwatchlogs.PutLogEventsOutput, error) {
	out := &cloudwatchlogs.PutLogEventsOutput{
		NextSequenceToken:     aws.String("abcd1234"),
	}
	return out, nil
}
