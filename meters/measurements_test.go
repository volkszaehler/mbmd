package meters

import (
	"testing"
	"time"
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
		withPrometheusDescription("My custom description for my measurement"),
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

func TestMeasurementCreation_WithUnitInPrometheusSpecified(t *testing.T) {
	measurement := newInternalMeasurement(
		withDescription("My Test Measurement"),
		withPrometheusDescription("My custom description for my measurement"),
		withPrometheusName("my_custom_name_for_my_test_measurement"),
		withUnit(KiloWattHour),
		withUnitInPrometheus(Joule),
		withMetricType(Counter),
	)

	expectedPrometheusName := "measurement_my_custom_name_for_my_test_measurement_joules_total"
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

func TestMeasurementResult_ConvertValueTo(t *testing.T) {
	measurementResult := &MeasurementResult{
		Measurement: Export,
		Value:       100.0,
		Timestamp:   time.Time{},
	}

	expected := 360_000_000.0
	actual := measurementResult.ConvertValueTo(Joule)

	if expected != actual {
		t.Errorf("Expected conversion value '%f' does not equal actual value '%f'",
			expected,
			actual)
	}
}

func TestMeasurementResult_ConvertValueTo_NonExistingConversionFunc(t *testing.T) {
	measurementResult := &MeasurementResult{
		Measurement: Export,
		Value:       100.0,
		Timestamp:   time.Time{},
	}

	expected := 0.0
	actual := measurementResult.ConvertValueTo(NoUnit)

	if expected != actual {
		t.Errorf("Expected conversion value '%f' does not equal actual value '%f'",
			expected,
			actual)
	}
}
