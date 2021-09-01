package internal

import (
	"github.com/prozz/aws-embedded-metrics-golang/emf"
)

// MetricsLogger is the metrics logger.
type MetricsLogger struct {
	logger *emf.Logger
}

// NewMetricsLogger creates a new metrics logger.
func NewMetricsLogger(logger *emf.Logger) *MetricsLogger {
	return &MetricsLogger{
		logger: logger,
	}
}

// Log ze dimensions.
func (l *MetricsLogger) Log(metrics *MetricSet) {
	for _, metric := range metrics.items {
		l.logger.DimensionSet(
			emf.NewDimension("Kind", metric.Kind),
			emf.NewDimension("Namespace", metric.Namespace),
			emf.NewDimension("Phase", string(metric.Phase)),
		).Metric("Count", metric.Count)
	}
	l.logger.Log()
}
