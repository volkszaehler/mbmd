package sunspec

import (
	sunspec "github.com/andig/gosunspec"
	"github.com/andig/gosunspec/models/model101"
	"github.com/andig/gosunspec/models/model103"
	"github.com/andig/gosunspec/models/model124"
	"github.com/andig/gosunspec/models/model160"
	"github.com/andig/gosunspec/models/model203"

	"github.com/volkszaehler/mbmd/meters"
)

// model- and block-specific opcode measurement mapping
var modelMap = map[sunspec.ModelId]map[int]map[string]meters.Measurement{
	// single phase inverter
	model101.ModelID: {
		0: {
			model101.A:      meters.Current,
			model101.AphA:   meters.CurrentL1,
			model101.PhVphA: meters.VoltageL1,
			model101.Hz:     meters.Frequency,
			model101.W:      meters.Power,
			model101.VA:     meters.ApparentPower,
			model101.VAr:    meters.ReactivePower,
			model101.PF:     meters.Cosphi,
			model101.WH:     meters.Export,
			model101.DCA:    meters.DCCurrent,
			model101.DCV:    meters.DCVoltage,
			model101.DCW:    meters.DCPower,
			model101.TmpCab: meters.HeatSinkTemp,
		},
	},
	// three phase inverter
	model103.ModelID: {
		0: {
			model103.A:      meters.Current,
			model103.AphA:   meters.CurrentL1,
			model103.AphB:   meters.CurrentL2,
			model103.AphC:   meters.CurrentL3,
			model103.PhVphA: meters.VoltageL1,
			model103.PhVphB: meters.VoltageL2,
			model103.PhVphC: meters.VoltageL3,
			model103.Hz:     meters.Frequency,
			model103.W:      meters.Power,
			model103.VA:     meters.ApparentPower,
			model103.VAr:    meters.ReactivePower,
			model103.PF:     meters.Cosphi,
			model103.WH:     meters.Export,
			model103.DCA:    meters.DCCurrent,
			model103.DCV:    meters.DCVoltage,
			model103.DCW:    meters.DCPower,
			model103.TmpCab: meters.HeatSinkTemp,
		},
	},
	model160.ModelID: {
		1: {
			model160.DCA:  meters.DCCurrentS1,
			model160.DCV:  meters.DCVoltageS1,
			model160.DCW:  meters.DCPowerS1,
			model160.DCWH: meters.DCEnergyS1,
		},
		2: {
			model160.DCA:  meters.DCCurrentS2,
			model160.DCV:  meters.DCVoltageS2,
			model160.DCW:  meters.DCPowerS2,
			model160.DCWH: meters.DCEnergyS2,
		},
		3: {
			model160.DCA:  meters.DCCurrentS3,
			model160.DCV:  meters.DCVoltageS3,
			model160.DCW:  meters.DCPowerS3,
			model160.DCWH: meters.DCEnergyS3,
		},
	},
	// wye-connect three phase (abcn) meter
	model203.ModelID: {
		0: {
			model203.A:           meters.Current,
			model203.AphA:        meters.CurrentL1,
			model203.AphB:        meters.CurrentL2,
			model203.AphC:        meters.CurrentL3,
			model203.Hz:          meters.Frequency,
			model203.PF:          meters.Cosphi,
			model203.PFphA:       meters.CosphiL1,
			model203.PFphB:       meters.CosphiL2,
			model203.PFphC:       meters.CosphiL3,
			model203.PhV:         meters.Voltage,
			model203.PhVphA:      meters.VoltageL1,
			model203.PhVphB:      meters.VoltageL2,
			model203.PhVphC:      meters.VoltageL3,
			model203.TotWhExp:    meters.Export,
			model203.TotWhExpPhA: meters.ExportL1,
			model203.TotWhExpPhB: meters.ExportL2,
			model203.TotWhExpPhC: meters.ExportL3,
			model203.TotWhImp:    meters.Import,
			model203.TotWhImpPhA: meters.ImportL1,
			model203.TotWhImpPhB: meters.ImportL2,
			model203.TotWhImpPhC: meters.ImportL3,
			model203.VA:          meters.ApparentPower,
			model203.VAR:         meters.ReactivePower,
			model203.VARphA:      meters.ReactivePowerL1,
			model203.VARphB:      meters.ReactivePowerL2,
			model203.VARphC:      meters.ReactivePowerL3,
			model203.VAphA:       meters.CurrentL1,
			model203.VAphB:       meters.CurrentL2,
			model203.VAphC:       meters.CurrentL3,
			model203.W:           meters.Power,
			model203.WphA:        meters.PowerL1,
			model203.WphB:        meters.PowerL2,
			model203.WphC:        meters.PowerL3,
		},
	},
	// storage
	model124.ModelID: {
		0: {
			model124.ChaState: meters.ChargeState,
			model124.InBatV:   meters.BatteryVoltage,
		},
	},
}

var dividerMap = map[meters.Measurement]float64{
	meters.Export:     1000,
	meters.ExportL1:   1000,
	meters.ExportL2:   1000,
	meters.ExportL3:   1000,
	meters.Import:     1000,
	meters.ImportL1:   1000,
	meters.ImportL2:   1000,
	meters.ImportL3:   1000,
	meters.DCEnergyS1: 1000,
	meters.DCEnergyS2: 1000,
	meters.DCEnergyS3: 1000,
}
