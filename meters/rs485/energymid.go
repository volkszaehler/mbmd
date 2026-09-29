package rs485

import (
	"math"

	"github.com/volkszaehler/mbmd/encoding"
	. "github.com/volkszaehler/mbmd/meters"
)

func init() {
	Register("ENERGYMID", NewEnergyMIDProducer)
}

const (
	energyMIDVoltageExponent = 12
	energyMIDCurrentExponent = 108
	energyMIDPowerExponent   = 212
	energyMIDEnergyFactor    = 308
)

type EnergyMIDProducer struct {
	Opcodes
}

func NewEnergyMIDProducer() Producer {
	/**
	 * Opcodes for Gossen Metrawatt ENERGYMID EM2281/EM2289/EM2381/EM2387/EM2389.
	 * Input registers, zero based.
	 * http://datenblatt.stark-elektronik.de/em2281-em2389-modbus-rtu.pdf
	 */
	ops := Opcodes{
		VoltageL1_L2: 0,
		VoltageL2_L3: 1,
		VoltageL3_L1: 2,
		VoltageL1:    4,
		VoltageL2:    5,
		VoltageL3:    6,
		Voltage:      7,

		Frequency: 11,

		CurrentL1: 100,
		CurrentL2: 101,
		CurrentL3: 102,
		Current:   103,

		PowerL1: 200,
		PowerL2: 201,
		PowerL3: 202,
		Power:   203,

		ReactivePowerL1: 204,
		ReactivePowerL2: 205,
		ReactivePowerL3: 206,
		ReactivePower:   207,

		CosphiL1: 208,
		CosphiL2: 209,
		CosphiL3: 210,
		Cosphi:   211,

		Import:         300,
		Export:         302,
		ReactiveImport: 304,
		ReactiveExport: 306,
	}
	return &EnergyMIDProducer{Opcodes: ops}
}

// Description implements Producer interface
func (p *EnergyMIDProducer) Description() string {
	return "Gossen Metrawatt ENERGYMID EM228x/EM238x"
}

// energyMIDMantissaTransform decodes a SINT16 mantissa with the SINT8 exponent in the low byte of the last register
func energyMIDMantissaTransform(b []byte) float64 {
	mantissa := encoding.Int16(b)
	if mantissa == math.MinInt16 {
		return math.NaN()
	}

	exponent := int8(b[len(b)-1])
	return float64(mantissa) * math.Pow10(int(exponent))
}

// energyMIDEnergyTransform decodes a UINT32 mantissa in Wh with the UINT32 primary energy factor in the last two registers
func energyMIDEnergyTransform(b []byte) float64 {
	mantissa := encoding.Uint32(b)
	factor := encoding.Uint32(b[len(b)-4:])
	return float64(mantissa) * float64(factor) / 1e3
}

func (p *EnergyMIDProducer) snip(iec Measurement, readlen uint16, transform RTUTransform) Operation {
	return Operation{
		FuncCode:  ReadInputReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   readlen,
		IEC61850:  iec,
		Transform: transform,
	}
}

// snipMantissa reads the mantissa register up to and including its exponent register
func (p *EnergyMIDProducer) snipMantissa(iec Measurement, exponent uint16) Operation {
	return p.snip(iec, exponent-p.Opcode(iec)+1, energyMIDMantissaTransform)
}

// snipEnergy reads the energy mantissa up to and including the primary energy factor
func (p *EnergyMIDProducer) snipEnergy(iec Measurement) Operation {
	return p.snip(iec, energyMIDEnergyFactor-p.Opcode(iec)+2, energyMIDEnergyTransform)
}

// Probe implements Producer interface
func (p *EnergyMIDProducer) Probe() Operation {
	return p.snipMantissa(VoltageL1, energyMIDVoltageExponent)
}

// Produce implements Producer interface
func (p *EnergyMIDProducer) Produce() (res []Operation) {
	for op := range p.Opcodes {
		switch op {
		case VoltageL1_L2, VoltageL2_L3, VoltageL3_L1, VoltageL1, VoltageL2, VoltageL3, Voltage:
			res = append(res, p.snipMantissa(op, energyMIDVoltageExponent))

		case CurrentL1, CurrentL2, CurrentL3, Current:
			res = append(res, p.snipMantissa(op, energyMIDCurrentExponent))

		case PowerL1, PowerL2, PowerL3, Power, ReactivePowerL1, ReactivePowerL2, ReactivePowerL3, ReactivePower:
			res = append(res, p.snipMantissa(op, energyMIDPowerExponent))

		case Frequency:
			res = append(res, p.snip(op, 1, MakeScaledTransform(RTUUint16ToFloat64, 100)))

		case CosphiL1, CosphiL2, CosphiL3, Cosphi:
			res = append(res, p.snip(op, 1, MakeScaledTransform(RTUInt16ToFloat64, 1000)))

		default:
			res = append(res, p.snipEnergy(op))
		}
	}

	return res
}
