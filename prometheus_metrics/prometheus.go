package prometheus_metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/volkszaehler/mbmd/meters"
)

const NAMESPACE = "mbmd"
const SSN_MISSING = "NOT_AVAILABLE"

// TODO remove?
var (
	ConnectionAttemptTotal = prometheus.NewCounterVec(
		newCounterOpts(
			"smart_meter_connection_attempt_total",
			"Total amount of a smart meter's connection attempts",
		),
		[]string{"model", "sub_device"},
	)

	ConnectionAttemptFailedTotal = prometheus.NewCounterVec(
		newCounterOpts(
			"smart_meter_connection_attempt_failed_total",
			"Amount of a smart meter's connection failures",
		),
		[]string{"model", "sub_device"},
	)

	ConnectionPartiallySuccessfulTotal = prometheus.NewCounterVec(
		newCounterOpts(
			"smart_meter_connection_partially_successful_total",
			"Number of connections that are partially open",
		),
		[]string{"model", "sub_device"},
	)

	DevicesCreatedTotal = prometheus.NewCounterVec(
		newCounterOpts(
			"smart_meter_devices_created_total",
			"Number of smart meter devices created/registered",
		),
		[]string{"meter_type", "sub_device"},
	)

	BusScanStartedTotal = prometheus.NewCounterVec(
		newCounterOpts(
			"bus_scan_started_total",
			"Total started bus scans",
		),
		[]string{"device_id"},
	)

	BusScanDeviceInitializationErrorTotal = prometheus.NewCounterVec(
		newCounterOpts(
			"bus_scan_device_initialization_error_total",
			"Total errors upon initialization of a device during bus scan",
		),
		[]string{"device_id"},
	)

	BusScanTotal = prometheus.NewCounter(
		newCounterOpts(
		"bus_scan_total",
		"Amount of bus scans done",
		),
	)

	BusScanDeviceProbeSuccessfulTotal = prometheus.NewCounterVec(
		newCounterOpts(
			"bus_scan_device_probe_successful_total",
			"Amount of successfully found devices during bus scan",
		),
		[]string{"device_id", "serial_number"},
	)

	BusScanDeviceProbeFailedTotal = prometheus.NewCounterVec(
		newCounterOpts(
			"bus_scan_device_probe_failed_total",
			"Amount of devices failed to be found during bus scan",
		),
		[]string{"device_id"},
	)

	MeasurementElectricCurrent = prometheus.NewGaugeVec(
		newGaugeOpts(
			"measurement_electric_current_ampere",
			"Last electric current measured",
		),
		[]string{"device_id", "serial_number"},
	)

	ReadDeviceDetailsFailedTotal = prometheus.NewCounterVec(
		newCounterOpts(
			"smart_meter_read_device_details_failed_total",
			"Reading additional details of a smart meter failed",
		),
		[]string{"model"},
	)

	DeviceQueriesTotal = prometheus.NewCounterVec(
		newCounterOpts(
			"smart_meter_queries_total",
			"Amount of queries/requests done for a smart meter",
		),
		[]string{"device_id", "serial_number"},
	)

	DeviceQueriesErrorTotal = prometheus.NewCounterVec(
		newCounterOpts(
			"smart_meter_queries_error_total",
			"Errors occured during smart meter query",
		),
		[]string{"device_id", "serial_number"},
	)

	DeviceQueriesSuccessTotal = prometheus.NewCounterVec(
		newCounterOpts(
			"smart_meter_queries_success_total",
			"Successful smart meter query",
		),
		[]string{"device_id", "serial_number"},
	)

	DeviceQueryMeasurementValueSkippedTotal = prometheus.NewCounterVec(
		newCounterOpts(
			"smart_meter_queries_measurement_value_skipped_total",
			"NaN measurement values found and skipped during smart meter query",
		),
		[]string{"device_id", "serial_number"},
	)

	MeasurementL1Current = prometheus.NewGaugeVec(
		newGaugeOpts(
			"measurement_l1_current_ampere",
			"Measurement of L1 current in ampere",
		),
		[]string{"device_id", "serial_number"},
	)

	MeasurementL2Current = prometheus.NewGaugeVec(
		newGaugeOpts(
			"measurement_l2_current_ampere",
			"Measurement of L2 current in ampere",
		),
		[]string{"device_id", "serial_number"},
	)

	MeasurementL3Current = prometheus.NewGaugeVec(
		newGaugeOpts(
			"measurement_l3_current_ampere",
			"Measurement of L3 current in ampere",
		),
		[]string{"device_id", "serial_number"},
	)

	MeasurementFrequency = prometheus.NewGaugeVec(
		newGaugeOpts(
			"measurement_frequency_hertz",
			"Last measurement of frequency in Hz",
		),
		[]string{"device_id", "serial_number"},
	)

	MeasurementVoltage = prometheus.NewGaugeVec(
		newGaugeOpts(
			"measurement_voltage_volt",
			"Last measurement of voltage in V",
		),
		[]string{"device_id", "serial_number"},
	)
)

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

// Init registers all globally defined metrics to Prometheus library's default registry
func Init() {
	collectors := make([]prometheus.Collector, 0, len(meters.MeasurementValues()))

	for _, measurement := range meters.MeasurementValues() {
		switch measurement.PrometheusMetricType() {
		case meters.Gauge:
			newGauge := prometheus.NewGaugeVec(
				newGaugeOpts(
					measurement.PrometheusName(),
					measurement.PrometheusDescription(),
				),
				[]string{"device_id", "serial_number"},
			)
			gaugeVecMap[measurement] = newGauge
			collectors = append(collectors, newGauge)
		case meters.Counter:
			newCounter := prometheus.NewCounterVec(
				newCounterOpts(
					measurement.PrometheusName(),
					measurement.PrometheusDescription(),
				),
				[]string{"device_id", "serial_number"},
			)
			counterVecMap[measurement] = newCounter
			collectors = append(collectors, newCounter)
		}
	}

	prometheus.MustRegister(collectors...)
}

// UpdateMeasurementMetric updates a counter or gauge based by passed measurement
func UpdateMeasurementMetric(
	deviceId 	 string,
	deviceSerial string,
	measurement  meters.MeasurementResult,
) {
	// TODO Remove when development is finished or think about a solution handling mocked devices
	if deviceSerial == "" {
		deviceSerial = SSN_MISSING
	}

	fmt.Printf("prometheus> [%s] deviceSerial: %s, measurement: %s\n", deviceId, deviceSerial, measurement.Value)
	if gauge, ok := gaugeVecMap[measurement.Measurement]; ok {
		fmt.Printf("prometheus> [%s] Setting gauge value of %s to %s\n", deviceId, gauge.WithLabelValues(deviceId, deviceSerial).Desc(), measurement.Value)
		gauge.WithLabelValues(deviceId, deviceSerial).Set(measurement.Value)
	} else if counter, ok := counterVecMap[measurement.Measurement]; ok {
		fmt.Printf("prometheus> [%s] Setting counter value of %s to %s\n", deviceId, counter.WithLabelValues(deviceId, deviceSerial).Desc(), measurement.Value)
		counter.WithLabelValues(deviceId, deviceSerial).Add(measurement.Value)
	}
}

// newCounterOpts creates a CounterOpts object, but with a predefined namespace
func newCounterOpts(name string, help string) prometheus.CounterOpts {
	return prometheus.CounterOpts{
		Namespace: NAMESPACE,
		Name: name,
		Help: help,
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
