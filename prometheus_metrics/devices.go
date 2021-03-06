package prometheus_metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/volkszaehler/mbmd/meters"
	"strconv"
)

type deviceCollectors struct {}

var (
	DevicesCreatedTotal = prometheus.NewCounterVec(
		*newCounterOpts(
			"devices_created_total",
			"Number of smart meter devices created/registered",
		),
		[]string{"meter_type"},
	)

	CurrentDevicesActive = prometheus.NewGaugeVec(
		*newGaugeOpts(
			"devices_currently_active",
			"",
		),
		[]string{"meter_type", "sub_device"},
	)

	DeviceModbusConnectionAttemptTotal = prometheus.NewCounterVec(
		*newCounterOpts(
			"device_connection_attempts_total",
			"Total amount of a smart meter's connection attempts",
		),
		[]string{"model_type", "sub_device"},
	)

	DeviceModbusConnectionFailure = prometheus.NewCounterVec(
		*newCounterOpts(
			"device_connection_failures_total",
			"Amount of a smart meter's connection failures",
		),
		[]string{"model_type", "sub_device"},
	)

	DeviceModbusConnectionSuccess = prometheus.NewCounterVec(
		*newCounterOpts(
			"device_connection_successes_total",
			"Amount of a smart meter's successful connection ",
		),
		[]string{"model_type", "sub_device"},
	)

	DeviceModbusConnectionPartialSuccess = prometheus.NewCounterVec(
		*newCounterOpts(
			"device_connection_partial_successes_total",
			"Amount of a smart meter's partial opened connection",
		),
		[]string{"model_type", "sub_device"},
	)

	DeviceByConfigNotFound = prometheus.NewCounterVec(
		*newCounterOpts(
			"device_not_found_by_config_total",
			"Amount of devices defined by config yaml not found",
		),
		[]string{"model_type", "sub_device"},
	)

	SunSpecDeviceModbusCommonBlockReadsSuccess = prometheus.NewCounterVec(
		*newCounterOptsWithSubsystem(
			"sunspec",
			"connection_common_block_read_successes_total",
			"Total amount of successful common reads on SunSpec smart meters",
		),
		[]string{"sub_device"},
	)

	SunSpecDeviceModbusCommonBlockReadsFailures = prometheus.NewCounterVec(
		*newCounterOptsWithSubsystem(
			"sunspec",
			"connection_common_block_read_failures_total",
			"Total amount of failed common reads on SunSpec smart meters",
		),
		[]string{"sub_device"},
	)

	SunSpecDeviceModbusModelCollectionSuccess = prometheus.NewCounterVec(
		*newCounterOptsWithSubsystem(
			"sunspec",
			"connection_model_collection_successes_total",
			"Total amount of successful model collection tasks on SunSpec smart meters",
		),
		[]string{"sub_device"},
	)

	SunSpecDeviceModbusModelCollectionFailure = prometheus.NewCounterVec(
		*newCounterOptsWithSubsystem(
			"sunspec",
			"connection_model_collection_failures_total",
			"Total amount of failed model collection tasks on SunSpec smart meters",
		),
		[]string{"sub_device"},
	)

	DeviceInfoDetails = prometheus.NewGaugeVec(
		*newGaugeOpts(
			"device_info",
			"Registered smart meter devices of which Prometheus metrics will be collected",
		),
		deviceInfoMetricLabels,
	)
)

func (d deviceCollectors) Collect() []prometheus.Collector {
	return []prometheus.Collector{
		DevicesCreatedTotal,
		CurrentDevicesActive,
		DeviceModbusConnectionAttemptTotal,
		DeviceModbusConnectionFailure,
		DeviceModbusConnectionSuccess,
		DeviceModbusConnectionPartialSuccess,

		DeviceByConfigNotFound,

		SunSpecDeviceModbusCommonBlockReadsSuccess,
		SunSpecDeviceModbusCommonBlockReadsFailures,
		SunSpecDeviceModbusModelCollectionSuccess,
		SunSpecDeviceModbusModelCollectionFailure,

		DeviceInfoDetails,
	}
}

var deviceInfoMetricLabels = []string{
	"type",
	"manufacturer",
	"model",
	"options",
	"version",
	"serial_number",
	"sub_device",
}

// RegisterDevice takes a meters.DeviceDescriptor struct, creates and registers a generalized info metric
// that can be used to join with other metrics that have same label values
func RegisterDevice(deviceDescriptor *meters.DeviceDescriptor) {
	DeviceInfoDetails.WithLabelValues(
		deviceDescriptor.Type,
		deviceDescriptor.Manufacturer,
		deviceDescriptor.Model,
		deviceDescriptor.Options,
		deviceDescriptor.Version,
		deviceDescriptor.Serial,
		strconv.Itoa(deviceDescriptor.SubDevice),
	).Add(1.0)
}
