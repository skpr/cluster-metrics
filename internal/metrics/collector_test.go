package metrics

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
	assert.Equal(t, 1, metrics.Items["Deployment-def-NotReady"].Value)

	assert.Equal(t, 3, stateSet["Deployment"]["Ready"])
	assert.Equal(t, 1, stateSet["Deployment"]["NotReady"])

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

		var suspended *bool
		if val["suspend"] == "true" {
			*suspended = true
		}

		cronjob.Spec.Suspend = suspended

		cronjobs = append(cronjobs, cronjob)
	}

	metrics, phaseSet := CollectCronJobs(cronjobs)

	assert.Equal(t, 2, metrics.Items["CronJob-def-Active"].Value)
	assert.Equal(t, 3, metrics.Items["CronJob-def-Suspended"].Value)

	assert.Equal(t, 2, phaseSet["CronJob"]["Active"])
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

	assert.Equal(t, 3, metrics.Items["Job-def-Active"].Value)
	assert.Equal(t, 2, metrics.Items["Job-def-Failed"].Value)
	assert.Equal(t, 4, metrics.Items["Job-def-Succeeded"].Value)

	assert.Equal(t, 3, phaseSet["Job"]["Active"])
	assert.Equal(t, 2, phaseSet["Job"]["Failed"])
	assert.Equal(t, 4, phaseSet["Job"]["Succeeded"])

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
