package units

// ConvertValueToElementaryUnit converts a sourceUnit and a sourceValue to their elementary unit if possible
// Otherwise, sourceUnit and sourceValue are returned again.
func ConvertValueToElementaryUnit(sourceUnit Unit, sourceValue float64) (Unit, float64) {
	switch sourceUnit {
	case KiloWattHour:
		fallthrough
	case KiloVarHour:
		return Joule, sourceValue * 1_000 * 3_600
	}

	return sourceUnit, sourceValue
}
