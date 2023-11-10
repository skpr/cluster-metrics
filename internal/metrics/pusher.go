package metrics

import (
	"context"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	awstypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"

	"github.com/skpr/cluster-metrics/internal/metrics/types"
)

const (
	dimensionKind      = "kind"
	dimensionNamespace = "namespace"
	dimensionState     = "phase"
	dimensionCluster   = "cluster"
	metricTotal        = "total"
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
func ConvertToMetricData(timestamp time.Time, cluster string, states StateSet) []awstypes.MetricDatum {

	// Sort keys for a consistent result order.
	keys := make([]string, 0, len(states))
	for k := range states {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var data []awstypes.MetricDatum
	for kind, stateData := range states {
		for state, _ := range stateData {
			datum := awstypes.MetricDatum{
				MetricName: aws.String(metricTotal),
				Dimensions: []awstypes.Dimension{
					{
						Name:  aws.String(dimensionKind),
						Value: aws.String(kind),
					},
					{
						Name:  aws.String(dimensionState),
						Value: aws.String(state),
					},
					{
						Name:  aws.String(dimensionCluster),
						Value: aws.String(cluster),
					},
				},
				Timestamp: aws.Time(timestamp),
				Value:     aws.Float64(float64(states[kind][state])),
			}
			data = append(data, datum)
		}
	}
	return data
}
