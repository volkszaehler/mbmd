package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register("SDM72", NewSDM72Producer)
}

type SDM72Producer struct {
	Opcodes
}

func NewSDM72Producer() Producer {
	/**
	 * Opcodes as defined by Eastron SDM72.
	 * See https://data.stromzähler.eu/eastron/SDM72DM-manual.pdf
	 */
	ops := Opcodes{
		Power:  0x0034,
		Import: 0x0048,
		Export: 0x004a,
		Sum:    0x0156,
	}
	return &SDM72Producer{Opcodes: ops}
}

func (p *SDM72Producer) Description() string {
	return "Eastron SDM72"
}

func (p *SDM72Producer) snip(iec Measurement) Operation {
	operation := Operation{
		FuncCode:  ReadInputReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   2,
		IEC61850:  iec,
		Transform: RTUIeee754ToFloat64,
	}
	return operation
}

// This device does not provide voltage data
// so it is not possible to automatically detect the device
func (p *SDM72Producer) Probe() Operation {
	return Operation{}
}

func (p *SDM72Producer) Produce() (res []Operation) {
	for op := range p.Opcodes {
		res = append(res, p.snip(op))
	}

	return res
}
