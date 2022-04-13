package metrics

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

// IncrementSelect will selectively increment the metric based upon uuid input.
func (s *MetricSet) IncrementSelect(kind, namespace string, uuid string, fields map[string]string) {
	if metric, found := s.Items[uuid]; found {
		metric.Value++
	} else {
		metric := &Metric{
			Labels: map[string]string{
				dimensionKind:      kind,
				dimensionNamespace: namespace,
			},
			Value: 1,
		}
		for i, field := range fields {
			metric.Labels[i] = field
		}
		s.Items[uuid] = metric
	}
}
