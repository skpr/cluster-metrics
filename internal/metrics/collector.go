package metrics

import (
	corev1 "k8s.io/api/core/v1"
)

// Collect ze metrics.
func Collect(pods []corev1.Pod) MetricSet {
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
