package types

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

// CloudwatchLogsInterface provides an interface for the cloudwatch logs cwLogsClient.
type CloudwatchLogsInterface interface {
	CreateLogGroup(ctx context.Context, params *cloudwatchlogs.CreateLogGroupInput, optFns ...func(options *cloudwatchlogs.Options)) (*cloudwatchlogs.CreateLogGroupOutput, error)
	CreateLogStream(ctx context.Context, params *cloudwatchlogs.CreateLogStreamInput, optFns ...func(options *cloudwatchlogs.Options)) (*cloudwatchlogs.CreateLogStreamOutput, error)
	PutLogEvents(ctx context.Context, params *cloudwatchlogs.PutLogEventsInput, optFns ...func(options *cloudwatchlogs.Options)) (*cloudwatchlogs.PutLogEventsOutput, error)
}
