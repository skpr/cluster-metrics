package metrics

import (
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// KindPod is the cost reference to the Kubernetes Pod object.
	KindPod = "Pod"
	// KindDeployment is the cost reference to the Kubernetes Deployment object.
	KindDeployment = "Deployment"
	// KindReplicaSet is the cost reference to the Kubernetes ReplicaSet object.
	KindReplicaSet = "ReplicaSet"
	// KindStatefulSet is the cost reference to the Kubernetes StatefulSet object.
	KindStatefulSet = "StatefulSet"
	// KindCronJob is the cost reference to the Kubernetes CronJob object.
	KindCronJob = "CronJob"
	// KindJob is the cost reference to the Kubernetes Job object.
	KindJob = "Job"

	// StateScaledDown indicates the desired replicas for the object is 0.
	StateScaledDown = "ScaledDown"
	// StateReady is the const representing the state for an object being Ready
	StateReady = string(corev1.PodReady)
	// StateNotReady is the const representing the state for an object not being Ready
	StateNotReady = string(corev1.PodPending)
	// StateSuspended is the const representing the state for a CronJob being Suspended
	StateSuspended = string(batchv1.JobSuspended)
	// StateActive is the const representing the state for a CronJob being not Suspended
	StateActive = string(corev1.PodRunning)
	// StateFailed is the const representing the state for an object having been Failed
	StateFailed = string(batchv1.JobFailed)
	// StateSucceeded is the const representing the state for an object having been completed successfully
	StateSucceeded = string(batchv1.JobComplete)
)

// CollectPods will collect ze metrics for pods.
func CollectPods(pods []corev1.Pod) (*MetricSet, StateSet) {
	metrics := NewMetricSet()
	stateSet := make(StateSet)
	stateSet[KindPod] = map[string]int{}
	for _, pod := range pods {

		for _, containerStatus := range pod.Status.ContainerStatuses {

			if !containerStatus.Ready {
				// Loop over terminated containers and stash reason.
				if containerStatus.State.Terminated != nil {
					stateSet[KindPod][containerStatus.State.Terminated.Reason]++
				}

				// Loop over waiting containers and stash reason.
				if containerStatus.State.Waiting != nil {
					stateSet[KindPod][containerStatus.State.Waiting.Reason]++
				}
			}
		}

		// Metrics for Logging
		metrics.Increment(findOwnerKind(pod.ObjectMeta), pod.ObjectMeta.Namespace, string(pod.Status.Phase))

		// Metrics for Pushing
		stateSet[KindPod][string(pod.Status.Phase)]++

	}
	return metrics, stateSet
}

// CollectDeployments will collect ze metrics for deployments.
func CollectDeployments(deployments []appsv1.Deployment) (*MetricSet, StateSet) {
	metrics := NewMetricSet()
	stateSet := make(StateSet)
	stateSet[KindDeployment] = map[string]int{}
	for _, deployment := range deployments {
		if deployment.Spec.Replicas != nil && *deployment.Spec.Replicas == 0 {
			metrics.Increment(findOwnerKind(deployment.ObjectMeta), deployment.ObjectMeta.Namespace, StateScaledDown)
			stateSet[KindDeployment][StateScaledDown]++
		} else if deployment.Status.ReadyReplicas > 0 {
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
		if statefulset.Spec.Replicas != nil && *statefulset.Spec.Replicas == 0 {
			metrics.Increment(findOwnerKind(statefulset.ObjectMeta), statefulset.ObjectMeta.Namespace, StateScaledDown)
			stateSet[KindStatefulSet][StateScaledDown]++
		} else if statefulset.Status.AvailableReplicas > 0 {
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
	for _, replicaset := range replicasets {
		if replicaset.Spec.Replicas != nil && *replicaset.Spec.Replicas == 0 {
			metrics.Increment(findOwnerKind(replicaset.ObjectMeta), replicaset.ObjectMeta.Namespace, StateScaledDown)
			stateSet[KindReplicaSet][StateScaledDown]++
		} else if replicaset.Status.AvailableReplicas > 0 {
			metrics.Increment(findOwnerKind(replicaset.ObjectMeta), replicaset.ObjectMeta.Namespace, StateReady)
			stateSet[KindReplicaSet][StateReady]++
		} else {
			metrics.Increment(findOwnerKind(replicaset.ObjectMeta), replicaset.ObjectMeta.Namespace, StateNotReady)
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
