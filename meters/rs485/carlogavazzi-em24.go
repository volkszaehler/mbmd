package rs485

import (
	"github.com/grid-x/modbus"
	"github.com/volkszaehler/mbmd/encoding"
	. "github.com/volkszaehler/mbmd/meters"
)

func init() {
	Register("CGEM24", NewCarloGavazziEM24Producer)
}

type CarloGavazziEM24Producer struct {
	ops_em24    Opcodes
	ops_em24_e1 Opcodes
	t           int
}

func NewCarloGavazziEM24Producer() Producer {
	/***
	 * Note: Carlo Gavazzi EM24 (RS-485)
	 * Doc for EM24: https://www.ccontrols.com/support/dp/CarloGavazziEM24.pdf
	 */
	ops_em24 := Opcodes{
		VoltageL1: 0x00,
		VoltageL2: 0x02,
		VoltageL3: 0x04,
		CurrentL1: 0x0C,
		CurrentL2: 0x0E,
		CurrentL3: 0x10,
		PowerL1:   0x12,
		PowerL2:   0x14,
		PowerL3:   0x16,
		Power:     0x28,
		CosphiL1:  0x32,
		CosphiL2:  0x33,
		CosphiL3:  0x34,
		Cosphi:    0x35,
		Frequency: 0x37,
		Import:    0x42,
		ImportL1:  0x46,
		ImportL2:  0x48,
		ImportL3:  0x4A,
		Export:    0x5C,
	}

	ops_em24_e1 := Opcodes{
		VoltageL1: 0x00,
		VoltageL2: 0x02,
		VoltageL3: 0x04,
		CurrentL1: 0x0C,
		CurrentL2: 0x0E,
		CurrentL3: 0x10,
		PowerL1:   0x12,
		PowerL2:   0x14,
		PowerL3:   0x16,
		Power:     0x28,
		CosphiL1:  0x2E,
		CosphiL2:  0x2F,
		CosphiL3:  0x30,
		Cosphi:    0x31,
		Frequency: 0x33,
		Import:    0x34,
		ImportL1:  0x40,
		ImportL2:  0x42,
		ImportL3:  0x44,
		Export:    0x4E,
	}
	return &CarloGavazziEM24Producer{ops_em24, ops_em24_e1, 0}
}

// Initialize implements Producer interface
func (p *CarloGavazziEM24Producer) Initialize(client modbus.Client) {
	var bytes []byte
	var err error
	bytes, err = client.ReadHoldingRegisters(0x00b, 1)
	if err == nil {
		t := encoding.Uint16(bytes)
		if (t >= 45) && (t <= 48) {
			// EM24
			p.t = 0
		}
		if (t >= 1648) && (t <= 1653) {
			// EM24_E1
			p.t = 1
		}
	}
}

// Description implements Producer interface
func (p *CarloGavazziEM24Producer) Description() string {
	if p.t == 0 {
		return "Carlo Gavazzi EM24"
	} else {
		return "Carlo Gavazzi EM24_E1"
	}
}

func (p *CarloGavazziEM24Producer) opCode(iec Measurement) uint16 {
	if p.t == 0 {
		return p.ops_em24.Opcode(iec)
	} else {
		return p.ops_em24_e1.Opcode(iec)
	}
}

func (p *CarloGavazziEM24Producer) snip16(iec Measurement, scaler ...float64) Operation {
	transform := RTUInt16ToFloat64 // default conversion
	if len(scaler) > 0 {
		transform = MakeScaledTransform(transform, scaler[0])
	}
	operation := Operation{
		FuncCode:  ReadInputReg,
		OpCode:    p.opCode(iec),
		ReadLen:   1,
		IEC61850:  iec,
		Transform: transform,
	}

	return operation
}

func (p *CarloGavazziEM24Producer) snip32(iec Measurement, scaler ...float64) Operation {
	transform := RTUInt32ToFloat64Swapped // default conversion
	if len(scaler) > 0 {
		transform = MakeScaledTransform(transform, scaler[0])
	}
	operation := Operation{
		FuncCode:  ReadInputReg,
		OpCode:    p.opCode(iec),
		ReadLen:   2,
		IEC61850:  iec,
		Transform: transform,
	}

	return operation
}

// Probe implements Producer interface
func (p *CarloGavazziEM24Producer) Probe() Operation {
	return p.snip32(VoltageL1, 10)
}

// Produce implements Producer interface
func (p *CarloGavazziEM24Producer) Produce() (res []Operation) {
	for _, op := range []Measurement{
		VoltageL1, VoltageL2, VoltageL3,
	} {
		res = append(res, p.snip32(op, 10))
	}

	for _, op := range []Measurement{
		CurrentL1, CurrentL2, CurrentL3,
	} {
		res = append(res, p.snip32(op, 1000))
	}

	for _, op := range []Measurement{
		Cosphi, CosphiL1, CosphiL2, CosphiL3,
	} {
		res = append(res, p.snip16(op, 1000))
	}

	for _, op := range []Measurement{
		Frequency,
	} {
		res = append(res, p.snip16(op, 10))
	}

	for _, op := range []Measurement{
		Power, PowerL1, PowerL2, PowerL3,
	} {
		res = append(res, p.snip32(op, 10))
	}

	for _, op := range []Measurement{
		Import, ImportL1, ImportL2, ImportL3,
		Export,
	} {
		res = append(res, p.snip32(op, 10))
	}

	return res
}
