package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register(NewSDM220Producer)
}

const (
	METERTYPE_SDM220 = "SDM220"
)

type SDM220Producer struct {
	Opcodes
}

func NewSDM220Producer() Producer {
	/**
	 * Opcodes as defined by Eastron SDM220.
	 * See https://bg-etech.de/download/manual/SDM220StandardDE.pdf
	 */
	ops := Opcodes{
		Voltage:        0x0000, // 220, 230
		Current:        0x0006, // 220, 230
		Import:         0x0048, // 220, 230
		Export:         0x004a, // 220, 230
		Sum:            0x0156, // 220, 230
		ReactiveSum:    0x0158, // 220
		ReactiveImport: 0x4C,   // 220, 230
		ReactiveExport: 0x4E,   // 220, 230
	}
	return &SDM220Producer{Opcodes: ops}
}

func (p *SDM220Producer) Type() string {
	return METERTYPE_SDM220
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
