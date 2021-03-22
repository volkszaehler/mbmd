// Package prometheus_metrics
//
// These functions take care of registering and updating counters/gauges.
// General naming convention: `mbmd_<NAME_IN_LOWER_CASE_WITH_UNDERSCORES>[_<UNIT>]['_total']`

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
// => Metric type: Gauge
// => Check if unit `W` can be converted to its elementary unit (see units.ConvertValueToElementaryUnit)
// => After prometheus.Metric creation:
//	- Name: `mbmd_measurement_l1_export_power_watts`
//	- Labels: {"device_name", "serial_number"}
//
// Measurement metrics are treated slightly differently and are maintained in prometheus_metrics/measurement.go
// It ensures that extensibility and customization of Prometheus names and help texts is easy.
// By default, if no custom Prometheus name is given, the measurement's description is transformed to lower case
// and whitespaces are replaced with underscores.
// Afterwards, the measurement's elementary units.Unit and `total` (if meters.MetricType equals meters.Counter) are appended
// in a snake-case fashion.
//
// Besides dynamic measurement metrics, some static metrics have been introduced
// and can be found in the file respectively, for instance: metrics for devices -> devices.go
//
// If you want to add new metrics, make sure your metric details comply to the usual Prometheus naming conventions
// e. g.: Amount of connection attempts
//	-> Most fitting metric type: Counter
//	-> Name (in newCounterOpts): `smart_meter_connection_attempt_total`
// For more information regarding naming conventions and best practices, see https://prometheus.io/docs/practices/naming/
package prometheus_metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/volkszaehler/mbmd/meters"
	"github.com/volkszaehler/mbmd/meters/units"
	"log"
)

// SSN_MISSING is used for mocked smart meters
const SSN_MISSING = "NOT_AVAILABLE"

// measurementMetricsLabels are the Prometheus labels commonly used for meters.Measurement metrics
var measurementMetricsLabels = []string{"device_name", "serial_number"}

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

// UpdateMeasurementMetric updates a counter or gauge based on passed measurement
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

	_, elementaryValue := units.ConvertValueToElementaryUnit(*measurement.Unit(), measurement.Value)

	if gauge, ok := gaugeVecMap[measurement.Measurement]; ok {
		gauge.WithLabelValues(deviceName, deviceSerial).Set(elementaryValue)
		return ok
	} else if counter, ok := counterVecMap[measurement.Measurement]; ok {
		counter.WithLabelValues(deviceName, deviceSerial).Add(elementaryValue)
		return ok
	} else {
		return ok
	}
}

// CreateMeasurementMetrics initializes all existing meters.Measurement
//
// If a prometheus.Metric could not be registered (see prometheus.Register),
// the affected prometheus.Metric will be omitted.
func CreateMeasurementMetrics() {
	for _, measurement := range meters.MeasurementValues() {
		fmt.Printf("%s - %s\n", measurement.PrometheusName(), measurement.PrometheusHelpText())
		switch measurement.PrometheusMetricType() {
		case meters.Gauge:
			newGauge := prometheus.NewGaugeVec(
				*newGaugeOpts(
					measurement.PrometheusName(),
					measurement.PrometheusHelpText(),
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
					measurement.PrometheusHelpText(),
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
