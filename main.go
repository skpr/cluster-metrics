package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/prozz/aws-embedded-metrics-golang/emf"
	"gopkg.in/alecthomas/kingpin.v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/skpr/cluster-metrics/internal"
)

var (
	kubeConfig string
	frequency  time.Duration
	namespace  string
)

func main() {
	defaultConfig := ""
	if home := homedir.HomeDir(); home != "" {
		defaultConfig = filepath.Join(home, ".kube", "config")
	}
	kingpin.Flag("kubeconfig", "The path to the kube config file.").Default(defaultConfig).StringVar(&kubeConfig)
	kingpin.Flag("frequency", "How often to poll for items data").Default("60s").DurationVar(&frequency)
	kingpin.Flag("namespace", "The metrics namespace").Required().StringVar(&namespace)
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

	collector := internal.NewMetricsCollector(clientset)

	// Format items
	l := emf.New(emf.WithWriter(os.Stderr)).Namespace(namespace)
	logger := internal.NewMetricsLogger(l)

	for range time.Tick(frequency) {
		// Get the pods
		pods, err := collector.ListPods()
		if err != nil {
			panic(err.Error())
		}

		// CollectMetrics the metrics.
		metrics := collector.CollectMetrics(pods)

		// Log the metrics.
		logger.Log(metrics)
	}
}
