package metrics

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	awstypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"

	"github.com/skpr/cluster-metrics/internal/metrics/types"
)

type LogsPusher struct {
	cwLogsClient types.CloudwatchLogsInterface
}

func NewLogsPusher(cwLogsClient types.CloudwatchLogsInterface) *LogsPusher {
	return &LogsPusher{
		cwLogsClient: cwLogsClient,
	}
}

// CreateLogGroup will attempt to create a log group and not return an error if it already exists.
func (p *LogsPusher) CreateLogGroup(ctx context.Context, group string) error {
	_, err := p.cwLogsClient.CreateLogGroup(ctx, &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: aws.String(group),
	})
	if err != nil {
		var awsErr *awstypes.ResourceAlreadyExistsException
		if errors.As(err, &awsErr) {
			return nil
		}
		return err
	}

	return nil
}

// CreateLogStream will attempt to create a log stream and not return an error if it already exists.
func (p *LogsPusher) CreateLogStream(ctx context.Context, group, stream string) error {
	_, err := p.cwLogsClient.CreateLogStream(ctx, &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(group),
		LogStreamName: aws.String(stream),
	})
	if err != nil {
		var awsErr *awstypes.ResourceAlreadyExistsException
		if errors.As(err, &awsErr) {
			return nil
		}
		return err
	}

	return nil
}
