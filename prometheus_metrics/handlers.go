package prometheus_metrics

import "github.com/prometheus/client_golang/prometheus"

// handlerCollectors contains all Prometheus metrics about modbus connection handlers
//
// Implements collectable interface
type handlerCollectors struct {}

var handlerCollectorsLabels = []string{"rtu_tcp_addr"}

var (
	// TODO Remove?
	ConnectionHandlerCreated = prometheus.NewCounterVec(
		*newCounterOpts(
			"connection_handler_created_total",
			"// TODO Remove?",
		),
		handlerCollectorsLabels,
	)

	ConnectionHandlerDeviceInitializationRoutineStarted = prometheus.NewCounter(
		*newCounterOpts(
			"connection_handler_device_initialization_routine_starts_total",
			"Total starts of routines where a device is initialized (e. g. initial ModBus connection and retrieving device information)",
		),
	)

	ConnectionHandlerDeviceInitializationFailure = prometheus.NewCounter(
		*newCounterOpts(
			"connection_handler_device_initialization_failures_total",
			"Total failures of routines where a device is initialized (e. g. initial ModBus connection and retrieving device information)",
		),
	)

	ConnectionHandlerDeviceInitializationSuccess = prometheus.NewCounter(
		*newCounterOpts(
			"connection_handler_device_initialization_successes_total",
			"Total successful routines where a device is initialized (e. g. initial ModBus connection and retrieving device information)",
		),
	)

	ConnectionHandlerDeviceQueriesTotal = prometheus.NewCounterVec(
		*newCounterOpts(
			"connection_handler_device_queries_total",
			"Amount of queries/requests done for a smart meter",
		),
		[]string{"serial_number"},
	)

	ConnectionHandlerDeviceQueriesErrorTotal = prometheus.NewCounterVec(
		*newCounterOpts(
			"connection_handler_device_queries_error_total",
			"Errors occurred during smart meter query",
		),
		[]string{"serial_number"},
	)

	ConnectionHandlerDeviceQueriesSuccessTotal = prometheus.NewCounterVec(
		*newCounterOpts(
			"connection_handler_device_queries_success_total",
			"Successful smart meter query",
		),
		[]string{"serial_number"},
	)

	ConnectionHandlerDeviceQueryMeasurementValueSkippedTotal = prometheus.NewCounterVec(
		*newCounterOpts(
			"connection_handler_device_queries_measurement_value_skipped_total",
			"NaN measurement values found and skipped during smart meter query",
		),
		[]string{"serial_number"},
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
