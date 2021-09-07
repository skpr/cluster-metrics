package metrics

import (
	"encoding/json"
	"io"
)

const (
	typeGauge = "gauge"
	metricName = "PodStatus"
)

// Log logs the metrics to the writer.
func Log(writer io.Writer, metrics *MetricSet) error {
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
