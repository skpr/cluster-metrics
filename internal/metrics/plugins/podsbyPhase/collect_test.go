package podsbyPhase

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestClient_Collect(t *testing.T) {
	values := []map[string]string{
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

	var pods []corev1.Pod
	for _, val := range values {
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

	// TODO Collect() no longer takes input so need to work around this...
	//metrics, phaseSet := client.Collect(pods)
	//assert.Equal(t, 3, metrics.Items["abc-def-Pending"].Value)
	//assert.Equal(t, 1, metrics.Items["abc-def-Succeeded"].Value)
	//assert.Equal(t, 1, metrics.Items["xyz-def-Succeeded"].Value)
	//assert.Equal(t, 2, metrics.Items["abc-def-Failed"].Value)
	//assert.Equal(t, 2, metrics.Items["abc-ghj-Running"].Value)
	//assert.Equal(t, 1, metrics.Items["xyz-ghj-Running"].Value)
	//
	//assert.Equal(t, 3, phaseSet["Running"])
}
