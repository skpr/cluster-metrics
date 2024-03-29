package metrics

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/stretchr/testify/assert"

	"github.com/skpr/cluster-metrics/internal/metrics/mock"
)

// TestPusher_Push tests the push function.
func TestPusher_Push(t *testing.T) {

	phases := StateSet{}

	phases["Pod"] = map[string]int{}

	phases["Pod"]["Pending"]++
	phases["Pod"]["Pending"]++
	phases["Pod"]["Running"]++
	phases["Pod"]["Running"]++
	phases["Pod"]["Running"]++

	namespace := "Skpr/Cluster"
	timestamp := time.Date(2020, time.September, 2, 9, 2, 0, 0, time.UTC)

	cloudwatch := mock.NewCloudwatch()
	cluster := "foo"
	metricData := ConvertToMetricData(timestamp, cluster, phases)

	for i, v := range metricData {
		if *v.Dimensions[1].Value == "Pending" {
			datum1 := metricData[i]
			assert.Equal(t, timestamp, *datum1.Timestamp)
			assert.Len(t, datum1.Dimensions, 3)
			assert.Equal(t, "phase", *datum1.Dimensions[1].Name)
			assert.Equal(t, "Pending", *datum1.Dimensions[1].Value)
			assert.Equal(t, int(*datum1.Value), 2)
			assert.Equal(t, metricTotal, *datum1.MetricName)
		} else if *v.Dimensions[1].Value == "Running" {
			datum2 := metricData[i]
			assert.Equal(t, timestamp, *datum2.Timestamp)
			assert.Len(t, datum2.Dimensions, 3)
			assert.Equal(t, "phase", *datum2.Dimensions[1].Name)
			assert.Equal(t, "Running", *datum2.Dimensions[1].Value)
			assert.Equal(t, int(*datum2.Value), 3)
			assert.Equal(t, metricTotal, *datum2.MetricName)
		} else {
			t.Fail()
		}
	}

	pusher := NewPusher(cloudwatch)
	err := pusher.Push(context.TODO(), namespace, metricData)
	assert.NoError(t, err)

}

func TestPusher_ConvertToMetricData(t *testing.T) {

	input := StateSet{
		"Pod": map[string]int{
			"1": 3,
		},
		"Deployment": map[string]int{
			"2": 5,
		},
	}

	timestamp := aws.Time(time.Now())
	expected := []types.MetricDatum{
		{
			MetricName: aws.String(metricTotal),
			Dimensions: []types.Dimension{
				{
					Name:  aws.String(dimensionKind),
					Value: aws.String("Pod"),
				},
				{
					Name:  aws.String(dimensionState),
					Value: aws.String("1"),
				},
				{
					Name:  aws.String(dimensionCluster),
					Value: aws.String("skpr-test"),
				},
			},
			Timestamp: timestamp,
			Value:     aws.Float64(float64(3)),
		},
		{
			MetricName: aws.String(metricTotal),
			Dimensions: []types.Dimension{
				{
					Name:  aws.String(dimensionKind),
					Value: aws.String("Deployment"),
				},
				{
					Name:  aws.String(dimensionState),
					Value: aws.String("2"),
				},
				{
					Name:  aws.String(dimensionCluster),
					Value: aws.String("skpr-test"),
				},
			},
			Timestamp: timestamp,
			Value:     aws.Float64(float64(5)),
		},
	}

	data := ConvertToMetricData(*timestamp, "skpr-test", input)

	for _, v := range expected {

		itemOne, errOne := getDataItemWithValue(*v.Dimensions[0].Value, data)
		itemTwo, errTwo := getDataItemWithValue(*v.Dimensions[1].Value, data)
		itemThree, errThree := getDataItemWithValue(*v.Dimensions[2].Value, data)

		assert.Nil(t, errOne)
		assert.Nil(t, errTwo)
		assert.Nil(t, errThree)

		assert.Equal(t, *v.Dimensions[0].Value, *itemOne.Dimensions[0].Value)
		assert.Equal(t, *v.Dimensions[1].Value, *itemTwo.Dimensions[1].Value)
		assert.Equal(t, *v.Dimensions[2].Value, *itemThree.Dimensions[2].Value)
	}
}

// Helper function for tests to prevent flake.
func getDataItemWithValue(want string, input []types.MetricDatum) (types.MetricDatum, error) {
	for _, v := range input {
		for x := range v.Dimensions {
			if *v.Dimensions[x].Value == want {
				return v, nil
			}
		}
	}
	return types.MetricDatum{}, errors.New("not found")
}
