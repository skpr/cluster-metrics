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
	DimensionKind      = "Kind"
	DimensionNamespace = "Namespace"
	DimensionPhase     = "Phase"
	MetricTotal        = "Total"
)

type Pusher struct {
	cloudwatchClient types.CloudwatchInterface
}

func NewPusher(cloudwatchClient types.CloudwatchInterface) *Pusher {
	return &Pusher{
		cloudwatchClient: cloudwatchClient,
	}
}

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

func ConvertToMetricData(timestamp time.Time, metrics MetricSet) []awstypes.MetricDatum {
	var data []awstypes.MetricDatum
	for _, metric := range metrics.Items {
		datum := awstypes.MetricDatum{
			MetricName: aws.String(MetricTotal),
			Dimensions: []awstypes.Dimension{
				{
					Name:  aws.String(DimensionKind),
					Value: aws.String(metric.Kind),
				},
				{
					Name:  aws.String(DimensionNamespace),
					Value: aws.String(metric.Namespace),
				},
				{
					Name:  aws.String(DimensionPhase),
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
