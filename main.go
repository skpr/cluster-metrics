package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/google/uuid"
	"gopkg.in/alecthomas/kingpin.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/skpr/cluster-metrics/internal/metrics"
)

var (
	cliKubeConfig = kingpin.Flag("kubeconfig", "The path to the kube config file.").Envar("KUBECONFIG").String()
	cliFrequency  = kingpin.Flag("frequency", "How often to poll for items data").Default("60s").Duration()
	cliNamespace  = kingpin.Flag("namespace", "The metrics namespace").Default("Skpr/Cluster").String()
	cliLogGroup   = kingpin.Flag("log-group", "The log group to use").Envar("CLUSTER_METRICS_LOG_GROUP").String()
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

	cwlogsClient := cloudwatchlogs.NewFromConfig(cfg)

	pusher := metrics.NewLogsPusher(cwlogsClient)
	ctx := context.TODO()
	err = pusher.CreateLogGroup(ctx, *cliLogGroup)
	if err != nil {
		panic(fmt.Sprintf("failed to create log group %s client: %d", *cliLogGroup, err))
	}

	for range time.Tick(*cliFrequency) {
		// Get the pods
		pods, err := clientset.CoreV1().Pods(corev1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		// Collect the metrics.
		mts := metrics.Collect(pods.Items)

		// Convert to log events.
		evts := metrics.ConvertToEvents(time.Now().UTC(), *cliNamespace, mts)

		// Generate a unique log stream name.
		logStream := uuid.NewString()

		// Create the log stream.
		err = pusher.CreateLogStream(ctx, *cliLogGroup, logStream)
		if err != nil {
			panic(err.Error())
		}

		// Push the logs.
		err = pusher.PushLogs(ctx, *cliLogGroup, logStream, evts, nil)
		if err != nil {
			panic(err.Error())
		}

	}
}
