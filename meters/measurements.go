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

	// Battery
	ChargeState
	BatteryVoltage

	PhaseAngle
)

var iec = map[Measurement]*measurement{
	Frequency:        newInternalMeasurement(withDescription("Frequency"), withUnit(Hertz), withMetricType(Gauge)),
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
	Cosphi:           newInternalMeasurement(withDescription("Cosphi"), withMetricType(Gauge)),
	CosphiL1:         newInternalMeasurement(withDescription("L1 Cosphi"), withMetricType(Gauge)),
	CosphiL2:         newInternalMeasurement(withDescription("L2 Cosphi"), withMetricType(Gauge)),
	CosphiL3:         newInternalMeasurement(withDescription("L3 Cosphi"), withMetricType(Gauge)),
	THD:              newInternalMeasurement(withDescription("Average voltage to neutral THD"), withUnit(Percent), withMetricType(Gauge)),
	THDL1:            newInternalMeasurement(withDescription("L1 Voltage to neutral THD"), withUnit(Percent), withMetricType(Gauge)),
	THDL2:            newInternalMeasurement(withDescription("L2 Voltage to neutral THD"), withUnit(Percent), withMetricType(Gauge)),
	THDL3:            newInternalMeasurement(withDescription("L3 Voltage to neutral THD"), withUnit(Percent), withMetricType(Gauge)),
	Sum:              newInternalMeasurement(withDescription("Total Sum"), withUnit(KiloWattHour), withUnitInPrometheus(Joule), withUnitInPrometheus(Joule), withMetricType(Counter)),
	SumT1:            newInternalMeasurement(withDescription("Tariff 1 Sum"), withUnit(KiloWattHour), withUnitInPrometheus(Joule), withMetricType(Counter)),
	SumT2:            newInternalMeasurement(withDescription("Tariff 2 Sum"), withUnit(KiloWattHour), withUnitInPrometheus(Joule), withMetricType(Counter)),
	SumL1:            newInternalMeasurement(withDescription("L1 Sum"), withUnit(KiloWattHour), withUnitInPrometheus(Joule), withMetricType(Counter)),
	SumL2:            newInternalMeasurement(withDescription("L2 Sum"), withUnit(KiloWattHour), withUnitInPrometheus(Joule), withMetricType(Counter)),
	SumL3:            newInternalMeasurement(withDescription("L3 Sum"), withUnit(KiloWattHour), withUnitInPrometheus(Joule), withMetricType(Counter)),
	Import:           newInternalMeasurement(withDescription("Total Import"), withPrometheusName("energy_import"), withUnit(KiloWattHour), withUnitInPrometheus(Joule), withMetricType(Counter)),
	ImportT1:         newInternalMeasurement(withDescription("Tariff 1 Import"), withPrometheusName("tariff_1_imported"), withUnit(KiloWattHour), withUnitInPrometheus(Joule), withMetricType(Counter)),
	ImportT2:         newInternalMeasurement(withDescription("Tariff 2 Import"), withPrometheusName("tariff_2_imported"), withUnit(KiloWattHour), withUnitInPrometheus(Joule), withMetricType(Counter)),
	ImportL1:         newInternalMeasurement(withDescription("L1 Import"), withPrometheusName("l1_imported"), withUnit(KiloWattHour), withUnitInPrometheus(Joule), withMetricType(Counter)),
	ImportL2:         newInternalMeasurement(withDescription("L2 Import"), withPrometheusName("l2_imported"), withUnit(KiloWattHour), withUnitInPrometheus(Joule), withMetricType(Counter)),
	ImportL3:         newInternalMeasurement(withDescription("L3 Import"), withPrometheusName("l3_imported"), withUnit(KiloWattHour), withUnitInPrometheus(Joule), withMetricType(Counter)),
	Export:           newInternalMeasurement(withDescription("Total Export"), withPrometheusName("energy_export"), withUnit(KiloWattHour), withUnitInPrometheus(Joule), withMetricType(Counter)),
	ExportT1:         newInternalMeasurement(withDescription("Tariff 1 Export"), withPrometheusName("tariff_1_exported"), withUnit(KiloWattHour), withUnitInPrometheus(Joule), withMetricType(Counter)),
	ExportT2:         newInternalMeasurement(withDescription("Tariff 2 Export"), withPrometheusName("tariff_2_exported"), withUnit(KiloWattHour), withUnitInPrometheus(Joule), withMetricType(Counter)),
	ExportL1:         newInternalMeasurement(withDescription("L1 Export"), withPrometheusName("l1_exported"), withUnit(KiloWattHour), withUnitInPrometheus(Joule), withMetricType(Counter)),
	ExportL2:         newInternalMeasurement(withDescription("L2 Export"), withPrometheusName("l2_exported"), withUnit(KiloWattHour), withUnitInPrometheus(Joule), withMetricType(Counter)),
	ExportL3:         newInternalMeasurement(withDescription("L3 Export"), withPrometheusName("l3_exported"), withUnit(KiloWattHour), withUnitInPrometheus(Joule), withMetricType(Counter)),
	ReactiveSum:      newInternalMeasurement(withDescription("Total Reactive"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveSumT1:    newInternalMeasurement(withDescription("Tariff 1 Reactive"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveSumT2:    newInternalMeasurement(withDescription("Tariff 2 Reactive"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveSumL1:    newInternalMeasurement(withDescription("L1 Reactive"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveSumL2:    newInternalMeasurement(withDescription("L2 Reactive"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveSumL3:    newInternalMeasurement(withDescription("L3 Reactive"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveImport:   newInternalMeasurement(withDescription("Reactive Import"), withPrometheusName("reactive_imported"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveImportT1: newInternalMeasurement(withDescription("Tariff 1 Reactive Import"), withPrometheusName("tariff_2_reactive_imported"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveImportT2: newInternalMeasurement(withDescription("Tariff 2 Reactive Import"), withPrometheusName("tariff_1_reactive_imported"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveImportL1: newInternalMeasurement(withDescription("L1 Reactive Import"), withPrometheusName("l1_reactive_imported"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveImportL2: newInternalMeasurement(withDescription("L2 Reactive Import"), withPrometheusName("l2_reactive_imported"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveImportL3: newInternalMeasurement(withDescription("L3 Reactive Import"), withPrometheusName("l3_reactive_imported"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveExport:   newInternalMeasurement(withDescription("Reactive Export"), withPrometheusName("reactive_exported"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveExportT1: newInternalMeasurement(withDescription("Tariff 1 Reactive Export"), withPrometheusName("tariff_1_reactive_exported"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveExportT2: newInternalMeasurement(withDescription("Tariff 2 Reactive Export"), withPrometheusName("tariff_2_reactive_exported"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveExportL1: newInternalMeasurement(withDescription("L1 Reactive Export"), withPrometheusName("l1_reactive_exported"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveExportL2: newInternalMeasurement(withDescription("L2 Reactive Export"), withPrometheusName("l2_reactive_exported"), withUnit(KiloVarHour), withMetricType(Counter)),
	ReactiveExportL3: newInternalMeasurement(withDescription("L3 Reactive Export"), withPrometheusName("l3_reactive_exported"), withUnit(KiloVarHour), withMetricType(Counter)),
	DCCurrent:        newInternalMeasurement(withDescription("DC Current"), withUnit(Ampere), withMetricType(Gauge)),
	DCVoltage:        newInternalMeasurement(withDescription("DC Voltage"), withUnit(Volt), withMetricType(Gauge)),
	DCPower:          newInternalMeasurement(withDescription("DC Power"), withUnit(Watt), withMetricType(Gauge)),
	HeatSinkTemp:     newInternalMeasurement(withDescription("Heat Sink Temperature"), withUnit(DegreeCelsius), withMetricType(Gauge)),
	DCCurrentS1:      newInternalMeasurement(withDescription("String 1 Current"), withUnit(Ampere), withMetricType(Gauge)),
	DCVoltageS1:      newInternalMeasurement(withDescription("String 1 Voltage"), withUnit(Volt), withMetricType(Gauge)),
	DCPowerS1:        newInternalMeasurement(withDescription("String 1 Power"), withUnit(Watt), withMetricType(Gauge)),
	DCEnergyS1:       newInternalMeasurement(withDescription("String 1 Generation"), withUnit(KiloWattHour), withUnitInPrometheus(Joule), withMetricType(Counter)),
	DCCurrentS2:      newInternalMeasurement(withDescription("String 2 Current"), withUnit(Ampere), withMetricType(Gauge)),
	DCVoltageS2:      newInternalMeasurement(withDescription("String 2 Voltage"), withUnit(Volt), withMetricType(Gauge)),
	DCPowerS2:        newInternalMeasurement(withDescription("String 2 Power"), withUnit(Watt), withMetricType(Gauge)),
	DCEnergyS2:       newInternalMeasurement(withDescription("String 2 Generation"), withUnit(KiloWattHour), withUnitInPrometheus(Joule), withMetricType(Counter)),
	DCCurrentS3:      newInternalMeasurement(withDescription("String 3 Current"), withUnit(Ampere), withMetricType(Gauge)),
	DCVoltageS3:      newInternalMeasurement(withDescription("String 3 Voltage"), withUnit(Volt), withMetricType(Gauge)),
	DCPowerS3:        newInternalMeasurement(withDescription("String 3 Power"), withUnit(Watt), withMetricType(Gauge)),
	DCEnergyS3:       newInternalMeasurement(withDescription("String 3 Generation"), withUnit(KiloWattHour), withUnitInPrometheus(Joule), withMetricType(Counter)),
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
	}

	if m.PrometheusInfo.Description == "" {
		withGenericPrometheusDescription()(m)
	}

	if m.PrometheusInfo.Name == "" {
		withGenericPrometheusName()(m)
	} else {
		var displayUnit *Unit
		if m.PrometheusInfo.Unit != nil {
			displayUnit = m.PrometheusInfo.Unit
		} else if m.Unit != nil {
			displayUnit = m.Unit
		}

		if displayUnit != nil {
			if m.PrometheusInfo.MetricType == Counter {
				m.PrometheusInfo.Name = generatePrometheusName(m.PrometheusInfo.Name, displayUnit.PrometheusName()+"_total")
			} else {
				m.PrometheusInfo.Name = generatePrometheusName(m.PrometheusInfo.Name, displayUnit.PrometheusName())
			}
		} else {
			m.PrometheusInfo.Name = generatePrometheusName(m.PrometheusInfo.Name, "")
		}
	}

	return m
}

// withPrometheusDescription enables setting a Prometheus description of a Measurement
func withPrometheusDescription(description string) measurementOptions {
	return func(m *measurement) {
		m.PrometheusInfo.Description = description
	}
}

// withGenericPrometheusDescription sets the Prometheus description to a generated, more generic format
func withGenericPrometheusDescription() measurementOptions {
	return func(m *measurement) {
		var displayUnit *Unit
		if m.PrometheusInfo.Unit != nil {
			displayUnit = m.PrometheusInfo.Unit
		} else if m.Unit != nil {
			displayUnit = m.Unit
		}

		if displayUnit != nil {
			m.PrometheusInfo.Description = generatePrometheusDescription(m.Description, displayUnit.FullName())
		} else {
			m.PrometheusInfo.Description = generatePrometheusDescription(m.Description, "")
		}
	}
}

// withUnit sets the Unit of a Measurement
// If u is nil, the unit will be set to NoUnit
func withUnit(u Unit) measurementOptions {
	return func(m *measurement) {
		m.Unit = &u
	}
}

// withUnitInPrometheus set the Unit to be displayed in Prometheus
// If u is set, any incoming measurements are automatically converted to specified u
func withUnitInPrometheus(u Unit) measurementOptions {
	return func(m *measurement) {
		m.PrometheusInfo.Unit = &u
	}
}

func withPrometheusName(name string) measurementOptions {
	return func(m *measurement) {
		m.PrometheusInfo.Name = name
	}
}

func withGenericPrometheusName() measurementOptions {
	return func(m *measurement) {
		var displayUnit *Unit
		if m.PrometheusInfo.Unit != nil {
			displayUnit = m.PrometheusInfo.Unit
		} else if m.Unit != nil {
			displayUnit = m.Unit
		}

		if displayUnit != nil {
			if m.PrometheusInfo.MetricType == Counter {
				m.PrometheusInfo.Name = generatePrometheusName(m.Description, displayUnit.PrometheusName()+"_total")
			} else {
				m.PrometheusInfo.Name = generatePrometheusName(m.Description, displayUnit.PrometheusName())
			}
		} else {
			m.PrometheusInfo.Name = generatePrometheusName(m.Description, "")
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

func generatePrometheusDescription(description string, unit string) string {
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

// ConvertValueTo queries the targetUnit's conversion func from a conversion map
// and converts it to the defined Unit
//
// If a conversionFunc cannot be found, this func will return 0.0
// as the source Unit cannot be converted to target Unit via non-defined conversion function!
func (r *MeasurementResult) ConvertValueTo(targetUnit Unit) float64 {
	if r.Value == 0.0 {
		return 0.0
	}

	if conversionFunc, ok := conversionMap[targetUnit]; ok {
		return conversionFunc(r.Value)
	}

	return 0.0
}
