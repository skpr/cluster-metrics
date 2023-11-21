package metrics

import (
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CollectPods will collect ze metrics for pods.
func CollectPods(pods []corev1.Pod) (*MetricSet, StateSet) {
	metrics := NewMetricSet()
	stateSet := make(StateSet)
	stateSet["Pod"] = map[string]int{}
	for _, pod := range pods {

		// Determine the correct PodStatus to report on.
		var PodStatus string
		if pod.Status.Phase == "Pending" {
			for _, container := range pod.Status.ContainerStatuses {
				PodStatus = container.State.Waiting.Reason
			}
		} else if pod.Status.Phase == "Failed" {
			for _, container := range pod.Status.ContainerStatuses {
				PodStatus = container.State.Terminated.Reason
			}
		} else {
			PodStatus = string(pod.Status.Phase)
		}

		// Metrics for Logging
		metrics.Increment(findOwnerKind(pod.ObjectMeta), pod.ObjectMeta.Namespace, pod.Status.Phase)

		// Metrics for Pushing
		stateSet["Pod"][PodStatus]++

	}
	return metrics, stateSet
}

// CollectDeployments will collect ze metrics for deployments.
func CollectDeployments(deployments []appsv1.Deployment) (*MetricSet, StateSet) {
	metrics := NewMetricSet()
	stateSet := make(StateSet)
	stateSet["Deployment"] = map[string]int{}
	for _, deployment := range deployments {
		if deployment.Status.ReadyReplicas > 0 {
			metrics.Increment(findOwnerKind(deployment.ObjectMeta), deployment.ObjectMeta.Namespace, "Ready")
			stateSet["Deployment"]["Ready"]++
		} else {
			metrics.Increment(findOwnerKind(deployment.ObjectMeta), deployment.ObjectMeta.Namespace, "NotReady")
			stateSet["Deployment"]["NotReady"]++
		}
	}
	return metrics, stateSet
}

// CollectStatefulSets will collect ze metrics for statefulsets.
func CollectStatefulSets(statefulsets []appsv1.StatefulSet) (*MetricSet, StateSet) {
	metrics := NewMetricSet()
	stateSet := make(StateSet)
	stateSet["StatefulSet"] = map[string]int{}
	for _, statefulset := range statefulsets {
		if statefulset.Status.AvailableReplicas > 0 {
			metrics.Increment(findOwnerKind(statefulset.ObjectMeta), statefulset.ObjectMeta.Namespace, "Ready")
			stateSet["StatefulSet"]["Ready"]++
		} else {
			metrics.Increment(findOwnerKind(statefulset.ObjectMeta), statefulset.ObjectMeta.Namespace, "NotReady")
			stateSet["StatefulSet"]["NotReady"]++
		}
	}
	return metrics, stateSet
}

// CollectCronJobs will collect ze metrics for cronjobs.
func CollectCronJobs(cronjobs []batchv1.CronJob) (*MetricSet, StateSet) {
	metrics := NewMetricSet()
	stateSet := make(StateSet)
	stateSet["CronJob"] = map[string]int{}
	for _, cronjob := range cronjobs {
		if cronjob.Spec.Suspend != nil && *cronjob.Spec.Suspend {
			metrics.Increment(findOwnerKind(cronjob.ObjectMeta), cronjob.ObjectMeta.Namespace, "Suspended")
			stateSet["CronJob"]["Suspended"]++
		} else if cronjob.Spec.Suspend == nil || !*cronjob.Spec.Suspend {
			metrics.Increment(findOwnerKind(cronjob.ObjectMeta), cronjob.ObjectMeta.Namespace, "Active")
			stateSet["CronJob"]["Active"]++
		}
	}
	return metrics, stateSet
}

// CollectJobs will collect ze metrics for jobs.
func CollectJobs(jobs []batchv1.Job) (*MetricSet, StateSet) {
	metrics := NewMetricSet()
	stateSet := make(StateSet)
	stateSet["Job"] = map[string]int{}
	for _, job := range jobs {
		if job.Status.Failed > 0 {
			metrics.Increment(findOwnerKind(job.ObjectMeta), job.ObjectMeta.Namespace, "Failed")
			stateSet["Job"]["Failed"]++
		}
		if job.Status.Active > 0 {
			metrics.Increment(findOwnerKind(job.ObjectMeta), job.ObjectMeta.Namespace, "Active")
			stateSet["Job"]["Active"]++
		}
		if job.Status.Succeeded > 0 {
			metrics.Increment(findOwnerKind(job.ObjectMeta), job.ObjectMeta.Namespace, "Succeeded")
			stateSet["Job"]["Succeeded"]++
		}
	}
	return metrics, stateSet
}

// TODO ADD SKPR-SPECIFIC (PROJECT AND ENVIRONMENT) RESOURCES

// findOwnerKind find the owner kind.
func findOwnerKind(meta v1.ObjectMeta) string {
	for _, ref := range meta.OwnerReferences {
		return ref.Kind
	}
	return ""
}
