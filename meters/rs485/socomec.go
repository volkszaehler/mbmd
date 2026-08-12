package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register("SOCOMEC", NewSocomecProducer)
}

type SocomecProducer struct {
	Opcodes
}

func NewSocomecProducer() Producer {
	/**
	 * Opcodes for Socomec Countis E3x/E4x, JBUS common table version 1.01.
	 * Metrology (0xC550) and energy (0xC650) tables, 32 bit registers, MSW first.
	 */
	ops := Opcodes{
		VoltageL1_L2: 0xC552, // V/100
		VoltageL2_L3: 0xC554,
		VoltageL3_L1: 0xC556,
		VoltageL1:    0xC558,
		VoltageL2:    0xC55A,
		VoltageL3:    0xC55C,
		Voltage:      0xC58C, // V sys

		Frequency: 0xC55E, // Hz/100

		CurrentL1: 0xC560, // mA
		CurrentL2: 0xC562,
		CurrentL3: 0xC564,
		Current:   0xC588, // I sys

		Power:   0xC568, // kW/100, signed
		PowerL1: 0xC570,
		PowerL2: 0xC572,
		PowerL3: 0xC574,

		ReactivePower:   0xC56A, // kvar/100, signed
		ReactivePowerL1: 0xC576,
		ReactivePowerL2: 0xC578,
		ReactivePowerL3: 0xC57A,

		ApparentPower:   0xC56C, // kVA/100
		ApparentPowerL1: 0xC57C,
		ApparentPowerL2: 0xC57E,
		ApparentPowerL3: 0xC580,

		Cosphi:   0xC56E, // 0.001, signed
		CosphiL1: 0xC582,
		CosphiL2: 0xC584,
		CosphiL3: 0xC586,

		Import:         0xC652, // kWh
		Export:         0xC658,
		ReactiveImport: 0xC654, // kvarh
		ReactiveExport: 0xC65A,
	}
	return &SocomecProducer{Opcodes: ops}
}

// Description implements Producer interface
func (p *SocomecProducer) Description() string {
	return "Socomec Countis E3x/E4x series"
}

func (p *SocomecProducer) snip(iec Measurement, transform RTUTransform, scaler ...float64) Operation {
	snip := Operation{
		FuncCode:  ReadHoldingReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   2,
		IEC61850:  iec,
		Transform: transform,
	}

	if len(scaler) > 0 {
		snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
	}

	return snip
}

// Probe implements Producer interface
func (p *SocomecProducer) Probe() Operation {
	return p.snip(VoltageL1, RTUUint32ToFloat64, 100)
}

// Produce implements Producer interface
func (p *SocomecProducer) Produce() (res []Operation) {
	for op := range p.Opcodes {
		switch op {
		case VoltageL1_L2, VoltageL2_L3, VoltageL3_L1, VoltageL1, VoltageL2, VoltageL3, Voltage, Frequency:
			res = append(res, p.snip(op, RTUUint32ToFloat64, 100))

		case CurrentL1, CurrentL2, CurrentL3, Current:
			res = append(res, p.snip(op, RTUUint32ToFloat64, 1000))

		// hundredths of kW/kvar/kVA, reported in W/var/VA
		case Power, PowerL1, PowerL2, PowerL3, ReactivePower, ReactivePowerL1, ReactivePowerL2, ReactivePowerL3:
			res = append(res, p.snip(op, RTUInt32ToFloat64, 0.1))

		case ApparentPower, ApparentPowerL1, ApparentPowerL2, ApparentPowerL3:
			res = append(res, p.snip(op, RTUUint32ToFloat64, 0.1))

		case Cosphi, CosphiL1, CosphiL2, CosphiL3:
			res = append(res, p.snip(op, RTUInt32ToFloat64, 1000))

		default:
			res = append(res, p.snip(op, RTUUint32ToFloat64))
		}
	}

	return res
}
