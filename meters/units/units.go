package units

// Unit is used to represent measurements
type Unit int

//go:generate go run github.com/dmarkham/enumer -type=Unit -transform=snake
const (
	KiloVarHour Unit = iota + 1
	KiloWattHour
	Joule

	Ampere
	Volt

	Watt
	Voltampere
	Var

	Degree
	DegreeCelsius

	Hertz

	Percent

	NoUnit // max value
)

type unit struct {
	abbreviation unitAbbreviation
	name         unitName
}

type internalUnitOption func(*unit)

// unitAbbreviation defines the default abbreviation and - if needed - an alternative.
type unitAbbreviation struct {
	Default     string
	Alternative string
}

type unitName struct {
	Singular       string
	Plural         string
	PrometheusForm string
}

var units = map[Unit]*unit{
	KiloVarHour: newInternalUnit(
		KiloVarHour,
		withName("Kilovoltampere-hour (reactive)", "Kilovoltampere-hours (reactive)"),
		withAbbreviation("kvarh", ""),
		// Unit is automatically converted to Joules in Prometheus context!
	),

	KiloWattHour: newInternalUnit(
		KiloWattHour,
		withName("Kilowatt-hour", ""),
		withAbbreviation("kWh", ""),
		// Unit is automatically converted to Joules in Prometheus context!
	),

	Var: newInternalUnit(
		Var,
		withName("Voltampere (reactive)", "Voltamperes (reactive)"),
		withAbbreviation(Var.String(), ""),
		withNameInPrometheusForm("voltamperes"),
	),

	Watt: newInternalUnit(
		Watt,
		withName("Watt", ""),
		withAbbreviation("W", ""),
	),

	Ampere: newInternalUnit(
		Ampere,
		withName("Ampere", ""),
		withAbbreviation("A", ""),
	),

	Volt: newInternalUnit(
		Volt,
		withName("Volt", ""),
		withAbbreviation("V", ""),
	),

	Voltampere: newInternalUnit(
		Voltampere,
		withName("Voltampere", ""),
		withAbbreviation("VA", ""),
	),

	Degree: newInternalUnit(
		Degree,
		withName("Degree", ""),
		withAbbreviation("°", "degree"),
	),

	DegreeCelsius: newInternalUnit(
		DegreeCelsius,
		withName("Degree Celsius", "Degrees Celsius"),
		withAbbreviation("°C", "degree celsius"),
		withNameInPrometheusForm("degrees_celsius"),
	),

	Hertz: newInternalUnit(
		Hertz,
		withName("Hertz", "Hertz"),
		withAbbreviation("hz", ""),
		withNameInPrometheusForm("hertz"),
	),

	Percent: newInternalUnit(
		Percent,
		withName("Percent", "Percent"),
		withAbbreviation("%", "percent"),
		withNameInPrometheusForm("percent"),
	),

	Joule: newInternalUnit(
		Joule,
		withName("Joule", ""),
		withAbbreviation("J", ""),
	),

	NoUnit: newInternalUnit(NoUnit),
}

// PrometheusForm returns Prometheus form of the Unit's associated elementary unit
func (u Unit) PrometheusForm() string {
	if u == 0 || u == NoUnit {
		return ""
	}

	elementaryUnit, _ := ConvertValueToElementaryUnit(u, 0.0)

	if unit, ok := units[elementaryUnit]; ok {
		return unit.name.PrometheusForm
	}

	return ""
}

// Abbreviation returns the matching abbreviation of a Unit if it exists
func (u Unit) Abbreviation() string {
	if unit, ok := units[u]; ok {
		return unit.abbreviation.Default
	}
	return ""
}

// AlternativeAbbreviation returns the matching alternative abbreviation of a Unit if it exists
func (u Unit) AlternativeAbbreviation() string {
	if unit, ok := units[u]; ok {
		alternative := unit.abbreviation.Alternative
		if alternative != "" {
			return alternative
		}
	}
	return ""
}

// Name returns the singular and plural form of a Unit's name
func (u Unit) Name() (string, string) {
	if unit, ok := units[u]; ok {
		return unit.name.Singular, unit.name.Plural
	}
	return "", ""
}

// newInternalUnit is a factory method for instantiating internal unit struct
func newInternalUnit(associatedUnit Unit, opts ...internalUnitOption) *unit {
	unit := &unit{}

	// Early return for NoUnit as we do not want to specify anything for "nothing"
	if associatedUnit == NoUnit {
		return unit
	}

	for _, opt := range opts {
		opt(unit)
	}

	if unit.name.PrometheusForm == "" {
		unit.name.PrometheusForm = associatedUnit.String() + "s"
	}

	if unit.name.Plural == "" {
		unit.name.Plural = unit.name.Singular + "s"
	}

	return unit
}

// withAbbreviation defines a default and an alternative representation of the unit's abbreviation
//
// Leave `alternative` empty if not needed
func withAbbreviation(defaultAbbreviation string, alternative string) internalUnitOption {
	return func(u *unit) {
		u.abbreviation.Default = defaultAbbreviation
		u.abbreviation.Alternative = alternative
	}
}

// withName defines the name of a unit - singular and plural
//
// Leave `plural` empty if an `s` can be appended to the literal `singular` string
func withName(singular string, plural string) internalUnitOption {
	return func(u *unit) {
		u.name.Singular = singular
		u.name.Plural = plural
	}
}

// withNameInPrometheusForm sets an explicit Prometheus form for a unit
//
// You don't have to use this option if the associated Unit key with an appended `s` can be used
func withNameInPrometheusForm(prometheusForm string) internalUnitOption {
	return func(u *unit) {
		u.name.PrometheusForm = prometheusForm
	}
}
