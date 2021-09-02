# Cluster Metrics

This is a simple app to periodically query a Kubernetes cluster
to get some metrics, and log them in AWS Embedded Metrics format.

## Usage

```
usage: cluster-metrics [<flags>]

Flags:
  --help                   Show context-sensitive help (also try --help-long and --help-man).
  --kubeconfig=KUBECONFIG  The path to the kube config file.
  --frequency=60s          How often to poll for items data
  --namespace="Skpr/ClusterMetrics"  
                           The metrics namespace
```
