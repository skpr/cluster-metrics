package main

import (
	"context"
	"fmt"
	"github.com/skpr/cluster-metrics/internal/metrics/plugins"
	pods2 "github.com/skpr/cluster-metrics/internal/metrics/plugins/podsbyPhase"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"gopkg.in/alecthomas/kingpin.v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/skpr/cluster-metrics/internal/metrics"
)

var (
	cliKubeConfig  = kingpin.Flag("kubeconfig", "The path to the kube config file.").Envar("KUBECONFIG").String()
	cliFrequency   = kingpin.Flag("frequency", "How often to poll for items data").Envar("CLUSTER_METRICS_FREQUENCY").Default("60s").Duration()
	cliClusterName = kingpin.Flag("cluster", "The cluster name").Envar("CLUSTER_NAME").String()
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

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Sprintf("failed to setup client: %d", err))
	}

	cloudwatchClient := cloudwatch.NewFromConfig(cfg)
	pusher := metrics.NewPusher(cloudwatchClient)

	for range time.Tick(*cliFrequency) {

		// PodsByPhase Client
		var (
			podsByPhase pods2.Client
		)

		// Consolidate Clients so we can loop over them
		clients := []plugins.ClusterMetricsPluginInterface{
			podsByPhase,
		}

		for _, client := range clients {

			// Collect the metrics.
			mts, _ := client.Collect(clientset)

			// Log the detailed metrics.
			err = client.Log(os.Stdout, mts)
			if err != nil {
				panic(err.Error())
			}

			// Convert to metric data.
			metricData := client.Convert(time.Now().UTC(), *cliClusterName, map[string]interface{}{})

			// Push the phase metrics.
			err = pusher.Push(context.TODO(), "Skpr/Cluster", metricData)
			if err != nil {
				panic(err.Error())
			}
		}

	}
}
