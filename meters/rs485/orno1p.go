package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register("ORNO1p", NewORNO1PProducer)
}

type ORNO1PProducer struct {
	Opcodes
}

func NewORNO1PProducer() Producer {
	/***
	 * Opcodes for ORNO WE-514 and WE-515
	 * https://github.com/gituser-rk/orno-modbus-mqtt/blob/master/Register%20description%20OR-WE-514%26OR-WE-515.pdf
	 */
	ops := Opcodes{
		Frequency: 0x130, // 16 bit, 0.01Hz

		VoltageL1:       0x131, // 16 bit, 0.01V
		CurrentL1:       0x139, // 16 bit, 0.001A
		PowerL1:         0x140, // 32 bit, 0.001kW
		ReactivePowerL1: 0x148, // 32 bit, 0.001kvar
		ApparentPowerL1: 0x150, // 32 bit, 0.001kva
		CosphiL1:        0x158, // 16 bit, 0,001

		VoltageL2:       0x132, // 16 bit, 0.01V
		CurrentL2:       0x13B, // 32 bit, 0.001A
		PowerL2:         0x142, // 32 bit, 0.001kW
		ReactivePowerL2: 0x14A, // 32 bit, 0.001kvar
		ApparentPowerL2: 0x152, // 32 bit, 0.001kva
		CosphiL2:        0x159, // 16 bit, 0,001

		VoltageL3:       0x133, // 16 bit, 0.01V
		CurrentL3:       0x13D, // 32 bit, 0.001A
		PowerL3:         0x144, // 32 bit, 0.001kW
		ReactivePowerL3: 0x14C, // 32 bit, 0.001kvar
		ApparentPowerL3: 0x154, // 32 bit, 0.001kva
		CosphiL3:        0x15A, // 16 bit, 0,001

		Power:         0x146, // 32 bit, 0.001kW
		ReactivePower: 0x14E, // 32 bit, 0.001kvar
		ApparentPower: 0x156, // 32 bit, 0.001kva
		Cosphi:        0x15B, // 16 bit, 0.001

		Sum:   0xA000, //32 Bit, 0.01kwh
		SumT1: 0xA002, //32 Bit, 0.01kwh
		SumT2: 0xA004, //32 Bit, 0.01kwh
		//		SumT3:           0xA006, //32 Bit, 0.01kwh // currently not supported
		//		SumT4:           0xA008, //32 Bit, 0.01kwh // currently not supported
		ReactiveSum:   0xA01E, //32 Bit, 0.01kvarh
		ReactiveSumT1: 0xA020, //32 Bit, 0.01kvarh
		ReactiveSumT2: 0xA022, //32 Bit, 0.01kvarh
		//		ReactiveSumT3:   0xA024, //32 Bit, 0.01kvarh // currently not supported
		//		ReactiveSumT4:   0xA026, //32 Bit, 0.01kvarh // currently not supported
	}

	return &ORNO1PProducer{Opcodes: ops}
}

// Description implements Producer interface
func (p *ORNO1PProducer) Description() string {
	return "ORNO WE-514 & WE-515"
}

// snip creates modbus operation
func (p *ORNO1PProducer) snip(iec Measurement, readlen uint16) Operation {
	return Operation{
		FuncCode: ReadHoldingReg,
		OpCode:   p.Opcode(iec), // adjust according to docs
		ReadLen:  readlen,
		IEC61850: iec,
	}
}

// snip16 creates modbus operation for single register
func (p *ORNO1PProducer) snip16(iec Measurement, scaler ...float64) Operation {
	snip := p.snip(iec, 1)

	snip.Transform = RTUUint16ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
	}

	return snip
}

// snip32 creates modbus operation for double register
func (p *ORNO1PProducer) snip32(iec Measurement, scaler ...float64) Operation {
	snip := p.snip(iec, 2)

	snip.Transform = RTUUint32ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
	}

	return snip
}

func (p *ORNO1PProducer) Probe() Operation {
	return p.snip16(VoltageL1, 100)
}

// Produce implements Producer interface
func (p *ORNO1PProducer) Produce() (res []Operation) {

	for _, op := range []Measurement{
		VoltageL1,
		Frequency,
	} {
		res = append(res, p.snip16(op, 100))
	}

	for _, op := range []Measurement{
		CurrentL1,
	} {
		res = append(res, p.snip32(op, 1000))
	}

	for _, op := range []Measurement{
		PowerL1, ReactivePowerL1, ApparentPowerL1,
	} {
		res = append(res, p.snip32(op, 1))
	}

	for _, op := range []Measurement{
		CosphiL1,
	} {
		res = append(res, p.snip16(op, 1000))
	}

	for _, op := range []Measurement{
		Sum, ReactiveSum,
	} {
		res = append(res, p.snip32(op, 100))
	}
	return res
}
