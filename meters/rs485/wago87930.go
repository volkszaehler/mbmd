package rs485

import (
	. "github.com/volkszaehler/mbmd/meters"
)

func init() {
	Register("WAGO87930", NewWago87930Producer)
}

type Wago87930Producer struct {
	Opcodes
}

func NewWago87930Producer() Producer {
	/**
	 * Registers as defined by Wago 879-30X0
	 * See https://www.wago.com/de-en/d/5937710
	 */
	ops := Opcodes{
		Voltage:         0x5000,
		VoltageL1:       0x5002,
		VoltageL2:       0x5004,
		VoltageL3:       0x5006,
		Frequency:       0x5008,
		Current:         0x500a,
		CurrentL1:       0x500c,
		CurrentL2:       0x500e,
		CurrentL3:       0x5010,
		Power:           0x5012,
		PowerL1:         0x5014,
		PowerL2:         0x5016,
		PowerL3:         0x5018,
		ReactivePower:   0x501a,
		ReactivePowerL1: 0x501c,
		ReactivePowerL2: 0x501e,
		ReactivePowerL3: 0x5020,
		ApparentPower:   0x5022,
		ApparentPowerL1: 0x5024,
		ApparentPowerL2: 0x5026,
		ApparentPowerL3: 0x5028,
		Import:          0x600c,
		ImportL1:        0x6012,
		ImportL2:        0x6014,
		ImportL3:        0x6016,
		Export:          0x6018,
		ExportL1:        0x601e,
		ExportL2:        0x6020,
		ExportL3:        0x6022,
		CosphiL1:        0x502c,
		CosphiL2:        0x502e,
		CosphiL3:        0x5030,
		Cosphi:          0x502a,
	}
	return &Wago87930Producer{Opcodes: ops}
}

func (p *Wago87930Producer) Description() string {
	return "Wago 879-30XX"
}

func (p *Wago87930Producer) snip(iec Measurement) Operation {
	operation := Operation{
		FuncCode:  ReadHoldingReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   2,
		IEC61850:  iec,
		Transform: RTUIeee754ToFloat64,
	}
	return operation
}

func (p *Wago87930Producer) Probe() Operation {
	return p.snip(Voltage)
}

func (p *Wago87930Producer) Produce() (res []Operation) {
	for op := range p.Opcodes {
		res = append(res, p.snip(op))
	}

	return res
}
