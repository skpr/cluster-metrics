package metrics

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestMetricsLogger_Log(t *testing.T) {

	mts := NewMetricSet()
	mts.Increment("ReplicaSet", "foo", corev1.PodPending)
	mts.Increment("ReplicaSet", "foo", corev1.PodPending)
	mts.Increment("ReplicaSet", "foo", corev1.PodRunning)
	mts.Increment("ReplicaSet", "bar", corev1.PodPending)
	mts.Increment("ReplicaSet", "bar", corev1.PodRunning)

	namespace := "Skpr/Cluster"
	timestamp := time.Date(2020, time.September, 2, 9, 2, 0, 0, time.UTC)
	var lines []string
	Process(namespace, mts, timestamp, func(line string) {
		lines = append(lines, line)
	})

	fmt.Println(lines)

	json1 := `{"Kind":"ReplicaSet","Namespace":"bar","Phase":"Pending","Total":1,"_aws":{"Timestamp":1599037320000,"CloudWatchMetrics":[{"Namespace":"Skpr/Cluster","Dimensions":[["Kind","Namespace","Phase"]],"Metrics":[{"Name":"Total","Unit":"None"}]}]}}`
	json2 := `{"Kind":"ReplicaSet","Namespace":"bar","Phase":"Running","Total":1,"_aws":{"Timestamp":1599037320000,"CloudWatchMetrics":[{"Namespace":"Skpr/Cluster","Dimensions":[["Kind","Namespace","Phase"]],"Metrics":[{"Name":"Total","Unit":"None"}]}]}}`
	json3 := `{"Kind":"ReplicaSet","Namespace":"foo","Phase":"Pending","Total":2,"_aws":{"Timestamp":1599037320000,"CloudWatchMetrics":[{"Namespace":"Skpr/Cluster","Dimensions":[["Kind","Namespace","Phase"]],"Metrics":[{"Name":"Total","Unit":"None"}]}]}}`
	json4 := `{"Kind":"ReplicaSet","Namespace":"foo","Phase":"Running","Total":1,"_aws":{"Timestamp":1599037320000,"CloudWatchMetrics":[{"Namespace":"Skpr/Cluster","Dimensions":[["Kind","Namespace","Phase"]],"Metrics":[{"Name":"Total","Unit":"None"}]}]}}`

	assert.Len(t, lines, 4)
	assert.Contains(t, lines, json3)
	assert.Contains(t, lines, json4)
	assert.Contains(t, lines, json1)
	assert.Contains(t, lines, json2)

}
