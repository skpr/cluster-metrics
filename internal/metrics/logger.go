package metrics

import (
	"encoding/json"
	"io"
)

const typeGauge = "gauge"

// Log logs the metrics to the writer.
func Log(writer io.Writer, name string, metrics *MetricSet) error {
	encoder := json.NewEncoder(writer)
	for _, metric := range metrics.Items {
		metric.Name = name
		metric.Type = typeGauge
		err := encoder.Encode(&metric)
		if err != nil {
			return err
		}
	}
	return nil
}
