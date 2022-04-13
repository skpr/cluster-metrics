package podsbyPhase

import (
	"bytes"
	metrics2 "github.com/skpr/cluster-metrics/internal/metrics"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_Log(t *testing.T) {
	metrics := metrics2.NewMetricSet()

	metrics.IncrementSelect("ReplicaSet", "foo", "replicaset-foo-pending", map[string]string{"phase": "Pending"})
	metrics.IncrementSelect("ReplicaSet", "foo", "replicaset-foo-pending", map[string]string{"phase": "Pending"})
	metrics.IncrementSelect("ReplicaSet", "foo", "replicaset-foo-running", map[string]string{"phase": "Running"})
	metrics.IncrementSelect("ReplicaSet", "bar", "replicaset-bar-pending", map[string]string{"phase": "Pending"})
	metrics.IncrementSelect("ReplicaSet", "bar", "replicaset-bar-running", map[string]string{"phase": "Running"})
	metrics.IncrementSelect("ReplicaSet", "bar", "replicaset-bar-running", map[string]string{"phase": "Running"})
	metrics.IncrementSelect("ReplicaSet", "bar", "replicaset-bar-running", map[string]string{"phase": "Running"})
	metrics.IncrementSelect("ReplicaSet", "bar", "replicaset-bar-running", map[string]string{"phase": "Running"})

	var buf bytes.Buffer
	var client Client
	err := client.Log(&buf, metrics)
	assert.NoError(t, err)
	s := buf.String()

	json1 := `{"name":"PodStatus","value":2,"type":"gauge","labels":{"kind":"ReplicaSet","namespace":"foo","phase":"Pending"}}`
	json2 := `{"name":"PodStatus","value":1,"type":"gauge","labels":{"kind":"ReplicaSet","namespace":"foo","phase":"Running"}}`
	json3 := `{"name":"PodStatus","value":1,"type":"gauge","labels":{"kind":"ReplicaSet","namespace":"bar","phase":"Pending"}}`
	json4 := `{"name":"PodStatus","value":4,"type":"gauge","labels":{"kind":"ReplicaSet","namespace":"bar","phase":"Running"}}`

	assert.Contains(t, s, json1)
	assert.Contains(t, s, json2)
	assert.Contains(t, s, json3)
	assert.Contains(t, s, json4)
}
