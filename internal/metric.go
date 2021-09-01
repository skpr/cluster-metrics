package internal

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

// MetricSet is the metric set.
type MetricSet struct {
	items map[string]*Metric
}

// Metric represents an individual metric.
type Metric struct {
	Kind      string
	Namespace string
	Phase     corev1.PodPhase
	Count     int
}

// Increment the metric.
func (s *MetricSet) Increment(kind, namespace string, phase corev1.PodPhase) {
	key := fmt.Sprintf("%s-%s-%s", kind, namespace, phase)
	if metric, found := s.items[key]; found {
		metric.Count++
	} else {
		metric := &Metric{
			Kind:      kind,
			Namespace: namespace,
			Phase:     phase,
			Count:     1,
		}
		s.items[key] = metric
	}
}
