package meters

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
	 * Opcodes for SunSpec- compatible Inverters like SolarEdge
	 * https://www.solaredge.com/sites/default/files/sunspec-implementation-technical-note.pdf
	 */
	ops := Measurements{
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

func (p *SUNProducer) Probe() Operation {
	return p.snip16uint(VoltageL1, 10)
}

func (p *SUNProducer) Produce() (res []Operation) {
	for _, op := range []Measurement{
		VoltageL1, VoltageL2, VoltageL1,
	} {
		res = append(res, p.snip16uint(op, 10))
	}

	for _, op := range []Measurement{
		CurrentL1, CurrentL2, CurrentL1,
		Frequency,
	} {
		res = append(res, p.snip16uint(op, 100))
	}

	for _, op := range []Measurement{
		Power, Cosphi, DCPower, HeatSinkTemp,
	} {
		res = append(res, p.snip16int(op, 100))
	}

	for _, op := range []Measurement{
		DCCurrent, DCVoltage,
	} {
		res = append(res, p.snip16uint(op, 10))
	}

	for _, op := range []Measurement{
		Export,
	} {
		res = append(res, p.snip32(op, 1000))
	}

	return res
}
