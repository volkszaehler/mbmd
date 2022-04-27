package rs485

import (
	. "github.com/volkszaehler/mbmd/meters"
)

func init() {
	Register("CarloGavazzi", NewCarloGavazziProducer)
}

type CarloGavazziProducer struct {
	Opcodes
}

func NewCarloGavazziProducer() Producer {
	/***
	 * https://www.aggsoft.com/serial-data-logger/tutorials/modbus-data-logging/carlo-gavazzi-em24.htm
	 */
	ops := Opcodes{
		VoltageL1: 00,
		VoltageL2: 02,
		VoltageL3: 04,

		CurrentL1: 12,
		CurrentL2: 14,
		CurrentL3: 16,

		PowerL1: 18,
		PowerL2: 20,
		PowerL3: 22,
		Power:   40,

		ImportL1: 70,
		ImportL2: 72,
		ImportL3: 74,
		Import:   62,

		Cosphi:   53,
		CosphiL1: 50,
		CosphiL2: 51,
		CosphiL3: 52,

		Frequency: 55,
	}
	return &CarloGavazziProducer{Opcodes: ops}
}

// Description implements Producer interface
func (p *CarloGavazziProducer) Description() string {
	return "Carlo Gavazzi EM24"
}

func (p *CarloGavazziProducer) snip16(iec Measurement, scaler ...float64) Operation {
	transform := RTUInt16ToFloat64 // default conversion
	if len(scaler) > 0 {
		transform = MakeScaledTransform(transform, scaler[0])
	}

	operation := Operation{
		FuncCode:  ReadInputReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   1,
		IEC61850:  iec,
		Transform: transform,
	}
	return operation
}

func (p *CarloGavazziProducer) snip32(iec Measurement, scaler ...float64) Operation {
	transform := RTUInt32ToFloat64 // default conversion
	if len(scaler) > 0 {
		transform = MakeScaledTransform(transform, scaler[0])
	}

	operation := Operation{
		FuncCode:  ReadInputReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   2,
		IEC61850:  iec,
		Transform: transform,
	}
	return operation
}

// Probe implements Producer interface
func (p *CarloGavazziProducer) Probe() Operation {
	return p.snip32(VoltageL1, 10)
}

// Produce implements Producer interface
func (p *CarloGavazziProducer) Produce() (res []Operation) {
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
		res = append(res, p.snip16(op, 100))
	}

	for _, op := range []Measurement{
		Frequency,
	} {
		res = append(res, p.snip16(op, 10))
	}

	for _, op := range []Measurement{
		Power, PowerL1, PowerL2, PowerL3,
	} {
		res = append(res, p.snip32(op, 100))
	}

	for _, op := range []Measurement{
		Import, ImportL1, ImportL2, ImportL3,
	} {
		res = append(res, p.snip32(op, 10))
	}

	return res
}
