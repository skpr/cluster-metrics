package metrics

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Collector is ze metrics collector.
type Collector struct {
	clientset *kubernetes.Clientset
}

// NewCollector creates ze new metrics collector.
func NewCollector(clientset *kubernetes.Clientset) *Collector {
	return &Collector{
		clientset: clientset,
	}
}

// Collect ze metrics.
func (c *Collector) Collect(pods []corev1.Pod) *MetricSet {
	metrics := NewMetricSet()
	for _, pod := range pods {
		metrics.Increment(findOwnerKind(pod), pod.ObjectMeta.Namespace, pod.Status.Phase)
	}
	return metrics
}

// findOwnerKind find the owner kind.
func findOwnerKind(pod corev1.Pod) string {
	for _, ref := range pod.ObjectMeta.OwnerReferences {
		return ref.Kind
	}
	return ""
}

// ListPods gets ze list of pods.
func (c *Collector) ListPods() ([]corev1.Pod, error) {
	podList, err := c.clientset.CoreV1().Pods(corev1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return []corev1.Pod{}, err
	}
	return podList.Items, nil
}
