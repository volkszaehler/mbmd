package meters

import (
	"testing"

	"github.com/volkszaehler/mbmd/meters/units"
)

func TestMeasurementCreation_WithRequiredOptions_WithMetricType_Counter(t *testing.T) {
	measurement := newInternalMeasurement(
		withDescription("My Test Measurement"),
		withUnit(units.Ampere),
		withMetricType(Counter),
	)

	expectedPrometheusName := "measurement_my_test_measurement_amperes_total"
	expectedDescription := "My Test Measurement in Amperes"

	if measurement.PrometheusInfo.Name != expectedPrometheusName {
		t.Errorf(
			"Prometheus metric name '%s' does not equal expected '%s'",
			measurement.PrometheusInfo.Name,
			expectedPrometheusName,
		)
	}

	if measurement.PrometheusInfo.HelpText != expectedDescription {
		t.Errorf("Prometheus description '%s' does not equal expected '%s'",
			measurement.PrometheusInfo.HelpText,
			expectedDescription,
		)
	}

	if measurement.Unit != units.Ampere {
		t.Errorf("Prometheus unit '%s' does not equal expected '%s'",
			measurement.Unit,
			units.Ampere,
		)
	}
}

func TestMeasurementCreation_WithCustomName_AndDescription(t *testing.T) {
	measurement := newInternalMeasurement(
		withDescription("My Test Measurement"),
		withPrometheusHelpText("My custom description for my measurement"),
		withPrometheusName("my_custom_name_for_my_test_measurement"),
		withUnit(units.Ampere),
		withMetricType(Gauge),
	)

	expectedPrometheusName := "measurement_my_custom_name_for_my_test_measurement_amperes"
	expectedDescription := "My custom description for my measurement"

	if measurement.PrometheusInfo.Name != expectedPrometheusName {
		t.Errorf(
			"Prometheus metric name '%s' does not equal expected '%s'",
			measurement.PrometheusInfo.Name,
			expectedPrometheusName,
		)
	}

	if measurement.PrometheusInfo.HelpText != expectedDescription {
		t.Errorf("Prometheus description '%s' does not equal expected '%s'",
			measurement.PrometheusInfo.HelpText,
			expectedDescription,
		)
	}
}

func TestInternalMeasurement_AutoConvertToElementaryUnit(t *testing.T) {
	measurementKwh := newInternalMeasurement(
		withDescription("My Test Measurement with kWh"),
		withPrometheusHelpText("My custom description for my measurement"),
		withPrometheusName("my_custom_name_for_my_test_measurement_energy"),
		withUnit(units.KiloWattHour),
		withMetricType(Gauge),
	)

	measurementKvarh := newInternalMeasurement(
		withDescription("My Test Measurement"),
		withPrometheusHelpText("My custom description for my measurement"),
		withPrometheusName("my_custom_name_for_my_test_measurement_energy"),
		withUnit(units.KiloWattHour),
		withMetricType(Gauge),
	)

	expectedConvertedUnit := units.Joule
	_, expectedConvertedUnitPluralForm := expectedConvertedUnit.Name()

	if measurementKwh.PrometheusInfo.Unit != expectedConvertedUnit {
		actualConvertedUnit := measurementKwh.PrometheusInfo.Unit
		_, actualConvertedUnitPluralForm := actualConvertedUnit.Name()

		t.Errorf(
			"measurement_kWh could not be converted to elementary unit '%s' automatically (actual: %s)",
			expectedConvertedUnitPluralForm,
			actualConvertedUnitPluralForm,
		)
	}

	if measurementKvarh.PrometheusInfo.Unit != expectedConvertedUnit {
		actualConvertedUnit := measurementKwh.PrometheusInfo.Unit
		_, actualConvertedUnitPluralForm := actualConvertedUnit.Name()

		t.Errorf("measurement_kvarh could not be converted to elementary unit '%s' automatically (actual: %s)",
			expectedConvertedUnitPluralForm,
			actualConvertedUnitPluralForm,
		)
	}
}
