package rs485

import (
	. "github.com/volkszaehler/mbmd/meters"
)

func init() {
	Register("DTSU666", NewDTSU666Producer)
}

type DTSU666Producer struct {
	Opcodes
}

func NewDTSU666Producer() Producer {
	// docs: https://www.chintglobal.com/content/dam/chint/global/product-center/instruments-meters/electricity-meter/din-rail-meter/dtsu666/manual/DTSU666%20DSSU666%20User%20Manual.pdf
	ops := Opcodes{
		Import:          0x101E,
		Export:          0x1028,
		VoltageL1:       0x2000,
		VoltageL2:       0x2002,
		VoltageL3:       0x2004,
		CurrentL1:       0x200C,
		CurrentL2:       0x200E,
		CurrentL3:       0x2010,
		Power:           0x2012,
		PowerL1:         0x2014,
		PowerL2:         0x2016,
		PowerL3:         0x2018,
		ReactivePower:   0x201A,
		ReactivePowerL1: 0x201C,
		ReactivePowerL2: 0x201E,
		ReactivePowerL3: 0x2020,
		Cosphi:          0x202A, // power factor, not cosine phi
		CosphiL1:        0x202C, // power factor, not cosine phi
		CosphiL2:        0x202E, // power factor, not cosine phi
		CosphiL3:        0x2030, // power factor, not cosine phi
		Frequency:       0x2044,
	}
	return &DTSU666Producer{Opcodes: ops}
}

func (p *DTSU666Producer) Description() string {
	return "Chint DTSU666"
}

func (p *DTSU666Producer) snip(iec Measurement, scaler ...float64) Operation {
	operation := Operation{
		FuncCode:  ReadHoldingReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   2,
		IEC61850:  iec,
		Transform: RTUIeee754ToFloat64,
	}

	if len(scaler) > 0 {
		operation.Transform = MakeScaledTransform(operation.Transform, scaler[0])
	}

	return operation
}

func (p *DTSU666Producer) Probe() Operation {
	return p.snip(VoltageL1, 10)
}

func (p *DTSU666Producer) Produce() (res []Operation) {
	for _, op := range []Measurement{
		VoltageL1, VoltageL2, VoltageL3,
		Power, ReactivePower, PowerL1, PowerL2, PowerL3,
	} {
		res = append(res, p.snip(op, 10))
	}

	for _, op := range []Measurement{
		CurrentL1, CurrentL2, CurrentL3,
		Cosphi, CosphiL1, CosphiL2, CosphiL3,
	} {
		res = append(res, p.snip(op, 1000))
	}

	for _, op := range []Measurement{
		Frequency,
	} {
		res = append(res, p.snip(op, 100))
	}

	for _, op := range []Measurement{
		Import, Export,
	} {
		res = append(res, p.snip(op, 1))
	}

	return res
}
