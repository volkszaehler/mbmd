package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register(NewIEM3000Producer)
}

const (
	METERTYPE_IEM3000 = "IEM3000"
)

type IEM3000Producer struct {
	Opcodes
}

func NewIEM3000Producer() Producer {
	/***
	 * https://download.schneider-electric.com/files?p_enDocType=User+guide&p_File_Name=DOCA0005DE-12.pdf&p_Doc_Ref=DOCA0005DE#page49
	 */
	ops := Opcodes{
		VoltageL1: 3028,
		VoltageL2: 3030,
		VoltageL3: 3032,
		Voltage:   3036,

		CurrentL1: 3000,
		CurrentL2: 3002,
		CurrentL3: 3004,
		Current:   3010,

		PowerL1: 3054,
		PowerL2: 3056,
		PowerL3: 3058,
		Power:   3060,

		ReactivePower: 3068,
		ApparentPower: 3076,

		// PowerFactor: 3084,
		Frequency: 3110,

		Import:   3204,
		ImportL1: 3518,
		ImportL2: 3522,
		ImportL3: 3526,
		Export:   3208,

		ReactiveImport: 3220,
		ReactiveExport: 3224,
	}
	return &IEM3000Producer{Opcodes: ops}
}

// Type implements Producer interface
func (p *IEM3000Producer) Type() string {
	return METERTYPE_IEM3000
}

// Description implements Producer interface
func (p *IEM3000Producer) Description() string {
	return "Schneider Electric iEM3000 series"
}

func (p *IEM3000Producer) snipF(iec Measurement, scaler ...float64) Operation {
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

// Probe implements Producer interface
func (p *IEM3000Producer) Probe() Operation {
	return p.snipF(VoltageL1)
}

// Produce implements Producer interface
func (p *IEM3000Producer) Produce() (res []Operation) {
	for op := range p.Opcodes {
		switch op {
		case PowerL1, PowerL2, PowerL3, Power, ReactivePower, ApparentPower:
			res = append(res, p.snipF(op, 0.001))
		case Import, ImportL1, ImportL2, ImportL3, Export, ReactiveImport, ReactiveExport:
			snip := Operation{
				FuncCode:  ReadHoldingReg,
				OpCode:    p.Opcodes[op],
				ReadLen:   2,
				IEC61850:  op,
				Transform: MakeScaledTransform(RTUInt64ToFloat64, 1000),
			}
			res = append(res, snip)
		default:
			res = append(res, p.snipF(op))
		}
	}

	return res
}
