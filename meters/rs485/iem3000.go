package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register("IEM3000", NewIEM3000Producer)
}

type IEM3000Producer struct {
	Opcodes
}

func NewIEM3000Producer() Producer {
	/***
	 * https://download.schneider-electric.com/files?p_enDocType=User+guide&p_File_Name=DOCA0005DE-12.pdf&p_Doc_Ref=DOCA0005DE#page49
	 */
	ops := Opcodes{
		VoltageL1: 0x0BD4,
		VoltageL2: 0x0BD6,
		VoltageL3: 0x0BD8,
		Voltage:   0x0BDC,

		CurrentL1: 0x0BB8,
		CurrentL2: 0x0BBA,
		CurrentL3: 0x0BBC,
		Current:   0x0BC2,

		PowerL1: 0x0BEE,
		PowerL2: 0x0BF0,
		PowerL3: 0x0BF2,
		Power:   0x0BF4,

		ReactivePower: 0x0BFC,
		ApparentPower: 0x0C04,

		// PowerFactor: 0x0C0C,
		Frequency: 0x0C26,

		Import:   0x0C84,
		ImportL1: 0x0DBE,
		ImportL2: 0x0DC2,
		ImportL3: 0x0DC6,
		Export:   0x0C88,

		ReactiveImport: 0x0C94,
		ReactiveExport: 0x0C98,
	}
	return &IEM3000Producer{Opcodes: ops}
}

// Description implements Producer interface
func (p *IEM3000Producer) Description() string {
	return "Schneider Electric iEM3000 series"
}

func (p *IEM3000Producer) snipFloat32(iec Measurement, scaler ...float64) Operation {
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

func (p *IEM3000Producer) snipInt64(iec Measurement, scaler ...float64) Operation {
	snip := Operation{
		FuncCode:  ReadHoldingReg,
		OpCode:    p.Opcodes[iec],
		ReadLen:   4,
		IEC61850:  iec,
		Transform: RTUInt64ToFloat64,
	}

	if len(scaler) > 0 {
		snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
	}

	return snip
}

// Probe implements Producer interface
func (p *IEM3000Producer) Probe() Operation {
	return p.snipFloat32(VoltageL1)
}

// Produce implements Producer interface
func (p *IEM3000Producer) Produce() (res []Operation) {
	for op := range p.Opcodes {
		switch op {
		case PowerL1, PowerL2, PowerL3, Power, ReactivePower, ApparentPower:
			res = append(res, p.snipFloat32(op, 0.001))
		case Import, ImportL1, ImportL2, ImportL3, Export, ReactiveImport, ReactiveExport:
			res = append(res, p.snipInt64(op, 1000))
		default:
			res = append(res, p.snipFloat32(op))
		}
	}

	return res
}
