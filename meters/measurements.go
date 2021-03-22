package meters

import (
	"fmt"
	"strings"
	"time"
)

type PrometheusInfo struct {
	Name        string
	Description string
	MetricType  MetricType
	Unit        *Unit
}

// MetricType is the type of a measurement's prometheus.Metric to be used
type MetricType int

const (
	_ MetricType = iota
	Gauge
	Counter
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

// measurement describes a Measurement itself, its unit and according prometheus.Metric type
// A measurement object is built by using the builder function newInternalMeasurement.
// Then, its fields can be set by using measurementOptions.
// Required fields are
// - Description
// - Unit
// - MetricType
// A Prometheus name and help text is "auto-generated". The format is:
// <Name>			::=	measurement_<Description>_<Unit>[_<CounterTotal>]
// <Description>	::= <measurementOption.withDescription()> | <measurementOption.WithCustomDescription()>
// <Unit>			::= <measurementOption.withUnit()>
// <CounterTotal>	::= "total" // if metric type is Counter
// E. g.:
//  Assuming a device's manufacturer is "myManufacturer":
//		newInternalMeasurement(withDescription("Frequency Test With Some Text"), withUnit(Hertz), withMetricType(Counter))
//	=> Name (before creating prometheus.Metric): "measurement_frequency_test_with_some_text_hertz_total"
//  => Description: "Measurement of Frequency Test With Some Text in Hertz"
//
// You can set custom Prometheus names and help texts by using the measurementOptions
// to override the "auto-generated" name and help text
// - WithCustomPrometheusName
// - WithCustomPrometheusDescription
// However, please make sure that the custom name conforms to Prometheus' naming conventions.
// (See https://prometheus.io/docs/practices/naming/)
// Please also note that PrometheusInfo.Name does not equal the actual name of prometheus.Metric;
// It's processed when initializing all prometheus.Metric in prometheus_metrics.UpdateMeasurementMetrics
// (see prometheus_metrics/measurements.go for more information).
type measurement struct {
	Description    string
	Unit           *Unit
	PrometheusInfo *PrometheusInfo
}

type measurementOptions func(*measurement)

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

var iec = map[Measurement]*measurement{
	Frequency:        newInternalMeasurement(withDescription("Frequency"), withPrometheusHelpText("Frequency of the power line in Hertz"), withUnit(Hertz), withMetricType(Gauge)),
	Current:          newInternalMeasurement(withDescription("Current"), withUnit(Ampere), withMetricType(Gauge)),
	CurrentL1:        newInternalMeasurement(withDescription("L1 Current"), withUnit(Ampere), withMetricType(Gauge)),
	CurrentL2:        newInternalMeasurement(withDescription("L2 Current"), withUnit(Ampere), withMetricType(Gauge)),
	CurrentL3:        newInternalMeasurement(withDescription("L3 Current"), withUnit(Ampere), withMetricType(Gauge)),
	Voltage:          newInternalMeasurement(withDescription("Voltage"), withUnit(Volt), withMetricType(Gauge)),
	VoltageL1:        newInternalMeasurement(withDescription("L1 Voltage"), withUnit(Volt), withMetricType(Gauge)),
	VoltageL2:        newInternalMeasurement(withDescription("L2 Voltage"), withUnit(Volt), withMetricType(Gauge)),
	VoltageL3:        newInternalMeasurement(withDescription("L3 Voltage"), withUnit(Volt), withMetricType(Gauge)),
	Power:            newInternalMeasurement(withDescription("Power"), withUnit(Watt), withMetricType(Gauge)),
	PowerL1:          newInternalMeasurement(withDescription("L1 Power"), withUnit(Watt), withMetricType(Gauge)),
	PowerL2:          newInternalMeasurement(withDescription("L2 Power"), withUnit(Watt), withMetricType(Gauge)),
	PowerL3:          newInternalMeasurement(withDescription("L3 Power"), withUnit(Watt), withMetricType(Gauge)),
	ImportPower:      newInternalMeasurement(withDescription("Import Power"), withUnit(Watt), withMetricType(Gauge)),
	ImportPowerL1:    newInternalMeasurement(withDescription("L1 Import Power"), withUnit(Watt), withMetricType(Gauge)),
	ImportPowerL2:    newInternalMeasurement(withDescription("L2 Import Power"), withUnit(Watt), withMetricType(Gauge)),
	ImportPowerL3:    newInternalMeasurement(withDescription("L3 Import Power"), withUnit(Watt), withMetricType(Gauge)),
	ExportPower:      newInternalMeasurement(withDescription("Export Power"), withUnit(Watt), withMetricType(Gauge)),
	ExportPowerL1:    newInternalMeasurement(withDescription("L1 Export Power"), withUnit(Watt), withMetricType(Gauge)),
	ExportPowerL2:    newInternalMeasurement(withDescription("L2 Export Power"), withUnit(Watt), withMetricType(Gauge)),
	ExportPowerL3:    newInternalMeasurement(withDescription("L3 Export Power"), withUnit(Watt), withMetricType(Gauge)),
	ReactivePower:    newInternalMeasurement(withDescription("Reactive Power"), withUnit(Var), withMetricType(Gauge)),
	ReactivePowerL1:  newInternalMeasurement(withDescription("L1 Reactive Power"), withUnit(Var), withMetricType(Gauge)),
	ReactivePowerL2:  newInternalMeasurement(withDescription("L2 Reactive Power"), withUnit(Var), withMetricType(Gauge)),
	ReactivePowerL3:  newInternalMeasurement(withDescription("L3 Reactive Power"), withUnit(Var), withMetricType(Gauge)),
	ApparentPower:    newInternalMeasurement(withDescription("Apparent Power"), withUnit(VoltAmpere), withMetricType(Gauge)),
	ApparentPowerL1:  newInternalMeasurement(withDescription("L1 Apparent Power"), withUnit(VoltAmpere), withMetricType(Gauge)),
	ApparentPowerL2:  newInternalMeasurement(withDescription("L2 Apparent Power"), withUnit(VoltAmpere), withMetricType(Gauge)),
	ApparentPowerL3:  newInternalMeasurement(withDescription("L3 Apparent Power"), withUnit(VoltAmpere), withMetricType(Gauge)),
	Cosphi:           newInternalMeasurement(withDescription("Power Factor Cosphi"), withMetricType(Gauge)),
	CosphiL1:         newInternalMeasurement(withDescription("L1 Power Factor Cosphi"), withMetricType(Gauge)),
	CosphiL2:         newInternalMeasurement(withDescription("L2 Power Factor Cosphi"), withMetricType(Gauge)),
	CosphiL3:         newInternalMeasurement(withDescription("L3 Power Factor Cosphi"), withMetricType(Gauge)),
	THD:              newInternalMeasurement(withDescription("Average voltage to neutral THD"), withUnit(Percent), withMetricType(Gauge)),
	THDL1:            newInternalMeasurement(withDescription("L1 Voltage to neutral THD"), withUnit(Percent), withMetricType(Gauge)),
	THDL2:            newInternalMeasurement(withDescription("L2 Voltage to neutral THD"), withUnit(Percent), withMetricType(Gauge)),
	THDL3:            newInternalMeasurement(withDescription("L3 Voltage to neutral THD"), withUnit(Percent), withMetricType(Gauge)),
	Sum:              newInternalMeasurement(withDescription("Total Energy Sum"), withPrometheusName("energy_sum"), withUnit(KiloWattHour), withMetricType(Counter)),
	SumT1:            newInternalMeasurement(withDescription("Tariff 1 Energy Sum"), withUnit(KiloWattHour), withMetricType(Counter)),
	SumT2:            newInternalMeasurement(withDescription("Tariff 2 Energy Sum"), withUnit(KiloWattHour), withMetricType(Counter)),
	SumL1:            newInternalMeasurement(withDescription("L1 Energy Sum"), withUnit(KiloWattHour), withMetricType(Counter)),
	SumL2:            newInternalMeasurement(withDescription("L2 Energy Sum"), withUnit(KiloWattHour), withMetricType(Counter)),
	SumL3:            newInternalMeasurement(withDescription("L3 Energy Sum"), withUnit(KiloWattHour), withMetricType(Counter)),
	Import:           newInternalMeasurement(withDescription("Total Import Energy"), withPrometheusName("energy_imported"), withUnit(KiloWattHour), withMetricType(Counter)),
	ImportT1:         newInternalMeasurement(withDescription("Tariff 1 Import Energy"), withPrometheusName("tariff_1_energy_imported"), withUnit(KiloWattHour), withMetricType(Counter)),
	ImportT2:         newInternalMeasurement(withDescription("Tariff 2 Import Energy"), withPrometheusName("tariff_2_energy_imported"), withUnit(KiloWattHour), withMetricType(Counter)),
	ImportL1:         newInternalMeasurement(withDescription("L1 Import Energy"), withPrometheusName("l1_energy_imported"), withUnit(KiloWattHour), withMetricType(Counter)),
	ImportL2:         newInternalMeasurement(withDescription("L2 Import Energy"), withPrometheusName("l2_energy_imported"), withUnit(KiloWattHour), withMetricType(Counter)),
	ImportL3:         newInternalMeasurement(withDescription("L3 Import Energy"), withPrometheusName("l3_energy_imported"), withUnit(KiloWattHour), withMetricType(Counter)),
	Export:           newInternalMeasurement(withDescription("Total Export Energy"), withPrometheusName("energy_exported"), withUnit(KiloWattHour), withMetricType(Counter)),
	ExportT1:         newInternalMeasurement(withDescription("Tariff 1 Export Energy"), withPrometheusName("tariff_1_energy_exported"), withUnit(KiloWattHour), withMetricType(Counter)),
	ExportT2:         newInternalMeasurement(withDescription("Tariff 2 Export Energy"), withPrometheusName("tariff_2_energy_exported"), withUnit(KiloWattHour), withMetricType(Counter)),
	ExportL1:         newInternalMeasurement(withDescription("L1 Export Energy"), withPrometheusName("l1_energy_exported"), withUnit(KiloWattHour), withMetricType(Counter)),
	ExportL2:         newInternalMeasurement(withDescription("L2 Export Energy"), withPrometheusName("l2_energy_exported"), withUnit(KiloWattHour), withMetricType(Counter)),
	ExportL3:         newInternalMeasurement(withDescription("L3 Export Energy"), withPrometheusName("l3_energy_exported"), withUnit(KiloWattHour), withMetricType(Counter)),
	ReactiveSum:      newInternalMeasurement(withDescription("Total Reactive Energy"), withPrometheusName("reactive_energy"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveSumT1:    newInternalMeasurement(withDescription("Tariff 1 Reactive Energy"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveSumT2:    newInternalMeasurement(withDescription("Tariff 2 Reactive Energy"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveSumL1:    newInternalMeasurement(withDescription("L1 Reactive Energy"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveSumL2:    newInternalMeasurement(withDescription("L2 Reactive Energy"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveSumL3:    newInternalMeasurement(withDescription("L3 Reactive Energy"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveImport:   newInternalMeasurement(withDescription("Reactive Import Energy"), withPrometheusName("reactive_energy_imported"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveImportT1: newInternalMeasurement(withDescription("Tariff 1 Reactive Import Energy"), withPrometheusName("tariff_2_reactive_energy_imported"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveImportT2: newInternalMeasurement(withDescription("Tariff 2 Reactive Import Energy"), withPrometheusName("tariff_1_reactive_energy_imported"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveImportL1: newInternalMeasurement(withDescription("L1 Reactive Import Energy"), withPrometheusName("l1_reactive_energy_imported"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveImportL2: newInternalMeasurement(withDescription("L2 Reactive Import Energy"), withPrometheusName("l2_reactive_energy_imported"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveImportL3: newInternalMeasurement(withDescription("L3 Reactive Import Energy"), withPrometheusName("l3_reactive_energy_imported"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveExport:   newInternalMeasurement(withDescription("Reactive Export Energy"), withPrometheusName("reactive_energy_exported"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveExportT1: newInternalMeasurement(withDescription("Tariff 1 Reactive Export Energy"), withPrometheusName("tariff_1_reactive_energy_exported"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveExportT2: newInternalMeasurement(withDescription("Tariff 2 Reactive Export Energy"), withPrometheusName("tariff_2_reactive_energy_exported"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveExportL1: newInternalMeasurement(withDescription("L1 Reactive Export Energy"), withPrometheusName("l1_reactive_energy_exported"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveExportL2: newInternalMeasurement(withDescription("L2 Reactive Export Energy"), withPrometheusName("l2_reactive_energy_exported"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveExportL3: newInternalMeasurement(withDescription("L3 Reactive Export Energy"), withPrometheusName("l3_reactive_energy_exported"), withUnit(KiloVarHour), withMetricType(Counter)),
	DCCurrent:        newInternalMeasurement(withDescription("DC Current"), withUnit(Ampere), withMetricType(Gauge)),
	DCVoltage:        newInternalMeasurement(withDescription("DC Voltage"), withUnit(Volt), withMetricType(Gauge)),
	DCPower:          newInternalMeasurement(withDescription("DC Power"), withUnit(Watt), withMetricType(Gauge)),
	HeatSinkTemp:     newInternalMeasurement(withDescription("Heat Sink Temperature"), withUnit(DegreeCelsius), withMetricType(Gauge)),
	DCCurrentS1:      newInternalMeasurement(withDescription("String 1 Current"), withUnit(Ampere), withMetricType(Gauge)),
	DCVoltageS1:      newInternalMeasurement(withDescription("String 1 Voltage"), withUnit(Volt), withMetricType(Gauge)),
	DCPowerS1:        newInternalMeasurement(withDescription("String 1 Power"), withUnit(Watt), withMetricType(Gauge)),
	DCEnergyS1:       newInternalMeasurement(withDescription("String 1 Generation"), withUnit(KiloWattHour), withMetricType(Counter)),
	DCCurrentS2:      newInternalMeasurement(withDescription("String 2 Current"), withUnit(Ampere), withMetricType(Gauge)),
	DCVoltageS2:      newInternalMeasurement(withDescription("String 2 Voltage"), withUnit(Volt), withMetricType(Gauge)),
	DCPowerS2:        newInternalMeasurement(withDescription("String 2 Power"), withUnit(Watt), withMetricType(Gauge)),
	DCEnergyS2:       newInternalMeasurement(withDescription("String 2 Generation"), withUnit(KiloWattHour), withMetricType(Counter)),
	DCCurrentS3:      newInternalMeasurement(withDescription("String 3 Current"), withUnit(Ampere), withMetricType(Gauge)),
	DCVoltageS3:      newInternalMeasurement(withDescription("String 3 Voltage"), withUnit(Volt), withMetricType(Gauge)),
	DCPowerS3:        newInternalMeasurement(withDescription("String 3 Power"), withUnit(Watt), withMetricType(Gauge)),
	DCEnergyS3:       newInternalMeasurement(withDescription("String 3 Generation"), withUnit(KiloWattHour), withMetricType(Counter)),
	DCCurrentS4:      {"String 4 Current", "A"},
	DCVoltageS4:      {"String 4 Voltage", "V"},
	DCPowerS4:        {"String 4 Power", "W"},
	DCEnergyS4:       {"String 4 Generation", "kWh"},
	ChargeState:      newInternalMeasurement(withDescription("Charge State"), withUnit(Percent), withMetricType(Gauge)),
	BatteryVoltage:   newInternalMeasurement(withDescription("Battery Voltage"), withUnit(Volt), withMetricType(Gauge)),
	PhaseAngle:       newInternalMeasurement(withDescription("Phase Angle"), withUnit(Degree), withMetricType(Gauge)),
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

func (m *Measurement) Unit() *Unit {
	if details, ok := iec[*m]; ok {
		return details.Unit
	}

	return nil
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
func (m *Measurement) PrometheusMetricType() MetricType {
	if measurement, ok := iec[*m]; ok {
		return measurement.PrometheusInfo.MetricType
	}
	return 0
}

// PrometheusDescription returns a description text appropriate for prometheus.Metric
func (m *Measurement) PrometheusDescription() string {
	if measurement, ok := iec[*m]; ok {
		return measurement.PrometheusInfo.Description
	}
	return ""
}

// PrometheusName returns a name and its associated unit for Prometheus counters
func (m *Measurement) PrometheusName() string {
	if details, ok := iec[*m]; ok {
		return details.PrometheusInfo.Name
	}
	return ""
}

// measurementOptions functional parameters for internal measurement object generation
//
// newInternalMeasurement generates an internal measurement object based on passed options
func newInternalMeasurement(opts ...measurementOptions) *measurement {
	promInfo := &PrometheusInfo{}
	m := &measurement{}
	m.PrometheusInfo = promInfo

	for _, opt := range opts {
		opt(m)
	}

	if m.Unit == nil {
		withUnit(NoUnit)(m)
	} else {
		elementaryUnit, _ := ConvertValueToElementaryUnit(*m.Unit, 0.0)
		m.PrometheusInfo.Unit = &elementaryUnit
	}

	if m.PrometheusInfo.Description == "" {
		withGenericPrometheusHelpText()(m)
	}

	if m.PrometheusInfo.Name == "" {
		withGenericPrometheusName()(m)
	} else {
		if m.PrometheusInfo.MetricType == Counter {
			m.PrometheusInfo.Name = generatePrometheusName(m.PrometheusInfo.Name, m.Unit.PrometheusName()+"_total")
		} else {
			m.PrometheusInfo.Name = generatePrometheusName(m.PrometheusInfo.Name, m.Unit.PrometheusName())
		}
	}

	return m
}

// withPrometheusHelpText enables setting a Prometheus description of a Measurement
func withPrometheusHelpText(description string) measurementOptions {
	return func(m *measurement) {
		m.PrometheusInfo.Description = description
	}
}

// withGenericPrometheusHelpText sets the Prometheus description to a generated, more generic format
func withGenericPrometheusHelpText() measurementOptions {
	return func(m *measurement) {
		m.PrometheusInfo.Description = generatePrometheusHelpText(m.Description, m.Unit.FullName())
	}
}

// withUnit sets the Unit of a Measurement
// If u is nil, the unit will be set to NoUnit
func withUnit(u Unit) measurementOptions {
	return func(m *measurement) {
		m.Unit = &u
	}
}

func withPrometheusName(name string) measurementOptions {
	return func(m *measurement) {
		m.PrometheusInfo.Name = name
	}
}

func withGenericPrometheusName() measurementOptions {
	return func(m *measurement) {
		if m.PrometheusInfo.MetricType == Counter {
			m.PrometheusInfo.Name = generatePrometheusName(m.Description, m.Unit.PrometheusName()+"_total")
		} else {
			m.PrometheusInfo.Name = generatePrometheusName(m.Description, m.Unit.PrometheusName())
		}
	}
}

func withMetricType(metricType MetricType) measurementOptions {
	return func(m *measurement) {
		m.PrometheusInfo.MetricType = metricType
	}
}

func withDescription(description string) measurementOptions {
	return func(m *measurement) {
		m.Description = description
	}
}

func generatePrometheusHelpText(description string, unit string) string {
	if unit != "" {
		return fmt.Sprintf("%s in %s", description, unit)
	} else {
		return fmt.Sprintf("%s", description)
	}
}

func generatePrometheusName(name string, unit string) string {
	measurementName := strings.ToLower(name)
	prometheusUnit := strings.ToLower(unit)

	measurementName = strings.Trim(strings.ReplaceAll(strings.ToLower(measurementName), " ", "_"), "_")

	return strings.Trim( // Trim trailing underscore (e. g. when unit string is empty)
		strings.Join(
			[]string{"measurement", measurementName, prometheusUnit},
			"_",
		),
		"_",
	)
}

// ConvertValueToElementaryUnit converts a sourceUnit and a sourceValue to their elementary unit if possible
// Otherwise, sourceUnit and sourceValue are returned again.
func ConvertValueToElementaryUnit(sourceUnit Unit, sourceValue float64) (Unit, float64) {
	switch sourceUnit {
	case KiloWattHour:
		fallthrough
	case KiloVarHour:
		return Joule, sourceValue * 1_000 * 3_600
	}

	return sourceUnit, sourceValue
}
