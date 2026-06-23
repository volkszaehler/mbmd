package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	// Both meters share the same (publicly documented) register layout, so a
	// single producer is registered under both type codes.
	Register("ELTAKODSZ15", NewEltakoProducer)
	Register("ELTAKODSZ16", NewEltakoProducer)
}

type EltakoProducer struct {
	Opcodes
}

func NewEltakoProducer() Producer {
	/**
	 * Opcodes for Eltako DSZ15DZMOD (and DSZ16 series in integer data format).
	 * https://www.eltako.com/fileadmin/downloads/de/_bedienung/Modbus-RTU_protocol_specification_for_DSZ15DZMOD_V1.6_English_version.pdf
	 *
	 * TODO: the DSZ16 series additionally exposes reactive/apparent power,
	 * frequency, per-phase energy and a float data format. These registers are
	 * not part of the public DSZ15DZMOD spec.
	 */
	ops := Opcodes{
		VoltageL1: 0x0000,
		VoltageL2: 0x0002,
		VoltageL3: 0x0004,

		CurrentL1: 0x0006,
		CurrentL2: 0x0008,
		CurrentL3: 0x000A,

		PowerL1: 0x000C,
		PowerL2: 0x000E,
		PowerL3: 0x0010,
		Power:   0x0034,

		CosphiL1: 0x001E,
		CosphiL2: 0x0020,
		CosphiL3: 0x0022,
		Cosphi:   0x003E,

		Import: 0x0048,
		Export: 0x004A,
	}
	return &EltakoProducer{Opcodes: ops}
}

// Description implements Producer interface
func (p *EltakoProducer) Description() string {
	return "Eltako DSZ15DZMOD / DSZ16"
}

func (p *EltakoProducer) snip(iec Measurement, transform RTUTransform, scaler ...float64) Operation {
	if len(scaler) > 0 {
		transform = MakeScaledTransform(transform, scaler[0])
	}

	return Operation{
		FuncCode:  ReadInputReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   2,
		IEC61850:  iec,
		Transform: transform,
	}
}

// snipU creates a modbus operation for an unsigned 32 bit register
func (p *EltakoProducer) snipU(iec Measurement, scaler ...float64) Operation {
	return p.snip(iec, RTUUint32ToFloat64, scaler...)
}

// snipI creates a modbus operation for a signed 32 bit register
func (p *EltakoProducer) snipI(iec Measurement, scaler ...float64) Operation {
	return p.snip(iec, RTUInt32ToFloat64, scaler...)
}

// Probe implements Producer interface
func (p *EltakoProducer) Probe() Operation {
	return p.snipU(VoltageL1, 100)
}

// Produce implements Producer interface
func (p *EltakoProducer) Produce() (res []Operation) {
	for _, op := range []Measurement{
		VoltageL1, VoltageL2, VoltageL3,
		CurrentL1, CurrentL2, CurrentL3,
	} {
		res = append(res, p.snipU(op, 100))
	}

	for _, op := range []Measurement{
		CosphiL1, CosphiL2, CosphiL3, Cosphi,
	} {
		res = append(res, p.snipI(op, 1000))
	}

	for _, op := range []Measurement{
		PowerL1, PowerL2, PowerL3, Power,
	} {
		res = append(res, p.snipI(op))
	}

	// energy registers report kWh with 2 decimals
	for _, op := range []Measurement{
		Import, Export,
	} {
		res = append(res, p.snipU(op, 100))
	}

	return res
}
