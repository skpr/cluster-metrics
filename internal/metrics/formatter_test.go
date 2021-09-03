package metrics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestFormat(t *testing.T) {
	namespace := "Skpr/Cluster"
	timestamp := time.Date(2020, time.September, 2, 9, 2, 0, 0, time.UTC)
	metric := Metric{
		Kind:      "ReplicaSet",
		Namespace: "skpr-test-project",
		Phase:     corev1.PodPending,
		Total:     5,
	}
	out := Format(timestamp, namespace, metric)

	json1 := `{"Kind":"ReplicaSet","Namespace":"skpr-test-project","Phase":"Pending","Total":5,"_aws":{"Timestamp":1599037320000,"CloudWatchMetrics":[{"Namespace":"Skpr/Cluster","Dimensions":[["Kind","Namespace","Phase"]],"Metrics":[{"Name":"Total","Unit":"None"}]}]}}`

	assert.Equal(t, json1, out)
}
