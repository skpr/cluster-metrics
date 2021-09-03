package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestMetricsCollector_CollectMetrics(t *testing.T) {

	vals := provideTestValues()

	pods := []corev1.Pod{}
	for _, val := range vals {
		pods = append(pods, corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: val["namespace"],
				OwnerReferences: []metav1.OwnerReference{
					{
						Kind: val["kind"],
					},
				},
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodPhase(val["phase"]),
			},
		})
	}

	metrics := Collect(pods)

	assert.Equal(t, 3, metrics.Items["abc-def-Pending"].Total)
	assert.Equal(t, 1, metrics.Items["abc-def-Succeeded"].Total)
	assert.Equal(t, 1, metrics.Items["xyz-def-Succeeded"].Total)
	assert.Equal(t, 2, metrics.Items["abc-def-Failed"].Total)
	assert.Equal(t, 2, metrics.Items["abc-ghj-Running"].Total)
	assert.Equal(t, 1, metrics.Items["xyz-ghj-Running"].Total)

}

func provideTestValues() []map[string]string {
	vals := []map[string]string{
		{"kind": "abc", "namespace": "def", "phase": string(corev1.PodPending)},
		{"kind": "abc", "namespace": "def", "phase": string(corev1.PodPending)},
		{"kind": "abc", "namespace": "def", "phase": string(corev1.PodPending)},
		{"kind": "abc", "namespace": "def", "phase": string(corev1.PodSucceeded)},
		{"kind": "xyz", "namespace": "def", "phase": string(corev1.PodSucceeded)},
		{"kind": "abc", "namespace": "def", "phase": string(corev1.PodFailed)},
		{"kind": "abc", "namespace": "def", "phase": string(corev1.PodFailed)},
		{"kind": "abc", "namespace": "ghj", "phase": string(corev1.PodRunning)},
		{"kind": "abc", "namespace": "ghj", "phase": string(corev1.PodRunning)},
		{"kind": "xyz", "namespace": "ghj", "phase": string(corev1.PodRunning)},
	}
	return vals
}
