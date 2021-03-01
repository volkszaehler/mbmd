// Package prometheus_metrics
//
// These functions take care of registering and updating counters/gauges.
// Whenever a prometheus.Metric is created, it complies with following convention "automatically".
// For instance: A counter for total connection attempts
// => Name for newCounterOpts: `smart_meter_connection_attempt_total`
// => Metric type: Counter
// => After prometheus.Metric creation: `mbmd_smart_meter_connection_attempt_total`
//
// If a new measurement prometheus.Metric is created, it follows this convention:
// For instance: "L1 Export Power" with unit `W`
// => Name: `l1_export_power`
// => Subsystem (device manufacturer): `sunspec`
// => Metric type: Gauge
// => After prometheus.Metric creation: `mbmd_sunspec_l1_export_power_watts`
// General naming convention: `mbmd_<DEVICE_MANUFACTURER>_<NAME_IN_LOWER_CASE_WITH_UNDERSCORES>_<UNIT>`
// Measurement metrics are treated slightly differently and are maintained in measurement.go
// It ensures that extensibility and customization of Prometheus names and descriptions is easy.
// By default, if no custom Prometheus name is given, the measurement name is transformed to lower case
// and whitespaces are replaced with underscores.
//
// Besides dynamic measurement metrics, some static metrics have been introduced
// in order to keep track of e. g. connection attempts, failed connection attempts, amount of bus scans, ...
//
// If you want to add new metrics, make sure your metric details comply to the usual Prometheus naming conventions
// e. g.: Amount of connection attempts -> Most fitting metric type: Counter -> Name (in newCounterOpts): `smart_meter_connection_attempt_total`
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
// TODO make?
var counterVecMap = map[meters.Measurement]*prometheus.CounterVec{}

// gaugeVecMap contains all meters.Measurement that are associated with a prometheus.Gauge
//
// If a new meters.Measurement is introduced, it needs to be added either to counterVecMap
// or to gaugeVecMap - Otherwise Prometheus won't keep track of the newly added meters.Measurement
// TODO make?
var gaugeVecMap = map[meters.Measurement]*prometheus.GaugeVec{}

// UpdateMeasurementMetric updates a counter or gauge based by passed measurement
//
// Returns false if the associated prometheus.Metric does not exist
func UpdateMeasurementMetric(
	deviceId string,
	deviceSerial string,
	measurement meters.MeasurementResult,
) (ok bool) {
	// TODO Remove when development is finished or think about a solution handling mocked devices
	if deviceSerial == "" {
		deviceSerial = SSN_MISSING
	}

	// 		fmt.Printf("prometheus> [%s] deviceSerial: %s, measurement: %s\n", deviceId, deviceSerial, measurement.Value)
	if gauge, ok := gaugeVecMap[measurement.Measurement]; ok {
		// 	fmt.Printf("prometheus> [%s] Setting gauge value of %s to %s\n", deviceId, gauge.WithLabelValues(deviceId, deviceSerial).Desc(), measurement.Value)
		gauge.WithLabelValues(deviceId, deviceSerial).Set(measurement.Value)
		return ok
	} else if counter, ok := counterVecMap[measurement.Measurement]; ok {
		// 	fmt.Printf("prometheus> [%s] Setting counter value of %s to %s\n", deviceId, counter.WithLabelValues(deviceId, deviceSerial).Desc(), measurement.Value)
		counter.WithLabelValues(deviceId, deviceSerial).Add(measurement.Value)
		return ok
	} else {
		return ok
	}
}

// CreateMeasurementMetrics initializes all existing meters.Measurement for a manufacturer
//
// If a prometheus.Metric could not be registered (see prometheus.Register),
// the affected prometheus.Metric will be omitted.
func CreateMeasurementMetrics(device meters.Device, labels ...string) {
	if !(len(labels) > 0) {
		labels = []string{"device_id", "serial_number"}
	}

	manufacturerName := device.Descriptor().Manufacturer

	for _, measurement := range meters.MeasurementValues() {
		switch measurement.PrometheusMetricType() {
		case meters.Gauge:
			newGauge := prometheus.NewGaugeVec(
				*newGaugeOptsWithSubsystem(
					manufacturerName,
					measurement.PrometheusName(),
					measurement.PrometheusDescription(),
				),
				labels,
			)

			if err := prometheus.Register(newGauge); err != nil {
				log.Fatalf(
					"Could not register gauge for measurement '%s' for devices with manufacturer '%s'. Error: %s\nOmitting...\n",
					measurement,
					manufacturerName,
					err,
				)
			} else {
				gaugeVecMap[measurement] = newGauge
			}
		case meters.Counter:
			newCounter := prometheus.NewCounterVec(
				*newCounterOptsWithSubsystem(
					manufacturerName,
					measurement.PrometheusName(),
					measurement.PrometheusDescription(),
				),
				labels,
			)

			if err := prometheus.Register(newCounter); err != nil {
				log.Fatalf(
					"Could not register counter for measurement '%s' for device '%s'. Error: %s\nOmitting...\n",
					measurement,
					manufacturerName,
					err,
				)
			} else {
				counterVecMap[measurement] = newCounter
			}
		}
	}
}
