# Cluster Metrics

This is a simple app to periodically query a Kubernetes cluster
to get some metrics. It pushes high level pod status metrics to CloudWatch, and logs more detailed metrics to stdout.

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

## Building

```
goreleaser build --snapshot --rm-dist
```
