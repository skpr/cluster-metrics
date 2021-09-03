package metrics

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

// MetricSet is the metric set.
type MetricSet struct {
	Items      map[string]*Metric
}

// NewMetricSet creates a new metric set.
func NewMetricSet() MetricSet {
	return MetricSet{
		Items: make(map[string]*Metric),
	}
}

// Metric represents an individual metric.
type Metric struct {
	Kind      string
	Namespace string
	Phase corev1.PodPhase
	Total int
}

// Increment the metric.
func (s *MetricSet) Increment(kind, namespace string, phase corev1.PodPhase) {
	key := fmt.Sprintf("%s-%s-%s", kind, namespace, phase)
	if metric, found := s.Items[key]; found {
		metric.Total++
	} else {
		metric := &Metric{
			Kind:      kind,
			Namespace: namespace,
			Phase:     phase,
			Total:     1,
		}
		s.Items[key] = metric
	}
}
