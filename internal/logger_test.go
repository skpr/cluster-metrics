package internal

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestMetricsLogger_Log(t *testing.T) {
	var buf bytes.Buffer

	logger := NewMetricsLogger(&buf)

	metrics := NewMetricSet()
	metrics.Increment("ReplicaSet", "foo", corev1.PodPending)
	metrics.Increment("ReplicaSet", "foo", corev1.PodPending)
	metrics.Increment("ReplicaSet", "foo", corev1.PodRunning)
	metrics.Increment("ReplicaSet", "bar", corev1.PodPending)
	metrics.Increment("ReplicaSet", "bar", corev1.PodRunning)

	timestamp := time.Date(2020, time.September, 2, 9, 2, 0, 0, time.UTC)
	logger.Log(metrics, timestamp)
	s := buf.String()
	fmt.Println(s)

	json1 := `{"Kind":"ReplicaSet","Namespace":"foo","Phase":"Pending","Total":2,"_aws":{"Timestamp":1599037320000,"CloudWatchMetrics":[{"Namespace":"aws-embedded-metrics","Dimensions":[["Kind","Namespace","Phase"]],"Metrics":[{"Name":"Total","Unit":"None"}]}]}}`
	json2 := `{"Kind":"ReplicaSet","Namespace":"foo","Phase":"Running","Total":1,"_aws":{"Timestamp":1599037320000,"CloudWatchMetrics":[{"Namespace":"aws-embedded-metrics","Dimensions":[["Kind","Namespace","Phase"]],"Metrics":[{"Name":"Total","Unit":"None"}]}]}}`
	json3 := `{"Kind":"ReplicaSet","Namespace":"bar","Phase":"Pending","Total":1,"_aws":{"Timestamp":1599037320000,"CloudWatchMetrics":[{"Namespace":"aws-embedded-metrics","Dimensions":[["Kind","Namespace","Phase"]],"Metrics":[{"Name":"Total","Unit":"None"}]}]}}`
	json4 := `{"Kind":"ReplicaSet","Namespace":"bar","Phase":"Running","Total":1,"_aws":{"Timestamp":1599037320000,"CloudWatchMetrics":[{"Namespace":"aws-embedded-metrics","Dimensions":[["Kind","Namespace","Phase"]],"Metrics":[{"Name":"Total","Unit":"None"}]}]}}`

	assert.Contains(t, s, json1)
	assert.Contains(t, s, json2)
	assert.Contains(t, s, json3)
	assert.Contains(t, s, json4)

}
