package metrics

import (
	"io"
	"time"

	"github.com/prozz/aws-embedded-metrics-golang/emf"
)

// Logger is the metrics logger.
type Logger struct {
	writer    io.Writer
	namespace string
}

// NewLogger creates a new metrics logger.
func NewLogger(writer io.Writer, namespace string) *Logger {
	return &Logger{
		writer:    writer,
		namespace: namespace,
	}
}

// Log ze dimensions.
func (l *Logger) Log(metrics *MetricSet, timestamp time.Time) {
	for _, metric := range metrics.Items {
		logger := emf.New(emf.WithWriter(l.writer), emf.WithTimestamp(timestamp)).Namespace(l.namespace)
		logger.DimensionSet(
			emf.NewDimension("Kind", metric.Kind),
			emf.NewDimension("Namespace", metric.Namespace),
			emf.NewDimension("Phase", string(metric.Phase)),
		).Metric("Total", metric.Total).Log()
	}
}
