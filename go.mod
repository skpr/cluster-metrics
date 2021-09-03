module github.com/skpr/cluster-metrics

go 1.16

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20210208195552-ff826a37aa15 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.8.0
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.7.0
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/prozz/aws-embedded-metrics-golang v1.0.0
	github.com/stretchr/testify v1.7.0
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	k8s.io/api v0.22.1
	k8s.io/apimachinery v0.22.1
	k8s.io/client-go v0.22.1
)
