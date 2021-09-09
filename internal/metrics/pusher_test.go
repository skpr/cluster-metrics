package metrics

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/skpr/cluster-metrics/internal/metrics/mock"
)

// TestPusher_Push tests the push function.
func TestPusher_Push(t *testing.T) {

	phases := make(PhaseSet)
	phases["Pending"]++
	phases["Pending"]++
	phases["Running"]++
	phases["Running"]++
	phases["Running"]++

	namespace := "Skpr/Cluster"
	timestamp := time.Date(2020, time.September, 2, 9, 2, 0, 0, time.UTC)

	cloudwatch := mock.NewCloudwatch()
	cluster := "foo"
	metricData := ConvertToMetricData(timestamp, cluster, phases)

	datum1 := metricData[0]
	assert.Equal(t, timestamp, *datum1.Timestamp)
	assert.Len(t, datum1.Dimensions, 2)
	assert.Equal(t, "phase", *datum1.Dimensions[0].Name)
	assert.Equal(t, "Pending", *datum1.Dimensions[0].Value)
	assert.Equal(t, int(*datum1.Value), 2)
	assert.Equal(t, metricTotal, *datum1.MetricName)

	datum2 := metricData[1]
	assert.Equal(t, timestamp, *datum2.Timestamp)
	assert.Len(t, datum2.Dimensions, 2)
	assert.Equal(t, "phase", *datum2.Dimensions[0].Name)
	assert.Equal(t, "Running", *datum2.Dimensions[0].Value)
	assert.Equal(t, int(*datum2.Value), 3)
	assert.Equal(t, metricTotal, *datum2.MetricName)

	pusher := NewPusher(cloudwatch)
	err := pusher.Push(context.TODO(), namespace, metricData)
	assert.NoError(t, err)

}
