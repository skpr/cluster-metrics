package metrics

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCombineRecords(t *testing.T) {
	starting := MetricSet{
		Items: map[string]*Metric{
			"Pod-1": &Metric{
				Name:  "Pod-default-Running",
				Value: 1,
				Labels: map[string]string{
					"kind":      "Pod",
					"namespace": "default",
					"phase":     "Running",
				},
			},
			"Pod-2": &Metric{
				Name:  "Pod-default-Failed",
				Value: 1,
				Labels: map[string]string{
					"kind":      "Pod",
					"namespace": "default",
					"phase":     "Failed",
				},
			},
		},
	}

	adding := MetricSet{
		Items: map[string]*Metric{
			"Pod-3": &Metric{
				Name:  "Pod-default-Running",
				Value: 1,
				Labels: map[string]string{
					"kind":      "Pod",
					"namespace": "default",
					"phase":     "Running",
				},
			},
			"Pod-4": &Metric{
				Name:  "Pod-default-Failed",
				Value: 1,
				Labels: map[string]string{
					"kind":      "Pod",
					"namespace": "default",
					"phase":     "Failed",
				},
			},
			"Pod-5": &Metric{
				Name:  "Pod-default-Succeeded",
				Value: 1,
				Labels: map[string]string{
					"kind":      "Pod",
					"namespace": "default",
					"phase":     "Succeeded",
				},
			},
		},
	}

	expectedOutput := MetricSet{
		Items: map[string]*Metric{
			"Pod-default-Running": &Metric{
				Name:  "Pod-default-Running",
				Value: 2,
				Labels: map[string]string{
					"kind":      "Pod",
					"namespace": "default",
					"phase":     "Running",
				},
			},
			"Pod-default-Failed": &Metric{
				Name:  "Pod-default-Failed",
				Value: 2,
				Labels: map[string]string{
					"kind":      "Pod",
					"namespace": "default",
					"phase":     "Failed",
				},
			},
			"Pod-default-Succeeded": &Metric{
				Name:  "Pod-default-Running",
				Value: 1,
				Labels: map[string]string{
					"kind":      "Pod",
					"namespace": "default",
					"phase":     "Running",
				},
			},
		},
	}

	output := CombineRecords(&starting, &adding)
	fmt.Sprint(output, expectedOutput)

	assert.Equal(t, expectedOutput.Items["Pod-default-Running"].Value, output.Items["Pod-default-Running"].Value)
	assert.Equal(t, expectedOutput.Items["Pod-default-Failed"].Value, output.Items["Pod-default-Failed"].Value)
	assert.Equal(t, expectedOutput.Items["Pod-default-Succeeded"].Value, output.Items["Pod-default-Succeeded"].Value)

	assert.Equal(t, len(expectedOutput.Items), len(output.Items))
}

func TestCombineStates(t *testing.T) {

	one := &StateSet{
		"Pod": {
			"y": 1,
		},
		"PretendObject": {
			"c": 3,
		},
		"Deployment": {
			"b": 1,
		},
	}

	two := &StateSet{
		"Deployment": {
			"b": 2,
		},
		"OtherPretendObject": {
			"c": 3,
		},
	}

	three := CombineStates(one, two)
	fmt.Sprint(three)
}
