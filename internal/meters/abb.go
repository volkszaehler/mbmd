package meters

const (
	METERTYPE_ABB = "ABB"
)

type ABBProducer struct {
	ops Measurements
}

func NewABBProducer() *ABBProducer {
	/***
	 * http://datenblatt.stark-elektronik.de/Energiezaehler_B-Serie_Handbuch.pdf
	 */
	ops := Measurements{
		VoltageL1: 0x5B00,
		VoltageL2: 0x5B02,
		VoltageL3: 0x5B04,

		CurrentL1: 0x5B0C,
		CurrentL2: 0x5B0E,
		CurrentL3: 0x5B10,

		Power:   0x5B24, // Apparent Power
		PowerL1: 0x5B26,
		PowerL2: 0x5B28,
		PowerL3: 0x5B2A,

		Cosphi:   0x5B3A,
		CosphiL1: 0x5B3B,
		CosphiL2: 0x5B3C,
		CosphiL3: 0x5B3D,

		Frequency: 0x5B2C,
	}
	return &ABBProducer{
		ops: ops,
	}
}

func (p *ABBProducer) GetMeterType() string {
	return METERTYPE_ABB
}

func (p *ABBProducer) snip(iec Measurement, readlen uint16) Operation {
	opcode := p.ops[iec]
	return Operation{
		FuncCode: ReadHoldingReg,
		OpCode:   opcode,
		ReadLen:  readlen,
		IEC61850: iec,
	}
}

// snip16 creates modbus operation for single register
func (p *ABBProducer) snip16(iec Measurement, scaler ...float64) Operation {
	snip := p.snip(iec, 1)

	snip.Transform = RTUUint16ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeRTUScaledUint16ToFloat64(scaler[0])
	}

	return snip
}

// snip32 creates modbus operation for double register
func (p *ABBProducer) snip32(iec Measurement, scaler ...float64) Operation {
	snip := p.snip(iec, 2)

	snip.Transform = RTUUint32ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeRTUScaledUint32ToFloat64(scaler[0])
	}

	return snip
}

func (p *ABBProducer) Probe() Operation {
	return p.snip16(VoltageL1)
}

func (p *ABBProducer) Produce() (res []Operation) {
	for _, op := range []Measurement{
		VoltageL1, VoltageL2, VoltageL3,
		CurrentL1, CurrentL2, CurrentL3,
		Power, PowerL1, PowerL2, PowerL3,
	} {
		res = append(res, p.snip32(op, 100))
	}

	for _, op := range []Measurement{
		Cosphi, CosphiL1, CosphiL2, CosphiL3,
		Frequency,
	} {
		res = append(res, p.snip16(op))
	}

	return res
}
