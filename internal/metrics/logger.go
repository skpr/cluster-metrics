package metrics

import (
	"bytes"
	"strings"
	"time"

	"github.com/prozz/aws-embedded-metrics-golang/emf"
)

// Process ze dimensions.
func Process(namespace string, metrics MetricSet, timestamp time.Time, processLine func(line string)) {
	for _, metric := range metrics.Items {
		var buf bytes.Buffer
		logger := emf.New(emf.WithWriter(&buf), emf.WithTimestamp(timestamp)).Namespace(namespace)
		logger.DimensionSet(
			emf.NewDimension("Kind", metric.Kind),
			emf.NewDimension("Namespace", metric.Namespace),
			emf.NewDimension("Phase", string(metric.Phase)),
		).Metric("Total", metric.Total).Log()
		// Remove the linebreak added by emf.
		processLine(strings.TrimSuffix(buf.String(), "\n"))
	}
}
