package rs485

import (
	"math"

	. "github.com/volkszaehler/mbmd/meters"
)

func init() {
	Register("ABB", NewABBProducer)
}

type signedness int

const (
	unsigned signedness = iota
	signed
)

type ABBProducer struct {
	Opcodes
}

func NewABBProducer() Producer {
	/***
	 * http://datenblatt.stark-elektronik.de/Energiezaehler_B-Serie_Handbuch.pdf
	 */
	ops := Opcodes{
		VoltageL1: 0x5B00,
		VoltageL2: 0x5B02,
		VoltageL3: 0x5B04,

		CurrentL1: 0x5B0C,
		CurrentL2: 0x5B0E,
		CurrentL3: 0x5B10,

		Power:   0x5B14,
		PowerL1: 0x5B16,
		PowerL2: 0x5B18,
		PowerL3: 0x5B1A,

		ImportL1: 0x5460,
		ImportL2: 0x5464,
		ImportL3: 0x5468,
		Import:   0x5000,

		ExportL1: 0x546C,
		ExportL2: 0x5470,
		ExportL3: 0x5474,
		Export:   0x5004,

		Cosphi:   0x5B3A,
		CosphiL1: 0x5B3B,
		CosphiL2: 0x5B3C,
		CosphiL3: 0x5B3D,

		Frequency: 0x5B2C,
	}
	return &ABBProducer{Opcodes: ops}
}

// Description implements Producer interface
func (p *ABBProducer) Description() string {
	return "ABB A/B-Series meters"
}

// wrapTransform validates if reading result is undefined and returns NaN in that case
func wrapTransform(byteCount uint16, sign signedness, transform RTUTransform) RTUTransform {
	var nan []byte
	if sign == signed {
		nan = []byte{0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	} else {
		nan = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	}

	return func(b []byte) float64 {
		var i uint16
		for i = 0; i < byteCount; i++ {
			if b[i] != nan[i] {
				return transform(b)
			}
		}
		return math.NaN()
	}
}

func (p *ABBProducer) snip(iec Measurement, readlen uint16, sign signedness, transform RTUTransform, scaler ...float64) Operation {
	// wrap the transformation inside a NaN check
	nanAwareTransform := wrapTransform(2*readlen, sign, transform)

	snip := Operation{
		FuncCode:  ReadHoldingReg,
		OpCode:    p.Opcodes[iec],
		ReadLen:   readlen,
		Transform: nanAwareTransform,
		IEC61850:  iec,
	}

	if len(scaler) > 0 {
		snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
	}

	return snip
}

// snip16u creates modbus operation for single register
func (p *ABBProducer) snip16u(iec Measurement, scaler ...float64) Operation {
	return p.snip(iec, 1, unsigned, RTUUint16ToFloat64, scaler...)
}

// snip16i creates modbus operation for single register
func (p *ABBProducer) snip16i(iec Measurement, scaler ...float64) Operation {
	return p.snip(iec, 1, signed, RTUInt16ToFloat64, scaler...)
}

// snip32u creates modbus operation for double register
func (p *ABBProducer) snip32u(iec Measurement, scaler ...float64) Operation {
	return p.snip(iec, 2, unsigned, RTUUint32ToFloat64, scaler...)
}

// snip32i creates modbus operation for double register
func (p *ABBProducer) snip32i(iec Measurement, scaler ...float64) Operation {
	return p.snip(iec, 2, signed, RTUInt32ToFloat64, scaler...)
}

// snip64u creates modbus operation for double register
func (p *ABBProducer) snip64u(iec Measurement, scaler ...float64) Operation {
	return p.snip(iec, 4, unsigned, RTUUint64ToFloat64, scaler...)
}

// Probe implements Producer interface
func (p *ABBProducer) Probe() Operation {
	return p.snip32u(VoltageL1, 10)
}

// Produce implements Producer interface
func (p *ABBProducer) Produce() (res []Operation) {
	for _, op := range []Measurement{
		VoltageL1, VoltageL2, VoltageL3,
	} {
		res = append(res, p.snip32u(op, 10))
	}

	for _, op := range []Measurement{
		CurrentL1, CurrentL2, CurrentL3,
	} {
		res = append(res, p.snip32u(op, 100))
	}

	for _, op := range []Measurement{
		Cosphi, CosphiL1, CosphiL2, CosphiL3,
	} {
		res = append(res, p.snip16i(op, 1000))
	}

	for _, op := range []Measurement{
		Frequency,
	} {
		res = append(res, p.snip16u(op, 100))
	}

	for _, op := range []Measurement{
		Power, PowerL1, PowerL2, PowerL3,
	} {
		res = append(res, p.snip32i(op, 100))
	}

	for _, op := range []Measurement{
		Import, ImportL1, ImportL2, ImportL3,
		Export, ExportL1, ExportL2, ExportL3,
	} {
		res = append(res, p.snip64u(op, 100))
	}

	return res
}
