package podsbyPhase

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awstypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"sort"
	"time"
)

func (c Client) Convert(timestamp time.Time, cluster string, labels map[string]interface{}) []awstypes.MetricDatum {
	// Sort keys for a consistent result order.
	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var data []awstypes.MetricDatum
	for _, label := range labels {
		value := labels[label.(string)]
		datum := awstypes.MetricDatum{
			MetricName: aws.String(metricTotal),
			Dimensions: []awstypes.Dimension{
				{
					Name:  aws.String(dimensionCluster),
					Value: aws.String(cluster),
				},
			},
			Timestamp: aws.Time(timestamp),
			Value:     aws.Float64(value.(float64)),
		}
		for i, value := range labels {
			datum.Dimensions = append(datum.Dimensions, awstypes.Dimension{
				Name:  aws.String(i),
				Value: aws.String(fmt.Sprintf("%v", value)),
			})
		}
		data = append(data, datum)
	}
	return data
}
