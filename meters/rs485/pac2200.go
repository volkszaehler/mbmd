package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register("PAC2200", NewPacProducer)
}

type PacProducer struct {
	Opcodes
}

func NewPacProducer() Producer {
	/**
	 * https://github.com/evcc-io/evcc/files/8588580/PAC2200.pdf
	 */
	ops := Opcodes{
		VoltageL1:     1,
		VoltageL2:     3,
		VoltageL3:     5,
		CurrentL1:     13,
		CurrentL2:     15,
		CurrentL3:     17,
		PowerL1:       25,
		PowerL2:       27,
		PowerL3:       29,
		Power:         65,
		ApparentPower: 63,
		ReactivePower: 67,
		CosphiL1:      37,
		CosphiL2:      39,
		CosphiL3:      41,
		Cosphi:        69,
		Frequency:     46,
		Import:        801, // 2w
		Export:        809, // 2w
	}
	return &PacProducer{Opcodes: ops}
}

func (p *PacProducer) Description() string {
	return "Siemens PAC2200"
}

func (p *PacProducer) snip32(iec Measurement) Operation {
	operation := Operation{
		FuncCode:  ReadInputReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   2,
		IEC61850:  iec,
		Transform: RTUIeee754ToFloat64,
	}
	return operation
}

func (p *PacProducer) snip64(iec Measurement, scaler float64) Operation {
	operation := Operation{
		FuncCode:  ReadInputReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   4,
		IEC61850:  iec,
		Transform: MakeScaledTransform(RTUFloat64ToFloat64, scaler),
	}
	return operation
}

func (p *PacProducer) Probe() Operation {
	return p.snip32(VoltageL1)
}

func (p *PacProducer) Produce() (res []Operation) {
	for op := range p.Opcodes {
		switch op {
		case Import, Export:
			res = append(res, p.snip64(op, 1e3))
		default:
			res = append(res, p.snip32(op))
		}
	}

	return res
}
