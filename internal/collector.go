package internal

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// MetricsCollector is ze metrics collector.
type MetricsCollector struct {
	clientset *kubernetes.Clientset
}

// NewMetricsCollector creates ze new metrics collector.
func NewMetricsCollector(clientset *kubernetes.Clientset) *MetricsCollector {
	return &MetricsCollector{
		clientset: clientset,
	}
}

// CollectMetrics collects ze metrics.
func (c *MetricsCollector) CollectMetrics(pods []v1.Pod) *MetricSet {
	metrics := NewMetricSet()
	for _, pod := range pods {
		metrics.Increment(pod.ObjectMeta.OwnerReferences[0].Kind, pod.ObjectMeta.Namespace, pod.Status.Phase)
	}
	return metrics
}

// ListPods gets ze list of pods.
func (c *MetricsCollector) ListPods() ([]v1.Pod, error) {
	podList, err := c.clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return []v1.Pod{}, err
	}
	return podList.Items, nil
}
