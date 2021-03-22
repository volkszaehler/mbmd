package meters

import (
	"testing"
)

func TestMeasurementCreation_WithRequiredOptions_WithMetricType_Counter(t *testing.T) {
	measurement := newInternalMeasurement(
		withDescription("My Test Measurement"),
		withUnit(Ampere),
		withMetricType(Counter),
	)

	expectedPrometheusName := "measurement_my_test_measurement_amperes_total"
	expectedDescription := "Measurement of My Test Measurement in A"

	if measurement.PrometheusInfo.Name != expectedPrometheusName {
		t.Errorf(
			"Prometheus metric name '%s' does not equal expected '%s'",
			measurement.PrometheusInfo.Name,
			expectedPrometheusName,
		)
	}

	if measurement.PrometheusInfo.Description != expectedDescription {
		t.Errorf("Prometheus description '%s' does not equal expected '%s'",
			measurement.PrometheusInfo.Description,
			expectedDescription,
		)
	}

	if *measurement.Unit != Ampere {
		t.Errorf("Prometheus unit '%s' does not equal expected '%s'",
			measurement.Unit,
			Ampere,
		)
	}
}

func TestMeasurementCreation_WithCustomName_AndDescription(t *testing.T) {
	measurement := newInternalMeasurement(
		withDescription("My Test Measurement"),
		withPrometheusHelpText("My custom description for my measurement"),
		withPrometheusName("my_custom_name_for_my_test_measurement"),
		withUnit(Ampere),
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

	if measurement.PrometheusInfo.Description != expectedDescription {
		t.Errorf("Prometheus description '%s' does not equal expected '%s'",
			measurement.PrometheusInfo.Description,
			expectedDescription,
		)
	}
}

func TestInternalMeasurement_AutoConvertToElementaryUnit(t *testing.T) {
	measurementKwh := newInternalMeasurement(
		withDescription("My Test Measurement with kWh"),
		withPrometheusHelpText("My custom description for my measurement"),
		withPrometheusName("my_custom_name_for_my_test_measurement_energy"),
		withUnit(KiloWattHour),
		withMetricType(Gauge),
	)

	measurementKvarh := newInternalMeasurement(
		withDescription("My Test Measurement"),
		withPrometheusHelpText("My custom description for my measurement"),
		withPrometheusName("my_custom_name_for_my_test_measurement_energy"),
		withUnit(KiloWattHour),
		withMetricType(Gauge),
	)

	expectedConvertedUnit := Joule

	if *measurementKwh.PrometheusInfo.Unit != expectedConvertedUnit {
		actualConvertedUnit := measurementKwh.PrometheusInfo.Unit
		t.Errorf(
			"measurement_kWh could not be converted to elementary unit %s automatically (actual: %s)",
			expectedConvertedUnit.FullName(),
			actualConvertedUnit.FullName(),
		)
	}

	if *measurementKvarh.PrometheusInfo.Unit != expectedConvertedUnit {
		actualConvertedUnit := measurementKwh.PrometheusInfo.Unit
		t.Errorf("measurement_kvarh could not be converted to elementary unit %s automatically (actual: %s)",
			expectedConvertedUnit.FullName(),
			actualConvertedUnit.FullName(),
		)
	}
}

func TestConvertValueToElementaryUnit(t *testing.T) {
	measurementResult := &MeasurementResult{
		Measurement: Export,
		Value:       100.0,
	}

	expectedConvertedUnit, expectedConvertedValue := Joule, 360_000_000.0
	actualConvertedUnit, actualConvertedValue := ConvertValueToElementaryUnit(*measurementResult.Measurement.Unit(), measurementResult.Value)

	if actualConvertedUnit != expectedConvertedUnit {
		t.Errorf(
			"Actual converted unit '%s' does not equal expected unit '%s'",
			actualConvertedUnit,
			expectedConvertedUnit,
		)
	}

	if actualConvertedValue != expectedConvertedValue {
		t.Errorf(
			"Actual converted value '%f' does not equal expected value '%f'",
			actualConvertedValue,
			expectedConvertedValue,
		)
	}
}
