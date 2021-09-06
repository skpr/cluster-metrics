package metrics

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestLog(t *testing.T) {
	metrics := NewMetricSet()
	metrics.Increment("ReplicaSet", "foo", corev1.PodPending)
	metrics.Increment("ReplicaSet", "foo", corev1.PodPending)
	metrics.Increment("ReplicaSet", "foo", corev1.PodRunning)
	metrics.Increment("ReplicaSet", "bar", corev1.PodPending)
	metrics.Increment("ReplicaSet", "bar", corev1.PodRunning)
	metrics.Increment("ReplicaSet", "bar", corev1.PodRunning)
	metrics.Increment("ReplicaSet", "bar", corev1.PodRunning)
	metrics.Increment("ReplicaSet", "bar", corev1.PodRunning)

	var buf bytes.Buffer
	err := Log(&buf, "pod_count_total", metrics)
	assert.NoError(t, err)
	s := buf.String()
	fmt.Println(s)

	json1 := `{"name":"pod_count_total","value":2,"type":"gauge","labels":{"kind":"ReplicaSet","namespace":"foo","phase":"Pending"}}`
	json2 := `{"name":"pod_count_total","value":1,"type":"gauge","labels":{"kind":"ReplicaSet","namespace":"foo","phase":"Running"}}`
	json3 := `{"name":"pod_count_total","value":1,"type":"gauge","labels":{"kind":"ReplicaSet","namespace":"bar","phase":"Pending"}}`
	json4 := `{"name":"pod_count_total","value":4,"type":"gauge","labels":{"kind":"ReplicaSet","namespace":"bar","phase":"Running"}}`

	assert.Contains(t, s, json1)
	assert.Contains(t, s, json2)
	assert.Contains(t, s, json3)
	assert.Contains(t, s, json4)
}
