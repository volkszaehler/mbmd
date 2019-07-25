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

//go:generate go run golang.org/x/tools/cmd/stringer -type=Measurement
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

	Power
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

	ActivePower
	ActivePowerL1
	ActivePowerL2
	ActivePowerL3

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
	Net
	NetL1
	NetL2
	NetL3

	Active
	ActiveNet
	ActiveNetL1
	ActiveNetL2
	ActiveNetL3

	Reactive
	ReactiveNet
	ReactiveNetL1
	ReactiveNetL2
	ReactiveNetL3

	Import
	ImportL1
	ImportL2
	ImportL3

	Export
	ExportL1
	ExportL2
	ExportL3

	ActiveImportT1
	ActiveImportT2

	ReactiveImportT1
	ReactiveImportT2

	ActiveExportT1
	ActiveExportT2

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
	ActivePower:      {"Active Power", "W"},
	ActivePowerL1:    {"L1 Active Power", "W"},
	ActivePowerL2:    {"L2 Active Power", "W"},
	ActivePowerL3:    {"L3 Active Power", "W"},
	ReactivePower:    {"Reactive Power", "W"},
	ReactivePowerL1:  {"L1 Reactive Power", "W"},
	ReactivePowerL2:  {"L2 Reactive Power", "W"},
	ReactivePowerL3:  {"L3 Reactive Power", "W"},
	ApparentPower:    {"Apparent Power", "W"},
	ApparentPowerL1:  {"L1 Apparent Power", "W"},
	ApparentPowerL2:  {"L2 Apparent Power", "W"},
	ApparentPowerL3:  {"L3 Apparent Power", "W"},
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
	Net:              {"Total  Net", "kWh"},
	NetL1:            {"L1 Net", "kWh"},
	NetL2:            {"L2 Net", "kWh"},
	NetL3:            {"L3 Net", "kWh"},
	Active:           {"Total Active", "kWh"},
	ActiveNet:        {"Total Active Net", "kWh"},
	ActiveNetL1:      {"L1 Active Net", "kWh"},
	ActiveNetL2:      {"L2 Active Net", "kWh"},
	ActiveNetL3:      {"L3 Active Net", "kWh"},
	Reactive:         {"Total Reactive", "kWh"},
	ReactiveNet:      {"Reactive Net", "kWh"},
	ReactiveNetL1:    {"L1 Reactive Net", "kWh"},
	ReactiveNetL2:    {"L2 Reactive Net", "kWh"},
	ReactiveNetL3:    {"L3 Reactive Net", "kWh"},
	ActiveImportT1:   {"Tariff 1 Active Import", "kWh"},
	ActiveImportT2:   {"Tariff 2 Active Import", "kWh"},
	ReactiveImportT1: {"Tariff 1 Reactive Import", "kWh"},
	ReactiveImportT2: {"Tariff 2 Reactive Import", "kWh"},
	ActiveExportT1:   {"Tariff 1 Active Export", "kWh"},
	ActiveExportT2:   {"Tariff 2 Active Export", "kWh"},
	ReactiveExportT1: {"Tariff 1 Reactive Export", "kWh"},
	ReactiveExportT2: {"Tariff 2 Reactive Export", "kWh"},
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
