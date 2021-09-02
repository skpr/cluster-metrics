package main

import (
	"os"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/skpr/cluster-metrics/internal/metrics"
)

var (
	kubeConfig string
	frequency  time.Duration
	namespace  string
)

func main() {
	kingpin.Flag("kubeconfig", "The path to the kube config file.").StringVar(&kubeConfig)
	kingpin.Flag("frequency", "How often to poll for items data").Default("60s").DurationVar(&frequency)
	kingpin.Flag("namespace", "The metrics namespace").Default("Skpr/ClusterMetrics").StringVar(&namespace)
	kingpin.Parse()

	// use the current context in kubeConfig
	c, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(c)
	if err != nil {
		panic(err.Error())
	}

	collector := metrics.NewCollector(clientset)

	// Format items
	logger := metrics.NewLogger(os.Stderr, namespace)

	for range time.Tick(frequency) {
		// Get the pods
		pods, err := collector.ListPods()
		if err != nil {
			panic(err.Error())
		}

		// Collect the metrics.
		metrics := collector.Collect(pods)

		// Log the metrics.
		logger.Log(metrics, time.Now().UTC())
	}
}
