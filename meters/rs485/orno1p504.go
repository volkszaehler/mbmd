package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register("ORNO1P504", NewORNO1P504Producer)
}

var ops1p504 Opcodes = Opcodes{

	Frequency:     0x02, // 16 bit, Hz
	Voltage:       0x00, // 16 bit, V
	Current:       0x01, // 16 bit, A
	Power:         0x03, // 16 bit, W
	ReactivePower: 0x04, // 16 bit, var
	ApparentPower: 0x05, // 16 bit, va
	Cosphi:        0x06, // 16 bit,

	Sum:         0x07, //32 Bit, wh
	ReactiveSum: 0x09, //32 Bit, varh
}

type ORNO1P504Producer struct {
	Opcodes
}

func NewORNO1P504Producer() Producer {
	return &ORNO1P504Producer{Opcodes: ops1p504}
}

// Description implements Producer interface
func (p *ORNO1P504Producer) Description() string {
	return "ORNO WE-504"
}

// snip creates modbus operation
func (p *ORNO1P504Producer) snip(iec Measurement, readlen uint16) Operation {
	return Operation{
		FuncCode: ReadHoldingReg,
		OpCode:   p.Opcode(iec), // adjust according to docs
		ReadLen:  readlen,
		IEC61850: iec,
	}
}

// snip16 creates modbus operation for single register
func (p *ORNO1P504Producer) snip16(iec Measurement, scaler ...float64) Operation {
	snip := p.snip(iec, 1)

	snip.Transform = RTUUint16ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
	}

	return snip
}

// snip32 creates modbus operation for double register
func (p *ORNO1P504Producer) snip32(iec Measurement, scaler ...float64) Operation {
	snip := p.snip(iec, 2)

	snip.Transform = RTUUint32ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
	}

	return snip
}

func (p *ORNO1P504Producer) Probe() Operation {
	return p.snip32(Voltage, 1)
}

// Produce implements Producer interface
func (p *ORNO1P504Producer) Produce() (res []Operation) {
	for _, op := range []Measurement{
		Power, ReactivePower, ApparentPower,
	} {
		res = append(res, p.snip16(op, 1))
	}

	for _, op := range []Measurement{
		Frequency, Voltage, Current,
	} {
		res = append(res, p.snip16(op, 10))
	}

	for _, op := range []Measurement{
		Cosphi,
	} {
		res = append(res, p.snip16(op, 1000))
	}

	for _, op := range []Measurement{
		ReactiveSum, Sum,
	} {
		res = append(res, p.snip32(op, 1000))
	}

	return res
}
