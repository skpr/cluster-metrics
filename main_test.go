package main

import (
	"testing"

	"github.com/skpr/cluster-metrics/internal/metrics"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestWorkflow(t *testing.T) {
	mts, phases := &metrics.MetricSet{}, metrics.StateSet{}

	{
		pods := corev1.PodList{
			Items: []corev1.Pod{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "abc",
						Namespace: "project-a",
						OwnerReferences: []metav1.OwnerReference{
							{
								Kind: "Pod",
							},
						},
					},
					Status: corev1.PodStatus{
						Phase: "Running",
					},
				},
				{
					TypeMeta: metav1.TypeMeta{
						Kind: "Pod",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "def",
						Namespace: "project-b",
						OwnerReferences: []metav1.OwnerReference{
							{
								Kind: "Pod",
							},
						},
					},
					Status: corev1.PodStatus{
						Phase: "Pending",
					},
				},
				{
					TypeMeta: metav1.TypeMeta{
						Kind: "Pod",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "ghi",
						Namespace: "project-c",
						OwnerReferences: []metav1.OwnerReference{
							{
								Kind: "Pod",
							},
						},
					},
					Status: corev1.PodStatus{
						Phase: "Failed",
					},
				},
				{
					TypeMeta: metav1.TypeMeta{
						Kind: "Pod",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "jkl",
						Namespace: "project-d",
						OwnerReferences: []metav1.OwnerReference{
							{
								Kind: "Pod",
							},
						},
					},
					Status: corev1.PodStatus{
						Phase: "Succeeded",
					},
				},
				{
					TypeMeta: metav1.TypeMeta{
						Kind: "Pod",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "jkl",
						Namespace: "project-d",
						OwnerReferences: []metav1.OwnerReference{
							{
								Kind: "Pod",
							},
						},
					},
					Status: corev1.PodStatus{
						Phase: "Succeeded",
					},
				},
			},
		}

		// Collect the metrics.
		metricSetAddition, stateSetAddition := metrics.CollectPods(pods.Items)
		mts = metrics.CombineRecords(mts, metricSetAddition)

		for i, v := range stateSetAddition {
			phases[i] = v
		}
	}

	assert.Len(t, mts.Items, 4)
	assert.Len(t, phases["Pod"], 4)

	assert.Equal(t, mts.Items["Pod-project-a-Running"].Value, 1)
	assert.Equal(t, mts.Items["Pod-project-a-Running"].Labels["kind"], "Pod")
	assert.Equal(t, mts.Items["Pod-project-a-Running"].Labels["namespace"], "project-a")
	assert.Equal(t, mts.Items["Pod-project-a-Running"].Labels["phase"], "Running")

	assert.Equal(t, mts.Items["Pod-project-b-Pending"].Value, 1)
	assert.Equal(t, mts.Items["Pod-project-b-Pending"].Labels["kind"], "Pod")
	assert.Equal(t, mts.Items["Pod-project-b-Pending"].Labels["namespace"], "project-b")
	assert.Equal(t, mts.Items["Pod-project-b-Pending"].Labels["phase"], "Pending")

	assert.Equal(t, mts.Items["Pod-project-c-Failed"].Value, 1)
	assert.Equal(t, mts.Items["Pod-project-c-Failed"].Labels["kind"], "Pod")
	assert.Equal(t, mts.Items["Pod-project-c-Failed"].Labels["namespace"], "project-c")
	assert.Equal(t, mts.Items["Pod-project-c-Failed"].Labels["phase"], "Failed")

	assert.Equal(t, mts.Items["Pod-project-d-Succeeded"].Value, 2)
	assert.Equal(t, mts.Items["Pod-project-d-Succeeded"].Labels["kind"], "Pod")
	assert.Equal(t, mts.Items["Pod-project-d-Succeeded"].Labels["namespace"], "project-d")
	assert.Equal(t, mts.Items["Pod-project-d-Succeeded"].Labels["phase"], "Succeeded")

	assert.Equal(t, phases["Pod"]["Running"], 1)
	assert.Equal(t, phases["Pod"]["Pending"], 1)
	assert.Equal(t, phases["Pod"]["Failed"], 1)
	assert.Equal(t, phases["Pod"]["Succeeded"], 2)
}
