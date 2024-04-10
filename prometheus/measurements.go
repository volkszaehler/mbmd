package prometheus

import (
	"github.com/volkszaehler/mbmd/meters"
)

// These functions take care of registering and updating counters/gauges.
// General naming convention: `mbmd_<NAME_IN_LOWER_CASE_WITH_UNDERSCORES>['_total']`
//
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
// => After prometheus.Metric creation:
//	- Name: `mbmd_measurement_l1_export_power`
//	- Labels: {"device_name", "serial_number", "unit"}
//
// Measurement metrics are treated slightly differently and are maintained in
// prometheus/measurement.go. It ensures that extensibility and customization of
// Prometheus names and help texts is easy. By default, if no custom Prometheus
// name is given, the measurement's description is transformed to lower case and
// whitespaces are replaced with underscores. Afterward, the measurement's
// elementary units.Unit and `total` (if meters.MetricType equals
// meters.Counter) are appended in a snake-case fashion.
//
// Besides dynamic measurement metrics, some static metrics have been introduced
// and can be found in the file respectively, for instance: metrics for devices
// -> devices.go
//
// If you want to add new metrics, make sure your metric details comply to the
// usual Prometheus naming conventions e.g.: Amount of connection attempts
//	-> Most fitting metric type: Counter
//	-> Name (in newCounterOpts): `smart_meter_connection_attempt_total`
// For more information regarding naming conventions and best practices, see
// https://prometheus.io/docs/practices/naming/

// SSN_MISSING is used for mocked smart meters
const SSN_MISSING = "NOT_AVAILABLE"

// counterVecMap contains all meters.Measurement that are associated with a
// prometheus.Counter
//
// If a new meters.Measurement is introduced, it needs to be added either to
// counterVecMap or to gaugeVecMap - Otherwise Prometheus won't keep track of
// the newly added meters.Measurement
var counterVecMap = map[meters.Measurement]*MeasurementCounterCollector{}

// gaugeVecMap contains all meters.Measurement that are associated with a
// prometheus.Gauge
//
// If a new meters.Measurement is introduced, it needs to be added either to
// counterVecMap or to gaugeVecMap - Otherwise Prometheus won't keep track of
// the newly added meters.Measurement
var gaugeVecMap = map[meters.Measurement]*MeasurementGaugeCollector{}

// UpdateMeasurementMetric updates a counter or gauge based on passed measurement.
func UpdateMeasurementMetric(
	deviceName string,
	deviceSerial string,
	measurement meters.MeasurementResult,
) {
	// Handle empty device serial numbers (e.g. on mocks)
	if deviceSerial == "" {
		deviceSerial = SSN_MISSING
	}

	if gauge, ok := gaugeVecMap[measurement.Measurement]; ok {
		gauge.Set(measurement.Timestamp, measurement.Value, deviceName, deviceSerial, measurement.Unit().Abbreviation())
		return
	}
	if counter, ok := counterVecMap[measurement.Measurement]; ok {
		counter.Set(measurement.Timestamp, measurement.Value, deviceName, deviceSerial, measurement.Unit().Abbreviation())
	}
}
