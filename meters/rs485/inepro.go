package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register(NewIneproProducer)
}

const (
	METERTYPE_INEPRO = "INEPRO"
)

type IneproProducer struct {
	Opcodes
}

func NewIneproProducer() Producer {
	/***
	 * https://ineprometering.com/wp-content/uploads/2018/09/PRO380-user-manual-V2-18.pdf
	 */
	ops := Opcodes{
		Voltage:   0x5000,
		VoltageL1: 0x5002,
		VoltageL2: 0x5004,
		VoltageL3: 0x5006,

		Frequency: 0x5008,

		Current:   0x500A,
		CurrentL1: 0x500C,
		CurrentL2: 0x500E,
		CurrentL3: 0x5010,

		ActivePower:   0x5012,
		ActivePowerL1: 0x5014,
		ActivePowerL2: 0x5016,
		ActivePowerL3: 0x5018,

		ReactivePower:   0x501A,
		ReactivePowerL1: 0x501C,
		ReactivePowerL2: 0x501E,
		ReactivePowerL3: 0x5020,

		ApparentPower:   0x5022,
		ApparentPowerL1: 0x5024,
		ApparentPowerL2: 0x5026,
		ApparentPowerL3: 0x5028,

		Cosphi:   0x502A,
		CosphiL1: 0x502C,
		CosphiL2: 0x502E,
		CosphiL3: 0x5030,

		Active:   0x6000,
		Reactive: 0x6024,
	}
	return &IneproProducer{Opcodes: ops}
}

func (p *IneproProducer) Type() string {
	return METERTYPE_INEPRO
}

func (p *IneproProducer) Description() string {
	return "Inepro Metering Pro 380 (experimental)"
}

func (p *IneproProducer) snip(iec Measurement) Operation {
	opcode := p.Opcodes[iec]
	return Operation{
		FuncCode:  ReadHoldingReg,
		OpCode:    opcode,
		ReadLen:   2,
		IEC61850:  iec,
		Transform: RTUIeee754ToFloat64,
	}
}

func (p *IneproProducer) Probe() Operation {
	return p.snip(VoltageL1)
}

func (p *IneproProducer) Produce() (res []Operation) {
	for op := range p.Opcodes {
		res = append(res, p.snip(op))
	}

	return res
}
