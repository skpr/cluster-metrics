package metrics

import (
	"bytes"
	"strings"
	"time"

	"github.com/prozz/aws-embedded-metrics-golang/emf"
)

func Format(timestamp time.Time, namespace string, metric Metric) string {
	var buf bytes.Buffer
	formatter := emf.New(emf.WithWriter(&buf), emf.WithTimestamp(timestamp)).Namespace(namespace)
	formatter.DimensionSet(
		emf.NewDimension("Kind", metric.Kind),
		emf.NewDimension("Namespace", metric.Namespace),
		emf.NewDimension("Phase", string(metric.Phase)),
	).Metric("Total", metric.Total).Log()
	// Remove the linebreak added by emf.
	return strings.TrimSuffix(buf.String(), "\n")
}
