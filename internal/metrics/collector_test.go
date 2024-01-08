package metrics

import (
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestMetricsCollector_CollectMetrics_PodMultipleStatuses(t *testing.T) {

	values := provideTestPodValuesMultipleErrors()

	var pods []corev1.Pod
	for _, val := range values {

		Pod := corev1.Pod{
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
		}

		if val["namespace"] == "abc" {
			Pod.Status.ContainerStatuses = []corev1.ContainerStatus{
				{
					Name: "ContainerTerminated",
					State: corev1.ContainerState{
						Terminated: &corev1.ContainerStateTerminated{
							Reason: "CustomErrImageForbidden",
						},
					},
				},
				{
					Name: "ContainerWaiting",
					State: corev1.ContainerState{
						Waiting: &corev1.ContainerStateWaiting{
							Reason: "CustomErrImagePull",
						},
					},
				},
				{
					Name:  "ContainerReady",
					Ready: true,
				},
			}
		}

		if val["namespace"] == "def" {
			Pod.Status.ContainerStatuses = []corev1.ContainerStatus{
				{
					Name: "ContainerTerminated",
					State: corev1.ContainerState{
						Terminated: &corev1.ContainerStateTerminated{
							Reason: "ErrCrashloopBackoff",
						},
					},
				},
				{
					Name: "ContainerWaiting",
					State: corev1.ContainerState{
						Waiting: &corev1.ContainerStateWaiting{
							Reason: "ErrImgPull",
						},
					},
				},
			}
		}

		if val["namespace"] == "ghi" {
			Pod.Status.ContainerStatuses = []corev1.ContainerStatus{
				{
					Name: "ContainerTerminated",
					State: corev1.ContainerState{
						Terminated: &corev1.ContainerStateTerminated{
							Reason: "CustomErrImageForbidden",
						},
					},
				},
				{
					Name: "ContainerWaiting",
					State: corev1.ContainerState{
						Waiting: &corev1.ContainerStateWaiting{
							Reason: "CrashLoopBackoff",
						},
					},
				},
			}
		}

		if val["namespace"] == "hjk" {
			Pod.Status.ContainerStatuses = []corev1.ContainerStatus{
				{
					Name: "ContainerTerminated",
					State: corev1.ContainerState{
						Terminated: &corev1.ContainerStateTerminated{
							Reason: "CustomUnknown",
						},
					},
				},
				{
					Name: "ContainerWaiting",
					State: corev1.ContainerState{
						Waiting: &corev1.ContainerStateWaiting{
							Reason: "Unknown",
						},
					},
				},
			}
		}

		if val["namespace"] == "lmn" {
			Pod.Status.ContainerStatuses = []corev1.ContainerStatus{
				{
					Name: "ContainerTerminated",
					State: corev1.ContainerState{
						Terminated: &corev1.ContainerStateTerminated{
							Reason: "CustomUnknown",
						},
					},
				},
				{
					Name:  "ContainerReady",
					Ready: true,
				},
			}
		}
		pods = append(pods, Pod)
	}

	assert.Equal(t, 11, len(pods))
	metrics, stateSet := CollectPods(pods)

	assert.Equal(t, 10, len(metrics.Items))
	assert.Equal(t, 11, len(pods))

	assert.Equal(t, 1, metrics.Items["Pod-abc-Pending"].Value)
	assert.Equal(t, 1, metrics.Items["Pod-abc-Failed"].Value)
	assert.Equal(t, 1, metrics.Items["Pod-def-Pending"].Value)
	assert.Equal(t, 1, metrics.Items["Pod-def-Failed"].Value)
	assert.Equal(t, 1, metrics.Items["Pod-ghi-Pending"].Value)
	assert.Equal(t, 1, metrics.Items["Pod-ghi-Failed"].Value)
	assert.Equal(t, 2, metrics.Items["Pod-hjk-Pending"].Value)
	assert.Equal(t, 1, metrics.Items["Pod-hjk-Failed"].Value)
	assert.Equal(t, 1, metrics.Items["Pod-lmn-Pending"].Value)
	assert.Equal(t, 1, metrics.Items["Pod-lmn-Failed"].Value)

	assert.Equal(t, 2, stateSet["Pod"]["CustomErrImageForbidden,CustomErrImagePull"])
	assert.Equal(t, 2, stateSet["Pod"]["ErrCrashloopBackoff,ErrImgPull"])
	assert.Equal(t, 2, stateSet["Pod"]["CustomErrImageForbidden,CrashLoopBackoff"])
	assert.Equal(t, 3, stateSet["Pod"]["CustomUnknown,Unknown"])
	assert.Equal(t, 2, stateSet["Pod"]["CustomUnknown"])
}

func TestMetricsCollector_CollectMetrics_Pods(t *testing.T) {

	values := provideTestPodValues()

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

	metrics, stateSet := CollectPods(pods)

	assert.Equal(t, 3, metrics.Items["Pod-def-Pending"].Value)
	assert.Equal(t, 1, metrics.Items["Pod-def-Succeeded"].Value)
	assert.Equal(t, 2, metrics.Items["Pod-def-Failed"].Value)
	assert.Equal(t, 3, metrics.Items["Pod-ghj-Running"].Value)

	assert.Equal(t, 3, stateSet["Pod"]["Running"])
}

func provideTestPodValues() []map[string]string {
	vals := []map[string]string{
		{"kind": "Pod", "namespace": "def", "phase": string(corev1.PodPending)},
		{"kind": "Pod", "namespace": "def", "phase": string(corev1.PodPending)},
		{"kind": "Pod", "namespace": "def", "phase": string(corev1.PodPending)},
		{"kind": "Pod", "namespace": "def", "phase": string(corev1.PodSucceeded)},
		{"kind": "Pod", "namespace": "def", "phase": string(corev1.PodFailed)},
		{"kind": "Pod", "namespace": "def", "phase": string(corev1.PodFailed)},
		{"kind": "Pod", "namespace": "ghj", "phase": string(corev1.PodRunning)},
		{"kind": "Pod", "namespace": "ghj", "phase": string(corev1.PodRunning)},
		{"kind": "Pod", "namespace": "ghj", "phase": string(corev1.PodRunning)},
	}
	return vals
}

func provideTestPodValuesMultipleErrors() []map[string]string {
	vals := []map[string]string{
		{"kind": "Pod", "namespace": "abc", "phase": string(corev1.PodPending)},
		{"kind": "Pod", "namespace": "abc", "phase": string(corev1.PodFailed)},
		{"kind": "Pod", "namespace": "def", "phase": string(corev1.PodPending)},
		{"kind": "Pod", "namespace": "def", "phase": string(corev1.PodFailed)},
		{"kind": "Pod", "namespace": "ghi", "phase": string(corev1.PodPending)},
		{"kind": "Pod", "namespace": "ghi", "phase": string(corev1.PodFailed)},
		{"kind": "Pod", "namespace": "hjk", "phase": string(corev1.PodPending)},
		{"kind": "Pod", "namespace": "hjk", "phase": string(corev1.PodFailed)},
		{"kind": "Pod", "namespace": "hjk", "phase": string(corev1.PodPending)},
		{"kind": "Pod", "namespace": "lmn", "phase": string(corev1.PodFailed)},
		{"kind": "Pod", "namespace": "lmn", "phase": string(corev1.PodPending)},
	}
	return vals
}

func TestMetricsCollector_CollectMetrics_Deployments(t *testing.T) {

	values := provideTestDeploymentValues()

	var deployments []appsv1.Deployment
	for _, val := range values {
		deployment := appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: val["namespace"],
				OwnerReferences: []metav1.OwnerReference{
					{
						Kind: val["kind"],
					},
				},
			},
			Status: appsv1.DeploymentStatus{},
		}

		if i, _ := strconv.ParseInt(val["replicas"], 10, 10); i > 0 {
			deployment.Status.ReadyReplicas = int32(i)
		}

		deployments = append(deployments, deployment)
	}

	metrics, stateSet := CollectDeployments(deployments)

	assert.Equal(t, 3, metrics.Items["Deployment-def-Ready"].Value)
	assert.Equal(t, 1, metrics.Items["Deployment-def-Pending"].Value)

	assert.Equal(t, 3, stateSet["Deployment"]["Ready"])
	assert.Equal(t, 1, stateSet["Deployment"]["Pending"])

}

func provideTestDeploymentValues() []map[string]string {
	vals := []map[string]string{
		{"kind": "Deployment", "namespace": "def", "replicas": "1"},
		{"kind": "Deployment", "namespace": "def", "replicas": "1"},
		{"kind": "Deployment", "namespace": "def", "replicas": "1"},
		{"kind": "Deployment", "namespace": "def", "replicas": "0"},
	}
	return vals
}

func TestMetricsCollector_CollectMetrics_CronJobs(t *testing.T) {

	values := provideTestCronJobsValues()

	var cronjobs []batchv1.CronJob
	for _, val := range values {
		cronjob := batchv1.CronJob{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: val["namespace"],
				OwnerReferences: []metav1.OwnerReference{
					{
						Kind: val["kind"],
					},
				},
			},
			Spec: batchv1.CronJobSpec{},
		}

		if val["suspend"] == "true" {
			cronjob.Spec.Suspend = aws.Bool(true)
		}

		cronjobs = append(cronjobs, cronjob)
	}

	metrics, phaseSet := CollectCronJobs(cronjobs)

	assert.Equal(t, 2, metrics.Items["CronJob-def-Running"].Value)
	assert.Equal(t, 3, metrics.Items["CronJob-def-Suspended"].Value)

	assert.Equal(t, 2, phaseSet["CronJob"]["Running"])
	assert.Equal(t, 3, phaseSet["CronJob"]["Suspended"])

}

func provideTestCronJobsValues() []map[string]string {
	vals := []map[string]string{
		{"kind": "CronJob", "namespace": "def", "suspend": "true"},
		{"kind": "CronJob", "namespace": "def", "suspend": "true"},
		{"kind": "CronJob", "namespace": "def", "suspend": "true"},
		{"kind": "CronJob", "namespace": "def", "suspend": "false"},
		{"kind": "CronJob", "namespace": "def", "suspend": "false"},
	}
	return vals
}

func TestMetricsCollector_CollectMetrics_Jobs(t *testing.T) {

	values := provideTestJobsValues()

	var jobs []batchv1.Job
	for _, val := range values {
		job := batchv1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: val["namespace"],
				OwnerReferences: []metav1.OwnerReference{
					{
						Kind: val["kind"],
					},
				},
			},
			Status: batchv1.JobStatus{},
		}
		switch val["status"] {
		case "Active":
			job.Status.Active++
		case "Failed":
			job.Status.Failed++
		case "Succeeded":
			job.Status.Succeeded++
		}

		jobs = append(jobs, job)
	}

	metrics, phaseSet := CollectJobs(jobs)

	assert.Equal(t, 3, metrics.Items["Job-def-Running"].Value)
	assert.Equal(t, 2, metrics.Items["Job-def-Failed"].Value)
	assert.Equal(t, 4, metrics.Items["Job-def-Complete"].Value)

	assert.Equal(t, 3, phaseSet["Job"]["Running"])
	assert.Equal(t, 2, phaseSet["Job"]["Failed"])
	assert.Equal(t, 4, phaseSet["Job"]["Complete"])

}

func provideTestJobsValues() []map[string]string {
	vals := []map[string]string{
		{"kind": "Job", "namespace": "def", "status": "Succeeded"},
		{"kind": "Job", "namespace": "def", "status": "Succeeded"},
		{"kind": "Job", "namespace": "def", "status": "Succeeded"},
		{"kind": "Job", "namespace": "def", "status": "Succeeded"},
		{"kind": "Job", "namespace": "def", "status": "Active"},
		{"kind": "Job", "namespace": "def", "status": "Active"},
		{"kind": "Job", "namespace": "def", "status": "Active"},
		{"kind": "Job", "namespace": "def", "status": "Failed"},
		{"kind": "Job", "namespace": "def", "status": "Failed"},
	}
	return vals
}
