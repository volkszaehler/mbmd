package units

import "testing"

func TestConvertValueToElementaryUnit(t *testing.T) {
	unit, value := KiloWattHour, 100.0

	expectedConvertedUnit, expectedConvertedValue := Joule, 360_000_000.0
	actualConvertedUnit, actualConvertedValue := ConvertValueToElementaryUnit(unit, value)

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
