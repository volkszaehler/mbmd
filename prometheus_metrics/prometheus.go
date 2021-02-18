package prometheus_metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/volkszaehler/mbmd/meters"
)

var (
	ConnectionAttemptTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "smart_meter_connection_attempt_total",
			Help: "Total amount of a smart meter's connection attempts",
		},
		[]string{"model", "sub_device"},
	)

	ConnectionAttemptFailedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "smart_meter_connection_attempt_failed_total",
			Help: "Amount of a smart meter's connection failures",
		},
		[]string{"model", "sub_device"},
	)

	ConnectionPartiallySuccessfulTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "smart_meter_connection_partially_successful_total",
			Help: "Number of connections that are partially open",
		},
		[]string{"model", "sub_device"},
	)

	DevicesCreatedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "smart_meter_devices_created_total",
			Help: "Number of smart meter devices created/registered",
		},
		[]string{"meter_type", "sub_device"},
	)

	BusScanStartedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bus_scan_started_total",
			Help: "Total started bus scans",
		},
		[]string{"device_id"},
	)

	BusScanDeviceInitializationErrorTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bus_scan_device_initialization_error_total",
			Help: "Total errors upon initialization of a device during bus scan",
		},
		[]string{"device_id"},
	)

	BusScanTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "bus_scan_total",
			Help: "Amount of bus scans done",
		},
	)

	BusScanDeviceProbeSuccessfulTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bus_scan_device_probe_successful_total",
			Help: "Amount of successfully found devices during bus scan",
		},
		[]string{"device_id", "serial_number"},
	)

	BusScanDeviceProbeFailedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bus_scan_device_probe_failed_total",
			Help: "Amount of devices failed to be found during bus scan",
		},
		[]string{"device_id"},
	)

	MeasurementElectricCurrent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "measurement_electric_current_ampere",
			Help: "Last electric current measured",
		},
		[]string{"device_id", "serial_number"},
	)

	ReadDeviceDetailsFailedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "smart_meter_read_device_details_failed_total",
			Help: "Reading additional details of a smart meter failed",
		},
		[]string{"model"},
	)

	DeviceQueriesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "smart_meter_queries_total",
			Help: "Amount of queries/requests done for a smart meter",
		},
		[]string{"device_id", "serial_number"},
	)

	DeviceQueriesErrorTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "smart_meter_queries_error_total",
			Help: "Errors occured during smart meter query",
		},
		[]string{"device_id", "serial_number"},
	)

	DeviceQueriesSuccessTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "smart_meter_queries_success_total",
			Help: "Successful smart meter query",
		},
		[]string{"device_id", "serial_number"},
	)

	DeviceQueryMeasurementValueSkippedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "smart_meter_queries_measurement_value_skipped_total",
			Help: "NaN measurement values found and skipped during smart meter query",
		},
		[]string{"device_id", "serial_number"},
	)
)

var counterVecMap = map[meters.Measurement]*prometheus.CounterVec{}

var gaugeVecMap = map[meters.Measurement]*prometheus.GaugeVec{
	//meters.Frequency:        {"Frequency", "Hz"},
	meters.Current:          MeasurementElectricCurrent,
	//meters.CurrentL1:        {"L1 Current", "A"},
	//meters.CurrentL2:        {"L2 Current", "A"},
	//meters.CurrentL3:        {"L3 Current", "A"},
	//meters.Voltage:          {"Voltage", "V"},
	//meters.VoltageL1:        {"L1 Voltage", "V"},
	//meters.VoltageL2:        {"L2 Voltage", "V"},
	//meters.VoltageL3:        {"L3 Voltage", "V"},
	//meters.Power:            {"Power", "W"},
	//meters.PowerL1:          {"L1 Power", "W"},
	//meters.PowerL2:          {"L2 Power", "W"},
	//meters.PowerL3:          {"L3 Power", "W"},
	//meters.ImportPower:      {"Import Power", "W"},
	//meters.ImportPowerL1:    {"L1 Import Power", "W"},
	//meters.ImportPowerL2:    {"L2 Import Power", "W"},
	//meters.ImportPowerL3:    {"L3 Import Power", "W"},
	//meters.ExportPower:      {"Export Power", "W"},
	//meters.ExportPowerL1:    {"L1 Export Power", "W"},
	//meters.ExportPowerL2:    {"L2 Export Power", "W"},
	//meters.ExportPowerL3:    {"L3 Export Power", "W"},
	//meters.ReactivePower:    {"Reactive Power", "var"},
	//meters.ReactivePowerL1:  {"L1 Reactive Power", "var"},
	//meters.ReactivePowerL2:  {"L2 Reactive Power", "var"},
	//meters.ReactivePowerL3:  {"L3 Reactive Power", "var"},
	//meters.ApparentPower:    {"Apparent Power", "VA"},
	//meters.ApparentPowerL1:  {"L1 Apparent Power", "VA"},
	//meters.ApparentPowerL2:  {"L2 Apparent Power", "VA"},
	//meters.ApparentPowerL3:  {"L3 Apparent Power", "VA"},
	//meters.Cosphi:           {"Cosphi"},
	//meters.CosphiL1:         {"L1 Cosphi"},
	//meters.CosphiL2:         {"L2 Cosphi"},
	//meters.CosphiL3:         {"L3 Cosphi"},
	//meters.THD:              {"Average voltage to neutral THD", "%"},
	//meters.THDL1:            {"L1 Voltage to neutral THD", "%"},
	//meters.THDL2:            {"L2 Voltage to neutral THD", "%"},
	//meters.THDL3:            {"L3 Voltage to neutral THD", "%"},
	//meters.Sum:              {"Total Sum", "kWh"},
	//meters.SumT1:            {"Tariff 1 Sum", "kWh"},
	//meters.SumT2:            {"Tariff 2 Sum", "kWh"},
	//meters.SumL1:            {"L1 Sum", "kWh"},
	//meters.SumL2:            {"L2 Sum", "kWh"},
	//meters.SumL3:            {"L3 Sum", "kWh"},
	//meters.Import:           {"Total Import", "kWh"},
	//meters.ImportT1:         {"Tariff 1 Import", "kWh"},
	//meters.ImportT2:         {"Tariff 2 Import", "kWh"},
	//meters.ImportL1:         {"L1 Import", "kWh"},
	//meters.ImportL2:         {"L2 Import", "kWh"},
	//meters.ImportL3:         {"L3 Import", "kWh"},
	//meters.Export:           {"Total Export", "kWh"},
	//meters.ExportT1:         {"Tariff 1 Export", "kWh"},
	//meters.ExportT2:         {"Tariff 2 Export", "kWh"},
	//meters.ExportL1:         {"L1 Export", "kWh"},
	//meters.ExportL2:         {"L2 Export", "kWh"},
	//meters.ExportL3:         {"L3 Export", "kWh"},
	//meters.ReactiveSum:      {"Total Reactive", "kvarh"},
	//meters.ReactiveSumT1:    {"Tariff 1 Reactive", "kvarh"},
	//meters.ReactiveSumT2:    {"Tariff 2 Reactive", "kvarh"},
	//meters.ReactiveSumL1:    {"L1 Reactive", "kvarh"},
	//meters.ReactiveSumL2:    {"L2 Reactive", "kvarh"},
	//meters.ReactiveSumL3:    {"L3 Reactive", "kvarh"},
	//meters.ReactiveImport:   {"Reactive Import", "kvarh"},
	//meters.ReactiveImportT1: {"Tariff 1 Reactive Import", "kvarh"},
	//meters.ReactiveImportT2: {"Tariff 2 Reactive Import", "kvarh"},
	//meters.ReactiveImportL1: {"L1 Reactive Import", "kvarh"},
	//meters.ReactiveImportL2: {"L2 Reactive Import", "kvarh"},
	//meters.ReactiveImportL3: {"L3 Reactive Import", "kvarh"},
	//meters.ReactiveExport:   {"Reactive Export", "kvarh"},
	//meters.ReactiveExportT1: {"Tariff 1 Reactive Export", "kvarh"},
	//meters.ReactiveExportT2: {"Tariff 2 Reactive Export", "kvarh"},
	//meters.ReactiveExportL1: {"L1 Reactive Export", "kvarh"},
	//meters.ReactiveExportL2: {"L2 Reactive Export", "kvarh"},
	//meters.ReactiveExportL3: {"L3 Reactive Export", "kvarh"},
	//meters.DCCurrent:        {"DC Current", "A"},
	//meters.DCVoltage:        {"DC Voltage", "V"},
	//meters.DCPower:          {"DC Power", "W"},
	//meters.HeatSinkTemp:     {"Heat Sink Temperature", "°C"},
	//meters.DCCurrentS1:      {"String 1 Current", "A"},
	//meters.DCVoltageS1:      {"String 1 Voltage", "V"},
	//meters.DCPowerS1:        {"String 1 Power", "W"},
	//meters.DCEnergyS1:       {"String 1 Generation", "kWh"},
	//meters.DCCurrentS2:      {"String 2 Current", "A"},
	//meters.DCVoltageS2:      {"String 2 Voltage", "V"},
	//meters.DCPowerS2:        {"String 2 Power", "W"},
	//meters.DCEnergyS2:       {"String 2 Generation", "kWh"},
	//meters.DCCurrentS3:      {"String 3 Current", "A"},
	//meters.DCVoltageS3:      {"String 3 Voltage", "V"},
	//meters.DCPowerS3:        {"String 3 Power", "W"},
	//meters.DCEnergyS3:       {"String 3 Generation", "kWh"},
	//meters.ChargeState:      {"Charge State", "%"},
	//meters.BatteryVoltage:   {"Battery Voltage", "V"},
	//meters.PhaseAngle:       {"Phase Angle", "°"},
}

// Init registers all globally defined metrics to Prometheus library's default registry
func Init() {
	prometheus.MustRegister(
		//ConnectionAttemptTotal,
		//ConnectionAttemptFailedTotal,
		//ConnectionPartiallySuccessfulTotal,
		//ReadDeviceDetailsFailedTotal,
		//DevicesCreatedTotal,

		// Device specific actions
		DeviceQueriesTotal,
		DeviceQueriesErrorTotal,
		DeviceQueriesSuccessTotal,
		DeviceQueryMeasurementValueSkippedTotal,

		// Bus scan metrics of cmd.scan
		BusScanTotal,
		BusScanDeviceProbeSuccessfulTotal,
		BusScanDeviceProbeFailedTotal,

		// Measurement gauges
		MeasurementElectricCurrent,
	)
}

// UpdateMeasurementMetric updates a counter or gauge based by passed measurement
func UpdateMeasurementMetric(
	deviceId string,
	deviceSerial string,
	measurement meters.MeasurementResult,
) {
	if gauge, err := gaugeVecMap[measurement.Measurement]; !err {
		gauge.WithLabelValues(deviceId, deviceSerial).Set(measurement.Value)
	} else if counter, err := counterVecMap[measurement.Measurement]; !err {
		counter.WithLabelValues(deviceId, deviceSerial).Add(measurement.Value)
	}
}
