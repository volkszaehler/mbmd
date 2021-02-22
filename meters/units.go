package meters

// Unit is used to represent measurements
type Unit int

//go:generate enumer -type=Unit -transform=snake
const (
	_ Unit = iota

	KiloVarHour
	Var
	KiloWattHour
	Watt

	Ampere
	Volt
	VoltAmpere

	Degree
	DegreeCelsius

	Hertz

	Percent

	NoUnit
)

type unit struct {
	Abbreviation 	unitAbbreviation
	PluralForm		string
}

// unitAbbreviation defines the default abbreviation and - if needed - an alternative.
type unitAbbreviation struct {
	Default     string
	Alternative string
}

var units = map[Unit]*unit{
	KiloVarHour: 	{unitAbbreviation{"kvarh", ""}, ""},
	Var:			{unitAbbreviation{"var", ""}, ""},
	KiloWattHour:   {unitAbbreviation{"kWh", ""}, ""},
	Watt: 			{unitAbbreviation{"W", ""}, ""},
	Ampere: 		{unitAbbreviation{"A", ""}, ""},
	Volt: 			{unitAbbreviation{"V", ""}, ""},
	VoltAmpere: 	{unitAbbreviation{"VA", ""}, ""},
	Degree: 		{unitAbbreviation{"°", "degree"}, ""},
	DegreeCelsius: 	{unitAbbreviation{"°C", "degree_celsius"}, "degrees_celsius"},
	Hertz:			{unitAbbreviation{"Hz", "hertz"}, "hertz"},
	Percent:		{unitAbbreviation{"%", "percent"}, "percent"},

	NoUnit:			{unitAbbreviation{"", ""}, ""},
}

func (u *Unit) PrometheusName() string {
	if u == nil || *u == NoUnit {
		return ""
	}

	if unit, ok := units[*u]; ok {
		plural := unit.PluralForm
		if plural != "" {
			return plural
		}
	}

	return u.String() + "s"
}

// Abbreviation returns the matching abbreviation of a Unit if it exists
func (u *Unit) Abbreviation() string {
	if unit, ok := units[*u]; ok {
		return unit.Abbreviation.Default
	}
	return ""
}

// AlternativeAbbreviation returns the matching alternative abbreviation of a Unit if it exists
func (u *Unit) AlternativeAbbreviation() string {
	if unit, ok := units[*u]; ok {
		alternative := unit.Abbreviation.Alternative
		if alternative != "" {
			return alternative
		}
	}
	return ""
}

