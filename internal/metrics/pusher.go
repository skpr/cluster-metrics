package metrics

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	awstypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"

	"github.com/skpr/cluster-metrics/internal/metrics/types"
)

const (
	dimensionKind      = "Kind"
	dimensionNamespace = "Namespace"
	dimensionPhase     = "Phase"
	metricTotal        = "Total"
)

// Pusher the metrics pusher.
type Pusher struct {
	cloudwatchClient types.CloudwatchInterface
}

// NewPusher creates a new metrics pusher.
func NewPusher(cloudwatchClient types.CloudwatchInterface) *Pusher {
	return &Pusher{
		cloudwatchClient: cloudwatchClient,
	}
}

// Push the metrics.
func (p *Pusher) Push(ctx context.Context, namespace string, metricData []awstypes.MetricDatum) error {
	_, err := p.cloudwatchClient.PutMetricData(ctx, &cloudwatch.PutMetricDataInput{
		MetricData: metricData,
		Namespace:  aws.String(namespace),
	})
	if err != nil {
		return err
	}
	return nil
}

// ConvertToMetricData converts our metrics to aws metric data.
func ConvertToMetricData(timestamp time.Time, metrics MetricSet) []awstypes.MetricDatum {
	var data []awstypes.MetricDatum
	for _, metric := range metrics.Items {
		datum := awstypes.MetricDatum{
			MetricName: aws.String(metricTotal),
			Dimensions: []awstypes.Dimension{
				{
					Name:  aws.String(dimensionKind),
					Value: aws.String(metric.Kind),
				},
				{
					Name:  aws.String(dimensionNamespace),
					Value: aws.String(metric.Namespace),
				},
				{
					Name:  aws.String(dimensionPhase),
					Value: aws.String(string(metric.Phase)),
				},
			},
			Timestamp: aws.Time(timestamp),
			Value:     aws.Float64(float64(metric.Total)),
		}
		data = append(data, datum)
	}
	return data
}
