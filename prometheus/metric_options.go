package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

const NAMESPACE = "mbmd"

// newCounterOpts creates a CounterOpts object, but with a predefined namespace.
func newCounterOpts(name string, help string) prometheus.CounterOpts {
	return prometheus.CounterOpts{
		Namespace: NAMESPACE,
		Name:      name,
		Help:      help,
	}
}

// newGaugeOpts creates a GaugeOpts object, but with a predefined namespace
func newGaugeOpts(name string, help string) prometheus.GaugeOpts {
	return prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Name:      name,
		Help:      help,
	}
}
