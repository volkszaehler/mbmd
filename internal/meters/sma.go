package meters

import (
	"encoding/binary"
	"math"
)

const (
	METERTYPE_SMA = "SMA"
)

type SMAProducer struct {
	MeasurementMapping
}

func NewSMAProducer() *SMAProducer {
	/***
	 * Opcodes for SMA SunSpec-compatible Inverters
	 * https://www.sma.de/fileadmin/content/landingpages/pl/FAQ/SunSpec_Modbus-TI-en-15.pdf
	 */
	ops := Measurements{
		Current:   188, // uint16
		CurrentL1: 189,
		CurrentL2: 190,
		CurrentL3: 191, // + scaler

		VoltageL1: 196, // uint16
		VoltageL2: 197,
		VoltageL3: 198, // + scaler

		Power: 200, // int16 + scaler
		// ApparentPower: 204, // int16 + scaler
		// ReactivePower: 206, // int16 + scaler
		Export: 210, // uint32 + scaler

		Cosphi:    208, // int16 + scaler
		Frequency: 202, // uint16 + scaler

		DCPower: 217, // int16 + scaler

		// DC block with global scale factors
		// DCCurrent1: 641,  // uint16
		// DCVoltage1: 642,  // uint16
		// DCPower1:   643, // uint16
		// DCCurrent2: 661,  // uint16
		// DCVoltage2: 662,  // uint16
		// DCPower2:   663, // uint16

		HeatSinkTemp: 219, // int16 + scaler
	}
	return &SMAProducer{
		MeasurementMapping{ops},
	}
}

func (p *SMAProducer) GetMeterType() string {
	return METERTYPE_SMA
}

func (p *SMAProducer) snip(iec Measurement, readlen uint16) Operation {
	return Operation{
		FuncCode: ReadHoldingReg,
		OpCode:   sunspecBase + p.Opcode(iec) - 1, // adjust according to docs
		ReadLen:  readlen,
		IEC61850: iec,
	}
}

func (p *SMAProducer) snip16uint(iec Measurement, scaler ...float64) Operation {
	snip := p.snip(iec, 1)

	snip.Transform = RTUUint16ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeRTUScaledUint16ToFloat64(scaler[0])
	}

	return snip
}

func (p *SMAProducer) snip16int(iec Measurement, scaler ...float64) Operation {
	snip := p.snip(iec, 1)

	snip.Transform = RTUInt16ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeRTUScaledInt16ToFloat64(scaler[0])
	}

	return snip
}

func (p *SMAProducer) snip32(iec Measurement, scaler ...float64) Operation {
	snip := p.snip(iec, 2)

	snip.Transform = RTUUint32ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeRTUScaledUint32ToFloat64(scaler[0])
	}

	return snip
}

func (p *SMAProducer) minMax(iec ...Measurement) (uint16, uint16) {
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
func (p *SMAProducer) scaleSnip16(splitter func(...Measurement) Splitter, iecs ...Measurement) Operation {
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

func (p *SMAProducer) scaleSnip32(splitter func(...Measurement) Splitter, iecs ...Measurement) Operation {
	op := p.scaleSnip16(splitter, iecs...)
	op.ReadLen = (op.ReadLen-1)*2 + 1 // read 4 bytes instead of 2 plus trailing scale factor
	return op
}

func (p *SMAProducer) mkSplitInt16(iecs ...Measurement) Splitter {
	return p.mkBlockSplitter(2, RTUInt16ToFloat64, iecs...)
}

func (p *SMAProducer) mkSplitUint16(iecs ...Measurement) Splitter {
	return p.mkBlockSplitter(2, RTUUint16ToFloat64WithNaN, iecs...)
}

func (p *SMAProducer) mkSplitUint32(iecs ...Measurement) Splitter {
	// use div 1000 for kWh conversion
	return p.mkBlockSplitter(4, MakeRTUScaledUint32ToFloat64(1000), iecs...)
}

func (p *SMAProducer) mkBlockSplitter(dataSize uint16, valFunc func([]byte) float64, iecs ...Measurement) Splitter {
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

			// filter results of RTUUint16ToFloat64WithNaN
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

func (p *SMAProducer) Probe() Operation {
	return p.snip16uint(VoltageL1, 10)
}

func (p *SMAProducer) Produce() (res []Operation) {
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
