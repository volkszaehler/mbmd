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

	Joule

	NoUnit
)

type unit struct {
	Abbreviation unitAbbreviation
	FullName     string
	PluralForm   string
}

// unitAbbreviation defines the default abbreviation and - if needed - an alternative.
type unitAbbreviation struct {
	Default     string
	Alternative string
}

var units = map[Unit]*unit{
	KiloVarHour:   {unitAbbreviation{"kvarh", ""}, "Kilovoltampere-hours (reactive)", ""},
	Var:           {unitAbbreviation{"var", "volt_ampere"}, "Voltamperes (reactive)", "volt_amperes"},
	KiloWattHour:  {unitAbbreviation{"kWh", ""}, "Kilowatt-hours", ""},
	Watt:          {unitAbbreviation{"W", ""}, "Watts", ""},
	Ampere:        {unitAbbreviation{"A", ""}, "Amperes", ""},
	Volt:          {unitAbbreviation{"V", ""}, "Volts", ""},
	VoltAmpere:    {unitAbbreviation{"VA", ""}, "Voltamperes", ""},
	Degree:        {unitAbbreviation{"°", "degree"}, "Degrees", ""},
	DegreeCelsius: {unitAbbreviation{"°C", "degree_celsius"}, "Degree Celsius", "degrees_celsius"},
	Hertz:         {unitAbbreviation{"Hz", "hertz"}, "Hertz", "hertz"},
	Percent:       {unitAbbreviation{"%", "percent"}, "Per cent", "percent"},
	Joule:         {unitAbbreviation{"J", ""}, "Joules", ""},

	NoUnit: {unitAbbreviation{"", ""}, "", ""},
}

var conversionMap = map[Unit]func(float64) float64{
	Joule: func(input float64) float64 {
		return (input * 3.6) * 1_000_000
	},
}

//var units = map[Unit]*unit{
//	KiloVarHour:   makeUnitDefinition(
//		"kvarh",
//		"",
//		"Kilovoltampere-hours (reactive)",
//		"",
//		),
//	Var:
//}

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

func (u *Unit) FullName() string {
	if unit, ok := units[*u]; ok {
		return unit.FullName
	}
	return ""
}

//func makeUnitDefinition(
//	abbreviation string,
//	abbreviationAlternative string,
//	fullName string,
//	pluralForm string,
//) *unit {
//	unitAbbr := &unitAbbreviation{
//		Default:		abbreviation,
//		Alternative: 	abbreviationAlternative,
//	}
//
//	unitDef := &unit{
//		Abbreviation: 	unitAbbr,
//		FullName: 		fullName,
//		PluralForm: 	pluralForm,
//	}
//
//	return unitDef
//}
//
//func makeUnitDefinitionWithConversion(
//	abbreviation string,
//	abbreviationAlternative string,
//	fullName string,
//	pluralForm string,
//	conversionFuncs map[string]*func(float64)float64,
//) *unit {
//	unitDef := makeUnitDefinition(abbreviation, abbreviationAlternative, fullName, pluralForm)
//
//	unitDef.Converter = conversionFuncs
//
//	return unitDef
//}
