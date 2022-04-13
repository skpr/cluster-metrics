package podsbyPhase

import (
	"encoding/json"
	"io"

	metricsS "github.com/skpr/cluster-metrics/internal/metrics"
)

func (c Client) Log(writer io.Writer, metrics *metricsS.MetricSet) error {
	encoder := json.NewEncoder(writer)
	for _, metric := range metrics.Items {
		metric.Name = metricName
		metric.Type = typeGauge
		err := encoder.Encode(&metric)
		if err != nil {
			return err
		}
	}
	return nil
}
