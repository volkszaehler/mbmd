package prometheus

import "github.com/prometheus/client_golang/prometheus"

// collectable is a generic interface that collects all "categorized" Prometheus metrics into a slice.
type collectable interface {
	Collect() []prometheus.Collector
}
