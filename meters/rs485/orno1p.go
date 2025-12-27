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
		Frequency:     0x130,  // 16 bit, 0.01Hz
		Voltage:       0x131,  // 16 bit, 0.01V
		Current:       0x139,  // 32 bit, 0.001A
		Power:         0x140,  // 32 bit, 0.001kW
		ReactivePower: 0x148,  // 32 bit, 0.001kvar
		ApparentPower: 0x150,  // 32 bit, 0.001kva
		Cosphi:        0x158,  // 16 bit, 0,001
		Sum:           0xA000, // 32 Bit, 0.01kwh
		SumT1:         0xA002, // 32 Bit, 0.01kwh
		SumT2:         0xA004, // 32 Bit, 0.01kwh
		ReactiveSum:   0xA01E, // 32 Bit, 0.01kvarh
		ReactiveSumT1: 0xA020, // 32 Bit, 0.01kvarh
		ReactiveSumT2: 0xA022, // 32 Bit, 0.01kvarh
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
	return p.snip16(Voltage, 100)
}

// Produce implements Producer interface
func (p *ORNO1PProducer) Produce() (res []Operation) {

	for _, op := range []Measurement{
		Voltage,
		Frequency,
	} {
		res = append(res, p.snip16(op, 100))
	}

	for _, op := range []Measurement{
		Current,
	} {
		res = append(res, p.snip32(op, 1000))
	}

	for _, op := range []Measurement{
		Power, ReactivePower, ApparentPower,
	} {
		res = append(res, p.snip32(op, 1))
	}

	for _, op := range []Measurement{
		Cosphi,
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
