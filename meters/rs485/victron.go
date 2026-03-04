package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register("VM3P", NewVM3PProducer)
}

// VM3PProducer implements the Producer interface for Victron Energy Meters
// VM-3P75CT (product ID 0xa1b1) and VM-3P5A (product ID 0xa1b2).
// Both models share an identical Modbus register map.
// Register map derived from:
// https://github.com/victronenergy/dbus-modbus-client/blob/master/victron_em.py
//
// The meters communicate via Modbus/UDP using Holding Registers (FC=3).
// All multi-word values are big-endian (most-significant word first).
type VM3PProducer struct {
	Opcodes
}

// NewVM3PProducer creates a new VM3PProducer.
func NewVM3PProducer() Producer {
	/**
	 * Opcodes for Victron Energy Meter VM-3P75CT and VM-3P5A.
	 * Register map: https://github.com/victronenergy/dbus-modbus-client/blob/master/victron_em.py
	 *
	 * Uses Holding Registers (Function Code 0x03), big-endian word order.
	 * Per-phase base address: 0x3040 + 8*(n-1)
	 * Per-phase power address: 0x3082 + 4*(n-1)
	 */
	ops := Opcodes{
		Frequency: 0x3032, // u16 /100 → Hz
		Import:    0x3034, // u32 /100 → kWh
		Export:    0x3036, // u32 /100 → kWh
		Power:     0x3080, // s32 /1   → W
		Cosphi:    0x303a, // s16 /1000

		VoltageL1: 0x3040, // s16 /100 → V
		CurrentL1: 0x3041, // s16 /100 → A
		ImportL1:  0x3042, // u32 /100 → kWh
		ExportL1:  0x3044, // u32 /100 → kWh
		CosphiL1:  0x3047, // s16 /1000
		PowerL1:   0x3082, // s32 /1   → W

		VoltageL2: 0x3048, // s16 /100 → V
		CurrentL2: 0x3049, // s16 /100 → A
		ImportL2:  0x304a, // u32 /100 → kWh
		ExportL2:  0x304c, // u32 /100 → kWh
		CosphiL2:  0x304f, // s16 /1000
		PowerL2:   0x3086, // s32 /1   → W

		VoltageL3: 0x3050, // s16 /100 → V
		CurrentL3: 0x3051, // s16 /100 → A
		ImportL3:  0x3052, // u32 /100 → kWh
		ExportL3:  0x3054, // u32 /100 → kWh
		CosphiL3:  0x3057, // s16 /1000
		PowerL3:   0x308a, // s32 /1   → W
	}
	return &VM3PProducer{Opcodes: ops}
}

// Description implements Producer interface.
func (p *VM3PProducer) Description() string {
	return "Victron VM-3P75CT/VM-3P5A"
}

// snip creates a Holding Register operation for the given measurement.
func (p *VM3PProducer) snip(iec Measurement, readlen uint16, transform RTUTransform) Operation {
	return Operation{
		FuncCode:  ReadHoldingReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   readlen,
		IEC61850:  iec,
		Transform: transform,
	}
}

// snip16s creates an operation for a signed 16-bit register with divisor.
func (p *VM3PProducer) snip16s(iec Measurement, divisor float64) Operation {
	return p.snip(iec, 1, MakeScaledTransform(RTUInt16ToFloat64, divisor))
}

// snip16u creates an operation for an unsigned 16-bit register with divisor.
func (p *VM3PProducer) snip16u(iec Measurement, divisor float64) Operation {
	return p.snip(iec, 1, MakeScaledTransform(RTUUint16ToFloat64, divisor))
}

// snip32s creates an operation for a signed 32-bit register.
func (p *VM3PProducer) snip32s(iec Measurement) Operation {
	return p.snip(iec, 2, RTUInt32ToFloat64)
}

// snip32u creates an operation for an unsigned 32-bit register with divisor.
func (p *VM3PProducer) snip32u(iec Measurement, divisor float64) Operation {
	return p.snip(iec, 2, MakeScaledTransform(RTUUint32ToFloat64, divisor))
}

// Probe implements Producer interface.
func (p *VM3PProducer) Probe() Operation {
	return p.snip16s(VoltageL1, 100)
}

// Produce implements Producer interface.
func (p *VM3PProducer) Produce() (res []Operation) {
	// Voltage: s16 /100 → V
	for _, op := range []Measurement{VoltageL1, VoltageL2, VoltageL3} {
		res = append(res, p.snip16s(op, 100))
	}

	// Current: s16 /100 → A
	for _, op := range []Measurement{CurrentL1, CurrentL2, CurrentL3} {
		res = append(res, p.snip16s(op, 100))
	}

	// Active power: s32 /1 → W
	for _, op := range []Measurement{Power, PowerL1, PowerL2, PowerL3} {
		res = append(res, p.snip32s(op))
	}

	// Power factor: s16 /1000
	for _, op := range []Measurement{Cosphi, CosphiL1, CosphiL2, CosphiL3} {
		res = append(res, p.snip16s(op, 1000))
	}

	// Frequency: u16 /100 → Hz
	res = append(res, p.snip16u(Frequency, 100))

	// Energy counters: u32 /100 → kWh
	for _, op := range []Measurement{Import, ImportL1, ImportL2, ImportL3, Export, ExportL1, ExportL2, ExportL3} {
		res = append(res, p.snip32u(op, 100))
	}

	return res
}
