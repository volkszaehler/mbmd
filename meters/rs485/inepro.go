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

		Power:   0x5012,
		PowerL1: 0x5014,
		PowerL2: 0x5016,
		PowerL3: 0x5018,

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

		Sum:    0x6000,
		SumL1:  0x6006,
		SumL2:  0x6008,
		SumL3:  0x600A,
		Import: 0x600C,
		Export: 0x6018,

		Reactive:       0x6024,
		ReactiveL1:     0x602A,
		ReactiveL2:     0x602C,
		ReactiveL3:     0x602E,
		ReactiveImport: 0x6030,
		ReactiveExport: 0x603C,
	}
	return &IneproProducer{Opcodes: ops}
}

func (p *IneproProducer) Type() string {
	return METERTYPE_INEPRO
}

func (p *IneproProducer) Description() string {
	return "Inepro Metering Pro 380"
}

func (p *IneproProducer) snip(iec Measurement, scaler ...float64) Operation {
	snip := Operation{
		FuncCode:  ReadHoldingReg,
		OpCode:    p.Opcodes[iec],
		ReadLen:   2,
		IEC61850:  iec,
		Transform: RTUIeee754ToFloat64,
	}

	if len(scaler) > 0 {
		snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
	}

	return snip
}

func (p *IneproProducer) Probe() Operation {
	return p.snip(VoltageL1)
}

func (p *IneproProducer) Produce() (res []Operation) {
	for op := range p.Opcodes {
		switch op {
		case Power, PowerL1, PowerL2, PowerL3,
			ReactivePower, ReactivePowerL1, ReactivePowerL2, ReactivePowerL3,
			ApparentPower, ApparentPowerL1, ApparentPowerL2, ApparentPowerL3:
			res = append(res, p.snip(op, 1000))
		default:
			res = append(res, p.snip(op))
		}
	}

	return res
}
