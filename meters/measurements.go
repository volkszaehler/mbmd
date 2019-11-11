package meters

import (
	"fmt"
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

// Measurement is the type of measurement, i.e. the physical property being measued in common notation
type Measurement int

//go:generate stringer -type=Measurement
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
	SumL1
	SumL2
	SumL3

	Import
	ImportL1
	ImportL2
	ImportL3

	Export
	ExportL1
	ExportL2
	ExportL3

	Reactive
	ReactiveL1
	ReactiveL2
	ReactiveL3

	Apparent
	ApparentL1
	ApparentL2
	ApparentL3

	ApparentImport
	ApparentImportT1
	ApparentImportT2

	ReactiveImport
	ReactiveImportT1
	ReactiveImportT2

	ApparentExport
	ApparentExportT1
	ApparentExportT2

	ReactiveExport
	ReactiveExportT1
	ReactiveExportT2

	// DC
	DCCurrent
	DCVoltage
	DCPower
	HeatSinkTemp

	// Battery
	ChargeState
	BatteryVoltage
)

var iec = map[Measurement][]string{
	Frequency:        {"Frequency", "Hz"},
	Current:          {"Current", "A"},
	CurrentL1:        {"L1 Current", "A"},
	CurrentL2:        {"L2 Current", "A"},
	CurrentL3:        {"L3 Current", "A"},
	Voltage:          {"Voltage", "V"},
	VoltageL1:        {"L1 Voltage", "V"},
	VoltageL2:        {"L2 Voltage", "V"},
	VoltageL3:        {"L3 Voltage", "V"},
	Power:            {"Power", "W"},
	PowerL1:          {"L1 Power", "W"},
	PowerL2:          {"L2 Power", "W"},
	PowerL3:          {"L3 Power", "W"},
	ImportPower:      {"Import Power", "W"},
	ImportPowerL1:    {"L1 Import Power", "W"},
	ImportPowerL2:    {"L2 Import Power", "W"},
	ImportPowerL3:    {"L3 Import Power", "W"},
	ExportPower:      {"Export Power", "W"},
	ExportPowerL1:    {"L1 Export Power", "W"},
	ExportPowerL2:    {"L2 Export Power", "W"},
	ExportPowerL3:    {"L3 Export Power", "W"},
	ReactivePower:    {"Reactive Power", "var"},
	ReactivePowerL1:  {"L1 Reactive Power", "var"},
	ReactivePowerL2:  {"L2 Reactive Power", "var"},
	ReactivePowerL3:  {"L3 Reactive Power", "var"},
	ApparentPower:    {"Apparent Power", "VA"},
	ApparentPowerL1:  {"L1 Apparent Power", "VA"},
	ApparentPowerL2:  {"L2 Apparent Power", "VA"},
	ApparentPowerL3:  {"L3 Apparent Power", "VA"},
	Cosphi:           {"Cosphi"},
	CosphiL1:         {"L1 Cosphi"},
	CosphiL2:         {"L2 Cosphi"},
	CosphiL3:         {"L3 Cosphi"},
	THD:              {"Average voltage to neutral THD", "%"},
	THDL1:            {"L1 Voltage to neutral THD", "%"},
	THDL2:            {"L2 Voltage to neutral THD", "%"},
	THDL3:            {"L3 Voltage to neutral THD", "%"},
	Import:           {"Total Import", "kWh"},
	ImportL1:         {"L1 Import", "kWh"},
	ImportL2:         {"L2 Import", "kWh"},
	ImportL3:         {"L3 Import", "kWh"},
	Export:           {"Total Export", "kWh"},
	ExportL1:         {"L1 Export", "kWh"},
	ExportL2:         {"L2 Export", "kWh"},
	ExportL3:         {"L3 Export", "kWh"},
	Sum:              {"Total Sum", "kWh"},
	SumL1:            {"L1 Sum", "kWh"},
	SumL2:            {"L2 Sum", "kWh"},
	SumL3:            {"L3 Sum", "kWh"},
	Reactive:         {"Total Reactive", "kvarh"},
	ReactiveL1:       {"L1 Reactive", "kvarh"},
	ReactiveL2:       {"L2 Reactive", "kvarh"},
	ReactiveL3:       {"L3 Reactive", "kvarh"},
	Apparent:         {"Total Apparent", "kVA"},
	ApparentL1:       {"L1 Apparent", "kVA"},
	ApparentL2:       {"L2 Apparent", "kVA"},
	ApparentL3:       {"L3 Apparent", "kVA"},
	ApparentImport:   {"Apparent Import", "kVA"},
	ApparentImportT1: {"Tariff 1 Apparent Import", "kVA"},
	ApparentImportT2: {"Tariff 2 Apparent Import", "kVA"},
	ReactiveImport:   {"Reactive Import", "kvarh"},
	ReactiveImportT1: {"Tariff 1 Reactive Import", "kvarh"},
	ReactiveImportT2: {"Tariff 2 Reactive Import", "kvarh"},
	ApparentExport:   {"Apparent Export", "kVA"},
	ApparentExportT1: {"Tariff 1 Apparent Export", "kVA"},
	ApparentExportT2: {"Tariff 2 Apparent Export", "kVA"},
	ReactiveExport:   {"Reactive Export", "kvarh"},
	ReactiveExportT1: {"Tariff 1 Reactive Export", "kvarh"},
	ReactiveExportT2: {"Tariff 2 Reactive Export", "kvarh"},
	DCCurrent:        {"DC Current", "A"},
	DCVoltage:        {"DC Voltage", "V"},
	DCPower:          {"DC Power", "W"},
	HeatSinkTemp:     {"Heat Sink Temperature", "°C"},
	ChargeState:      {"Charge State", "%"},
	BatteryVoltage:   {"Battery Voltage", "V"},
}

// MarshalText implements encoding.TextMarshaler
func (m *Measurement) MarshalText() (text []byte, err error) {
	return []byte(m.String()), nil
}

// DescriptionAndUnit returns a measurements human-readable name and its unit
func (m *Measurement) DescriptionAndUnit() (string, string) {
	if details, ok := iec[*m]; ok {
		unit := ""
		description := details[0]
		if len(details) > 1 {
			unit = details[1]
		}
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
