package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register("SDM54", NewSDM54Producer)
}

type SDM54Producer struct {
	Opcodes
}

func NewSDM54Producer() Producer {
	/**
	 * Opcodes as defined by Eastron SDM54.
	 * See https://www.eastrongroup.com/eastrongroup/2024/08/21/eastronsdm54seriesusermanualv1.2.pdf
	 */
	ops := Opcodes{
		VoltageL1:        0x0000, // Phase 1 line to neutral volts
		VoltageL2:        0x0002, // Phase 2 line to neutral volts
		VoltageL3:        0x0004, // Phase 3 line to neutral volts
		CurrentL1:        0x0006, // Phase 1 current
		CurrentL2:        0x0008, // Phase 2 current
		CurrentL3:        0x000A, // Phase 3 current
		PowerL1:          0x000C, // Phase 1 active power
		PowerL2:          0x000E, // Phase 2 active power
		PowerL3:          0x0010, // Phase 3 active power
		ApparentPowerL1:  0x0012, // Phase 1 apparent power
		ApparentPowerL2:  0x0014, // Phase 2 apparent power
		ApparentPowerL3:  0x0016, // Phase 3 apparent power
		ReactivePowerL1:  0x0018, // Phase 1 reactive power
		ReactivePowerL2:  0x001A, // Phase 2 reactive power
		ReactivePowerL3:  0x001C, // Phase 3 reactive power
		CosphiL1:         0x001E, // Phase 1 power factor
		CosphiL2:         0x0020, // Phase 2 power factor
		CosphiL3:         0x0022, // Phase 3 power factor
		Power:            0x0034, // Total system power
		ApparentPower:    0x0038, // Total system volt amps.
		ReactivePower:    0x003C, // Total system VAr
		Cosphi:           0x003E, // Total system power factor
		PhaseAngle:       0x0042, // Total system phase angle
		Frequency:        0x0046, // Frequency of supply voltages
		Import:           0x0048, // Total Import kWh
		Export:           0x004A, // Total Export kWh
		THDL1:            0x00EA, // Phase 1 L/N volts THD
		THDL2:            0x00EC, // Phase 2 L/N volts THD
		THDL3:            0x00EE, // Phase 3 L/N volts THD
		THD:              0x00F8, // Average line to neutral volts THD
		Sum:              0x0156, // Total kWh
		ReactiveSum:      0x0158, // Total kVArh
		ImportL1:         0x015A, // L1 import kWh
		ImportL2:         0x015C, // L2 import kWh
		ImportL3:         0x015E, // L3 import kWh
		ExportL1:         0x0160, // L1 export kWh
		ExportL2:         0x0162, // L2 export kWh
		ExportL3:         0x0164, // L3 export kWh
		SumL1:            0x0166, // L1 total kWh
		SumL2:            0x0168, // L2 total kWh
		SumL3:            0x016A, // L3 total kWh
		SumT1:            0x130C, // Tariff 1 total kWh
		SumT2:            0x130E, // Tariff 2 total kWh
		ImportT1:         0x1314, // Tariff 1 import kWh
		ImportT2:         0x1316, // Tariff 2 import kWh
		ExportT1:         0x131C, // Tariff 1 export kWh
		ExportT2:         0x131E, // Tariff 2 export kWh
		ReactiveSumT1:    0x1324, // Tariff 1 total kVArh
		ReactiveSumT2:    0x1326, // Tariff 2 total kVArh
		ReactiveImportT1: 0x132C, // Tariff 1 import kVArh
		ReactiveImportT2: 0x132E, // Tariff 2 import kVArh
		ReactiveExportT1: 0x1334, // Tariff 1 export kVArh
		ReactiveExportT2: 0x1336, // Tariff 2 export kVArh
	}
	return &SDM54Producer{Opcodes: ops}
}

func (p *SDM54Producer) Description() string {
	return "Eastron SDM54"
}

func (p *SDM54Producer) snip(iec Measurement) Operation {
	operation := Operation{
		FuncCode:  ReadInputReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   2,
		IEC61850:  iec,
		Transform: RTUIeee754ToFloat64,
	}
	return operation
}

func (p *SDM54Producer) Probe() Operation {
	return p.snip(VoltageL1)
}

func (p *SDM54Producer) Produce() (res []Operation) {
	for op := range p.Opcodes {
		res = append(res, p.snip(op))
	}

	return res
}
