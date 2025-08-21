package units

import (
	"testing"
)

func Test_makeInternalUnit(t *testing.T) {
	internalUnit := newInternalUnit(
		Ampere,
		withName("Ampere", ""),
		withAbbreviation("A", ""),
	)

	if internalUnit.name.Singular != "Ampere" {
		t.Errorf(
			"Actual defined singular name '%s' for Unit '%v' does not equal expected name '%s'",
			internalUnit.name.Singular,
			Ampere,
			"Ampere",
		)
	}

	if internalUnit.name.Plural != "Amperes" {
		t.Errorf(
			"Actual defined plural name '%s' for Unit '%v' does not equal expected name '%s'",
			internalUnit.name.Plural,
			Ampere,
			"Amperes",
		)
	}

	expectedPrometheusForm := Ampere.String() + "s"

	if internalUnit.name.PrometheusForm != expectedPrometheusForm {
		t.Errorf(
			"Actual defined Prometheus form '%s' for Unit '%v' does not equal expected name '%s'",
			internalUnit.name.PrometheusForm,
			Ampere,
			expectedPrometheusForm,
		)
	}
}
