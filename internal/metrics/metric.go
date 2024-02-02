package metrics

import (
	"fmt"
)

// StateSet is the state of the metrics.
// ie Pod Phase, Deployment Status or CronJob Suspended state.
type StateSet map[string]map[string]int

// MetricSet is the metric set.
type MetricSet struct {
	Items map[string]*Metric
}

// NewMetricSet creates a new metric set.
func NewMetricSet() *MetricSet {
	return &MetricSet{
		Items: map[string]*Metric{},
	}
}

// Metric represents an individual metric.
type Metric struct {
	Name   string            `json:"name"`
	Value  int               `json:"value"`
	Type   string            `json:"type"`
	Labels map[string]string `json:"labels"`
}

// Increment the metric.
func (s *MetricSet) Increment(kind, namespace string, phase string) {
	key := fmt.Sprintf("%s-%s-%s", kind, namespace, phase)
	if metric, found := s.Items[key]; found {
		metric.Value++
	} else {
		metric := &Metric{
			Name:  key,
			Value: 1,
			Labels: map[string]string{
				dimensionKind:      kind,
				dimensionNamespace: namespace,
				dimensionState:     string(phase),
			},
		}
		s.Items[key] = metric
	}
}

// CombineRecords will combine two metric sets.
func CombineRecords(recordsInput *MetricSet, recordsAppend *MetricSet) *MetricSet {

	output := NewMetricSet()

	for _, record := range recordsInput.Items {
		output.Items[record.Name] = record
	}

	for _, record := range recordsAppend.Items {
		if output.Items[record.Name] == nil {
			output.Items[record.Name] = record
		} else {
			output.Items[record.Name].Value++
		}
	}

	return output
}

// CombineStates will combine two state sets.
func CombineStates(input *StateSet, og *StateSet) StateSet {
	output := StateSet{}

	for i, v := range *input {
		for a, b := range v {
			if output[i] == nil {
				output[i] = make(map[string]int)
			}
			if output[i][a] == 0 {
				output[i][a] = b
			} else {
				output[i][a] += b
			}
		}

	}

	for i, v := range *og {
		for a, b := range v {
			if output[i] == nil {
				output[i] = make(map[string]int)
			}
			if output[i][a] == 0 {
				output[i][a] = b
			} else {
				output[i][a] += b
			}
		}
	}

	return output
}
