package metrics

import (
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	KindPod         = "Pod"
	KindDeployment  = "Deployment"
	KindReplicaSet  = "ReplicaSet"
	KindStatefulSet = "StatefulSet"
	KindCronJob     = "CronJob"
	KindJob         = "Job"

	StateReady     = string(corev1.PodReady)
	StateNotReady  = string(corev1.PodPending)
	StateSuspended = string(batchv1.JobSuspended)
	StateActive    = string(corev1.PodRunning)
	StateFailed    = string(batchv1.JobFailed)
	StateSucceeded = string(batchv1.JobComplete)
)

// CollectPods will collect ze metrics for pods.
func CollectPods(pods []corev1.Pod) (*MetricSet, StateSet) {
	metrics := NewMetricSet()
	stateSet := make(StateSet)
	stateSet[KindPod] = map[string]int{}
	for _, pod := range pods {

		// Determine the correct PodStatus to report on.
		var PodStatus string
		var Reasons []string

		for _, container := range pod.Status.ContainerStatuses {
			if !container.Ready {
				// Loop over terminated containers and stash reason.
				if container.State.Terminated != nil {
					Reasons = append(Reasons, container.State.Terminated.Reason)
				}

				// Loop over waiting containers and stash reason.
				if container.State.Waiting != nil {
					Reasons = append(Reasons, container.State.Waiting.Reason)
				}
			}
		}

		// Concat reasons into single string
		if len(Reasons) > 1 {
			PodStatus = strings.Join(Reasons, ",")
		} else if len(Reasons) == 1 {
			PodStatus = Reasons[0]
		}

		if PodStatus == "" {
			PodStatus = string(pod.Status.Phase)
		}

		// Metrics for Logging
		metrics.Increment(findOwnerKind(pod.ObjectMeta), pod.ObjectMeta.Namespace, string(pod.Status.Phase))

		// Metrics for Pushing
		stateSet[KindPod][PodStatus]++

	}
	return metrics, stateSet
}

// CollectDeployments will collect ze metrics for deployments.
func CollectDeployments(deployments []appsv1.Deployment) (*MetricSet, StateSet) {
	metrics := NewMetricSet()
	stateSet := make(StateSet)
	stateSet[KindDeployment] = map[string]int{}
	for _, deployment := range deployments {
		if deployment.Status.ReadyReplicas > 0 {
			metrics.Increment(findOwnerKind(deployment.ObjectMeta), deployment.ObjectMeta.Namespace, StateReady)
			stateSet[KindDeployment][StateReady]++
		} else {
			metrics.Increment(findOwnerKind(deployment.ObjectMeta), deployment.ObjectMeta.Namespace, StateNotReady)
			stateSet[KindDeployment][StateNotReady]++
		}
	}
	return metrics, stateSet
}

// CollectStatefulSets will collect ze metrics for statefulsets.
func CollectStatefulSets(statefulsets []appsv1.StatefulSet) (*MetricSet, StateSet) {
	metrics := NewMetricSet()
	stateSet := make(StateSet)
	stateSet[KindStatefulSet] = map[string]int{}
	for _, statefulset := range statefulsets {
		if statefulset.Status.AvailableReplicas > 0 {
			metrics.Increment(findOwnerKind(statefulset.ObjectMeta), statefulset.ObjectMeta.Namespace, StateReady)
			stateSet[KindStatefulSet][StateReady]++
		} else {
			metrics.Increment(findOwnerKind(statefulset.ObjectMeta), statefulset.ObjectMeta.Namespace, StateNotReady)
			stateSet[KindStatefulSet][StateNotReady]++
		}
	}
	return metrics, stateSet
}

// CollectReplicaSets will collect ze metrics for replicasets.
func CollectReplicaSets(replicasets []appsv1.ReplicaSet) (*MetricSet, StateSet) {
	metrics := NewMetricSet()
	stateSet := make(StateSet)
	stateSet[KindReplicaSet] = map[string]int{}
	for _, statefulset := range replicasets {
		if statefulset.Status.AvailableReplicas > 0 {
			metrics.Increment(findOwnerKind(statefulset.ObjectMeta), statefulset.ObjectMeta.Namespace, StateReady)
			stateSet[KindReplicaSet][StateReady]++
		} else {
			metrics.Increment(findOwnerKind(statefulset.ObjectMeta), statefulset.ObjectMeta.Namespace, StateNotReady)
			stateSet[KindReplicaSet][StateNotReady]++
		}
	}
	return metrics, stateSet
}

// CollectCronJobs will collect ze metrics for cronjobs.
func CollectCronJobs(cronjobs []batchv1.CronJob) (*MetricSet, StateSet) {
	metrics := NewMetricSet()
	stateSet := make(StateSet)
	stateSet[KindCronJob] = map[string]int{}
	for _, cronjob := range cronjobs {
		if cronjob.Spec.Suspend != nil && *cronjob.Spec.Suspend {
			metrics.Increment(findOwnerKind(cronjob.ObjectMeta), cronjob.ObjectMeta.Namespace, StateSuspended)
			stateSet[KindCronJob][StateSuspended]++
		} else if cronjob.Spec.Suspend == nil || !*cronjob.Spec.Suspend {
			metrics.Increment(findOwnerKind(cronjob.ObjectMeta), cronjob.ObjectMeta.Namespace, StateActive)
			stateSet[KindCronJob][StateActive]++
		}
	}
	return metrics, stateSet
}

// CollectJobs will collect ze metrics for jobs.
func CollectJobs(jobs []batchv1.Job) (*MetricSet, StateSet) {
	metrics := NewMetricSet()
	stateSet := make(StateSet)
	stateSet[KindJob] = map[string]int{}
	for _, job := range jobs {
		if job.Status.Failed > 0 {
			metrics.Increment(findOwnerKind(job.ObjectMeta), job.ObjectMeta.Namespace, StateFailed)
			stateSet[KindJob][StateFailed]++
		}
		if job.Status.Active > 0 {
			metrics.Increment(findOwnerKind(job.ObjectMeta), job.ObjectMeta.Namespace, StateActive)
			stateSet[KindJob][StateActive]++
		}
		if job.Status.Succeeded > 0 {
			metrics.Increment(findOwnerKind(job.ObjectMeta), job.ObjectMeta.Namespace, StateSucceeded)
			stateSet[KindJob][StateSucceeded]++
		}
	}
	return metrics, stateSet
}

// findOwnerKind find the owner kind.
func findOwnerKind(meta v1.ObjectMeta) string {
	for _, ref := range meta.OwnerReferences {
		return ref.Kind
	}
	return ""
}
