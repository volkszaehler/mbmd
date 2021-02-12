package server

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Prometheusd struct {
	Gauges	map[string]map[string]prometheus.Gauge
}

func NewPrometheusd() *Prometheusd {
	gauges := make(map[string]map[string]prometheus.Gauge)

	return &Prometheusd{Gauges: gauges}
}

// Run Prometheusd metric collection
func (p *Prometheusd) Run(in <-chan QuerySnip) {
	for snip := range in {
		gauge := p.GetOrCreate(snip)
		(*gauge).Set(snip.Value)
	}
}

// GetOrCreate returns a pointer to desired prometheus.Gauge
// If it doesn't exist, GetOrCreate will create a new prometheus.Gauge
// and add it to the Gauges list
func (p *Prometheusd) GetOrCreate(q QuerySnip) *prometheus.Gauge {
	gauges := p.Gauges[q.Device]

	if gauges == nil {
		gauges = make(map[string]prometheus.Gauge)
	}

	gauge := gauges[q.PrometheusName()]

	if gauge == nil {
		gauge = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Subsystem: q.Device,
				Name: q.PrometheusName(),
				Help: q.Description(),
			},
		)

		// Register gauge to Prometheus Registry
		prometheus.MustRegister(gauge)

		gauges[q.PrometheusName()] = gauge
	}

	return &gauge
}
