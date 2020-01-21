package rs485

// validator checks if value is in range of reference values
type validator struct {
	refs []float64
}

func (v validator) validate(f float64) bool {
	tolerance := 0.1 // 10%
	for _, ref := range v.refs {
		if f >= (1-tolerance)*ref && f <= (1+tolerance)*ref {
			return true
		}
	}
	return false
}
