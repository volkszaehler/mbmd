package prometheus

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/volkszaehler/mbmd/meters"
)

var MBMDRegistry = prometheus.NewRegistry()

type Config struct {
	Enable                 bool // defaults to yes
	EnableProcessCollector bool
	EnableGoCollector      bool
}

// RegisterAllMetrics registers all static metrics and dynamically created measurement metrics
// to the Prometheus Default registry.
func RegisterAllMetrics(c Config) {
	if !c.Enable {
		return
	}
	createMeasurementMetrics()
	if c.EnableProcessCollector {
		MBMDRegistry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	}
	if c.EnableGoCollector {
		MBMDRegistry.MustRegister(collectors.NewGoCollector())
	}
}

// createMeasurementMetrics initializes all existing meters.Measurement
//
// If a prometheus.Metric could not be registered, the affected prometheus.Metric will be omitted.
func createMeasurementMetrics() {
	for _, measurement := range meters.MeasurementValues() {
		switch measurement.PrometheusMetricType() {
		case meters.Gauge:
			newGauge := NewMeasurementGaugeCollector(
				newGaugeOpts(
					measurement.PrometheusName(),
					measurement.PrometheusHelpText(),
				),
			)

			if err := MBMDRegistry.Register(newGauge); err != nil {
				log.Printf(
					"Could not register gauge for measurement '%s'. Omitting... (Error: %s)\n",
					measurement,
					err,
				)
			} else {
				gaugeVecMap[measurement] = newGauge
			}
		case meters.Counter:
			measurementCollector := NewMeasurementCounterCollector(
				newCounterOpts(
					measurement.PrometheusName(),
					measurement.PrometheusHelpText(),
				),
			)

			if err := MBMDRegistry.Register(measurementCollector); err != nil {
				log.Printf("could not register counter for measurement '%s'. omitting... (Error: %s)\n",
					measurement,
					err,
				)
			} else {
				counterVecMap[measurement] = measurementCollector
			}
		}
	}
}
