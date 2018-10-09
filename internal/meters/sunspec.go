package meters

import (
	"encoding/binary"
	"math"
)

const (
	METERTYPE_SUN = "SUN"

	// MODBUS protocol address (base 0)
	base = 40000
)

type SUNProducer struct {
	MeasurementMapping
}

func NewSUNProducer() *SUNProducer {
	/***
	 * Opcodes for SunSpec-compatible Inverters like SolarEdge
	 * https://www.solaredge.com/sites/default/files/sunspec-implementation-technical-note.pdf
	 */
	ops := Measurements{
		Current:   72,
		CurrentL1: 73,
		CurrentL2: 74,
		CurrentL3: 75, // + scaler

		VoltageL1: 80,
		VoltageL2: 81,
		VoltageL3: 82, // + scaler

		Power: 84, // + scaler
		// ApparentPower: 88, // + scaler
		// ReactivePower: 90, // + scaler
		Export: 94, // + scaler

		Cosphi:    92, // + scaler
		Frequency: 86, // + scaler

		DCCurrent: 97,  // + scaler
		DCVoltage: 99,  // + scaler
		DCPower:   101, // + scaler

		HeatSinkTemp: 104, // + scaler
	}
	return &SUNProducer{
		MeasurementMapping{ops},
	}
}

func (p *SUNProducer) GetMeterType() string {
	return METERTYPE_SUN
}

func (p *SUNProducer) snip(iec Measurement, readlen uint16) Operation {
	return Operation{
		FuncCode: ReadHoldingReg,
		OpCode:   base + p.Opcode(iec) - 1, // adjust according to docs
		ReadLen:  readlen,
		IEC61850: iec,
	}
}

func (p *SUNProducer) snip16uint(iec Measurement, scaler ...float64) Operation {
	snip := p.snip(iec, 1)

	snip.Transform = RTUUint16ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeRTUScaledUint16ToFloat64(scaler[0])
	}

	return snip
}

func (p *SUNProducer) snip16int(iec Measurement, scaler ...float64) Operation {
	snip := p.snip(iec, 1)

	snip.Transform = RTUInt16ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeRTUScaledInt16ToFloat64(scaler[0])
	}

	return snip
}

func (p *SUNProducer) snip32(iec Measurement, scaler ...float64) Operation {
	snip := p.snip(iec, 2)

	snip.Transform = RTUUint32ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeRTUScaledUint32ToFloat64(scaler[0])
	}

	return snip
}

func (p *SUNProducer) minMax(iec ...Measurement) (uint16, uint16) {
	var min = uint16(0xFFFF)
	var max = uint16(0x0000)
	for _, i := range iec {
		op := p.Opcode(i)
		if op < min {
			min = op
		}
		if op > max {
			max = op
		}
	}
	return min, max
}

// create a block reading function the result of which is then split into measurements
func (p *SUNProducer) scaleSnip16(splitter func(...Measurement) Splitter, iecs ...Measurement) Operation {
	min, max := p.minMax(iecs...)

	// read register block
	op := Operation{
		FuncCode: ReadHoldingReg,
		OpCode:   base + min - 1, // adjust according to docs
		ReadLen:  max - min + 2,  // registers plus int16 scale factor
		IEC61850: Split,
		Splitter: splitter(iecs...),
	}

	return op
}

func (p *SUNProducer) scaleSnip32(splitter func(...Measurement) Splitter, iecs ...Measurement) Operation {
	op := p.scaleSnip16(splitter, iecs...)
	op.ReadLen = (op.ReadLen-1)*2 + 1 // read 4 bytes instead of 2 plus trailing scale factor
	return op
}

func (p *SUNProducer) mkSplitInt16(iecs ...Measurement) Splitter {
	return p.mkBlockSplitter(2, RTUInt16ToFloat64, iecs...)
}

func (p *SUNProducer) mkSplitUint16(iecs ...Measurement) Splitter {
	return p.mkBlockSplitter(2, RTUUint16ToFloat64, iecs...)
}

func (p *SUNProducer) mkSplitUint32(iecs ...Measurement) Splitter {
	return p.mkBlockSplitter(4, RTUUint32ToFloat64, iecs...)
}

func (p *SUNProducer) mkBlockSplitter(dataSize uint16, valFunc func([]byte) float64, iecs ...Measurement) Splitter {
	min, _ := p.minMax(iecs...)
	return func(b []byte) []SplitResult {
		// get scaler from last entry in result block
		exp := int(int16(binary.BigEndian.Uint16(b[len(b)-2:]))) // last int16
		scaler := math.Pow10(exp)

		res := make([]SplitResult, len(iecs))

		// split result block into individual readings
		for idx, iec := range iecs {
			opcode := p.Opcode(iec)
			val := valFunc(b[dataSize*(opcode-min):]) // 2 bytes per uint16, 4 bytes per uint32

			op := SplitResult{
				OpCode:   base + opcode - 1,
				IEC61850: iec,
				Value:    scaler * val,
			}

			res[idx] = op
		}

		return res
	}
}

func (p *SUNProducer) Probe() Operation {
	return p.snip16uint(VoltageL1, 10)
}

func (p *SUNProducer) Produce() (res []Operation) {
	res = []Operation{
		// uint16
		p.scaleSnip16(p.mkSplitUint16, VoltageL1, VoltageL2, VoltageL3),
		p.scaleSnip16(p.mkSplitUint16, Current, CurrentL1, CurrentL2, CurrentL3),

		p.scaleSnip16(p.mkSplitUint16, Frequency),
		p.scaleSnip16(p.mkSplitUint16, DCCurrent),
		p.scaleSnip16(p.mkSplitUint16, DCVoltage),

		// int16
		p.scaleSnip16(p.mkSplitInt16, Cosphi),
		p.scaleSnip16(p.mkSplitInt16, Power),
		p.scaleSnip16(p.mkSplitInt16, DCPower),
		p.scaleSnip16(p.mkSplitInt16, HeatSinkTemp),

		// uint32
		p.scaleSnip32(p.mkSplitUint32, Export),
	}

	return res
}
