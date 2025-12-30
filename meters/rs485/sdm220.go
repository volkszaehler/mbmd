package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register("SDM220", NewSDM220Producer)
}

type SDM220Producer struct {
	Opcodes
}

func NewSDM220Producer() Producer {
	/**
	 * Opcodes as defined by Eastron SDM220.
	 * See https://bg-etech.de/download/manual/SDM220StandardDE.pdf
	 */
	ops := Opcodes{
		Voltage:        0x0000, // Line to neutral volts
		Current:        0x0006, // Current
		Power:          0x000C, // Active power
		ApparentPower:  0x0012, // Apparent power
		ReactivePower:  0x0018, // Reactive power
		Cosphi:         0x001E, // Power factor
		PhaseAngle:     0x0024, // Phase angle
		Frequency:      0x0046, // Frequency of supply voltage
		Import:         0x0048, // Total Import kWh
		Export:         0x004A, // Total Export kWh
		ReactiveImport: 0x004C, // Total Import kVArh
		ReactiveExport: 0x004E, // Total Export kVArh
		Sum:            0x0156, // Total kWh
		ReactiveSum:    0x0158, // Total kVArh
	}
	return &SDM220Producer{Opcodes: ops}
}

func (p *SDM220Producer) Description() string {
	return "Eastron SDM220"
}

func (p *SDM220Producer) snip(iec Measurement) Operation {
	operation := Operation{
		FuncCode:  ReadInputReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   2,
		IEC61850:  iec,
		Transform: RTUIeee754ToFloat64,
	}
	return operation
}

func (p *SDM220Producer) Probe() Operation {
	return p.snip(Voltage)
}

func (p *SDM220Producer) Produce() (res []Operation) {
	for op := range p.Opcodes {
		res = append(res, p.snip(op))
	}

	return res
}
