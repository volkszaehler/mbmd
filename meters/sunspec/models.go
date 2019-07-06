package sunspec

import (
	"github.com/andig/gosunspec/models/model101"
	"github.com/andig/gosunspec/models/model103"
	"github.com/andig/gosunspec/models/model124"

	"github.com/volkszaehler/mbmd/meters"
)

var modelPoints = map[int][]string{
	// single phase inverter
	model101.ModelID: {
		model101.A,
		model101.AphA,
		model101.PhVphA,
		model101.Hz,
		model101.W,
		model101.VA,
		model101.VAr,
		model101.PF,
		model101.WH,
		model101.DCA,
		model101.DCV,
		model101.DCW,
		model101.TmpCab,
	},
	// three phase inverter
	model103.ModelID: {
		model103.A,
		model103.AphA,
		model103.AphB,
		model103.AphC,
		model103.PhVphA,
		model103.PhVphB,
		model103.PhVphC,
		model103.Hz,
		model103.W,
		model103.VA,
		model103.VAr,
		model103.PF,
		model103.WH,
		model103.DCA,
		model103.DCV,
		model103.DCW,
		model103.TmpCab,
	},
	// storage
	model124.ModelID: {
		model124.ChaState,
		model124.InBatV,
	},
}

var opcodeMap = map[string]meters.Measurement{
	model103.A:        meters.Current,
	model103.AphA:     meters.CurrentL1,
	model103.AphB:     meters.CurrentL2,
	model103.AphC:     meters.CurrentL3,
	model103.PhVphA:   meters.VoltageL1,
	model103.PhVphB:   meters.VoltageL2,
	model103.PhVphC:   meters.VoltageL3,
	model103.Hz:       meters.Frequency,
	model103.W:        meters.Power,
	model103.VA:       meters.ApparentPower,
	model103.VAr:      meters.ReactivePower,
	model103.PF:       meters.Cosphi,
	model103.WH:       meters.Export,
	model103.DCA:      meters.DCCurrent,
	model103.DCV:      meters.DCVoltage,
	model103.DCW:      meters.DCPower,
	model103.TmpCab:   meters.HeatSinkTemp,
	model124.ChaState: meters.ChargeState,
	model124.InBatV:   meters.BatteryVoltage,
}
