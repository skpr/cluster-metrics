package internal

import (
	"io"
	"time"

	"github.com/prozz/aws-embedded-metrics-golang/emf"
)

// MetricsLogger is the metrics logger.
type MetricsLogger struct {
	writer    io.Writer
}

// NewMetricsLogger creates a new metrics logger.
func NewMetricsLogger(writer io.Writer) *MetricsLogger {
	return &MetricsLogger{
		writer:    writer,
	}
}

// Log ze dimensions.
func (l *MetricsLogger) Log(metrics *MetricSet, timestamp time.Time) {
	for _, metric := range metrics.Items {
		logger := emf.New(emf.WithWriter(l.writer), emf.WithTimestamp(timestamp))
		logger.DimensionSet(
			emf.NewDimension("Kind", metric.Kind),
			emf.NewDimension("Namespace", metric.Namespace),
			emf.NewDimension("Phase", string(metric.Phase)),
		).Metric("Total", metric.Total).Log()
	}
}
