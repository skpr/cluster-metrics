package metrics

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

// PhaseSet is the phase set.
type PhaseSet map[string]int

// MetricSet is the metric set.
type MetricSet struct {
	Items map[string]*Metric
}

// NewMetricSet creates a new metric set.
func NewMetricSet() *MetricSet {
	return &MetricSet{
		Items: make(map[string]*Metric),
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
func (s *MetricSet) Increment(kind, namespace string, phase corev1.PodPhase) {
	key := fmt.Sprintf("%s-%s-%s", kind, namespace, phase)
	if metric, found := s.Items[key]; found {
		metric.Value++
	} else {
		metric := &Metric{
			Labels: map[string]string{
				dimensionKind:      kind,
				dimensionNamespace: namespace,
				dimensionPhase:     string(phase),
			},
			Value: 1,
		}
		s.Items[key] = metric
	}
}
