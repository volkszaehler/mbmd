package server

import (
	"testing"
	"time"

	. "github.com/volkszaehler/mbmd/meters"
)

func TestQuerySnipMerge(t *testing.T) {
	r := Readings{
		Timestamp: time.Now(),
		Unix:      time.Now().Unix(),
		DeviceId:  1,
		UniqueId:  "Instrument1",
		Power: ThreePhaseReadings{
			L1: F2fp(1.0), L2: F2fp(2.0), L3: F2fp(3.0),
		},
		Voltage: ThreePhaseReadings{
			L1: F2fp(1.0), L2: F2fp(2.0), L3: F2fp(3.0),
		},
		Current: ThreePhaseReadings{
			L1: F2fp(4.0), L2: F2fp(5.0), L3: F2fp(6.0),
		},
		Cosphi: ThreePhaseReadings{
			L1: F2fp(7.0), L2: F2fp(8.0), L3: F2fp(9.0),
		},
		Import: ThreePhaseReadings{
			L1: F2fp(10.0), L2: F2fp(11.0), L3: F2fp(12.0),
		},
		Export: ThreePhaseReadings{
			L1: F2fp(13.0), L2: F2fp(14.0), L3: F2fp(15.0),
		},
	}

	setvalue := float64(230.0)
	var sniptests = []struct {
		snip  QuerySnip
		param func(Readings) float64
	}{
		{
			QuerySnip{
				Device: "dev",
				MeasurementResult: MeasurementResult{
					Value:       setvalue,
					Measurement: VoltageL1,
				},
			},
			func(r Readings) float64 { return Fp2f(r.Voltage.L1) },
		},
		{
			QuerySnip{
				Device: "dev",
				MeasurementResult: MeasurementResult{
					Value:       setvalue,
					Measurement: VoltageL2,
				},
			},
			func(r Readings) float64 { return Fp2f(r.Voltage.L2) },
		},
		{
			QuerySnip{
				Device: "dev",
				MeasurementResult: MeasurementResult{
					Value:       setvalue,
					Measurement: VoltageL3,
				},
			},
			func(r Readings) float64 { return Fp2f(r.Voltage.L3) },
		},
		{
			QuerySnip{
				Device: "dev",
				MeasurementResult: MeasurementResult{
					Value:       setvalue,
					Measurement: CurrentL1,
				},
			},
			func(r Readings) float64 { return Fp2f(r.Current.L1) },
		},
		{
			QuerySnip{
				Device: "dev",
				MeasurementResult: MeasurementResult{
					Value:       setvalue,
					Measurement: CurrentL2,
				},
			},
			func(r Readings) float64 { return Fp2f(r.Current.L2) },
		},
		{
			QuerySnip{
				Device: "dev",
				MeasurementResult: MeasurementResult{
					Value:       setvalue,
					Measurement: CurrentL3,
				},
			},
			func(r Readings) float64 { return Fp2f(r.Current.L3) },
		},
		{
			QuerySnip{
				Device: "dev",
				MeasurementResult: MeasurementResult{
					Value:       setvalue,
					Measurement: PowerL1,
				},
			},
			func(r Readings) float64 { return Fp2f(r.Power.L1) },
		},
		{
			QuerySnip{
				Device: "dev",
				MeasurementResult: MeasurementResult{
					Value:       setvalue,
					Measurement: PowerL2,
				},
			},
			func(r Readings) float64 { return Fp2f(r.Power.L2) },
		},
		{
			QuerySnip{
				Device: "dev",
				MeasurementResult: MeasurementResult{
					Value:       setvalue,
					Measurement: PowerL3,
				},
			},
			func(r Readings) float64 { return Fp2f(r.Power.L3) },
		},
		{
			QuerySnip{
				Device: "dev",
				MeasurementResult: MeasurementResult{
					Value:       setvalue,
					Measurement: CosphiL1,
				},
			},
			func(r Readings) float64 { return Fp2f(r.Cosphi.L1) },
		},
		{
			QuerySnip{
				Device: "dev",
				MeasurementResult: MeasurementResult{
					Value:       setvalue,
					Measurement: CosphiL2,
				},
			},
			func(r Readings) float64 { return Fp2f(r.Cosphi.L2) },
		},
		{
			QuerySnip{
				Device: "dev",
				MeasurementResult: MeasurementResult{
					Value:       setvalue,
					Measurement: CosphiL3,
				},
			},
			func(r Readings) float64 { return Fp2f(r.Cosphi.L3) },
		},
		{
			QuerySnip{
				Device: "dev",
				MeasurementResult: MeasurementResult{
					Value:       setvalue,
					Measurement: ImportL1,
				},
			},
			func(r Readings) float64 { return Fp2f(r.Import.L1) },
		},
		{
			QuerySnip{
				Device: "dev",
				MeasurementResult: MeasurementResult{
					Value:       setvalue,
					Measurement: ImportL2,
				},
			},
			func(r Readings) float64 { return Fp2f(r.Import.L2) },
		},
		{
			QuerySnip{
				Device: "dev",
				MeasurementResult: MeasurementResult{
					Value:       setvalue,
					Measurement: ImportL3,
				},
			},
			func(r Readings) float64 { return Fp2f(r.Import.L3) },
		},
		{
			QuerySnip{
				Device: "dev",
				MeasurementResult: MeasurementResult{
					Value:       setvalue,
					Measurement: ExportL1,
				},
			},
			func(r Readings) float64 { return Fp2f(r.Export.L1) },
		},
		{
			QuerySnip{
				Device: "dev",
				MeasurementResult: MeasurementResult{
					Value:       setvalue,
					Measurement: ExportL2,
				},
			},
			func(r Readings) float64 { return Fp2f(r.Export.L2) },
		},
		{
			QuerySnip{
				Device: "dev",
				MeasurementResult: MeasurementResult{
					Value:       setvalue,
					Measurement: ExportL3,
				},
			},
			func(r Readings) float64 { return Fp2f(r.Export.L3) },
		},
	}

	for _, test := range sniptests {
		r.MergeSnip(test.snip)
		if test.param(r) != setvalue {
			t.Errorf("Merge of querysnip failed: Expected %.2f, got %.2f",
				setvalue, test.param(r))
		}
	}
}
