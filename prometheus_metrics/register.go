package prometheus_metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

// Static metrics
var (
	BusScanDeviceInitializationErrorTotal = prometheus.NewCounterVec(
		*newCounterOpts(
			"bus_scan_device_initialization_error_total",
			"Total errors upon initialization of a device during bus scan",
		),
		[]string{"device_name"},
	)

	BusScanDeviceProbeSuccessfulTotal = prometheus.NewCounterVec(
		*newCounterOpts(
			"bus_scan_device_probe_successful_total",
			"Amount of successfully found devices during bus scan",
		),
		[]string{"device_name", "serial_number"},
	)

	BusScanDeviceProbeFailedTotal = prometheus.NewCounterVec(
		*newCounterOpts(
			"bus_scan_device_probe_failed_total",
			"Amount of devices failed to be found during bus scan",
		),
		[]string{"device_name"},
	)
)

// RegisterStatics registers all globally defined static metrics to Prometheus library's default registry
func RegisterStatics() {
	collectables := getAllCollectors()

	for _, collectable := range collectables {
		for _, prometheusCollector := range collectable.Collect() {
			if err := prometheus.Register(prometheusCollector); err != nil {
				log.Printf("Could not register a metric '%s' (%s)", prometheusCollector, err)
			}
		}
	}
}

func getAllCollectors() []collectable {
	return []collectable{
		socketCollectors{},
		publisherCollectors{},
		handlerCollectors{},
		deviceCollectors{},
	}
}
