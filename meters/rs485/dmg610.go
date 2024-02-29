package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register("DMG610", NewDMG610Producer)
}

type DMG610Producer struct {
	Opcodes
}

func NewDMG610Producer() Producer {
	/**
	 * Opcodes as defined by Lovato DMG Series Protocol Manual.
	 * See https://www.manualslib.com/manual/2401832/Lovato-Dmg-Series.html
	 * As for Modbus® standard, the address in the query message must be decreased by one from the effective address reported in the table!
	 */
	ops := Opcodes{
		Frequency:       0x0031, // Frequency
		CurrentL1:       0x0007, // L1 Current
		CurrentL2:       0x0009, // L2 Current
		CurrentL3:       0x000B, // L3 Current
		VoltageL1:       0x0001, // L1 Phase Voltage
		VoltageL2:       0x0003, // L2 Phase Voltage
		VoltageL3:       0x0005, // L3 Phase Voltage
		PowerL1:         0x0013, // L1 Active Power
		PowerL2:         0x0015, // L2 Active Power
		PowerL3:         0x0017, // L3 Active Power
		ReactivePowerL1: 0x0019, // L1 Reactive Power
		ReactivePowerL2: 0x001B, // L2 Reactive Power
		ReactivePowerL3: 0x001D, // L3 Reactive Power
		ApparentPowerL1: 0x001F, // L1 Apparent Power
		ApparentPowerL2: 0x0021, // L2 Apparent Power
		ApparentPowerL3: 0x0023, // L3 Apparent Power
		CosphiL1:        0x0025, // L1 Power Factor
		CosphiL2:        0x0027, // L2 Power Factor
		CosphiL3:        0x0029, // L3 Power Factor
		THDL1:           0x0053, // L1 Voltage Thd
		THDL2:           0x0055, // L2 Voltage Thd
		THDL3:           0x0057, // L3 Voltage Thd
		Import:          0x1B1F, // Total imp. Active Energy
		ImportL1:        0x1E1F, // Energia attiva L1 importata
		ImportL2:        0x1E47, // Energia attiva L2 importata
		ImportL3:        0x1E6F, // Energia attiva L3 importata
		Export:          0x1B23, // Total exported Active Energy
		ExportL1:        0x1E23, // Energia attiva L1 esportata
		ExportL2:        0x1E4B, // Energia attiva L2 esportata
		ExportL3:        0x1E73, // Energia attiva L3 esportata
		ReactiveImport:  0x1B27, // Total imp. Reactive Energy
		ReactiveExport:  0x1B2B, // Total exp. Reactive Energy
	}
	return &DMG610Producer{Opcodes: ops}
}

func (p *DMG610Producer) Description() string {
	return "Lovato DMG610"
}

func (p *DMG610Producer) snip(iec Measurement, readlen uint16, transform RTUTransform, scaler ...float64) Operation {
	snip := Operation{
		FuncCode:  ReadInputReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   readlen,
		IEC61850:  iec,
		Transform: transform,
	}

	if len(scaler) > 0 {
		snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
	}

	return snip
}

func (p *DMG610Producer) snip64u(iec Measurement, scaler ...float64) Operation {
	return p.snip(iec, 4, RTUUint64ToFloat64, scaler...)
}

func (p *DMG610Producer) snip32u(iec Measurement, scaler ...float64) Operation {
	return p.snip(iec, 2, RTUUint32ToFloat64, scaler...)
}

func (p *DMG610Producer) snip32(iec Measurement, scaler ...float64) Operation {
	return p.snip(iec, 2, RTUInt32ToFloat64, scaler...)
}

func (p *DMG610Producer) Probe() Operation {
	return p.snip32u(VoltageL1, 100)
}

func (p *DMG610Producer) Produce() (res []Operation) {
	for op := range p.Opcodes {
		switch op {
		case CurrentL1, CurrentL2, CurrentL3:
			res = append(res, p.snip32u(op, 10000))
		case Frequency:
			res = append(res, p.snip32u(op, 1000))
		case VoltageL1, VoltageL2, VoltageL3,
			ApparentPowerL1, ApparentPowerL2, ApparentPowerL3,
			THDL1, THDL2, THDL3:
			res = append(res, p.snip32u(op, 100))
		case CosphiL1, CosphiL2, CosphiL3:
			res = append(res, p.snip32(op, 10000))
		case PowerL1, PowerL2, PowerL3,
			ReactivePowerL1, ReactivePowerL2, ReactivePowerL3:
			res = append(res, p.snip32(op, 100))
		default:
			res = append(res, p.snip64u(op, 100))
		}
	}

	return res
}
