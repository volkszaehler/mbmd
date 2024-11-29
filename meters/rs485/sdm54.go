package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register("SDM54", NewSDM54Producer)
}

type SDM54Producer struct {
	Opcodes
}

func NewSDM54Producer() Producer {
	/**
	 * Opcodes as defined by Eastron SDM54.
	 * See https://www.eastrongroup.com/eastrongroup/2024/08/21/eastronsdm54seriesusermanualv1.2.pdf
	 * This is to a large extent a superset of all SDM devices, however there are
	 * subtle differences (see 220, 230). Some opcodes might not work on some devices.
	 */
	ops := Opcodes{
		VoltageL1:       0x0000, // 220, 230
		VoltageL2:       0x0002,
		VoltageL3:       0x0004,
		CurrentL1:       0x0006, // 220, 230
		CurrentL2:       0x0008,
		CurrentL3:       0x000A,
		PowerL1:         0x000C, //      230
		PowerL2:         0x000E,
		PowerL3:         0x0010,
		ApparentPowerL1: 0x0012,
		ApparentPowerL2: 0x0014,
		ApparentPowerL3: 0x0016,
		ReactivePowerL1: 0x0018,
		ReactivePowerL2: 0x001A,
		ReactivePowerL3: 0x001C,
		CosphiL1:        0x001e, //      230
		CosphiL2:        0x0020,
		CosphiL3:        0x0022,
		Power:           0x0034,
		ApparentPower:   0x0038,
		ReactivePower:   0x003C,
		ImportPower:     0x0054,
		ImportL1:        0x015a,
		ImportL2:        0x015c,
		ImportL3:        0x015e,
		Import:          0x0048, // 220, 230
		ExportL1:        0x0160,
		ExportL2:        0x0162,
		ExportL3:        0x0164,
		Export:          0x004a, // 220, 230
		SumL1:           0x0166,
		SumL2:           0x0168,
		SumL3:           0x016a,
		Sum:             0x0156, // 220
		Cosphi:          0x003e,
		THDL1:           0x00ea, // voltage
		THDL2:           0x00ec, // voltage
		THDL3:           0x00ee, // voltage
		THD:             0x00F8, // voltage
		Frequency:       0x0046, //      230
	}
	return &SDM54Producer{Opcodes: ops}
}

func (p *SDM54Producer) Description() string {
	return "Eastron SDM54"
}

func (p *SDM54Producer) snip(iec Measurement) Operation {
	operation := Operation{
		FuncCode:  ReadInputReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   2,
		IEC61850:  iec,
		Transform: RTUIeee754ToFloat64,
	}
	return operation
}

func (p *SDM54Producer) Probe() Operation {
	return p.snip(VoltageL1)
}

func (p *SDM54Producer) Produce() (res []Operation) {
	for op := range p.Opcodes {
		res = append(res, p.snip(op))
	}

	return res
}
