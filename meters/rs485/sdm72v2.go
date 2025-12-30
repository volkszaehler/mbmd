package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register("SDM72V2", NewSDM72V2Producer)
}

type SDM72V2Producer struct {
	Opcodes
}

func NewSDM72V2Producer() Producer {
	/**
	 * Opcodes as defined by Eastron SDM72DMv2.
	 * https://stromzähler.eu/media/5f/b6/aa/1696582672/sdm72dm-v2.pdf
	 */
	ops := Opcodes{
		VoltageL1:       0x0000, // Phase 1 line to neutral volts
		VoltageL2:       0x0002, // Phase 2 line to neutral volts
		VoltageL3:       0x0004, // Phase 3 line to neutral volts
		CurrentL1:       0x0006, // Phase 1 current
		CurrentL2:       0x0008, // Phase 2 current
		CurrentL3:       0x000A, // Phase 3 current
		PowerL1:         0x000C, // Phase 1 active power
		PowerL2:         0x000E, // Phase 2 active power
		PowerL3:         0x0010, // Phase 3 active power
		ApparentPowerL1: 0x0012, // Phase 1 apparent power
		ApparentPowerL2: 0x0014, // Phase 2 apparent power
		ApparentPowerL3: 0x0016, // Phase 3 apparent power
		ReactivePowerL1: 0x0018, // Phase 1 reactive power
		ReactivePowerL2: 0x001A, // Phase 2 reactive power
		ReactivePowerL3: 0x001C, // Phase 3 reactive power
		CosphiL1:        0x001E, // Phase 1 power factor
		CosphiL2:        0x0020, // Phase 2 power factor
		CosphiL3:        0x0022, // Phase 3 power factor
		Power:           0x0034, // Total system power
		ApparentPower:   0x0038, // Total system volt amps.
		ReactivePower:   0x003C, // Total system VAr
		Cosphi:          0x003E, // Total system power factor
		Frequency:       0x0046, // Frequency of supply voltages
		Import:          0x0048, // Total Import kWh
		Export:          0x004A, // Total Export kWh
		Sum:             0x0156, // Total kWh
		ReactiveSum:     0x0158, // Total kVArh
	}
	return &SDM72V2Producer{Opcodes: ops}
}

func (p *SDM72V2Producer) Description() string {
	return "Eastron SDM72 v2"
}

func (p *SDM72V2Producer) snip(iec Measurement) Operation {
	operation := Operation{
		FuncCode:  ReadInputReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   2,
		IEC61850:  iec,
		Transform: RTUIeee754ToFloat64,
	}
	return operation
}

func (p *SDM72V2Producer) Probe() Operation {
	return p.snip(VoltageL1)
}

func (p *SDM72V2Producer) Produce() (res []Operation) {
	for op := range p.Opcodes {
		res = append(res, p.snip(op))
	}

	return res
}
