package prometheus_metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/volkszaehler/mbmd/meters"
)

const NAMESPACE = "mbmd"
const SSN_MISSING = "NOT_AVAILABLE"

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
var counterVecMap = map[meters.Measurement]*prometheus.CounterVec{
	meters.Sum:              nil,
	meters.SumT1:            nil,
	meters.SumT2:            nil,
	meters.SumL1:            nil,
	meters.SumL2:            nil,
	meters.SumL3:            nil,
	meters.Import:           nil,
	meters.ImportT1:         nil,
	meters.ImportT2:         nil,
	meters.ImportL1:         nil,
	meters.ImportL2:         nil,
	meters.ImportL3:         nil,
	meters.Export:           nil,
	meters.ExportT1:         nil,
	meters.ExportT2:         nil,
	meters.ExportL1:         nil,
	meters.ExportL2:         nil,
	meters.ExportL3:         nil,
	meters.ReactiveSum:      nil,
	meters.ReactiveSumT1:    nil,
	meters.ReactiveSumT2:    nil,
	meters.ReactiveSumL1:    nil,
	meters.ReactiveSumL2:    nil,
	meters.ReactiveSumL3:    nil,
	meters.ReactiveImport:   nil,
	meters.ReactiveImportT1: nil,
	meters.ReactiveImportT2: nil,
	meters.ReactiveImportL1: nil,
	meters.ReactiveImportL2: nil,
	meters.ReactiveImportL3: nil,
	meters.ReactiveExport:   nil,
	meters.ReactiveExportT1: nil,
	meters.ReactiveExportT2: nil,
	meters.ReactiveExportL1: nil,
	meters.ReactiveExportL2: nil,
	meters.ReactiveExportL3: nil,
	meters.DCEnergyS1:       nil,
	meters.DCEnergyS2:       nil,
	meters.DCEnergyS3:       nil,
}

// gaugeVecMap contains all meters.Measurement that are associated with a prometheus.Gauge
//
// If a new meters.Measurement is introduced, it needs to be added either to counterVecMap
// or to gaugeVecMap - Otherwise Prometheus won't keep track of the newly added meters.Measurement
var gaugeVecMap = map[meters.Measurement]*prometheus.GaugeVec{
	meters.Frequency:        MeasurementFrequency,
	meters.Current:          MeasurementElectricCurrent,
	meters.CurrentL1:        MeasurementL1Current,
	meters.CurrentL2:        MeasurementL2Current,
	meters.CurrentL3:        MeasurementL3Current,
	meters.Voltage:          MeasurementVoltage,
	meters.VoltageL1: 		 nil,
	meters.VoltageL2:        nil,
	meters.VoltageL3:        nil,
	meters.Power:            nil,
	meters.PowerL1:          nil,
	meters.PowerL2:          nil,
	meters.PowerL3:          nil,
	meters.ImportPower:      nil,
	meters.ImportPowerL1:    nil,
	meters.ImportPowerL2:    nil,
	meters.ImportPowerL3:    nil,
	meters.ExportPower:      nil,
	meters.ExportPowerL1:    nil,
	meters.ExportPowerL2:    nil,
	meters.ExportPowerL3:    nil,
	meters.ReactivePower:    nil,
	meters.ReactivePowerL1:  nil,
	meters.ReactivePowerL2:  nil,
	meters.ReactivePowerL3:  nil,
	meters.ApparentPower:    nil,
	meters.ApparentPowerL1:  nil,
	meters.ApparentPowerL2:  nil,
	meters.ApparentPowerL3:  nil,
	meters.Cosphi:           nil,
	meters.CosphiL1:         nil,
	meters.CosphiL2:         nil,
	meters.CosphiL3:         nil,
	meters.THD:              nil,
	meters.THDL1:            nil,
	meters.THDL2:            nil,
	meters.THDL3:            nil,
	meters.DCCurrent:        nil,
	meters.DCVoltage:        nil,
	meters.DCPower:          nil,
	meters.HeatSinkTemp:     nil,
	meters.DCCurrentS1:      nil,
	meters.DCVoltageS1:      nil,
	meters.DCPowerS1:        nil,
	meters.DCCurrentS2:      nil,
	meters.DCVoltageS2:      nil,
	meters.DCPowerS2:        nil,
	meters.DCCurrentS3:      nil,
	meters.DCVoltageS3:      nil,
	meters.DCPowerS3:        nil,
	meters.ChargeState:      nil,
	meters.BatteryVoltage:   nil,
	meters.PhaseAngle:       nil,
}

// Init registers all globally defined metrics to Prometheus library's default registry
func Init() {
	//prometheus.MustRegister(
	//	//ConnectionAttemptTotal,
	//	//ConnectionAttemptFailedTotal,
	//	//ConnectionPartiallySuccessfulTotal,
	//	//ReadDeviceDetailsFailedTotal,
	//	//DevicesCreatedTotal,
	//
	//	// Device specific actions
	//	DeviceQueriesTotal,
	//	DeviceQueriesErrorTotal,
	//	DeviceQueriesSuccessTotal,
	//	DeviceQueryMeasurementValueSkippedTotal,
	//
	//	// Bus scan metrics of cmd.scan
	//	BusScanTotal,
	//	BusScanDeviceProbeSuccessfulTotal,
	//	BusScanDeviceProbeFailedTotal,
	//
	//	// Measurement gauges
	//	MeasurementElectricCurrent,
	//	MeasurementL1Current,
	//	MeasurementL2Current,
	//	MeasurementL3Current,
	//)

	initAndRegisterGauges()
	initAndRegisterCounters()
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

func initAndRegisterCounters() {
	//collectors := make([]prometheus.Collector, len(counterVecMap))

	for measurement := range counterVecMap {
		newCounter := prometheus.NewCounterVec(
			newCounterOpts(
				measurement.PrometheusName(),
				prometheusDescription(measurement),
			),
			[]string{"device_id", "serial_number"},
		)
		counterVecMap[measurement] = newCounter
		//collectors = append(collectors, newCounter)
		prometheus.MustRegister(newCounter)
	}

	// fmt.Print("Registering new collectors of type CounterVec")
	//prometheus.MustRegister(collectors...)
}

func initAndRegisterGauges() {
	// gauges := make([]*prometheus.GaugeVec, len(gaugeVecMap))

	for measurement := range gaugeVecMap {
		newGauge := prometheus.NewGaugeVec(
			newGaugeOpts(
				measurement.PrometheusName(),
				prometheusDescription(measurement),
			),
			[]string{"device_id", "serial_number"},
		)
		gaugeVecMap[measurement] = newGauge
		// gauges = append(gauges, newGauge)
		prometheus.MustRegister(newGauge)
	}

	//fmt.Print("Registering new collectors of type GaugeVec")

}

func prometheusDescription(measurement meters.Measurement) string {
	description, unit := measurement.DescriptionAndUnit()

	return fmt.Sprintf("Measurement of %s in %s", description, unit)
}
