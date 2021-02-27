package prometheus_metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

// Static metrics
var (
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
)

// RegisterStatics registers all globally defined static metrics to Prometheus library's default registry
func RegisterStatics() {
	collectables := getAllCollectors()

	for _, collectable := range collectables {
		for _, prometheusCollector := range collectable.Collect() {
			if err := prometheus.Register(prometheusCollector); err != nil {
				log.Fatalf("Could not register a metric '%s' (%s)", prometheusCollector, err)
			}
		}
	}

	//prometheus.MustRegister(
	//	//DeviceModbusConnectionAttemptTotal,
	//	//DeviceModbusConnectionFailure,
	//	//ConnectionPartiallySuccessfulTotal,
	//	//DevicesCreatedTotal,
	//	//BusScanStartedTotal,
	//	//BusScanDeviceInitializationErrorTotal,
	//	//BusScanTotal,
	//	//BusScanDeviceProbeSuccessfulTotal,
	//	//BusScanDeviceProbeFailedTotal,
	//	//ReadDeviceDetailsFailedTotal,
	//	//ConnectionHandlerDeviceQueriesTotal,
	//	//ConnectionHandlerDeviceQueriesErrorTotal,
	//	//ConnectionHandlerDeviceQueriesSuccessTotal,
	//	//ConnectionHandlerDeviceQueryMeasurementValueSkippedTotal,
	//
	//	// Socket related metrics
	//	collectors...
	//)
}

func getAllCollectors() []collectable {
	return []collectable{
		socketCollectors{},
		publisherCollectors{},
		handlerCollectors{},
		deviceCollectors{},
	}
}
