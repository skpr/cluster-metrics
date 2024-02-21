package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"gopkg.in/alecthomas/kingpin.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

		mts, phases := &metrics.MetricSet{}, metrics.StateSet{}

		{
			// Get the pods
			pods, err := clientset.CoreV1().Pods(corev1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				panic(err.Error())
			}

			// Collect the metrics.
			metricSetAddition, stateSetAddition := metrics.CollectPods(pods.Items)
			mts = metrics.CombineRecords(mts, metricSetAddition)
			phases = metrics.CombineStates(&phases, &stateSetAddition)
		}

		// Get Namespaces so we can loop through them all.
		namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		for _, namespace := range namespaces.Items {
			{
				// Get the deployments
				deployments, err := clientset.AppsV1().Deployments(namespace.Name).List(context.TODO(), metav1.ListOptions{})
				if err != nil {
					panic(err.Error())
				}

				// Collect the metrics.
				metricSetAddition, stateSetAddition := metrics.CollectDeployments(deployments.Items)
				mts = metrics.CombineRecords(mts, metricSetAddition)
				phases = metrics.CombineStates(&phases, &stateSetAddition)
			}
			{
				// Get the statefulSets
				statefulsets, err := clientset.AppsV1().StatefulSets(namespace.Name).List(context.TODO(), metav1.ListOptions{})
				if err != nil {
					panic(err.Error())
				}

				// Collect the metrics.
				metricSetAddition, stateSetAddition := metrics.CollectStatefulSets(statefulsets.Items)
				mts = metrics.CombineRecords(mts, metricSetAddition)
				phases = metrics.CombineStates(&phases, &stateSetAddition)
			}
			{
				// Get the replicasets
				replicasets, err := clientset.AppsV1().ReplicaSets(namespace.Name).List(context.TODO(), metav1.ListOptions{})
				if err != nil {
					panic(err.Error())
				}

				// Collect the metrics.
				metricSetAddition, stateSetAddition := metrics.CollectReplicaSets(replicasets.Items)
				mts = metrics.CombineRecords(mts, metricSetAddition)
				phases = metrics.CombineStates(&phases, &stateSetAddition)
			}
			{
				// Get the daemonsets
				daemonsets, err := clientset.AppsV1().DaemonSets(namespace.Name).List(context.TODO(), metav1.ListOptions{})
				if err != nil {
					panic(err.Error())
				}

				// Collect the metrics.
				metricSetAddition, stateSetAddition := metrics.CollectDaemonSets(daemonsets.Items)
				mts = metrics.CombineRecords(mts, metricSetAddition)
				phases = metrics.CombineStates(&phases, &stateSetAddition)
			}
			{
				// Get the cronjobs
				cronjobs, err := clientset.BatchV1().CronJobs(namespace.Name).List(context.TODO(), metav1.ListOptions{})
				if err != nil {
					panic(err.Error())
				}

				// Collect the metrics.
				metricSetAddition, stateSetAddition := metrics.CollectCronJobs(cronjobs.Items)
				mts = metrics.CombineRecords(mts, metricSetAddition)
				phases = metrics.CombineStates(&phases, &stateSetAddition)
			}
			{
				// Get the jobs
				jobs, err := clientset.BatchV1().Jobs(namespace.Name).List(context.TODO(), metav1.ListOptions{})
				if err != nil {
					panic(err.Error())
				}

				// Collect the metrics.
				metricSetAddition, stateSetAddition := metrics.CollectJobs(jobs.Items)
				mts = metrics.CombineRecords(mts, metricSetAddition)
				phases = metrics.CombineStates(&phases, &stateSetAddition)
			}
		}

		// Log the detailed metrics.
		err = metrics.Log(os.Stdout, mts)
		if err != nil {
			panic(err.Error())
		}

		// Convert to metric data.
		metricData := metrics.ConvertToMetricData(time.Now().UTC(), *cliClusterName, phases)

		// Push the phase metrics.
		err = pusher.Push(context.TODO(), "Skpr/Cluster", metricData)
		if err != nil {
			panic(err.Error())
		}

	}
}
