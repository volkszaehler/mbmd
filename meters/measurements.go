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
// <Description>	::= <measurementOption.WithDescription()> | <measurementOption.WithCustomDescription()>
// <Unit>			::= <measurementOption.WithUnit()>
// <CounterTotal>	::= "total" // if metric type is Counter
// E. g.:
//  Assuming a device's manufacturer is "myManufacturer":
//		newInternalMeasurement(WithDescription("Frequency Test With Some Text"), WithUnit(Hertz), WithMetricType(Counter))
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
	Frequency:        newInternalMeasurement(WithDescription("Frequency"), WithUnit(Hertz), WithMetricType(Gauge)),
	Current:          newInternalMeasurement(WithDescription("Current"), WithUnit(Ampere), WithMetricType(Gauge)),
	CurrentL1:        newInternalMeasurement(WithDescription("L1 Current"), WithUnit(Ampere), WithMetricType(Gauge)),
	CurrentL2:        newInternalMeasurement(WithDescription("L2 Current"), WithUnit(Ampere), WithMetricType(Gauge)),
	CurrentL3:        newInternalMeasurement(WithDescription("L3 Current"), WithUnit(Ampere), WithMetricType(Gauge)),
	Voltage:          newInternalMeasurement(WithDescription("Voltage"), WithUnit(Volt), WithMetricType(Gauge)),
	VoltageL1:        newInternalMeasurement(WithDescription("L1 Voltage"), WithUnit(Volt), WithMetricType(Gauge)),
	VoltageL2:        newInternalMeasurement(WithDescription("L2 Voltage"), WithUnit(Volt), WithMetricType(Gauge)),
	VoltageL3:        newInternalMeasurement(WithDescription("L3 Voltage"), WithUnit(Volt), WithMetricType(Gauge)),
	Power:            newInternalMeasurement(WithDescription("Power"), WithUnit(Watt), WithMetricType(Gauge)),
	PowerL1:          newInternalMeasurement(WithDescription("L1 Power"), WithUnit(Watt), WithMetricType(Gauge)),
	PowerL2:          newInternalMeasurement(WithDescription("L2 Power"), WithUnit(Watt), WithMetricType(Gauge)),
	PowerL3:          newInternalMeasurement(WithDescription("L3 Power"), WithUnit(Watt), WithMetricType(Gauge)),
	ImportPower:      newInternalMeasurement(WithDescription("Import Power"), WithUnit(Watt), WithMetricType(Gauge)),
	ImportPowerL1:    newInternalMeasurement(WithDescription("L1 Import Power"), WithUnit(Watt), WithMetricType(Gauge)),
	ImportPowerL2:    newInternalMeasurement(WithDescription("L2 Import Power"), WithUnit(Watt), WithMetricType(Gauge)),
	ImportPowerL3:    newInternalMeasurement(WithDescription("L3 Import Power"), WithUnit(Watt), WithMetricType(Gauge)),
	ExportPower:      newInternalMeasurement(WithDescription("Export Power"), WithUnit(Watt), WithMetricType(Gauge)),
	ExportPowerL1:    newInternalMeasurement(WithDescription("L1 Export Power"), WithUnit(Watt), WithMetricType(Gauge)),
	ExportPowerL2:    newInternalMeasurement(WithDescription("L2 Export Power"), WithUnit(Watt), WithMetricType(Gauge)),
	ExportPowerL3:    newInternalMeasurement(WithDescription("L3 Export Power"), WithUnit(Watt), WithMetricType(Gauge)),
	ReactivePower:    newInternalMeasurement(WithDescription("Reactive Power"), WithUnit(Var), WithMetricType(Gauge)),
	ReactivePowerL1:  newInternalMeasurement(WithDescription("L1 Reactive Power"), WithUnit(Var), WithMetricType(Gauge)),
	ReactivePowerL2:  newInternalMeasurement(WithDescription("L2 Reactive Power"), WithUnit(Var), WithMetricType(Gauge)),
	ReactivePowerL3:  newInternalMeasurement(WithDescription("L3 Reactive Power"), WithUnit(Var), WithMetricType(Gauge)),
	ApparentPower:    newInternalMeasurement(WithDescription("Apparent Power"), WithUnit(VoltAmpere), WithMetricType(Gauge)),
	ApparentPowerL1:  newInternalMeasurement(WithDescription("L1 Apparent Power"), WithUnit(VoltAmpere), WithMetricType(Gauge)),
	ApparentPowerL2:  newInternalMeasurement(WithDescription("L2 Apparent Power"), WithUnit(VoltAmpere), WithMetricType(Gauge)),
	ApparentPowerL3:  newInternalMeasurement(WithDescription("L3 Apparent Power"), WithUnit(VoltAmpere), WithMetricType(Gauge)),
	Cosphi:           newInternalMeasurement(WithDescription("Cosphi"), WithMetricType(Gauge)),
	CosphiL1:         newInternalMeasurement(WithDescription("L1 Cosphi"), WithMetricType(Gauge)),
	CosphiL2:         newInternalMeasurement(WithDescription("L2 Cosphi"), WithMetricType(Gauge)),
	CosphiL3:         newInternalMeasurement(WithDescription("L3 Cosphi"), WithMetricType(Gauge)),
	THD:              newInternalMeasurement(WithDescription("Average voltage to neutral THD"), WithUnit(Percent), WithMetricType(Gauge)),
	THDL1:            newInternalMeasurement(WithDescription("L1 Voltage to neutral THD"), WithUnit(Percent), WithMetricType(Gauge)),
	THDL2:            newInternalMeasurement(WithDescription("L2 Voltage to neutral THD"), WithUnit(Percent), WithMetricType(Gauge)),
	THDL3:            newInternalMeasurement(WithDescription("L3 Voltage to neutral THD"), WithUnit(Percent), WithMetricType(Gauge)),
	Sum:              newInternalMeasurement(WithDescription("Total Sum"), WithUnit(KiloWattHour), WithMetricType(Counter)),
	SumT1:            newInternalMeasurement(WithDescription("Tariff 1 Sum"), WithUnit(KiloWattHour), WithMetricType(Counter)),
	SumT2:            newInternalMeasurement(WithDescription("Tariff 2 Sum"), WithUnit(KiloWattHour), WithMetricType(Counter)),
	SumL1:            newInternalMeasurement(WithDescription("L1 Sum"), WithUnit(KiloWattHour), WithMetricType(Counter)),
	SumL2:            newInternalMeasurement(WithDescription("L2 Sum"), WithUnit(KiloWattHour), WithMetricType(Counter)),
	SumL3:            newInternalMeasurement(WithDescription("L3 Sum"), WithUnit(KiloWattHour), WithMetricType(Counter)),
	Import:           newInternalMeasurement(WithDescription("Total Import"), WithPrometheusName("total_imported"), WithUnit(KiloWattHour), WithMetricType(Counter)),
	ImportT1:         newInternalMeasurement(WithDescription("Tariff 1 Import"), WithPrometheusName("tariff_1_imported"), WithUnit(KiloWattHour), WithMetricType(Counter)),
	ImportT2:         newInternalMeasurement(WithDescription("Tariff 2 Import"), WithPrometheusName("tariff_2_imported"), WithUnit(KiloWattHour), WithMetricType(Counter)),
	ImportL1:         newInternalMeasurement(WithDescription("L1 Import"), WithPrometheusName("l1_imported"), WithUnit(KiloWattHour), WithMetricType(Counter)),
	ImportL2:         newInternalMeasurement(WithDescription("L2 Import"), WithPrometheusName("l2_imported"), WithUnit(KiloWattHour), WithMetricType(Counter)),
	ImportL3:         newInternalMeasurement(WithDescription("L3 Import"), WithPrometheusName("l3_imported"), WithUnit(KiloWattHour), WithMetricType(Counter)),
	Export:           newInternalMeasurement(WithDescription("Total Export"), WithPrometheusName("total_exported"), WithUnit(KiloWattHour), WithMetricType(Counter)),
	ExportT1:         newInternalMeasurement(WithDescription("Tariff 1 Export"), WithPrometheusName("tariff_1_exported"), WithUnit(KiloWattHour), WithMetricType(Counter)),
	ExportT2:         newInternalMeasurement(WithDescription("Tariff 2 Export"), WithPrometheusName("tariff_2_exported"), WithUnit(KiloWattHour), WithMetricType(Counter)),
	ExportL1:         newInternalMeasurement(WithDescription("L1 Export"), WithPrometheusName("l1_exported"), WithUnit(KiloWattHour), WithMetricType(Counter)),
	ExportL2:         newInternalMeasurement(WithDescription("L2 Export"), WithPrometheusName("l2_exported"), WithUnit(KiloWattHour), WithMetricType(Counter)),
	ExportL3:         newInternalMeasurement(WithDescription("L3 Export"), WithPrometheusName("l3_exported"), WithUnit(KiloWattHour), WithMetricType(Counter)),
	ReactiveSum:      newInternalMeasurement(WithDescription("Total Reactive"), WithUnit(KiloVarHour), WithMetricType(Counter)),
	ReactiveSumT1:    newInternalMeasurement(WithDescription("Tariff 1 Reactive"), WithUnit(KiloVarHour), WithMetricType(Counter)),
	ReactiveSumT2:    newInternalMeasurement(WithDescription("Tariff 2 Reactive"), WithUnit(KiloVarHour), WithMetricType(Counter)),
	ReactiveSumL1:    newInternalMeasurement(WithDescription("L1 Reactive"), WithUnit(KiloVarHour), WithMetricType(Counter)),
	ReactiveSumL2:    newInternalMeasurement(WithDescription("L2 Reactive"), WithUnit(KiloVarHour), WithMetricType(Counter)),
	ReactiveSumL3:    newInternalMeasurement(WithDescription("L3 Reactive"), WithUnit(KiloVarHour), WithMetricType(Counter)),
	ReactiveImport:   newInternalMeasurement(WithDescription("Reactive Import"), WithPrometheusName("reactive_imported"), WithUnit(KiloVarHour), WithMetricType(Counter)),
	ReactiveImportT1: newInternalMeasurement(WithDescription("Tariff 1 Reactive Import"), WithPrometheusName("tariff_2_reactive_imported"), WithUnit(KiloVarHour), WithMetricType(Counter)),
	ReactiveImportT2: newInternalMeasurement(WithDescription("Tariff 2 Reactive Import"), WithPrometheusName("tariff_1_reactive_imported"), WithUnit(KiloVarHour), WithMetricType(Counter)),
	ReactiveImportL1: newInternalMeasurement(WithDescription("L1 Reactive Import"), WithPrometheusName("l1_reactive_imported"), WithUnit(KiloVarHour), WithMetricType(Counter)),
	ReactiveImportL2: newInternalMeasurement(WithDescription("L2 Reactive Import"), WithPrometheusName("l2_reactive_imported"), WithUnit(KiloVarHour), WithMetricType(Counter)),
	ReactiveImportL3: newInternalMeasurement(WithDescription("L3 Reactive Import"), WithPrometheusName("l3_reactive_imported"), WithUnit(KiloVarHour), WithMetricType(Counter)),
	ReactiveExport:   newInternalMeasurement(WithDescription("Reactive Export"), WithPrometheusName("reactive_exported"), WithUnit(KiloVarHour), WithMetricType(Counter)),
	ReactiveExportT1: newInternalMeasurement(WithDescription("Tariff 1 Reactive Export"), WithPrometheusName("tariff_1_reactive_exported"), WithUnit(KiloVarHour), WithMetricType(Counter)),
	ReactiveExportT2: newInternalMeasurement(WithDescription("Tariff 2 Reactive Export"), WithPrometheusName("tariff_2_reactive_exported"), WithUnit(KiloVarHour), WithMetricType(Counter)),
	ReactiveExportL1: newInternalMeasurement(WithDescription("L1 Reactive Export"), WithPrometheusName("l1_reactive_exported"), WithUnit(KiloVarHour), WithMetricType(Counter)),
	ReactiveExportL2: newInternalMeasurement(WithDescription("L2 Reactive Export"), WithPrometheusName("l2_reactive_exported"), WithUnit(KiloVarHour), WithMetricType(Counter)),
	ReactiveExportL3: newInternalMeasurement(WithDescription("L3 Reactive Export"), WithPrometheusName("l3_reactive_exported"), WithUnit(KiloVarHour), WithMetricType(Counter)),
	DCCurrent:        newInternalMeasurement(WithDescription("DC Current"), WithUnit(Ampere), WithMetricType(Gauge)),
	DCVoltage:        newInternalMeasurement(WithDescription("DC Voltage"), WithUnit(Volt), WithMetricType(Gauge)),
	DCPower:          newInternalMeasurement(WithDescription("DC Power"), WithUnit(Watt), WithMetricType(Gauge)),
	HeatSinkTemp:     newInternalMeasurement(WithDescription("Heat Sink Temperature"), WithUnit(DegreeCelsius), WithMetricType(Gauge)),
	DCCurrentS1:      newInternalMeasurement(WithDescription("String 1 Current"), WithUnit(Ampere), WithMetricType(Gauge)),
	DCVoltageS1:      newInternalMeasurement(WithDescription("String 1 Voltage"), WithUnit(Volt), WithMetricType(Gauge)),
	DCPowerS1:        newInternalMeasurement(WithDescription("String 1 Power"), WithUnit(Watt), WithMetricType(Gauge)),
	DCEnergyS1:       newInternalMeasurement(WithDescription("String 1 Generation"), WithUnit(KiloWattHour), WithMetricType(Counter)),
	DCCurrentS2:      newInternalMeasurement(WithDescription("String 2 Current"), WithUnit(Ampere), WithMetricType(Gauge)),
	DCVoltageS2:      newInternalMeasurement(WithDescription("String 2 Voltage"), WithUnit(Volt), WithMetricType(Gauge)),
	DCPowerS2:        newInternalMeasurement(WithDescription("String 2 Power"), WithUnit(Watt), WithMetricType(Gauge)),
	DCEnergyS2:       newInternalMeasurement(WithDescription("String 2 Generation"), WithUnit(KiloWattHour), WithMetricType(Counter)),
	DCCurrentS3:      newInternalMeasurement(WithDescription("String 3 Current"), WithUnit(Ampere), WithMetricType(Gauge)),
	DCVoltageS3:      newInternalMeasurement(WithDescription("String 3 Voltage"), WithUnit(Volt), WithMetricType(Gauge)),
	DCPowerS3:        newInternalMeasurement(WithDescription("String 3 Power"), WithUnit(Watt), WithMetricType(Gauge)),
	DCEnergyS3:       newInternalMeasurement(WithDescription("String 3 Generation"), WithUnit(KiloWattHour), WithMetricType(Counter)),
	ChargeState:      newInternalMeasurement(WithDescription("Charge State"), WithUnit(Percent), WithMetricType(Gauge)),
	BatteryVoltage:   newInternalMeasurement(WithDescription("Battery Voltage"), WithUnit(Volt), WithMetricType(Gauge)),
	PhaseAngle:       newInternalMeasurement(WithDescription("Phase Angle"), WithUnit(Degree), WithMetricType(Gauge)),
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
		WithUnit(NoUnit)(m)
	}

	if m.PrometheusInfo.Description == "" {
		WithGenericPrometheusDescription()(m)
	}

	if m.PrometheusInfo.Name == "" {
		WithGenericPrometheusName()(m)
	} else {
		if m.Unit != nil {
			if m.PrometheusInfo.MetricType == Counter {
				m.PrometheusInfo.Name = generatePrometheusName(m.PrometheusInfo.Name, m.Unit.PrometheusName()+"_total")
			} else {
				m.PrometheusInfo.Name = generatePrometheusName(m.PrometheusInfo.Name, m.Unit.PrometheusName())
			}
		} else {
			m.PrometheusInfo.Name = generatePrometheusName(m.PrometheusInfo.Name, "")
		}
	}

	return m
}

// WithPrometheusDescription enables setting a Prometheus description of a Measurement
func WithPrometheusDescription(description string) measurementOptions {
	return func(m *measurement) {
		m.PrometheusInfo.Description = description
	}
}

// WithGenericPrometheusDescription sets the Prometheus description to a generated, more generic format
func WithGenericPrometheusDescription() measurementOptions {
	return func(m *measurement) {
		if m.Unit != nil {
			m.PrometheusInfo.Description = generatePrometheusDescription(m.Description, m.Unit.FullName())
		} else {
			m.PrometheusInfo.Description = generatePrometheusDescription(m.Description, "")
		}
	}
}

// WithUnit sets the Unit of a Measurement
// If u is nil, the unit will be set to NoUnit
func WithUnit(u Unit) measurementOptions {
	return func(m *measurement) {
		m.Unit = &u
	}
}

func WithPrometheusName(name string) measurementOptions {
	return func(m *measurement) {
		m.PrometheusInfo.Name = name
	}
}

func WithGenericPrometheusName() measurementOptions {
	return func(m *measurement) {
		if m.Unit != nil {
			if m.PrometheusInfo.MetricType == Counter {
				m.PrometheusInfo.Name = generatePrometheusName(m.Description, m.Unit.PrometheusName()+"_total")
			} else {
				m.PrometheusInfo.Name = generatePrometheusName(m.Description, m.Unit.PrometheusName())
			}
		} else {
			m.PrometheusInfo.Name = generatePrometheusName(m.Description, "")
		}
	}
}

func WithMetricType(metricType MetricType) measurementOptions {
	return func(m *measurement) {
		m.PrometheusInfo.MetricType = metricType
	}
}

func WithDescription(description string) measurementOptions {
	return func(m *measurement) {
		m.Description = description
	}
}

func generatePrometheusDescription(description string, unit string) string {
	if unit != "" {
		return fmt.Sprintf("Measurement of %s in %s", description, unit)
	} else {
		return fmt.Sprintf("Measurement of %s", description)
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
