package rs485

import (
	. "github.com/volkszaehler/mbmd/meters"
)

func init() {
	Register("MPM", NewMPM3MPProducer)
}

type MPM3MPProducer struct {
	Opcodes
}

func NewMPM3MPProducer() Producer {
	/***
	 * http://www.qrck.info/MPM3PModBus.pdf
	 */
	ops := Opcodes{
		Sum:             0x00,
		Import:          0x02,
		Export:          0x04,
		ReactiveSum:     0x06,
		VoltageL1:       0x08,
		VoltageL2:       0x0A,
		VoltageL3:       0x0C,
		CurrentL1:       0x0E,
		CurrentL2:       0x10,
		CurrentL3:       0x12,
		PowerL1:         0x14,
		PowerL2:         0x16,
		PowerL3:         0x18,
		ReactivePowerL1: 0x1A,
		ReactivePowerL2: 0x1C,
		ReactivePowerL3: 0x1E,
		Cosphi:          0x2A,
		CosphiL1:        0x20,
		CosphiL2:        0x22,
		CosphiL3:        0x24,
		Power:           0x26,
		ReactivePower:   0x28,
		Frequency:       0x2C,
	}
	return &MPM3MPProducer{Opcodes: ops}
}

// Description implements Producer interface
func (p *MPM3MPProducer) Description() string {
	return "Bernecker Engineering MPM3PM meters"
}

func (p *MPM3MPProducer) snip(iec Measurement, readlen uint16, transform RTUTransform, scaler ...float64) Operation {
	snip := Operation{
		FuncCode:  ReadHoldingReg,
		OpCode:    p.Opcodes[iec],
		ReadLen:   readlen,
		Transform: transform,
		IEC61850:  iec,
	}

	if len(scaler) > 0 {
		snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
	}

	return snip
}

// snip32u creates modbus operation for double register
func (p *MPM3MPProducer) snip32u(iec Measurement, scaler ...float64) Operation {
	return p.snip(iec, 2, RTUUint32ToFloat64, scaler...)
}

// snip32i creates modbus operation for double register
func (p *MPM3MPProducer) snip32i(iec Measurement, scaler ...float64) Operation {
	return p.snip(iec, 2, RTUInt32ToFloat64, scaler...)
}

// Probe implements Producer interface
func (p *MPM3MPProducer) Probe() Operation {
	return p.snip32u(VoltageL1, 10)
}

// Produce implements Producer interface
func (p *MPM3MPProducer) Produce() (res []Operation) {
	for _, op := range []Measurement{
		VoltageL1, VoltageL2, VoltageL3,
	} {
		res = append(res, p.snip32u(op, 10))
	}

	for _, op := range []Measurement{
		CurrentL1, CurrentL2, CurrentL3,
		Import, Export, ReactiveSum, Frequency,
	} {
		res = append(res, p.snip32u(op, 100))
	}

	for _, op := range []Measurement{
		Sum,
		Power, PowerL1, PowerL2, PowerL3,
		ReactivePower, ReactivePowerL1, ReactivePowerL2, ReactivePowerL3,
	} {
		res = append(res, p.snip32i(op, 100))
	}

	for _, op := range []Measurement{
		Cosphi, CosphiL1, CosphiL2, CosphiL3,
	} {
		res = append(res, p.snip32u(op, 1000))
	}

	return res
}
