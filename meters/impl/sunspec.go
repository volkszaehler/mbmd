package impl

import (
	"encoding/binary"
	"errors"
	"math"
	"strings"

	. "github.com/gonium/gosdm630/meters"
)

const (
	// MODBUS protocol address (base 0)
	sunspecBase         = 40000
	sunspecID           = 1
	sunspecModelID      = 3
	sunspecManufacturer = 5
	sunspecModel        = 21
	sunspecOptions      = 37
	sunspecVersion      = 45
	sunspecSerial       = 53

	sunspecSignature = 0x53756e53 // SunS
)

// SunspecRTUUint16ToFloat64 converts 16 bit unsigned integer readings
// If byte sequence is 0xffff, NaN is returned for compatibility with SunSpec inverters
func SunspecRTUUint16ToFloat64(b []byte) float64 {
	u := binary.BigEndian.Uint16(b)
	if u == 0xffff {
		return math.NaN()
	}
	return float64(u)
}

// SunspecRTUUint16ToFloat64 converts 16 bit unsigned integer readings
// If byte sequence is 0xffff, NaN is returned for compatibility with SunSpec inverters
func SunspecRTUInt16ToFloat64(b []byte) float64 {
	u := binary.BigEndian.Uint16(b)
	if u == 0x8000 {
		return math.NaN()
	}
	return float64(int16(u))
}

// SunspecRTUUint32ToFloat64 converts 32 bit unsigned integer readings
// If byte sequence is 0xffffffff, NaN is returned for compatibility with SunSpec inverters
func SunspecRTUUint32ToFloat64(b []byte) float64 {
	u := binary.BigEndian.Uint32(b)
	if u == 0xffffffff {
		return math.NaN()
	}
	return float64(u)
}

type SunSpecCore struct {
	Opcodes
}

func (p *SunSpecCore) ConnectionType() ConnectionType {
	return TCP
}

func (p *SunSpecCore) GetSunSpecCommonBlock() Operation {
	// must return 0x53756e53 = SunS
	return Operation{
		FuncCode: ReadHoldingReg,
		OpCode:   sunspecBase,
		ReadLen:  sunspecSerial + 16,
	}
}

func (p *SunSpecCore) DecodeSunSpecCommonBlock(b []byte) (SunSpecDeviceDescriptor, error) {
	res := SunSpecDeviceDescriptor{}

	if len(b) < sunspecSerial+2*16 {
		return res, errors.New("Could not read SunSpec device descriptor")
	}

	u := binary.BigEndian.Uint32(b[sunspecID-1:])
	if u != sunspecSignature {
		return res, errors.New("Invalid SunSpec device signature")
	}

	res.Manufacturer = p.stringDecode(b, sunspecManufacturer, 16)
	res.Model = p.stringDecode(b, sunspecModel, 16)
	res.Options = p.stringDecode(b, sunspecOptions, 8)
	res.Version = p.stringDecode(b, sunspecVersion, 8)
	res.Serial = p.stringDecode(b, sunspecSerial, 16)

	return res, nil
}

func (p *SunSpecCore) stringDecode(b []byte, reg int, len int) string {
	start := 2 * (reg - 1)
	end := 2 * (reg + len - 1)
	// trim space and null
	return strings.TrimRight(string(b[start:end-1]), " \x00")
}

func (p *SunSpecCore) snip(iec Measurement, readlen uint16) Operation {
	return Operation{
		FuncCode: ReadHoldingReg,
		OpCode:   sunspecBase + p.Opcode(iec) - 1, // adjust according to docs
		ReadLen:  readlen,
		IEC61850: iec,
	}
}

func (p *SunSpecCore) snip16uint(iec Measurement, scaler ...float64) Operation {
	snip := p.snip(iec, 1)

	snip.Transform = SunspecRTUUint16ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
	}

	return snip
}

func (p *SunSpecCore) snip16int(iec Measurement, scaler ...float64) Operation {
	snip := p.snip(iec, 1)

	snip.Transform = SunspecRTUInt16ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
	}

	return snip
}

func (p *SunSpecCore) snip32(iec Measurement, scaler ...float64) Operation {
	snip := p.snip(iec, 2)

	snip.Transform = SunspecRTUUint32ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
	}

	return snip
}

func (p *SunSpecCore) minMax(iec ...Measurement) (uint16, uint16) {
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
func (p *SunSpecCore) scaleSnip16(splitter func(...Measurement) Splitter, iecs ...Measurement) Operation {
	min, max := p.minMax(iecs...)

	// read register block
	op := Operation{
		FuncCode: ReadHoldingReg,
		OpCode:   sunspecBase + min - 1, // adjust according to docs
		ReadLen:  max - min + 2,         // registers plus int16 scale factor
		IEC61850: Split,
		Splitter: splitter(iecs...),
	}

	return op
}

func (p *SunSpecCore) scaleSnip32(splitter func(...Measurement) Splitter, iecs ...Measurement) Operation {
	op := p.scaleSnip16(splitter, iecs...)
	op.ReadLen = (op.ReadLen-1)*2 + 1 // read 4 bytes instead of 2 plus trailing scale factor
	return op
}

func (p *SunSpecCore) mkSplitInt16(iecs ...Measurement) Splitter {
	return p.mkBlockSplitter(2, SunspecRTUInt16ToFloat64, iecs...)
}

func (p *SunSpecCore) mkSplitUint16(iecs ...Measurement) Splitter {
	return p.mkBlockSplitter(2, SunspecRTUUint16ToFloat64, iecs...)
}

func (p *SunSpecCore) mkSplitUint32(iecs ...Measurement) Splitter {
	// use div 1000 for kWh conversion
	return p.mkBlockSplitter(4, MakeScaledTransform(SunspecRTUUint32ToFloat64, 1000), iecs...)
}

func (p *SunSpecCore) mkBlockSplitter(dataSize uint16, valFunc func([]byte) float64, iecs ...Measurement) Splitter {
	min, _ := p.minMax(iecs...)
	return func(b []byte) []SplitResult {
		// get scaler from last entry in result block
		exp := int(int16(binary.BigEndian.Uint16(b[len(b)-2:]))) // last int16
		scaler := math.Pow10(exp)

		res := make([]SplitResult, 0, len(iecs))

		// split result block into individual readings
		for _, iec := range iecs {
			opcode := p.Opcode(iec)
			val := valFunc(b[dataSize*(opcode-min):]) // 2 bytes per uint16, 4 bytes per uint32

			// filter results of Sunspec transforms
			if math.IsNaN(val) {
				continue
			}

			op := SplitResult{
				OpCode:   sunspecBase + opcode - 1,
				IEC61850: iec,
				Value:    scaler * val,
			}

			res = append(res, op)
		}

		return res
	}
}
