package metrics

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

func CreateInput(logGroup, logStream string, metrics MetricSet, sequenceToken *string) *cloudwatchlogs.PutLogEventsInput {
	input := &cloudwatchlogs.PutLogEventsInput{
		LogEvents:     nil,
		LogGroupName:  aws.String(logGroup),
		LogStreamName: aws.String(logStream),
		SequenceToken: sequenceToken,
	}
	return input
}

func ConvertToEvents(timestamp time.Time, namespace string, metrics MetricSet) []types.InputLogEvent {
	var events []types.InputLogEvent
	for _, metric := range metrics.Items {
		mesg := Format(timestamp, namespace, *metric)
		ev := types.InputLogEvent{
			Message:   aws.String(mesg),
			Timestamp: aws.Int64(timestamp.UnixNano() / int64(time.Millisecond)),
		}
		events = append(events, ev)
	}
	return events
}
