package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register("SDM", NewSDMProducer)
}

type SDMProducer struct {
	Opcodes
}

func NewSDMProducer() Producer {
	/**
	 * Opcodes as defined by Eastron SDM630.
	 * See https://www.eastroneurope.com/images/uploads/products/protocol/SDM630_MODBUS_Protocol.pdf
	 * This is to a large extent a superset of all SDM devices, however there are
	 * subtle differences (see 220, 230). Some opcodes might not work on some devices.
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
		ReactiveImport:   0x004C, // Total Import kVArh
		ReactiveExport:   0x004E, // Total Export kVArh
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
		ReactiveImportL1: 0x016C, // L1 import kVArh
		ReactiveImportL2: 0x016E, // L2 import kVArh
		ReactiveImportL3: 0x0170, // L3 import kVArh
		ReactiveExportL1: 0x0172, // L1 export kVArh
		ReactiveExportL2: 0x0174, // L2 export kVArh
		ReactiveExportL3: 0x0176, // L3 export kVArh
		ReactiveSumL1:    0x0178, // L1 total kVArh
		ReactiveSumL2:    0x017A, // L2 total kVArh
		ReactiveSumL3:    0x017C, // L3 total kVArh
	}
	return &SDMProducer{Opcodes: ops}
}

func (p *SDMProducer) Description() string {
	return "Eastron SDM630"
}

func (p *SDMProducer) snip(iec Measurement) Operation {
	operation := Operation{
		FuncCode:  ReadInputReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   2,
		IEC61850:  iec,
		Transform: RTUIeee754ToFloat64,
	}
	return operation
}

func (p *SDMProducer) Probe() Operation {
	return p.snip(VoltageL1)
}

func (p *SDMProducer) Produce() (res []Operation) {
	for op := range p.Opcodes {
		res = append(res, p.snip(op))
	}

	return res
}
