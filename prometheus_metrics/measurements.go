// Package prometheus_metrics
//
// These functions take care of registering and updating counters/gauges.
// For instance: A counter for total connection attempts
// => Name for newCounterOpts: `smart_meter_connection_attempt_total`
// => Metric type: Counter
// => After prometheus.Metric creation:
//	- Name: `mbmd_smart_meter_connection_attempt_total`
//	- Metric: prometheus.Counter
//
// If a new measurement prometheus.Metric is created, it follows this convention:
// For instance: "L1 Export Power" with unit `W`
// => Name: `l1_export_power`
// => Subsystem (device manufacturer): `sunspec`
// => Metric type: Gauge
// => After prometheus.Metric creation:
//	- Name: `mbmd_l1_export_power_watts`
//	- Labels: {"serial_number"}
//
// General naming convention: `mbmd_<NAME_IN_LOWER_CASE_WITH_UNDERSCORES>_<UNIT>['_total']`
// Measurement metrics are treated slightly differently and are maintained in measurement.go
// It ensures that extensibility and customization of Prometheus names and descriptions is easy.
// By default, if no custom Prometheus name is given, the measurement name is transformed to lower case
// and whitespaces are replaced with underscores.
//
// Besides dynamic measurement metrics, some static metrics have been introduced
// and can be found in the file respectively, for instance: metrics for devices -> devices.go
//
// If you want to add new metrics, make sure your metric details comply to the usual Prometheus naming conventions
// e. g.: Amount of connection attempts
//	-> Most fitting metric type: Counter -> Name (in newCounterOpts): `smart_meter_connection_attempt_total`
// For more information regarding naming conventions and best practices, see https://prometheus.io/docs/practices/naming/
package prometheus_metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/volkszaehler/mbmd/meters"
	"log"
)

const SSN_MISSING = "NOT_AVAILABLE"

// counterVecMap contains all meters.Measurement that are associated with a prometheus.Counter
//
// If a new meters.Measurement is introduced, it needs to be added either to counterVecMap
// or to gaugeVecMap - Otherwise Prometheus won't keep track of the newly added meters.Measurement
var counterVecMap = map[meters.Measurement]*prometheus.CounterVec{}

// gaugeVecMap contains all meters.Measurement that are associated with a prometheus.Gauge
//
// If a new meters.Measurement is introduced, it needs to be added either to counterVecMap
// or to gaugeVecMap - Otherwise Prometheus won't keep track of the newly added meters.Measurement
var gaugeVecMap = map[meters.Measurement]*prometheus.GaugeVec{}

// UpdateMeasurementMetric updates a counter or gauge based by passed measurement
//
// Returns false if the associated prometheus.Metric does not exist
func UpdateMeasurementMetric(
	deviceName string,
	deviceSerial string,
	measurement meters.MeasurementResult,
) (ok bool) {
	// Handle empty device serial numbers (e. g. on mocks)
	// TODO Better handling??
	if deviceSerial == "" {
		deviceSerial = SSN_MISSING
	}

	if gauge, ok := gaugeVecMap[measurement.Measurement]; ok {
		gauge.WithLabelValues(deviceName, deviceSerial).Set(measurement.Value)
		return ok
	} else if counter, ok := counterVecMap[measurement.Measurement]; ok {
		if unit := measurement.Unit(); unit != nil && *unit == meters.KiloWattHour {
			counter.WithLabelValues(deviceName, deviceSerial).Add(measurement.ConvertValueTo(meters.Joule))
		} else {
			counter.WithLabelValues(deviceName, deviceSerial).Add(measurement.Value)
		}
		return ok
	} else {
		return ok
	}
}

var measurementMetricsLabels = []string{"device_name", "serial_number"}

// CreateMeasurementMetrics initializes all existing meters.Measurement
//
// If a prometheus.Metric could not be registered (see prometheus.Register),
// the affected prometheus.Metric will be omitted.
func CreateMeasurementMetrics() {
	for _, measurement := range meters.MeasurementValues() {
		switch measurement.PrometheusMetricType() {
		case meters.Gauge:
			newGauge := prometheus.NewGaugeVec(
				*newGaugeOpts(
					measurement.PrometheusName(),
					measurement.PrometheusDescription(),
				),
				measurementMetricsLabels,
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
			newCounter := prometheus.NewCounterVec(
				*newCounterOpts(
					measurement.PrometheusName(),
					measurement.PrometheusDescription(),
				),
				measurementMetricsLabels,
			)

			if err := prometheus.Register(newCounter); err != nil {
				log.Printf(
					"Could not register counter for measurement '%s'. Omitting... (Error: %s)\n",
					measurement,
					err,
				)
			} else {
				counterVecMap[measurement] = newCounter
			}
		}
	}
}
