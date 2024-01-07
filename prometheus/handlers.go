package prometheus

import (
	prometheusLib "github.com/prometheus/client_golang/prometheus"
)

// handlerCollectors contains all Prometheus metrics about modbus connection handlers
//
// Implements collectable interface
type handlerCollectors struct{}

var (
	ConnectionHandlerDeviceInitializationRoutineStarted = prometheusLib.NewCounter(
		newCounterOpts(
			"connection_handler_device_initialization_routine_starts_total",
			"Total starts of routines where a device is initialized (e. g. initial ModBus connection and retrieving device information)",
		),
	)

	ConnectionHandlerDeviceInitializationFailure = prometheusLib.NewCounter(
		newCounterOpts(
			"connection_handler_device_initialization_failures_total",
			"Total failures of routines where a device is initialized (e. g. initial ModBus connection and retrieving device information)",
		),
	)

	ConnectionHandlerDeviceInitializationSuccess = prometheusLib.NewCounter(
		newCounterOpts(
			"connection_handler_device_initialization_successes_total",
			"Total successful routines where a device is initialized (e. g. initial ModBus connection and retrieving device information)",
		),
	)

	ConnectionHandlerDeviceQueriesTotal = prometheusLib.NewCounterVec(
		newCounterOpts(
			"connection_handler_device_queries_total",
			"Number of queries to a meter",
		),
		[]string{"device_name", "serial_number"},
	)

	ConnectionHandlerDeviceQueriesErrorTotal = prometheusLib.NewCounterVec(
		newCounterOpts(
			"connection_handler_device_queries_error_total",
			"Errors occurred during smart meter query",
		),
		[]string{"device_name", "serial_number"},
	)

	ConnectionHandlerDeviceQueriesSuccessTotal = prometheusLib.NewCounterVec(
		newCounterOpts(
			"connection_handler_device_queries_success_total",
			"Successful smart meter query",
		),
		[]string{"device_name", "serial_number"},
	)

	ConnectionHandlerDeviceQueryMeasurementValueSkippedTotal = prometheusLib.NewCounterVec(
		newCounterOpts(
			"connection_handler_device_queries_measurement_value_skipped_total",
			"NaN measurement values found and skipped during smart meter query",
		),
		[]string{"device_name", "serial_number"},
	)
)

func (handlerCollectors) Collect() []prometheusLib.Collector {
	return []prometheusLib.Collector{
		ConnectionHandlerDeviceInitializationRoutineStarted,
		ConnectionHandlerDeviceInitializationFailure,
		ConnectionHandlerDeviceInitializationSuccess,

		ConnectionHandlerDeviceQueriesTotal,
		ConnectionHandlerDeviceQueriesErrorTotal,
		ConnectionHandlerDeviceQueriesSuccessTotal,
		ConnectionHandlerDeviceQueryMeasurementValueSkippedTotal,
	}
}
