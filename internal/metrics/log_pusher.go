package metrics

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	awstypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"

	"github.com/skpr/cluster-metrics/internal/metrics/types"
)

type LogPusher struct {
	log          log.Logger
	cwLogsClient types.CloudwatchLogsInterface
}

func NewLogsPusher(cwLogsClient types.CloudwatchLogsInterface) *LogPusher {
	return &LogPusher{
		cwLogsClient: cwLogsClient,
	}
}

func (p *LogPusher) PushLogs(ctx context.Context, logGroup, logStream string, events []awstypes.InputLogEvent, sequenceToken *string) error {
	_, err := p.cwLogsClient.PutLogEvents(ctx, &cloudwatchlogs.PutLogEventsInput{
		LogEvents:     events,
		LogGroupName:  aws.String(logGroup),
		LogStreamName: aws.String(logStream),
		SequenceToken: sequenceToken,
	}, func(options *cloudwatchlogs.Options) {
		options.Retryer = retry.AddWithMaxAttempts(options.Retryer, 0)
	})
	if err != nil {
		var seqTokenError *awstypes.InvalidSequenceTokenException
		if errors.As(err, &seqTokenError) {
			p.log.Println("Invalid token. Refreshing", logGroup, logStream)
			(&cloudwatchlogs.PutLogEventsInput{
				LogEvents:     events,
				LogGroupName:  aws.String(logGroup),
				LogStreamName: aws.String(logStream),
				SequenceToken: sequenceToken,
			}).SequenceToken = seqTokenError.ExpectedSequenceToken
			return p.PushLogs(ctx, logGroup, logStream, events, sequenceToken)
		}
		var alreadyAccErr *awstypes.DataAlreadyAcceptedException
		if errors.As(err, &alreadyAccErr) {
			p.log.Println("Data already accepted. Refreshing", logGroup, logStream)
			(&cloudwatchlogs.PutLogEventsInput{
				LogEvents:     events,
				LogGroupName:  aws.String(logGroup),
				LogStreamName: aws.String(logStream),
				SequenceToken: sequenceToken,
			}).SequenceToken = alreadyAccErr.ExpectedSequenceToken
			return p.PushLogs(ctx, logGroup, logStream, events, sequenceToken)
		}
		return err
	}
	return nil
}

// CreateLogGroup will attempt to create a log group and not return an error if it already exists.
func (p *LogPusher) CreateLogGroup(ctx context.Context, group string) error {
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
func (p *LogPusher) CreateLogStream(ctx context.Context, group, stream string) error {
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
