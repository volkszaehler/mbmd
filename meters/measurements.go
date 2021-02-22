package meters

import (
	"fmt"
	"strings"
	"time"
)

// MeasurementResult is the result of modbus read operation
type MeasurementResult struct {
	Measurement
	Value     float64
	Timestamp time.Time
}

func (r MeasurementResult) String() string {
	_, unit := r.Measurement.DescriptionAndUnit()
	return fmt.Sprintf("%s: %.2f%s", r.Measurement.String(), r.Value, unit)
}

// MeasurementDescription describes a Measurement itself, its unit and according prometheus.Metric type
type MeasurementDescription struct {
	Description string
	Unit        Unit
	MetricType  MeasurementMetricType
}

// MeasurementMetricType is the type of a Measurement's prometheus.Metric to be used
type MeasurementMetricType int

const (
	_ MeasurementMetricType = iota
	Gauge
	Counter
)

// Measurement is the type of measurement, i.e. the physical property being measured in common notation
type Measurement int

//go:generate enumer -type=Measurement
const (
	_ Measurement = iota

	Frequency

	Current
	CurrentL1
	CurrentL2
	CurrentL3

	// phases and sums
	Voltage
	VoltageL1
	VoltageL2
	VoltageL3

	Power // synonymous ActivePower
	PowerL1
	PowerL2
	PowerL3

	ImportPower
	ImportPowerL1
	ImportPowerL2
	ImportPowerL3

	ExportPower
	ExportPowerL1
	ExportPowerL2
	ExportPowerL3

	ReactivePower
	ReactivePowerL1
	ReactivePowerL2
	ReactivePowerL3

	ApparentPower
	ApparentPowerL1
	ApparentPowerL2
	ApparentPowerL3

	Cosphi
	CosphiL1
	CosphiL2
	CosphiL3

	THD
	THDL1
	THDL2
	THDL3

	// energy
	Sum // synonymous ActiveEnergy
	SumT1
	SumT2
	SumL1
	SumL2
	SumL3

	Import
	ImportT1
	ImportT2
	ImportL1
	ImportL2
	ImportL3

	Export
	ExportT1
	ExportT2
	ExportL1
	ExportL2
	ExportL3

	ReactiveSum
	ReactiveSumT1
	ReactiveSumT2
	ReactiveSumL1
	ReactiveSumL2
	ReactiveSumL3

	ReactiveImport
	ReactiveImportT1
	ReactiveImportT2
	ReactiveImportL1
	ReactiveImportL2
	ReactiveImportL3

	ReactiveExport
	ReactiveExportT1
	ReactiveExportT2
	ReactiveExportL1
	ReactiveExportL2
	ReactiveExportL3

	// DC
	DCCurrent
	DCVoltage
	DCPower
	HeatSinkTemp

	// Strings
	DCCurrentS1
	DCVoltageS1
	DCPowerS1
	DCEnergyS1
	DCCurrentS2
	DCVoltageS2
	DCPowerS2
	DCEnergyS2
	DCCurrentS3
	DCVoltageS3
	DCPowerS3
	DCEnergyS3

	// Battery
	ChargeState
	BatteryVoltage

	PhaseAngle
)

var iec = map[Measurement]MeasurementDescription{
	Frequency:        {"Frequency", Hertz, Gauge},
	Current:          {"Current", Ampere, Gauge},
	CurrentL1:        {"L1 Current", Ampere, Gauge},
	CurrentL2:        {"L2 Current", Ampere, Gauge},
	CurrentL3:        {"L3 Current", Ampere, Gauge},
	Voltage:          {"Voltage", Volt, Gauge},
	VoltageL1:        {"L1 Voltage", Volt, Gauge},
	VoltageL2:        {"L2 Voltage", Volt, Gauge},
	VoltageL3:        {"L3 Voltage", Volt, Gauge},
	Power:            {"Power", Watt, Gauge},
	PowerL1:          {"L1 Power", Watt, Gauge},
	PowerL2:          {"L2 Power", Watt, Gauge},
	PowerL3:          {"L3 Power", Watt, Gauge},
	ImportPower:      {"Import Power", Watt, Gauge},
	ImportPowerL1:    {"L1 Import Power", Watt, Gauge},
	ImportPowerL2:    {"L2 Import Power", Watt, Gauge},
	ImportPowerL3:    {"L3 Import Power", Watt, Gauge},
	ExportPower:      {"Export Power", Watt, Gauge},
	ExportPowerL1:    {"L1 Export Power", Watt, Gauge},
	ExportPowerL2:    {"L2 Export Power", Watt, Gauge},
	ExportPowerL3:    {"L3 Export Power", Watt, Gauge},
	ReactivePower:    {"Reactive Power", Var, Gauge},
	ReactivePowerL1:  {"L1 Reactive Power", Var, Gauge},
	ReactivePowerL2:  {"L2 Reactive Power", Var, Gauge},
	ReactivePowerL3:  {"L3 Reactive Power", Var, Gauge},
	ApparentPower:    {"Apparent Power", VoltAmpere, Gauge},
	ApparentPowerL1:  {"L1 Apparent Power", VoltAmpere, Gauge},
	ApparentPowerL2:  {"L2 Apparent Power", VoltAmpere, Gauge},
	ApparentPowerL3:  {"L3 Apparent Power", VoltAmpere, Gauge},
	Cosphi:           {"Cosphi", NoUnit, Gauge},
	CosphiL1:         {"L1 Cosphi", NoUnit,Gauge},
	CosphiL2:         {"L2 Cosphi", NoUnit,Gauge},
	CosphiL3:         {"L3 Cosphi", NoUnit, Gauge},
	THD:              {"Average voltage to neutral THD", Percent, Gauge},
	THDL1:            {"L1 Voltage to neutral THD", Percent, Gauge},
	THDL2:            {"L2 Voltage to neutral THD", Percent, Gauge},
	THDL3:            {"L3 Voltage to neutral THD", Percent, Gauge},
	Sum:              {"Total Sum", KiloWattHour, Counter},
	SumT1:            {"Tariff 1 Sum", KiloWattHour, Counter},
	SumT2:            {"Tariff 2 Sum", KiloWattHour, Counter},
	SumL1:            {"L1 Sum", KiloWattHour, Counter},
	SumL2:            {"L2 Sum", KiloWattHour, Counter},
	SumL3:            {"L3 Sum", KiloWattHour, Counter},
	Import:           {"Total Import", KiloWattHour, Counter},
	ImportT1:         {"Tariff 1 Import", KiloWattHour, Counter},
	ImportT2:         {"Tariff 2 Import", KiloWattHour, Counter},
	ImportL1:         {"L1 Import", KiloWattHour, Counter},
	ImportL2:         {"L2 Import", KiloWattHour, Counter},
	ImportL3:         {"L3 Import", KiloWattHour, Counter},
	Export:           {"Total Export", KiloWattHour, Counter},
	ExportT1:         {"Tariff 1 Export", KiloWattHour, Counter},
	ExportT2:         {"Tariff 2 Export", KiloWattHour, Counter},
	ExportL1:         {"L1 Export", KiloWattHour, Counter},
	ExportL2:         {"L2 Export", KiloWattHour, Counter},
	ExportL3:         {"L3 Export", KiloWattHour, Counter},
	ReactiveSum:      {"Total Reactive", KiloVarHour, Counter},
	ReactiveSumT1:    {"Tariff 1 Reactive", KiloVarHour, Counter},
	ReactiveSumT2:    {"Tariff 2 Reactive", KiloVarHour, Counter},
	ReactiveSumL1:    {"L1 Reactive", KiloVarHour, Counter},
	ReactiveSumL2:    {"L2 Reactive", KiloVarHour, Counter},
	ReactiveSumL3:    {"L3 Reactive", KiloVarHour, Counter},
	ReactiveImport:   {"Reactive Import", KiloVarHour, Counter},
	ReactiveImportT1: {"Tariff 1 Reactive Import", KiloVarHour, Counter},
	ReactiveImportT2: {"Tariff 2 Reactive Import", KiloVarHour, Counter},
	ReactiveImportL1: {"L1 Reactive Import", KiloVarHour, Counter},
	ReactiveImportL2: {"L2 Reactive Import", KiloVarHour, Counter},
	ReactiveImportL3: {"L3 Reactive Import", KiloVarHour, Counter},
	ReactiveExport:   {"Reactive Export", KiloVarHour, Counter},
	ReactiveExportT1: {"Tariff 1 Reactive Export", KiloVarHour, Counter},
	ReactiveExportT2: {"Tariff 2 Reactive Export", KiloVarHour, Counter},
	ReactiveExportL1: {"L1 Reactive Export", KiloVarHour, Counter},
	ReactiveExportL2: {"L2 Reactive Export", KiloVarHour, Counter},
	ReactiveExportL3: {"L3 Reactive Export", KiloVarHour, Counter},
	DCCurrent:        {"DC Current", Ampere, Gauge},
	DCVoltage:        {"DC Voltage", Volt, Gauge},
	DCPower:          {"DC Power", Watt, Gauge},
	HeatSinkTemp:     {"Heat Sink Temperature", DegreeCelsius, Gauge},
	DCCurrentS1:      {"String 1 Current", Ampere, Gauge},
	DCVoltageS1:      {"String 1 Voltage", Volt, Gauge},
	DCPowerS1:        {"String 1 Power", Watt, Gauge},
	DCEnergyS1:       {"String 1 Generation", KiloWattHour, Counter},
	DCCurrentS2:      {"String 2 Current", Ampere, Gauge},
	DCVoltageS2:      {"String 2 Voltage", Volt, Gauge},
	DCPowerS2:        {"String 2 Power", Watt, Gauge},
	DCEnergyS2:       {"String 2 Generation", KiloWattHour, Counter},
	DCCurrentS3:      {"String 3 Current", Ampere, Gauge},
	DCVoltageS3:      {"String 3 Voltage", Volt, Gauge},
	DCPowerS3:        {"String 3 Power", Watt, Gauge},
	DCEnergyS3:       {"String 3 Generation", KiloWattHour, Counter},
	ChargeState:      {"Charge State", Percent, Gauge},
	BatteryVoltage:   {"Battery Voltage", Volt, Gauge},
	PhaseAngle:       {"Phase Angle", Degree, Gauge},
}

// MarshalText implements encoding.TextMarshaler
func (m *Measurement) MarshalText() (text []byte, err error) {
	return []byte(m.String()), nil
}

// DescriptionAndUnit returns a measurements human-readable name and its unit
func (m *Measurement) DescriptionAndUnit() (string, string) {
	if details, ok := iec[*m]; ok {
		unit := details.Unit
		description := details.Description
		return description, unit.Abbreviation()
	}
	return m.String(), ""
}

// Description returns a measurements human-readable name
func (m *Measurement) Description() string {
	description, unit := m.DescriptionAndUnit()
	if unit != "" {
		description = description + " (" + unit + ")"
	}
	return description
}

// PrometheusMetricType returns the Measurement's associated prometheus.Metric
func (m *Measurement) PrometheusMetricType() MeasurementMetricType {
	if measurement, ok := iec[*m]; ok {
		return measurement.MetricType
	}
	return 0
}

// PrometheusDescription returns a description text appropriate for prometheus.Metric
func (m *Measurement) PrometheusDescription() string {
	description, unit := m.DescriptionAndUnit()

	if unit != "" {
		return fmt.Sprintf("Measurement of %s in %s", description, unit)
	} else {
		return fmt.Sprintf("Measurement of %s", description)
	}
}

// PrometheusName returns a name and its associated unit for Prometheus counters
func (m *Measurement) PrometheusName() string {
	var description string
	var unit string
	if details, ok := iec[*m]; ok {
		unit = details.Unit.PrometheusName()
		description = details.Description
	} else {
		description = m.String()
	}

	if unit != "" {
		unit = strings.ToLower(unit)
	}

	description = strings.ReplaceAll(strings.ToLower(description), " ", "_")

	return strings.Trim( // Trim trailing underscore (e. g. when unit string is empty)
		strings.Join(
			[]string{"measurement", description, unit},
			"_",
		),
		"_",
	)
}


