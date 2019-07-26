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

	res := lhs.add(lhs)
	if v, ok := res.Values[meters.VoltageL1]; v != 440 || !ok {
		t.Errorf("could not combine readings %v", res)
	}
	if v, ok := res.Values[meters.VoltageL2]; v != 760 || !ok {
		t.Errorf("could not combine readings %v", res)
	}

	div := res.divide(2)
	if v, ok := div.Values[meters.VoltageL1]; v != 220 || !ok {
		t.Errorf("could not divide reading %v", div)
	}
	if v, ok := div.Values[meters.VoltageL2]; v != 380 || !ok {
		t.Errorf("could not divide reading %v", div)
	}
}

func TestAverage(t *testing.T) {
	rs := &ReadingSlice{
		&Readings{
			Values: map[meters.Measurement]float64{
				meters.VoltageL1: 220,
				meters.VoltageL2: 380,
			},
		},
		&Readings{
			Values: map[meters.Measurement]float64{
				meters.VoltageL1: 240,
				meters.VoltageL2: 400,
			},
		},
	}

	avg, err := rs.Average()
	if err != nil {
		t.Errorf("could not average readings %v", err)
	}
	if v, ok := avg.Values[meters.VoltageL1]; v != 230 || !ok {
		t.Errorf("could not average readings %v", avg)
	}
	if v, ok := avg.Values[meters.VoltageL2]; v != 390 || !ok {
		t.Errorf("could not average readings %v", avg)
	}
}
