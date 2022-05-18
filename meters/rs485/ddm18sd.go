package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register("DDM", NewDDMProducer)
}

type DDMProducer struct {
	Opcodes
}

func NewDDMProducer() Producer {
	ops := Opcodes{
		VoltageL1:       0x0000,
		CurrentL1:       0x0008,
		PowerL1:         0x0012,
		ReactivePowerL1: 0x001A,
		Cosphi:          0x002A,
		Frequency:       0x0036,
		Sum:             0x0100,
		ReactiveSum:     0x0400,
	}
	return &DDMProducer{Opcodes: ops}
}

func (p *DDMProducer) Description() string {
	return "DDM18SD"
}

func (p *DDMProducer) snip(iec Measurement) Operation {
	operation := Operation{
		FuncCode:  ReadInputReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   2,
		IEC61850:  iec,
		Transform: RTUIeee754ToFloat64,
	}
	return operation
}

func (p *DDMProducer) Probe() Operation {
	return p.snip(VoltageL1)
}

func (p *DDMProducer) Produce() (res []Operation) {
	for op := range p.Opcodes {
		res = append(res, p.snip(op))
	}

	return res
}
