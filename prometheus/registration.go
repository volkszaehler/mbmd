package prometheus

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/volkszaehler/mbmd/meters"
)

// RegisterAllMetrics registers all static metrics and dynamically created measurement metrics
// to the Prometheus Default registry.
func RegisterAllMetrics() {
	registerStatics()
	createMeasurementMetrics()
}

// registerStatics registers all globally defined static metrics to Prometheus library's default registry
//
// If you add a new collectable, make sure to add it to getAllCollectors to have them registered on startup.
func registerStatics() {
	collectables := getAllCollectors()

	for _, collectable := range collectables {
		for _, prometheusCollector := range collectable.Collect() {
			if err := prometheus.Register(prometheusCollector); err != nil {
				log.Printf("Could not register a metric '%s' (Error: %s)", prometheusCollector, err)
			}
		}
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

			if err := prometheus.Register(newGauge); err != nil {
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

			if err := prometheus.Register(measurementCollector); err != nil {
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

func getAllCollectors() []collectable {
	return []collectable{
		socketCollectors{},
		publisherCollectors{},
		handlerCollectors{},
		deviceCollectors{},
	}
}
