package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register("FIND7M38", NewFinder7M38Producer)
}

type Finder7M38Producer struct {
	Opcodes
}

func NewFinder7M38Producer() Producer {
	/***
	 * Opcodes for Finder 7M.38 (3-phase energy meter)
	 * Based on Modbus protocol documentation
	 * https://cdn.findernet.com/app/uploads/2021/09/20090052/Modbus-7M24-7M38_v2_30062021.pdf
	 *
	 * 7M.38 is a 3-phase energy meter with MID certification
	 * Uses Input Registers (Function Code 0x04)
	 * All values are IEEE 754 Float (32-bit, 2 registers)
	 */
	ops := Opcodes{
		VoltageL1:       2500,
		VoltageL2:       2502,
		VoltageL3:       2504,
		CurrentL1:       2516,
		CurrentL2:       2518,
		CurrentL3:       2520,
		PowerL1:         2530,
		PowerL2:         2532,
		PowerL3:         2534,
		Power:           2536,
		ReactivePowerL1: 2538,
		ReactivePowerL2: 2540,
		ReactivePowerL3: 2542,
		ReactivePower:   2544,
		ApparentPowerL1: 2546,
		ApparentPowerL2: 2548,
		ApparentPowerL3: 2550,
		ApparentPower:   2552,
		CosphiL1:        2554,
		CosphiL2:        2556,
		CosphiL3:        2558,
		Cosphi:          2560,
		PhaseAngle:      2576,
		Frequency:       2584,
		THDL1:           2594,
		THDL2:           2596,
		THDL3:           2598,
		Import:          2752,
		ReactiveImport:  2754,
		Export:          2756,
		ReactiveExport:  2758,
	}

	return &Finder7M38Producer{Opcodes: ops}
}

// Description implements Producer interface
func (p *Finder7M38Producer) Description() string {
	return "Finder 7M.38"
}

// snip creates modbus operation
func (p *Finder7M38Producer) snip(iec Measurement, readlen uint16, transform RTUTransform) Operation {
	return Operation{
		FuncCode:  ReadInputReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   readlen,
		IEC61850:  iec,
		Transform: transform,
	}
}

// snip32 creates modbus operation for 32-bit register (2 registers)
func (p *Finder7M38Producer) snip32(iec Measurement, scaler ...float64) Operation {
	transform := RTUIeee754ToFloat64
	if len(scaler) > 0 {
		transform = MakeScaledTransform(RTUIeee754ToFloat64, scaler[0])
	}
	return p.snip(iec, 2, transform)
}

func (p *Finder7M38Producer) Probe() Operation {
	return p.snip32(VoltageL1)
}

// Produce implements Producer interface
func (p *Finder7M38Producer) Produce() (res []Operation) {
	// Instantaneous values (IEEE754 float32, no scaling)
	for _, op := range []Measurement{
		VoltageL1, VoltageL2, VoltageL3,
		CurrentL1, CurrentL2, CurrentL3,
		PowerL1, PowerL2, PowerL3, Power,
		ReactivePowerL1, ReactivePowerL2, ReactivePowerL3, ReactivePower,
		ApparentPowerL1, ApparentPowerL2, ApparentPowerL3, ApparentPower,
		CosphiL1, CosphiL2, CosphiL3, Cosphi,
		PhaseAngle, Frequency,
		THDL1, THDL2, THDL3,
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
