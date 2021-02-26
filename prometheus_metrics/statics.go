package prometheus_metrics

import "github.com/prometheus/client_golang/prometheus"

// Static metrics
var (
	ConnectionAttemptTotal = prometheus.NewCounterVec(
		*newCounterOpts(
			"smart_meter_connection_attempt_total",
			"Total amount of a smart meter's connection attempts",
		),
		[]string{"model", "sub_device"},
	)

	ConnectionAttemptFailedTotal = prometheus.NewCounterVec(
		*newCounterOpts(
			"smart_meter_connection_attempt_failed_total",
			"Amount of a smart meter's connection failures",
		),
		[]string{"model", "sub_device"},
	)

	ConnectionPartiallySuccessfulTotal = prometheus.NewCounterVec(
		*newCounterOpts(
			"smart_meter_connection_partially_successful_total",
			"Number of connections that are partially open",
		),
		[]string{"model", "sub_device"},
	)

	DevicesCreatedTotal = prometheus.NewCounterVec(
		*newCounterOpts(
			"smart_meter_devices_created_total",
			"Number of smart meter devices created/registered",
		),
		[]string{"meter_type", "sub_device"},
	)

	BusScanStartedTotal = prometheus.NewCounterVec(
		*newCounterOpts(
			"bus_scan_started_total",
			"Total started bus scans",
		),
		[]string{"device_id"},
	)

	BusScanDeviceInitializationErrorTotal = prometheus.NewCounterVec(
		*newCounterOpts(
			"bus_scan_device_initialization_error_total",
			"Total errors upon initialization of a device during bus scan",
		),
		[]string{"device_id"},
	)

	BusScanTotal = prometheus.NewCounter(
		*newCounterOpts(
			"bus_scan_total",
			"Amount of bus scans done",
		),
	)

	BusScanDeviceProbeSuccessfulTotal = prometheus.NewCounterVec(
		*newCounterOpts(
			"bus_scan_device_probe_successful_total",
			"Amount of successfully found devices during bus scan",
		),
		[]string{"device_id", "serial_number"},
	)

	BusScanDeviceProbeFailedTotal = prometheus.NewCounterVec(
		*newCounterOpts(
			"bus_scan_device_probe_failed_total",
			"Amount of devices failed to be found during bus scan",
		),
		[]string{"device_id"},
	)

	ReadDeviceDetailsFailedTotal = prometheus.NewCounterVec(
		*newCounterOpts(
			"smart_meter_read_device_details_failed_total",
			"Reading additional details of a smart meter failed",
		),
		[]string{"model"},
	)

	DeviceQueriesTotal = prometheus.NewCounterVec(
		*newCounterOpts(
			"smart_meter_queries_total",
			"Amount of queries/requests done for a smart meter",
		),
		[]string{"device_id", "serial_number"},
	)

	DeviceQueriesErrorTotal = prometheus.NewCounterVec(
		*newCounterOpts(
			"smart_meter_queries_error_total",
			"Errors occured during smart meter query",
		),
		[]string{"device_id", "serial_number"},
	)

	DeviceQueriesSuccessTotal = prometheus.NewCounterVec(
		*newCounterOpts(
			"smart_meter_queries_success_total",
			"Successful smart meter query",
		),
		[]string{"device_id", "serial_number"},
	)

	DeviceQueryMeasurementValueSkippedTotal = prometheus.NewCounterVec(
		*newCounterOpts(
			"smart_meter_queries_measurement_value_skipped_total",
			"NaN measurement values found and skipped during smart meter query",
		),
		[]string{"device_id", "serial_number"},
	)
)

// RegisterStatics registers all globally defined static metrics to Prometheus library's default registry
func RegisterStatics() {
	prometheus.MustRegister(
		ConnectionAttemptTotal,
		ConnectionAttemptFailedTotal,
		ConnectionPartiallySuccessfulTotal,
		DevicesCreatedTotal,
		BusScanStartedTotal,
		BusScanDeviceInitializationErrorTotal,
		BusScanTotal,
		BusScanDeviceProbeSuccessfulTotal,
		BusScanDeviceProbeFailedTotal,
		ReadDeviceDetailsFailedTotal,
		DeviceQueriesTotal,
		DeviceQueriesErrorTotal,
		DeviceQueriesSuccessTotal,
		DeviceQueryMeasurementValueSkippedTotal,
	)
}
