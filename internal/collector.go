package internal

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// MetricsCollector is the metrics collector.
type MetricsCollector struct {
	clientset *kubernetes.Clientset
}

// NewMetricsCollector creates a new metrics collector.
func NewMetricsCollector(clientset *kubernetes.Clientset) *MetricsCollector {
	return &MetricsCollector{
		clientset: clientset,
	}
}

// Collect ze metrics.
func (c *MetricsCollector) Collect() (*MetricSet, error) {
	pods, err := c.clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return &MetricSet{}, err
	}
	metrics := &MetricSet{}
	for _, pod := range pods.Items {
		metrics.Increment(pod.OwnerReferences[0].Kind, pod.ObjectMeta.Namespace, pod.Status.Phase)
	}
	return metrics, nil
}
