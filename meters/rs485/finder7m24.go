package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register("FIND7M24", NewFinder7M24Producer)
}

type Finder7M24Producer struct {
	Opcodes
}

func NewFinder7M24Producer() Producer {
	/***
	 * Opcodes for Finder 7M.24 (1-phase energy meter)
	 * Based on Modbus protocol documentation
	 * https://cdn.findernet.com/app/uploads/2021/09/20090052/Modbus-7M24-7M38_v2_30062021.pdf
	 *
	 * 7M.24 is a 1-phase energy meter with MID certification
	 * Uses Input Registers (Function Code 0x04)
	 * All values are IEEE 754 Float (32-bit, 2 registers)
	 */
	ops := Opcodes{
		Voltage:        2500,
		Current:        2516,
		Power:          2536,
		ReactivePower:  2544,
		ApparentPower:  2552,
		Cosphi:         2560,
		PhaseAngle:     2576,
		Frequency:      2584,
		THD:            2594,
		Import:         2752,
		ReactiveImport: 2754,
		Export:         2756,
		ReactiveExport: 2758,
	}

	return &Finder7M24Producer{Opcodes: ops}
}

// Description implements Producer interface
func (p *Finder7M24Producer) Description() string {
	return "Finder 7M.24"
}

// snip creates modbus operation
func (p *Finder7M24Producer) snip(iec Measurement, readlen uint16, transform RTUTransform) Operation {
	return Operation{
		FuncCode:  ReadInputReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   readlen,
		IEC61850:  iec,
		Transform: transform,
	}
}

// snip32 creates modbus operation for 32-bit register (2 registers)
func (p *Finder7M24Producer) snip32(iec Measurement, scaler ...float64) Operation {
	transform := RTUIeee754ToFloat64
	if len(scaler) > 0 {
		transform = MakeScaledTransform(RTUUint32ToFloat64, scaler[0])
	}
	return p.snip(iec, 2, transform)
}

func (p *Finder7M24Producer) Probe() Operation {
	return p.snip32(Voltage)
}

// Produce implements Producer interface
func (p *Finder7M24Producer) Produce() (res []Operation) {
	// Instantaneous values (IEEE754 float32, no scaling)
	for _, op := range []Measurement{
		Voltage, Current,
		Power, ReactivePower, ApparentPower,
		Cosphi, Frequency, PhaseAngle,
		THD,
	} {
		res = append(res, p.snip32(op))
	}

	// Energy counters (IEEE754 float32, in Wh - scale to kWh)
	for _, op := range []Measurement{
		Import, Export,
		ReactiveImport, ReactiveExport,
	} {
		res = append(res, p.snip32(op, 1000))
	}

	return res
}
