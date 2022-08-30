package sunspec

import (
	sunspec "github.com/andig/gosunspec"
	"github.com/andig/gosunspec/models/model101"
	"github.com/andig/gosunspec/models/model103"
	"github.com/andig/gosunspec/models/model111"
	"github.com/andig/gosunspec/models/model113"
	"github.com/andig/gosunspec/models/model124"
	"github.com/andig/gosunspec/models/model160"
	"github.com/andig/gosunspec/models/model201"
	"github.com/andig/gosunspec/models/model203"
	"github.com/andig/gosunspec/models/model211"
	"github.com/andig/gosunspec/models/model213"

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
	// single phase inverter - float
	model111.ModelID: {
		0: {
			model111.A:      meters.Current,
			model111.AphA:   meters.CurrentL1,
			model111.PhVphA: meters.VoltageL1,
			model111.Hz:     meters.Frequency,
			model111.W:      meters.Power,
			model111.VA:     meters.ApparentPower,
			model111.VAr:    meters.ReactivePower,
			model111.PF:     meters.Cosphi,
			model111.WH:     meters.Export,
			model111.DCA:    meters.DCCurrent,
			model111.DCV:    meters.DCVoltage,
			model111.DCW:    meters.DCPower,
			model111.TmpCab: meters.HeatSinkTemp,
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
	// three phase inverter - float
	model113.ModelID: {
		0: {
			model113.A:      meters.Current,
			model113.AphA:   meters.CurrentL1,
			model113.AphB:   meters.CurrentL2,
			model113.AphC:   meters.CurrentL3,
			model113.PhVphA: meters.VoltageL1,
			model113.PhVphB: meters.VoltageL2,
			model113.PhVphC: meters.VoltageL3,
			model113.Hz:     meters.Frequency,
			model113.W:      meters.Power,
			model113.VA:     meters.ApparentPower,
			model113.VAr:    meters.ReactivePower,
			model113.PF:     meters.Cosphi,
			model113.WH:     meters.Export,
			model113.DCA:    meters.DCCurrent,
			model113.DCV:    meters.DCVoltage,
			model113.DCW:    meters.DCPower,
			model113.TmpCab: meters.HeatSinkTemp,
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
		4: {
			model160.DCA:  meters.DCCurrentS4,
			model160.DCV:  meters.DCVoltageS4,
			model160.DCW:  meters.DCPowerS4,
			model160.DCWH: meters.DCEnergyS4,
		},
	},
	// single phase (AN or AB) meter
	model201.ModelID: {
		0: {
			model201.A:        meters.Current,
			model201.Hz:       meters.Frequency,
			model201.PF:       meters.Cosphi,
			model201.PhV:      meters.Voltage,
			model201.TotWhExp: meters.Export,
			model201.TotWhImp: meters.Import,
			model201.VA:       meters.ApparentPower,
			model201.VAR:      meters.ReactivePower,
			model201.W:        meters.Power,
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
			model203.VAphA:       meters.ApparentPowerL1,
			model203.VAphB:       meters.ApparentPowerL2,
			model203.VAphC:       meters.ApparentPowerL3,
			model203.VAR:         meters.ReactivePower,
			model203.VARphA:      meters.ReactivePowerL1,
			model203.VARphB:      meters.ReactivePowerL2,
			model203.VARphC:      meters.ReactivePowerL3,
			model203.W:           meters.Power,
			model203.WphA:        meters.PowerL1,
			model203.WphB:        meters.PowerL2,
			model203.WphC:        meters.PowerL3,
		},
	},
	// single phase (AN or AB) meter - float
	model211.ModelID: {
		0: {
			model211.A:        meters.Current,
			model211.Hz:       meters.Frequency,
			model211.PF:       meters.Cosphi,
			model211.PhV:      meters.Voltage,
			model211.TotWhExp: meters.Export,
			model211.TotWhImp: meters.Import,
			model211.VA:       meters.ApparentPower,
			model211.VAR:      meters.ReactivePower,
			model211.W:        meters.Power,
		},
	},
	// wye-connect three phase (abcn) meter - float
	model213.ModelID: {
		0: {
			model213.A:           meters.Current,
			model213.AphA:        meters.CurrentL1,
			model213.AphB:        meters.CurrentL2,
			model213.AphC:        meters.CurrentL3,
			model213.Hz:          meters.Frequency,
			model213.PF:          meters.Cosphi,
			model213.PFphA:       meters.CosphiL1,
			model213.PFphB:       meters.CosphiL2,
			model213.PFphC:       meters.CosphiL3,
			model213.PhV:         meters.Voltage,
			model213.PhVphA:      meters.VoltageL1,
			model213.PhVphB:      meters.VoltageL2,
			model213.PhVphC:      meters.VoltageL3,
			model213.TotWhExp:    meters.Export,
			model213.TotWhExpPhA: meters.ExportL1,
			model213.TotWhExpPhB: meters.ExportL2,
			model213.TotWhExpPhC: meters.ExportL3,
			model213.TotWhImp:    meters.Import,
			model213.TotWhImpPhA: meters.ImportL1,
			model213.TotWhImpPhB: meters.ImportL2,
			model213.TotWhImpPhC: meters.ImportL3,
			model213.VA:          meters.ApparentPower,
			model213.VAphA:       meters.ApparentPowerL1,
			model213.VAphB:       meters.ApparentPowerL2,
			model213.VAphC:       meters.ApparentPowerL3,
			model213.VAR:         meters.ReactivePower,
			model213.VARphA:      meters.ReactivePowerL1,
			model213.VARphB:      meters.ReactivePowerL2,
			model213.VARphC:      meters.ReactivePowerL3,
			model213.W:           meters.Power,
			model213.WphA:        meters.PowerL1,
			model213.WphB:        meters.PowerL2,
			model213.WphC:        meters.PowerL3,
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
