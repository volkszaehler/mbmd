package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register("JANITZA", NewJanitzaProducer)
}

type JanitzaProducer struct {
	Opcodes
}

func NewJanitzaProducer() Producer {
	/**
	 * Opcodes for Janitza B23.
	 * See https://www.janitza.de/betriebsanleitungen.html?file=files/download/manuals/current/B-Series/janitza-bhb-b2x-en.pdf
	 */
	ops := Opcodes{
		VoltageL1: 0x4A38,
		VoltageL2: 0x4A3A,
		VoltageL3: 0x4A3C,
		CurrentL1: 0x4A44,
		CurrentL2: 0x4A46,
		CurrentL3: 0x4A48,
		PowerL1:   0x4A4C,
		PowerL2:   0x4A4E,
		PowerL3:   0x4A50,
		Power:     0x4A52,
		ImportL1:  0x4A76,
		ImportL2:  0x4A78,
		ImportL3:  0x4A7A,
		Import:    0x4A7C,
		ExportL1:  0x4A7E,
		ExportL2:  0x4A80,
		ExportL3:  0x4A82,
		Export:    0x4A84,
		CosphiL1:  0x4A64,
		CosphiL2:  0x4A66,
		CosphiL3:  0x4A68,
		Frequency: 0x4A6A,
	}
	return &JanitzaProducer{Opcodes: ops}
}

// Description implements Producer interface
func (p *JanitzaProducer) Description() string {
	return "Janitza B-Series meters"
}

func (p *JanitzaProducer) snip(iec Measurement) Operation {
	snip := Operation{
		FuncCode:  ReadHoldingReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   2,
		IEC61850:  iec,
		Transform: RTUIeee754ToFloat64,
	}
	return snip
}

// Probe implements Producer interface
func (p *JanitzaProducer) Probe() Operation {
	return p.snip(VoltageL1)
}

// Produce implements Producer interface
func (p *JanitzaProducer) Produce() (res []Operation) {
	for op := range p.Opcodes {
		res = append(res, p.snip(op))
	}

	return res
}
