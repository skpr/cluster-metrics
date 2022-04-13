package plugins

import (
	"context"
	"io"
	"k8s.io/client-go/kubernetes"
	"time"

	awstypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/skpr/cluster-metrics/internal/metrics"
)

// ClusterMetricsPluginInterface defines an interface for cluster metrics to
// plug into. Given a plugin could include extremely specific to the given
// plugin, the signatures are as generic as possible and the functionality
// as high-level as possible.
type ClusterMetricsPluginInterface interface {
	Collect(*kubernetes.Clientset) (*metrics.MetricSet, metrics.PhaseSet)
	Convert(time.Time, string, map[string]interface{}) []awstypes.MetricDatum
	Log(io.Writer, *metrics.MetricSet) error
	Push(context.Context, *string, []awstypes.MetricDatum)
}
