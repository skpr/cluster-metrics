package podsbyPhase

import (
	"github.com/skpr/cluster-metrics/internal/metrics/plugins"
)

const (
	dimensionKind        = "kind"
	dimensionEnvironment = "environment"
	dimensionNamespace   = "namespace"
	dimensionPhase       = "phase"
	dimensionProject     = "project"
	dimensionCluster     = "cluster"
	metricTotal          = "total"
	typeGauge            = "gauge"
	metricName           = "PodStatus"
)

type Client struct {
	plugins.ClusterMetricsPluginInterface
}
