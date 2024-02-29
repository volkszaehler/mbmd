package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register("X961A", NewX961AProducer)
}

type X961AProducer struct {
	Opcodes
}

func NewX961AProducer() Producer {
	/**
	 * Opcodes as defined by Eastron SMART X96-1A.
	 * See https://www.eastroneurope.com/images/uploads/products/protocol/SMART_X96-1A_MODBUS_Protocol_v1.0.pdf
	 */
	ops := Opcodes{
		Frequency:        0x0046, // Frequency of supply voltages.
		Current:          0x0030, // Sum of line currents.
		CurrentL1:        0x0006, // Phase 1 current.
		CurrentL2:        0x0008, // Phase 2 current.
		CurrentL3:        0x000A, // Phase 3 current.
		VoltageL1:        0x0000, // Phase 1 line to neutral volts.
		VoltageL2:        0x0002, // Phase 2 line to neutral volts.
		VoltageL3:        0x0004, // Phase 3 line to neutral volts.
		Power:            0x0034, // Total system power.
		PowerL1:          0x000C, // Phase 1 active power.
		PowerL2:          0x000E, // Phase 2 active power.
		PowerL3:          0x0010, // Phase 3 active power
		ReactivePower:    0x003C, // Total system VAr.
		ReactivePowerL1:  0x0018, // Phase 1 reactive power.
		ReactivePowerL2:  0x001A, // Phase 2 reactive power.
		ReactivePowerL3:  0x001C, // Phase 3 reactive power.
		ApparentPowerL1:  0x0012, // Phase 1 apparent power.
		ApparentPowerL2:  0x0014, // Phase 2 apparent power.
		ApparentPowerL3:  0x0016, // Phase 3 apparent power.
		Cosphi:           0x003E, // Total system power factor.
		CosphiL1:         0x001E, // Phase 1 power factor.
		CosphiL2:         0x0020, // Phase 2 power factor.
		CosphiL3:         0x0022, // Phase 3 power factor.
		THD:              0x00F8, // Average line to neutral volts THD.
		THDL1:            0x00EA, // Phase 1 L/N volts THD
		THDL2:            0x00EC, // Phase 2 L/N volts THD
		THDL3:            0x00EE, // Phase 3 L/N volts THD
		Sum:              0x0156, // Total kWh
		SumL1:            0x0166, // L1 total kWh
		SumL2:            0x0168, // L2 total kWh
		SumL3:            0x016A, // L3 total kWh
		ImportL1:         0x015A, // L1 import kWh
		ImportL2:         0x015C, // L2 import kWh
		ImportL3:         0x015E, // L3 import kWh
		ExportL1:         0x0160, // L1 export kWh
		ExportL2:         0x0162, // L2 export kWh
		ExportL3:         0x0164, // L3 export kWh
		ReactiveSum:      0x0158, // Total kVArh
		ReactiveSumL1:    0x0178, // L1 total kVArh
		ReactiveSumL2:    0x017A, // L2 total kVArh
		ReactiveSumL3:    0x017C, // L3 total kVArh
		ReactiveImportL1: 0x016C, // L1 import kVArh
		ReactiveImportL2: 0x016E, // L2 import kVArh
		ReactiveImportL3: 0x0170, // L3 import kVArh
		ReactiveExportL1: 0x0172, // L1 export kVArh
		ReactiveExportL2: 0x0174, // L2 export kVArh
		ReactiveExportL3: 0x0176, // L3 export kVArh
		PhaseAngle:       0x0042, // Total system phase angle.
	}
	return &X961AProducer{Opcodes: ops}
}

func (p *X961AProducer) Description() string {
	return "Eastron SMART X96-1A"
}

func (p *X961AProducer) snip(iec Measurement) Operation {
	operation := Operation{
		FuncCode:  ReadInputReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   2,
		IEC61850:  iec,
		Transform: RTUIeee754ToFloat64,
	}
	return operation
}

func (p *X961AProducer) Probe() Operation {
	return p.snip(VoltageL1)
}

func (p *X961AProducer) Produce() (res []Operation) {
	for op := range p.Opcodes {
		res = append(res, p.snip(op))
	}

	return res
}
