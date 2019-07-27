package server

import (
	"testing"

	"github.com/volkszaehler/mbmd/meters"
)

func TestAddDivide(t *testing.T) {
	lhs := &Readings{
		Values: map[meters.Measurement]float64{},
	}
	lhs.Add(QuerySnip{
		MeasurementResult: meters.MeasurementResult{
			Measurement: meters.VoltageL1,
			Value:       220,
		},
	})
	if v, ok := lhs.Values[meters.VoltageL1]; v != 220 || !ok {
		t.Error("could not add reading")
	}

	lhs.Add(QuerySnip{
		MeasurementResult: meters.MeasurementResult{
			Measurement: meters.VoltageL2,
			Value:       380,
		},
	})
	if v, ok := lhs.Values[meters.VoltageL2]; v != 380 || !ok {
		t.Error("could not add reading")
	}
}
