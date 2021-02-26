package prometheus_metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"strings"
)

const NAMESPACE = "mbmd"

// newCounterOpts creates a CounterOpts object, but with a predefined namespace
func newCounterOpts(name string, help string) *prometheus.CounterOpts {
	return &prometheus.CounterOpts{
		Namespace: NAMESPACE,
		Name:      name,
		Help:      help,
	}
}

// newCounterOptsWithSubsystem acts the same as newCounterOpts, but specifies a subsystem for Prometheus fully qualified name
func newCounterOptsWithSubsystem(subsystem string, name string, help string) *prometheus.CounterOpts {
	opts := newCounterOpts(name, help)
	opts.Subsystem = strings.ToLower(subsystem)

	return opts
}

// newGaugeOpts creates a GaugeOpts object, but with a predefined namespace
func newGaugeOpts(name string, help string) *prometheus.GaugeOpts {
	return &prometheus.GaugeOpts{
		Namespace: NAMESPACE,
		Name:      name,
		Help:      help,
	}
}

// newGaugeOptsWithSubsystem acts the same as newGaugeOpts, but specifies a subsystem for Prometheus fully qualified name
func newGaugeOptsWithSubsystem(subsystem string, name string, help string) *prometheus.GaugeOpts {
	opts := newGaugeOpts(name, help)
	opts.Subsystem = strings.ToLower(subsystem)

	return opts
}

