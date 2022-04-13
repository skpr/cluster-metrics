package podsbyPhase

import (
	"context"
	"fmt"

	metricsS "github.com/skpr/cluster-metrics/internal/metrics"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func (c Client) Collect(clientset *kubernetes.Clientset) (*metricsS.MetricSet, metricsS.PhaseSet) {
	podList, err := clientset.CoreV1().Pods(corev1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic("could not list pods in all namespaces:" + err.Error())
	}
	metrics := metricsS.NewMetricSet()
	phaseSet := make(metricsS.PhaseSet)
	for _, pod := range podList.Items {
		for _, ref := range pod.ObjectMeta.OwnerReferences {
			if ref.Kind != "" {
				metrics.IncrementSelect(ref.Kind, pod.ObjectMeta.Namespace, fmt.Sprintf("%s-%s-%s", ref.Kind, dimensionNamespace, pod.Status.Phase), map[string]string{dimensionPhase: fmt.Sprint(pod.Status.Phase)})
				phaseSet[string(pod.Status.Phase)]++
			}
		}
	}
	return metrics, phaseSet
}
