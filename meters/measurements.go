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
	Unit        string
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
	DCCurrentS4
	DCVoltageS4
	DCPowerS4
	DCEnergyS4

	// Battery
	ChargeState
	BatteryVoltage

	PhaseAngle
)

var iec = map[Measurement]MeasurementDescription{
	Frequency:        {"Frequency", "Hz", Gauge},
	Current:          {"Current", "A", Gauge},
	CurrentL1:        {"L1 Current", "A", Gauge},
	CurrentL2:        {"L2 Current", "A", Gauge},
	CurrentL3:        {"L3 Current", "A", Gauge},
	Voltage:          {"Voltage", "V", Gauge},
	VoltageL1:        {"L1 Voltage", "V", Gauge},
	VoltageL2:        {"L2 Voltage", "V", Gauge},
	VoltageL3:        {"L3 Voltage", "V", Gauge},
	Power:            {"Power", "W", Gauge},
	PowerL1:          {"L1 Power", "W", Gauge},
	PowerL2:          {"L2 Power", "W", Gauge},
	PowerL3:          {"L3 Power", "W", Gauge},
	ImportPower:      {"Import Power", "W", Gauge},
	ImportPowerL1:    {"L1 Import Power", "W", Gauge},
	ImportPowerL2:    {"L2 Import Power", "W", Gauge},
	ImportPowerL3:    {"L3 Import Power", "W", Gauge},
	ExportPower:      {"Export Power", "W", Gauge},
	ExportPowerL1:    {"L1 Export Power", "W", Gauge},
	ExportPowerL2:    {"L2 Export Power", "W", Gauge},
	ExportPowerL3:    {"L3 Export Power", "W", Gauge},
	ReactivePower:    {"Reactive Power", "var", Gauge},
	ReactivePowerL1:  {"L1 Reactive Power", "var", Gauge},
	ReactivePowerL2:  {"L2 Reactive Power", "var", Gauge},
	ReactivePowerL3:  {"L3 Reactive Power", "var", Gauge},
	ApparentPower:    {"Apparent Power", "VA", Gauge},
	ApparentPowerL1:  {"L1 Apparent Power", "VA", Gauge},
	ApparentPowerL2:  {"L2 Apparent Power", "VA", Gauge},
	ApparentPowerL3:  {"L3 Apparent Power", "VA", Gauge},
	Cosphi:           {"Cosphi", "", Gauge},
	CosphiL1:         {"L1 Cosphi", "",Gauge},
	CosphiL2:         {"L2 Cosphi", "",Gauge},
	CosphiL3:         {"L3 Cosphi", "", Gauge},
	THD:              {"Average voltage to neutral THD", "%", Gauge},
	THDL1:            {"L1 Voltage to neutral THD", "%", Gauge},
	THDL2:            {"L2 Voltage to neutral THD", "%", Gauge},
	THDL3:            {"L3 Voltage to neutral THD", "%", Gauge},
	Sum:              {"Total Sum", "kWh", Counter},
	SumT1:            {"Tariff 1 Sum", "kWh", Counter},
	SumT2:            {"Tariff 2 Sum", "kWh", Counter},
	SumL1:            {"L1 Sum", "kWh", Counter},
	SumL2:            {"L2 Sum", "kWh", Counter},
	SumL3:            {"L3 Sum", "kWh", Counter},
	Import:           {"Total Import", "kWh", Counter},
	ImportT1:         {"Tariff 1 Import", "kWh", Counter},
	ImportT2:         {"Tariff 2 Import", "kWh", Counter},
	ImportL1:         {"L1 Import", "kWh", Counter},
	ImportL2:         {"L2 Import", "kWh", Counter},
	ImportL3:         {"L3 Import", "kWh", Counter},
	Export:           {"Total Export", "kWh", Counter},
	ExportT1:         {"Tariff 1 Export", "kWh", Counter},
	ExportT2:         {"Tariff 2 Export", "kWh", Counter},
	ExportL1:         {"L1 Export", "kWh", Counter},
	ExportL2:         {"L2 Export", "kWh", Counter},
	ExportL3:         {"L3 Export", "kWh", Counter},
	ReactiveSum:      {"Total Reactive", "kvarh", Counter},
	ReactiveSumT1:    {"Tariff 1 Reactive", "kvarh", Counter},
	ReactiveSumT2:    {"Tariff 2 Reactive", "kvarh", Counter},
	ReactiveSumL1:    {"L1 Reactive", "kvarh", Counter},
	ReactiveSumL2:    {"L2 Reactive", "kvarh", Counter},
	ReactiveSumL3:    {"L3 Reactive", "kvarh", Counter},
	ReactiveImport:   {"Reactive Import", "kvarh", Counter},
	ReactiveImportT1: {"Tariff 1 Reactive Import", "kvarh", Counter},
	ReactiveImportT2: {"Tariff 2 Reactive Import", "kvarh", Counter},
	ReactiveImportL1: {"L1 Reactive Import", "kvarh", Counter},
	ReactiveImportL2: {"L2 Reactive Import", "kvarh", Counter},
	ReactiveImportL3: {"L3 Reactive Import", "kvarh", Counter},
	ReactiveExport:   {"Reactive Export", "kvarh", Counter},
	ReactiveExportT1: {"Tariff 1 Reactive Export", "kvarh", Counter},
	ReactiveExportT2: {"Tariff 2 Reactive Export", "kvarh", Counter},
	ReactiveExportL1: {"L1 Reactive Export", "kvarh", Counter},
	ReactiveExportL2: {"L2 Reactive Export", "kvarh", Counter},
	ReactiveExportL3: {"L3 Reactive Export", "kvarh", Counter},
	DCCurrent:        {"DC Current", "A", Gauge},
	DCVoltage:        {"DC Voltage", "V", Gauge},
	DCPower:          {"DC Power", "W", Gauge},
	HeatSinkTemp:     {"Heat Sink Temperature", "°C", Gauge},
	DCCurrentS1:      {"String 1 Current", "A", Gauge},
	DCVoltageS1:      {"String 1 Voltage", "V", Gauge},
	DCPowerS1:        {"String 1 Power", "W", Gauge},
	DCEnergyS1:       {"String 1 Generation", "kWh", Counter},
	DCCurrentS2:      {"String 2 Current", "A", Gauge},
	DCVoltageS2:      {"String 2 Voltage", "V", Gauge},
	DCPowerS2:        {"String 2 Power", "W", Gauge},
	DCEnergyS2:       {"String 2 Generation", "kWh", Counter},
	DCCurrentS3:      {"String 3 Current", "A", Gauge},
	DCVoltageS3:      {"String 3 Voltage", "V", Gauge},
	DCPowerS3:        {"String 3 Power", "W", Gauge},
	DCEnergyS3:       {"String 3 Generation", "kWh", Counter},
	DCCurrentS4:      {"String 4 Current", "A"},
	DCVoltageS4:      {"String 4 Voltage", "V"},
	DCPowerS4:        {"String 4 Power", "W"},
	DCEnergyS4:       {"String 4 Generation", "kWh"},
	ChargeState:      {"Charge State", "%", Gauge},
	BatteryVoltage:   {"Battery Voltage", "V", Gauge},
	PhaseAngle:       {"Phase Angle", "°", Gauge},
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
		return description, unit
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
	description, unit := m.DescriptionAndUnit()
	if unit != "" {
		unit = strings.ToLower(unit)
	}

	description = strings.ReplaceAll(strings.ToLower(description), " ", "_")

	paraphrasedUnit := paraphraseChars(unit)

	return strings.Trim( // Trim trailing underscore (e. g. when unit string is empty)
		strings.Join(
			[]string{"measurement", description, paraphrasedUnit},
			"_",
		),
		"_",
	)
}

func paraphraseChars(text string) string {
	text = strings.ToLower(text)
	result := ""
	switch text {
	case "°":
		result = "degree"
	case "°c":
		result = "degree_celsius"
	case "%":
		result = "percent"
	case "a":
		result = "ampere"
	case "v":
		result = "volt"
	case "w":
		result = "watt"
	case "kw":
		result = "kilowatt"
	case "kwh":
		result = "kilowatt_per_hour"
	case "va":
		result = "voltampere"
	case "kvarh":
		result = "kilovar_per_hour"
	default:
		result = text
	}

	return result
}

