package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register("DS100", NewDS100Producer)
}

type DS100Producer struct {
	Opcodes
}

func NewDS100Producer() Producer {
	/**
	 * Opcodes as defined by B+G e-tech DS100.
	 * See https://data.xn--stromzhler-v5a.eu/manuals/bg_ds100serie_de.pdf
	 */
	ops := Opcodes{
		VoltageL1:       0x0400,
		VoltageL2:       0x0402,
		VoltageL3:       0x0404,
		Voltage:         0x040C,
		CurrentL1:       0x0410,
		CurrentL2:       0x0412,
		CurrentL3:       0x0414,
		Current:         0x0418,
		PowerL1:         0x041A,
		PowerL2:         0x041C,
		PowerL3:         0x041E,
		Power:           0x0420,
		ApparentPowerL1: 0x0422,
		ApparentPowerL2: 0x0424,
		ApparentPowerL3: 0x0426,
		ApparentPower:   0x0428,
		ReactivePowerL1: 0x042A,
		ReactivePowerL2: 0x042C,
		ReactivePowerL3: 0x042E,
		ReactivePower:   0x0430,
		Frequency:       0x0435,
		CosphiL1:        0x0436,
		CosphiL2:        0x0437,
		CosphiL3:        0x0438,
		Cosphi:          0x0439,
		ImportL1:        0x050A, //A phase forward active energy
		ImportL2:        0x056E, //B phase forward active energy
		ImportL3:        0x05D2, //C phase forward active energy
		Import:          0x010E, //  total forward active energy
		ExportL1:        0x0514, //A phase reverse active energy
		ExportL2:        0x0578, //B phase reverse active energy
		ExportL3:        0x05DC, //C phase reverse active energy
		Export:          0x118A, //  total reverse active energy
		SumL1:           0x0500, //A phase total active energy
		SumL2:           0x0564, //B phase total active energy
		SumL3:           0x05C8, //C phase total active energy
		Sum:             0x0122, //  total total active energy
	}
	return &DS100Producer{Opcodes: ops}
}

func (p *DS100Producer) Description() string {
	return "B+G e-tech DS100"
}

func (p *DS100Producer) snip(iec Measurement, readlen uint16, transform RTUTransform, scaler ...float64) Operation {
	snip := Operation{
		FuncCode:  ReadHoldingReg,
		OpCode:    p.Opcodes[iec],
		ReadLen:   readlen,
		Transform: transform,
		IEC61850:  iec,
	}

	if len(scaler) > 0 {
		snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
	}

	return snip
}

// snip16u creates modbus operation for single register
func (p *DS100Producer) snip16u(iec Measurement, scaler ...float64) Operation {
	return p.snip(iec, 1, RTUUint16ToFloat64, scaler...)
}

// snip32u creates modbus operation for double register
func (p *DS100Producer) snip32u(iec Measurement, scaler ...float64) Operation {
	return p.snip(iec, 2, RTUUint32ToFloat64, scaler...)
}

// snip16s creates modbus operation for single register (signed)
func (p *DS100Producer) snip16s(iec Measurement, scaler ...float64) Operation {
	return p.snip(iec, 1, RTUInt16ToFloat64, scaler...)
}

// snip32s creates modbus operation for double register (signed)
func (p *DS100Producer) snip32s(iec Measurement, scaler ...float64) Operation {
	return p.snip(iec, 2, RTUInt32ToFloat64, scaler...)
}

func (p *DS100Producer) Probe() Operation {
	return p.snip32u(Voltage, 1000)
}

func (p *DS100Producer) Produce() (res []Operation) {
	for _, op := range []Measurement{
		VoltageL1, VoltageL2, VoltageL3, Voltage,
	} {
		res = append(res, p.snip32u(op, 1000))
	}

	for _, op := range []Measurement{
		CurrentL1, CurrentL2, CurrentL3, Current,
	} {
		res = append(res, p.snip32s(op, 1000))
	}

	for _, op := range []Measurement{
		PowerL1, PowerL2, PowerL3, Power,
		ReactivePowerL1, ReactivePowerL2, ReactivePowerL3, ReactivePower,
	} {
		res = append(res, p.snip32s(op, 1))
	}

	for _, op := range []Measurement{
		ApparentPowerL1, ApparentPowerL2, ApparentPowerL3, ApparentPower,
	} {
		res = append(res, p.snip32u(op, 1))
	}

	for _, op := range []Measurement{
		ImportL1, ImportL2, ImportL3, Import,
		ExportL1, ExportL2, ExportL3, Export,
		SumL1, SumL2, SumL3, Sum,
	} {
		res = append(res, p.snip32u(op, 100))
	}

	for _, op := range []Measurement{
		Frequency,
	} {
		res = append(res, p.snip16u(op, 10))
	}

	for _, op := range []Measurement{
		CosphiL1, CosphiL2, CosphiL3, Cosphi,
	} {
		res = append(res, p.snip16s(op, 1000))
	}

	return res
}
