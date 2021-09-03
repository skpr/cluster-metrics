package metrics

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"

	"github.com/skpr/cluster-metrics/internal/metrics/mock"
)

// TestPusher_Push tests the push function.
func TestPusher_Push(t *testing.T) {

	mts := NewMetricSet()
	mts.Increment("ReplicaSet", "foo", corev1.PodPending)
	mts.Increment("ReplicaSet", "foo", corev1.PodPending)
	mts.Increment("ReplicaSet", "foo", corev1.PodRunning)
	mts.Increment("ReplicaSet", "bar", corev1.PodPending)
	mts.Increment("ReplicaSet", "bar", corev1.PodRunning)

	namespace := "Skpr/Cluster"
	timestamp := time.Date(2020, time.September, 2, 9, 2, 0, 0, time.UTC)

	cloudwatch := mock.NewCloudwatch()

	metricData := ConvertToMetricData(timestamp, mts)

	datum := metricData[0]
	assert.Equal(t, timestamp, *datum.Timestamp)
	assert.Len(t, datum.Dimensions, 3)
	assert.Greater(t, int(*datum.Value), 0)
	assert.Equal(t, metricTotal, *datum.MetricName)

	pusher := NewPusher(cloudwatch)
	err := pusher.Push(context.TODO(), namespace, metricData)
	assert.NoError(t, err)

}
