package prometheus_metrics

import "github.com/prometheus/client_golang/prometheus"

// handlerCollectors contains all Prometheus metrics about modbus connection handlers
//
// Implements collectable interface
type handlerCollectors struct {}

var handlerCollectorsLabels = []string{"rtu_tcp_addr"}

var (
	ConnectionHandlerCreated = prometheus.NewCounterVec(
		*newCounterOpts(
			"connection_handler_created_total",
			"",
		),
		handlerCollectorsLabels,
	)

	ConnectionHandlerDeviceInitializationRoutineStarted = prometheus.NewCounter(
		*newCounterOpts(
			"connection_handler_device_initialization_routine_starts_total",
			"",
		),
	)

	ConnectionHandlerDeviceInitializationFailure = prometheus.NewCounter(
		*newCounterOpts(
			"connection_handler_device_initialization_failures_total",
			"",
		),
	)

	ConnectionHandlerDeviceInitializationSuccess = prometheus.NewCounter(
		*newCounterOpts(
			"connection_handler_device_initialization_successes_total",
			"",
		),
	)

	ConnectionHandlerDeviceQueriesTotal = prometheus.NewCounterVec(
		*newCounterOpts(
			"connection_handler_device_queries_total",
			"Amount of queries/requests done for a smart meter",
		),
		[]string{"device_id", "serial_number"},
	)

	ConnectionHandlerDeviceQueriesErrorTotal = prometheus.NewCounterVec(
		*newCounterOpts(
			"connection_handler_device_queries_error_total",
			"Errors occurred during smart meter query",
		),
		[]string{"device_id", "serial_number"},
	)

	ConnectionHandlerDeviceQueriesSuccessTotal = prometheus.NewCounterVec(
		*newCounterOpts(
			"connection_handler_device_queries_success_total",
			"Successful smart meter query",
		),
		[]string{"device_id", "serial_number"},
	)

	ConnectionHandlerDeviceQueryMeasurementValueSkippedTotal = prometheus.NewCounterVec(
		*newCounterOpts(
			"connection_handler_device_queries_measurement_value_skipped_total",
			"NaN measurement values found and skipped during smart meter query",
		),
		[]string{"device_id", "serial_number"},
	)
)

func (h handlerCollectors) Collect() []prometheus.Collector {
	return []prometheus.Collector{
		ConnectionHandlerCreated,

		ConnectionHandlerDeviceInitializationRoutineStarted,
		ConnectionHandlerDeviceInitializationFailure,
		ConnectionHandlerDeviceInitializationSuccess,

		ConnectionHandlerDeviceQueriesTotal,
		ConnectionHandlerDeviceQueriesErrorTotal,
		ConnectionHandlerDeviceQueriesSuccessTotal,
		ConnectionHandlerDeviceQueryMeasurementValueSkippedTotal,
	}
}
