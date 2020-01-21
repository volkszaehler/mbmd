package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register(NewSBCProducer)
}

const (
	METERTYPE_SBC = "SBC"
)

type SBCProducer struct {
	typ    string
	phases int
	Opcodes
}

func NewSBCProducer() Producer {
	/**
	 * Opcodes for Saia Burgess ALE3
	 * https://www.sbc-support.com/uploads/tx_srcproducts/26-527_ENG_DS_EnergyMeter-ALE3-with-Modbus_01.pdf
	 * http://datenblatt.stark-elektronik.de/saia_burgess/DE_DS_Energymeter-ALE3-with-Modbus.pdf
	 */
	ops := Opcodes{
		Import: 28, // double, scaler 100
		Export: 32, // double, scaler 100
		// PartialImport: 30, // double, scaler 100
		// PartialExport: 34, // double, scaler 100

		VoltageL1:       36,
		CurrentL1:       37, // scaler 10
		PowerL1:         38, // scaler 100
		ReactivePowerL1: 39, // scaler 100
		CosphiL1:        40, // scaler 100

		VoltageL2:       41,
		CurrentL2:       42, // scaler 10
		PowerL2:         43, // scaler 100
		ReactivePowerL2: 44, // scaler 100
		CosphiL2:        45, // scaler 100

		VoltageL3:       46,
		CurrentL3:       47, // scaler 10
		PowerL3:         48, // scaler 100
		ReactivePowerL3: 49, // scaler 100
		CosphiL3:        50, // scaler 100

		Power:         51, // scaler 100
		ReactivePower: 52, // scaler 100
	}
	return &SBCProducer{
		typ:     "ALE3", // assume ALE3
		phases:  3,      // assume 3 phase device
		Opcodes: ops,
	}
}

// Type implements Producer interface
func (p *SBCProducer) Type() string {
	return METERTYPE_SBC
}

// Description implements Producer interface
func (p *SBCProducer) Description() string {
	return "Saia Burgess " + p.typ
}

// snip creates modbus operation
func (p *SBCProducer) snip(iec Measurement, readlen uint16) Operation {
	return Operation{
		FuncCode: ReadHoldingReg,
		OpCode:   p.Opcode(iec) - 1, // adjust according to docs
		ReadLen:  readlen,
		IEC61850: iec,
	}
}

// snip16 creates modbus operation for single register
func (p *SBCProducer) snip16(iec Measurement, scaler ...float64) Operation {
	snip := p.snip(iec, 1)

	snip.Transform = RTUUint16ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
	}

	return snip
}

// snip32 creates modbus operation for double register
func (p *SBCProducer) snip32(iec Measurement, scaler ...float64) Operation {
	snip := p.snip(iec, 2)

	snip.Transform = RTUUint32ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
	}

	return snip
}

// Identify implements Identifier interface
func (p *SBCProducer) Identify(bytes []byte) bool {
	if len(bytes) < 4 {
		return false
	}

	switch string(bytes[:4]) {
	case "ALD1":
		// single phase
		p.phases = 1
	case "ALE3", "AWE3":
		// three phase direct/ converter
		p.phases = 3
	default:
		return false
	}

	p.typ = string(bytes[:4])
	return true
}

// Probe implements Producer interface
func (p *SBCProducer) Probe() Operation {
	// return p.snip16(VoltageL1)
	return Operation{
		FuncCode: ReadHoldingReg,
		OpCode:   6,
		ReadLen:  4,
	}
}

// Produce implements Producer interface
func (p *SBCProducer) Produce() (res []Operation) {
	for _, op := range []Measurement{
		VoltageL1, VoltageL2, VoltageL3,
	} {
		res = append(res, p.snip16(op))
	}

	for _, op := range []Measurement{
		CurrentL1, CurrentL2, CurrentL3,
	} {
		res = append(res, p.snip16(op, 10))
	}

	for _, op := range []Measurement{
		PowerL1, PowerL2, PowerL3,
		CosphiL1, CosphiL2, CosphiL3,
	} {
		res = append(res, p.snip16(op, 100))
	}

	res = append(res, p.snip32(Import, 100))
	res = append(res, p.snip32(Export, 100))

	return res
}
