package meters

import (
	"fmt"
	"github.com/volkszaehler/mbmd/meters/units"
	"log"
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

var iec = map[Measurement]*measurement{
	Frequency:        newInternalMeasurement(withDescription("Frequency"), withPrometheusHelpText("Frequency of the power line in Hertz"), withUnit(units.Hertz), withMetricType(Gauge)),
	Current:          newInternalMeasurement(withDescription("Current"), withUnit(units.Ampere), withMetricType(Gauge)),
	CurrentL1:        newInternalMeasurement(withDescription("L1 Current"), withUnit(units.Ampere), withMetricType(Gauge)),
	CurrentL2:        newInternalMeasurement(withDescription("L2 Current"), withUnit(units.Ampere), withMetricType(Gauge)),
	CurrentL3:        newInternalMeasurement(withDescription("L3 Current"), withUnit(units.Ampere), withMetricType(Gauge)),
	Voltage:          newInternalMeasurement(withDescription("Voltage"), withUnit(units.Volt), withMetricType(Gauge)),
	VoltageL1:        newInternalMeasurement(withDescription("L1 Voltage"), withUnit(units.Volt), withMetricType(Gauge)),
	VoltageL2:        newInternalMeasurement(withDescription("L2 Voltage"), withUnit(units.Volt), withMetricType(Gauge)),
	VoltageL3:        newInternalMeasurement(withDescription("L3 Voltage"), withUnit(units.Volt), withMetricType(Gauge)),
	Power:            newInternalMeasurement(withDescription("Power"), withUnit(units.Watt), withMetricType(Gauge)),
	PowerL1:          newInternalMeasurement(withDescription("L1 Power"), withUnit(units.Watt), withMetricType(Gauge)),
	PowerL2:          newInternalMeasurement(withDescription("L2 Power"), withUnit(units.Watt), withMetricType(Gauge)),
	PowerL3:          newInternalMeasurement(withDescription("L3 Power"), withUnit(units.Watt), withMetricType(Gauge)),
	ImportPower:      newInternalMeasurement(withDescription("Import Power"), withUnit(units.Watt), withMetricType(Gauge)),
	ImportPowerL1:    newInternalMeasurement(withDescription("L1 Import Power"), withUnit(units.Watt), withMetricType(Gauge)),
	ImportPowerL2:    newInternalMeasurement(withDescription("L2 Import Power"), withUnit(units.Watt), withMetricType(Gauge)),
	ImportPowerL3:    newInternalMeasurement(withDescription("L3 Import Power"), withUnit(units.Watt), withMetricType(Gauge)),
	ExportPower:      newInternalMeasurement(withDescription("Export Power"), withUnit(units.Watt), withMetricType(Gauge)),
	ExportPowerL1:    newInternalMeasurement(withDescription("L1 Export Power"), withUnit(units.Watt), withMetricType(Gauge)),
	ExportPowerL2:    newInternalMeasurement(withDescription("L2 Export Power"), withUnit(units.Watt), withMetricType(Gauge)),
	ExportPowerL3:    newInternalMeasurement(withDescription("L3 Export Power"), withUnit(units.Watt), withMetricType(Gauge)),
	ReactivePower:    newInternalMeasurement(withDescription("Reactive Power"), withUnit(units.Var), withMetricType(Gauge)),
	ReactivePowerL1:  newInternalMeasurement(withDescription("L1 Reactive Power"), withUnit(units.Var), withMetricType(Gauge)),
	ReactivePowerL2:  newInternalMeasurement(withDescription("L2 Reactive Power"), withUnit(units.Var), withMetricType(Gauge)),
	ReactivePowerL3:  newInternalMeasurement(withDescription("L3 Reactive Power"), withUnit(units.Var), withMetricType(Gauge)),
	ApparentPower:    newInternalMeasurement(withDescription("Apparent Power"), withUnit(units.Voltampere), withMetricType(Gauge)),
	ApparentPowerL1:  newInternalMeasurement(withDescription("L1 Apparent Power"), withUnit(units.Voltampere), withMetricType(Gauge)),
	ApparentPowerL2:  newInternalMeasurement(withDescription("L2 Apparent Power"), withUnit(units.Voltampere), withMetricType(Gauge)),
	ApparentPowerL3:  newInternalMeasurement(withDescription("L3 Apparent Power"), withUnit(units.Voltampere), withMetricType(Gauge)),
	Cosphi:           newInternalMeasurement(withDescription("Power Factor Cosphi"), withMetricType(Gauge)),
	CosphiL1:         newInternalMeasurement(withDescription("L1 Power Factor Cosphi"), withMetricType(Gauge)),
	CosphiL2:         newInternalMeasurement(withDescription("L2 Power Factor Cosphi"), withMetricType(Gauge)),
	CosphiL3:         newInternalMeasurement(withDescription("L3 Power Factor Cosphi"), withMetricType(Gauge)),
	THD:              newInternalMeasurement(withDescription("Average voltage to neutral THD"), withUnit(units.Percent), withMetricType(Gauge)),
	THDL1:            newInternalMeasurement(withDescription("L1 Voltage to neutral THD"), withUnit(units.Percent), withMetricType(Gauge)),
	THDL2:            newInternalMeasurement(withDescription("L2 Voltage to neutral THD"), withUnit(units.Percent), withMetricType(Gauge)),
	THDL3:            newInternalMeasurement(withDescription("L3 Voltage to neutral THD"), withUnit(units.Percent), withMetricType(Gauge)),
	Sum:              newInternalMeasurement(withDescription("Total Energy Sum"), withPrometheusName("energy_sum"), withUnit(units.KiloWattHour), withMetricType(Counter)),
	SumT1:            newInternalMeasurement(withDescription("Tariff 1 Energy Sum"), withUnit(units.KiloWattHour), withMetricType(Counter)),
	SumT2:            newInternalMeasurement(withDescription("Tariff 2 Energy Sum"), withUnit(units.KiloWattHour), withMetricType(Counter)),
	SumL1:            newInternalMeasurement(withDescription("L1 Energy Sum"), withUnit(units.KiloWattHour), withMetricType(Counter)),
	SumL2:            newInternalMeasurement(withDescription("L2 Energy Sum"), withUnit(units.KiloWattHour), withMetricType(Counter)),
	SumL3:            newInternalMeasurement(withDescription("L3 Energy Sum"), withUnit(units.KiloWattHour), withMetricType(Counter)),
	Import:           newInternalMeasurement(withDescription("Total Import Energy"), withPrometheusName("energy_imported"), withUnit(units.KiloWattHour), withMetricType(Counter)),
	ImportT1:         newInternalMeasurement(withDescription("Tariff 1 Import Energy"), withPrometheusName("tariff_1_energy_imported"), withUnit(units.KiloWattHour), withMetricType(Counter)),
	ImportT2:         newInternalMeasurement(withDescription("Tariff 2 Import Energy"), withPrometheusName("tariff_2_energy_imported"), withUnit(units.KiloWattHour), withMetricType(Counter)),
	ImportL1:         newInternalMeasurement(withDescription("L1 Import Energy"), withPrometheusName("l1_energy_imported"), withUnit(units.KiloWattHour), withMetricType(Counter)),
	ImportL2:         newInternalMeasurement(withDescription("L2 Import Energy"), withPrometheusName("l2_energy_imported"), withUnit(units.KiloWattHour), withMetricType(Counter)),
	ImportL3:         newInternalMeasurement(withDescription("L3 Import Energy"), withPrometheusName("l3_energy_imported"), withUnit(units.KiloWattHour), withMetricType(Counter)),
	Export:           newInternalMeasurement(withDescription("Total Export Energy"), withPrometheusName("energy_exported"), withUnit(units.KiloWattHour), withMetricType(Counter)),
	ExportT1:         newInternalMeasurement(withDescription("Tariff 1 Export Energy"), withPrometheusName("tariff_1_energy_exported"), withUnit(units.KiloWattHour), withMetricType(Counter)),
	ExportT2:         newInternalMeasurement(withDescription("Tariff 2 Export Energy"), withPrometheusName("tariff_2_energy_exported"), withUnit(units.KiloWattHour), withMetricType(Counter)),
	ExportL1:         newInternalMeasurement(withDescription("L1 Export Energy"), withPrometheusName("l1_energy_exported"), withUnit(units.KiloWattHour), withMetricType(Counter)),
	ExportL2:         newInternalMeasurement(withDescription("L2 Export Energy"), withPrometheusName("l2_energy_exported"), withUnit(units.KiloWattHour), withMetricType(Counter)),
	ExportL3:         newInternalMeasurement(withDescription("L3 Export Energy"), withPrometheusName("l3_energy_exported"), withUnit(units.KiloWattHour), withMetricType(Counter)),
	ReactiveSum:      newInternalMeasurement(withDescription("Total Reactive Energy"), withPrometheusName("reactive_energy"), withUnit(units.KiloVarHour), withMetricType(Counter)),
	ReactiveSumT1:    newInternalMeasurement(withDescription("Tariff 1 Reactive Energy"), withUnit(units.KiloVarHour), withMetricType(Counter)),
	ReactiveSumT2:    newInternalMeasurement(withDescription("Tariff 2 Reactive Energy"), withUnit(units.KiloVarHour), withMetricType(Counter)),
	ReactiveSumL1:    newInternalMeasurement(withDescription("L1 Reactive Energy"), withUnit(units.KiloVarHour), withMetricType(Counter)),
	ReactiveSumL2:    newInternalMeasurement(withDescription("L2 Reactive Energy"), withUnit(units.KiloVarHour), withMetricType(Counter)),
	ReactiveSumL3:    newInternalMeasurement(withDescription("L3 Reactive Energy"), withUnit(units.KiloVarHour), withMetricType(Counter)),
	ReactiveImport:   newInternalMeasurement(withDescription("Reactive Import Energy"), withPrometheusName("reactive_energy_imported"), withUnit(units.KiloVarHour), withMetricType(Counter)),
	ReactiveImportT1: newInternalMeasurement(withDescription("Tariff 1 Reactive Import Energy"), withPrometheusName("tariff_2_reactive_energy_imported"), withUnit(units.KiloVarHour), withMetricType(Counter)),
	ReactiveImportT2: newInternalMeasurement(withDescription("Tariff 2 Reactive Import Energy"), withPrometheusName("tariff_1_reactive_energy_imported"), withUnit(units.KiloVarHour), withMetricType(Counter)),
	ReactiveImportL1: newInternalMeasurement(withDescription("L1 Reactive Import Energy"), withPrometheusName("l1_reactive_energy_imported"), withUnit(units.KiloVarHour), withMetricType(Counter)),
	ReactiveImportL2: newInternalMeasurement(withDescription("L2 Reactive Import Energy"), withPrometheusName("l2_reactive_energy_imported"), withUnit(units.KiloVarHour), withMetricType(Counter)),
	ReactiveImportL3: newInternalMeasurement(withDescription("L3 Reactive Import Energy"), withPrometheusName("l3_reactive_energy_imported"), withUnit(units.KiloVarHour), withMetricType(Counter)),
	ReactiveExport:   newInternalMeasurement(withDescription("Reactive Export Energy"), withPrometheusName("reactive_energy_exported"), withUnit(units.KiloVarHour), withMetricType(Counter)),
	ReactiveExportT1: newInternalMeasurement(withDescription("Tariff 1 Reactive Export Energy"), withPrometheusName("tariff_1_reactive_energy_exported"), withUnit(units.KiloVarHour), withMetricType(Counter)),
	ReactiveExportT2: newInternalMeasurement(withDescription("Tariff 2 Reactive Export Energy"), withPrometheusName("tariff_2_reactive_energy_exported"), withUnit(units.KiloVarHour), withMetricType(Counter)),
	ReactiveExportL1: newInternalMeasurement(withDescription("L1 Reactive Export Energy"), withPrometheusName("l1_reactive_energy_exported"), withUnit(units.KiloVarHour), withMetricType(Counter)),
	ReactiveExportL2: newInternalMeasurement(withDescription("L2 Reactive Export Energy"), withPrometheusName("l2_reactive_energy_exported"), withUnit(units.KiloVarHour), withMetricType(Counter)),
	ReactiveExportL3: newInternalMeasurement(withDescription("L3 Reactive Export Energy"), withPrometheusName("l3_reactive_energy_exported"), withUnit(units.KiloVarHour), withMetricType(Counter)),
	DCCurrent:        newInternalMeasurement(withDescription("DC Current"), withUnit(units.Ampere), withMetricType(Gauge)),
	DCVoltage:        newInternalMeasurement(withDescription("DC Voltage"), withUnit(units.Volt), withMetricType(Gauge)),
	DCPower:          newInternalMeasurement(withDescription("DC Power"), withUnit(units.Watt), withMetricType(Gauge)),
	HeatSinkTemp:     newInternalMeasurement(withDescription("Heat Sink Temperature"), withUnit(units.DegreeCelsius), withMetricType(Gauge)),
	DCCurrentS1:      newInternalMeasurement(withDescription("String 1 Current"), withUnit(units.Ampere), withMetricType(Gauge)),
	DCVoltageS1:      newInternalMeasurement(withDescription("String 1 Voltage"), withUnit(units.Volt), withMetricType(Gauge)),
	DCPowerS1:        newInternalMeasurement(withDescription("String 1 Power"), withUnit(units.Watt), withMetricType(Gauge)),
	DCEnergyS1:       newInternalMeasurement(withDescription("String 1 Generation"), withPrometheusName("string_1_energy_generated"), withUnit(units.KiloWattHour), withMetricType(Counter)),
	DCCurrentS2:      newInternalMeasurement(withDescription("String 2 Current"), withUnit(units.Ampere), withMetricType(Gauge)),
	DCVoltageS2:      newInternalMeasurement(withDescription("String 2 Voltage"), withUnit(units.Volt), withMetricType(Gauge)),
	DCPowerS2:        newInternalMeasurement(withDescription("String 2 Power"), withUnit(units.Watt), withMetricType(Gauge)),
	DCEnergyS2:       newInternalMeasurement(withDescription("String 2 Generation"), withPrometheusName("string_2_energy_generated"), withUnit(units.KiloWattHour), withMetricType(Counter)),
	DCCurrentS3:      newInternalMeasurement(withDescription("String 3 Current"), withUnit(units.Ampere), withMetricType(Gauge)),
	DCVoltageS3:      newInternalMeasurement(withDescription("String 3 Voltage"), withUnit(units.Volt), withMetricType(Gauge)),
	DCPowerS3:        newInternalMeasurement(withDescription("String 3 Power"), withUnit(units.Watt), withMetricType(Gauge)),
	DCEnergyS3:       newInternalMeasurement(withDescription("String 3 Generation"), withPrometheusName("string_3_energy_generated"), withUnit(units.KiloWattHour), withMetricType(Counter)),
	ChargeState:      newInternalMeasurement(withDescription("Charge State"), withUnit(units.Percent), withMetricType(Gauge)),
	BatteryVoltage:   newInternalMeasurement(withDescription("Battery Voltage"), withUnit(units.Volt), withMetricType(Gauge)),
	PhaseAngle:       newInternalMeasurement(withDescription("Phase Angle"), withUnit(units.Degree), withMetricType(Gauge)),
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

func (m *Measurement) Unit() *units.Unit {
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

// PrometheusMetricType returns the Measurement's associated prometheus.Metric type
func (m *Measurement) PrometheusMetricType() MetricType {
	if measurement, ok := iec[*m]; ok {
		return measurement.PrometheusInfo.MetricType
	}
	return 0
}

// PrometheusHelpText returns a description text appropriate for prometheus.Metric
func (m *Measurement) PrometheusHelpText() string {
	if measurement, ok := iec[*m]; ok {
		return measurement.PrometheusInfo.HelpText
	}
	return ""
}

// PrometheusName returns a name and its associated unit for Prometheus counters0
func (m *Measurement) PrometheusName() string {
	if details, ok := iec[*m]; ok {
		return details.PrometheusInfo.Name
	}
	return ""
}

// measurement describes a Measurement itself, its unit and according prometheus.Metric type
// A measurement object is built by using the builder function newInternalMeasurement.
//
// A Prometheus name and help text is "auto-generated". The format is:
// <Name>			::=	measurement_<HelpText>_<Unit>[_<CounterTotal>]
// <HelpText>		::= <measurementOption.withDescription()> | <measurementOption.WithCustomDescription()>
// <Unit>			::= <measurementOption.withUnit()> // Elementary unit!
// <CounterTotal>	::= "total" // if metric type is Counter
// E. g.:
//		newInternalMeasurement(withDescription("Frequency Test With Some Text"), withUnit(Hertz), withMetricType(Counter))
//	=> Name (before creating prometheus.Metric): "measurement_frequency_test_with_some_text_hertz_total"
//  => Description: "Frequency Test With Some Text in Hertz"
//
// In Prometheus context: If Unit is set, then it will be automatically converted to its elementary unit.
//							(see units.ConvertValueToElementaryUnit)
//
// You can set custom Prometheus names and help texts by using the measurementOptions
// to override the "auto-generated" name and help text
// - withPrometheusName
// - withPrometheusHelpText
// However, please make sure that the custom name conforms to Prometheus' naming conventions.
// (See https://prometheus.io/docs/practices/naming/)
// Please also note that PrometheusInfo.Name does not equal the actual name of prometheus.Metric;
// It's a partial name that will be concatenated together with a globally defined namespace (and for measurements with `measurement`)
// (see also prometheus_metrics.CreateMeasurementMetrics and generatePrometheusName)
type measurement struct {
	Description    string
	Unit           *units.Unit
	PrometheusInfo *PrometheusInfo
}

// measurementOptions are used in newInternalMeasurement
type measurementOptions func(*measurement)

// PrometheusInfo carries Prometheus relevant information for e. g. creating metrics
type PrometheusInfo struct {
	Name       string
	HelpText   string
	MetricType MetricType
	Unit       *units.Unit
}

// MetricType is the type of a measurement's prometheus.Metric to be used
type MetricType int

const (
	_ MetricType = iota
	Gauge
	Counter
)

// newInternalMeasurement generates an internal measurement object based on passed options
//
// If one of the following options are not passed:
//	- withDescription
//	- withMetricType
// the app will panic!
func newInternalMeasurement(opts ...measurementOptions) *measurement {
	promInfo := &PrometheusInfo{}
	m := &measurement{}
	m.PrometheusInfo = promInfo

	for _, opt := range opts {
		opt(m)
	}

	if m.Description == "" || m.PrometheusInfo.MetricType == 0 {
		log.Fatalf(
			"Cannot create internal `measurement` because either Description or MetricType is empty."+
				"(Description: %v, MetricType: %v)",
			m.Description,
			m.PrometheusInfo.MetricType,
		)
	}

	if m.Unit == nil {
		withUnit(units.NoUnit)(m)
	}

	if m.PrometheusInfo.HelpText == "" {
		withGenericPrometheusHelpText()(m)
	}

	if m.PrometheusInfo.Name == "" {
		m.PrometheusInfo.Name = generatePrometheusName(m.Description, m.PrometheusInfo.Unit, m.PrometheusInfo.MetricType)
	} else {
		m.PrometheusInfo.Name = generatePrometheusName(m.PrometheusInfo.Name, m.PrometheusInfo.Unit, m.PrometheusInfo.MetricType)
	}

	return m
}

// withPrometheusHelpText enables setting a Prometheus description of a Measurement
func withPrometheusHelpText(description string) measurementOptions {
	return func(m *measurement) {
		m.PrometheusInfo.HelpText = description
	}
}

// withGenericPrometheusHelpText sets the Prometheus description to a generated, more generic format
func withGenericPrometheusHelpText() measurementOptions {
	return func(m *measurement) {
		m.PrometheusInfo.HelpText = generatePrometheusHelpText(m.Description, m.PrometheusInfo.Unit)
	}
}

// withUnit sets the Unit of a Measurement
// If u is nil, the unit will be set to NoUnit
func withUnit(u units.Unit) measurementOptions {
	return func(m *measurement) {
		m.Unit = &u

		elementaryUnit, _ := units.ConvertValueToElementaryUnit(*m.Unit, 0.0)
		m.PrometheusInfo.Unit = &elementaryUnit
	}
}

func withPrometheusName(name string) measurementOptions {
	return func(m *measurement) {
		m.PrometheusInfo.Name = name
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

func generatePrometheusHelpText(description string, unit *units.Unit) string {
	if unit != nil && *unit != units.NoUnit {
		_, pluralForm := unit.Name()
		return fmt.Sprintf("%s in %s", description, pluralForm)
	} else {
		return fmt.Sprintf("%s", description)
	}
}

func generatePrometheusName(name string, unit *units.Unit, metricType MetricType) string {
	measurementName := strings.ToLower(name)
	prometheusUnit := strings.ToLower(unit.PrometheusForm())

	measurementName = strings.Trim(strings.ReplaceAll(strings.ToLower(measurementName), " ", "_"), "_")

	var counterSuffix string
	if metricType == Counter {
		counterSuffix = "total"
	}

	return strings.Trim( // Trim trailing underscore (e. g. when unit string is empty)
		strings.Join(
			[]string{"measurement", measurementName, prometheusUnit, counterSuffix},
			"_",
		),
		"_",
	)
}
