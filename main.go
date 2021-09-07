package main

import (
	"context"
	"os"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/skpr/cluster-metrics/internal/metrics"
)

var (
	cliKubeConfig = kingpin.Flag("kubeconfig", "The path to the kube config file.").Envar("KUBECONFIG").String()
	cliFrequency  = kingpin.Flag("frequency", "How often to poll for items data").Envar("CLUSTER_METRICS_FREQUENCY").Default("60s").Duration()
)

func main() {
	kingpin.Parse()

	// use the current context in kubeConfig
	c, err := clientcmd.BuildConfigFromFlags("", *cliKubeConfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(c)
	if err != nil {
		panic(err.Error())
	}

	for range time.Tick(*cliFrequency) {
		// Get the pods
		pods, err := clientset.CoreV1().Pods(corev1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		// Collect the metrics.
		mts := metrics.Collect(pods.Items)

		// Log the metrics.
		err = metrics.Log(os.Stdout, mts)
		if err != nil {
			panic(err.Error())
		}
	}
}
