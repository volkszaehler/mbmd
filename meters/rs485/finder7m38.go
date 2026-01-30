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
		VoltageL1:       2500, // Phase 1 line to neutral volts
		VoltageL2:       2502, // Phase 2 line to neutral volts
		VoltageL3:       2504, // Phase 3 line to neutral volts
		CurrentL1:       2516, // Phase 1 current
		CurrentL2:       2518, // Phase 2 current
		CurrentL3:       2520, // Phase 3 current
		PowerL1:         2530, // Phase 1 active power
		PowerL2:         2532, // Phase 2 active power
		PowerL3:         2534, // Phase 3 active power
		Power:           2536, // Total system power
		ReactivePowerL1: 2538, // Phase 1 reactive power
		ReactivePowerL2: 2540, // Phase 2 reactive power
		ReactivePowerL3: 2542, // Phase 3 reactive power
		ReactivePower:   2544, // Total system VAr
		ApparentPowerL1: 2546, // Phase 1 apparent power
		ApparentPowerL2: 2548, // Phase 2 apparent power
		ApparentPowerL3: 2550, // Phase 3 apparent power
		ApparentPower:   2552, // Total system volt amps.
		CosphiL1:        2554, // Phase 1 power factor
		CosphiL2:        2556, // Phase 2 power factor
		CosphiL3:        2558, // Phase 3 power factor
		Cosphi:          2560, // Total system power factor
		PhaseAngle:      2576, // Total system phase angle
		Frequency:       2584, // Frequency of supply voltages
		THDL1:           2594, // Phase 1 L/N volts THD
		THDL2:           2596, // Phase 2 L/N volts THD
		THDL3:           2598, // Phase 3 L/N volts THD
		Import:          2752, // Total Import kWh
		ReactiveImport:  2754, // Total Import kVArh
		Export:          2756, // Total Export kWh
		ReactiveExport:  2758, // Total Export kVArh
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
		transform = MakeScaledTransform(RTUUint32ToFloat64, scaler[0])
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
