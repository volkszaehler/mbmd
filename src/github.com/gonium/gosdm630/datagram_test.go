package sdm630_test

import (
	"github.com/gonium/gosdm630"
	"testing"
	"time"
)

func TestQuerySnipMerge(t *testing.T) {
	r := sdm630.Readings{
		Timestamp:      time.Now(),
		Unix:           time.Now().Unix(),
		ModbusDeviceId: 1,
		UniqueId: "Instrument1",
		Power: sdm630.ThreePhaseReadings{
			L1: 1, L2: 2, L3: 3,
		},
		Voltage: sdm630.ThreePhaseReadings{
			L1: 1, L2: 2, L3: 3,
		},
		Current: sdm630.ThreePhaseReadings{
			L1: 4, L2: 5, L3: 6,
		},
		Cosphi: sdm630.ThreePhaseReadings{
			L1: 7, L2: 8, L3: 9,
		},
		Import: sdm630.ThreePhaseReadings{
			L1: 10, L2: 11, L3: 12,
		},
		Export: sdm630.ThreePhaseReadings{
			L1: 13, L2: 14, L3: 15,
		},
	}

	setvalue := float64(230.0)
	var sniptests = []struct {
		snip  sdm630.QuerySnip
		param func(sdm630.Readings) float64
	}{
		{sdm630.QuerySnip{DeviceId: 1, OpCode: sdm630.OpCodeL1Voltage,
			Value: setvalue},
			func(r sdm630.Readings) float64 { return r.Voltage.L1 },
		},
		{sdm630.QuerySnip{DeviceId: 1, OpCode: sdm630.OpCodeL2Voltage,
			Value: setvalue},
			func(r sdm630.Readings) float64 { return r.Voltage.L2 },
		},
		{sdm630.QuerySnip{DeviceId: 1, OpCode: sdm630.OpCodeL3Voltage,
			Value: setvalue},
			func(r sdm630.Readings) float64 { return r.Voltage.L3 },
		},
		{sdm630.QuerySnip{DeviceId: 1, OpCode: sdm630.OpCodeL1Current,
			Value: setvalue},
			func(r sdm630.Readings) float64 { return r.Current.L1 },
		},
		{sdm630.QuerySnip{DeviceId: 1, OpCode: sdm630.OpCodeL2Current,
			Value: setvalue},
			func(r sdm630.Readings) float64 { return r.Current.L2 },
		},
		{sdm630.QuerySnip{DeviceId: 1, OpCode: sdm630.OpCodeL3Current,
			Value: setvalue},
			func(r sdm630.Readings) float64 { return r.Current.L3 },
		},
		{sdm630.QuerySnip{DeviceId: 1, OpCode: sdm630.OpCodeL1Power,
			Value: setvalue},
			func(r sdm630.Readings) float64 { return r.Power.L1 },
		},
		{sdm630.QuerySnip{DeviceId: 1, OpCode: sdm630.OpCodeL2Power,
			Value: setvalue},
			func(r sdm630.Readings) float64 { return r.Power.L2 },
		},
		{sdm630.QuerySnip{DeviceId: 1, OpCode: sdm630.OpCodeL3Power,
			Value: setvalue},
			func(r sdm630.Readings) float64 { return r.Power.L3 },
		},
		{sdm630.QuerySnip{DeviceId: 1, OpCode: sdm630.OpCodeL1Cosphi,
			Value: setvalue},
			func(r sdm630.Readings) float64 { return r.Cosphi.L1 },
		},
		{sdm630.QuerySnip{DeviceId: 1, OpCode: sdm630.OpCodeL2Cosphi,
			Value: setvalue},
			func(r sdm630.Readings) float64 { return r.Cosphi.L2 },
		},
		{sdm630.QuerySnip{DeviceId: 1, OpCode: sdm630.OpCodeL3Cosphi,
			Value: setvalue},
			func(r sdm630.Readings) float64 { return r.Cosphi.L3 },
		},
		{sdm630.QuerySnip{DeviceId: 1, OpCode: sdm630.OpCodeL1Import,
			Value: setvalue},
			func(r sdm630.Readings) float64 { return r.Import.L1 },
		},
		{sdm630.QuerySnip{DeviceId: 1, OpCode: sdm630.OpCodeL2Import,
			Value: setvalue},
			func(r sdm630.Readings) float64 { return r.Import.L2 },
		},
		{sdm630.QuerySnip{DeviceId: 1, OpCode: sdm630.OpCodeL3Import,
			Value: setvalue},
			func(r sdm630.Readings) float64 { return r.Import.L3 },
		},
		{sdm630.QuerySnip{DeviceId: 1, OpCode: sdm630.OpCodeL1Export,
			Value: setvalue},
			func(r sdm630.Readings) float64 { return r.Export.L1 },
		},
		{sdm630.QuerySnip{DeviceId: 1, OpCode: sdm630.OpCodeL2Export,
			Value: setvalue},
			func(r sdm630.Readings) float64 { return r.Export.L2 },
		},
		{sdm630.QuerySnip{DeviceId: 1, OpCode: sdm630.OpCodeL3Export,
			Value: setvalue},
			func(r sdm630.Readings) float64 { return r.Export.L3 },
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
