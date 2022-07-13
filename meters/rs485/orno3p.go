package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register("ORNO3p", NewORNO3PProducer)
}

type ORNO3PProducer struct {
	Opcodes
}

func NewORNO3PProducer() Producer {
	/***
	 * Opcodes for ORNO WE-514 and WE-515
	 * https://github.com/gituser-rk/orno-modbus-mqtt/blob/master/Register%20description%20OR-WE-514%26OR-WE-515.pdf
	 */
	ops := Opcodes{
		Frequency: 0x0014, // 32 bit, Hz

		VoltageL1:       0x000E, // 32 bit, V
		CurrentL1:       0x0016, // 32 bit, A
		PowerL1:         0x001E, // 32 bit, kW
		ReactivePowerL1: 0x0026, // 32 bit, kvar
		ApparentPowerL1: 0x002E, // 32 bit, kva
		CosphiL1:        0x0036, // 32 bit, XX,X(literal)

		VoltageL2:       0x0010, // 32 bit, V
		CurrentL2:       0x0018, // 32 bit, A
		PowerL2:         0x0020, // 32 bit, kW
		ReactivePowerL2: 0x0028, // 32 bit, kvar
		ApparentPowerL2: 0x0030, // 32 bit, kva
		CosphiL2:        0x0038, // 32 bit, XX,X(literal)

		VoltageL3:       0x0012, // 32 bit, V
		CurrentL3:       0x001A, // 32 bit, A
		PowerL3:         0x0022, // 32 bit, kW
		ReactivePowerL3: 0x002A, // 32 bit, kvar
		ApparentPowerL3: 0x0032, // 32 bit, kva
		CosphiL3:        0x003A, // 32 bit, XX,X(literal)

		Power:         0x001C, // 32 bit, kW
		ReactivePower: 0x0024, // 32 bit, kvar
		ApparentPower: 0x002C, // 32 bit, kva
		Cosphi:        0x0034, // 32 bit, XX,X(literal)

		Sum:   0x0100, //32 Bit, kwh
		SumL1: 0x0102, //32 Bit, kwh
		SumL2: 0x0104, //32 Bit, kwh
		SumL3: 0x0106, //32 Bit, kwh

		Import:   0x0108, //32 Bit, kwh
		ImportL1: 0x010A, //32 Bit, kwh
		ImportL2: 0x010C, //32 Bit, kwh
		ImportL3: 0x010E, //32 Bit, kwh

		Export:   0x0110, //32 Bit, kwh
		ExportL1: 0x0112, //32 Bit, kwh
		ExportL2: 0x0114, //32 Bit, kwh
		ExportL3: 0x0116, //32 Bit, kwh

		ReactiveSum:   0x0118, //32 Bit, kvarh
		ReactiveSumL1: 0x011A, //32 Bit, kvarh
		ReactiveSumL2: 0x011C, //32 Bit, kvarh
		ReactiveSumL3: 0x011E, //32 Bit, kvarh

		ReactiveImport:   0x0120, //32 Bit, kvarh
		ReactiveImportL1: 0x0122, //32 Bit, kvarh
		ReactiveImportL2: 0x0124, //32 Bit, kvarh
		ReactiveImportL3: 0x0126, //32 Bit, kvarh

		ReactiveExport:   0x0128, //32 Bit, kvarh
		ReactiveExportL1: 0x012A, //32 Bit, kvarh
		ReactiveExportL2: 0x012C, //32 Bit, kvarh
		ReactiveExportL3: 0x012E, //32 Bit, kvarh

		SumT1:            0x0130, //32 Bit, kwh
		ImportT1:         0x0132, //32 Bit, kwh
		ExportT1:         0x0134, //32 Bit, kwh
		ReactiveSumT1:    0x0136, //32 Bit, kvarh
		ReactiveImportT1: 0x0138, //32 Bit, kvarh
		ReactiveExportT1: 0x013A, //32 Bit, kvarh

		SumT2:            0x013C, //32 Bit, kwh
		ImportT2:         0x013E, //32 Bit, kwh
		ExportT2:         0x0140, //32 Bit, kwh
		ReactiveSumT2:    0x0142, //32 Bit, kvarh
		ReactiveImportT2: 0x0144, //32 Bit, kvarh
		ReactiveExportT2: 0x0146, //32 Bit, kvarh

		/* // Curently not supported
		SumT3:           0x0148, //32 Bit, kwh
		ImportT3:        0x014A, //32 Bit, kwh
		ExportT3:        0x014C, //32 Bit, kwh
		ReactiveSumT3:   0x015E, //32 Bit, kvarh
		ReactiveImportT3:0x0150, //32 Bit, kvarh
		ReactiveExportT3:0x0152, //32 Bit, kvarh

		SumT4:           0x0154, //32 Bit, kwh
		ImportT4:        0x0156, //32 Bit, kwh
		ExportT4:        0x0158, //32 Bit, kwh
		ReactiveSumT4:   0x015A, //32 Bit, kvarh
		ReactiveImportT4:0x015C, //32 Bit, kvarh
		ReactiveExportT4:0x015E, //32 Bit, kvarh
		*/
	}

	return &ORNO3PProducer{Opcodes: ops}
}

// Description implements Producer interface
func (p *ORNO3PProducer) Description() string {
	return "ORNO WE-516 & WE-517"
}

// snip creates modbus operation
func (p *ORNO3PProducer) snip(iec Measurement, readlen uint16) Operation {
	return Operation{
		FuncCode: ReadHoldingReg,
		OpCode:   p.Opcode(iec), // adjust according to docs
		ReadLen:  readlen,
		IEC61850: iec,
	}
}

// snip32 creates modbus operation for double register
func (p *ORNO3PProducer) snip32(iec Measurement, scaler ...float64) Operation {
	snip := p.snip(iec, 2)

	snip.Transform = RTUIeee754ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
	}

	return snip
}

func (p *ORNO3PProducer) Probe() Operation {
	return p.snip32(VoltageL1, 1)
}

// Produce implements Producer interface
func (p *ORNO3PProducer) Produce() (res []Operation) {
	// These values are stored as literals
	for _, op := range []Measurement{
		Frequency,
		VoltageL1, CurrentL1, CosphiL1,
		VoltageL2, CurrentL2, CosphiL2,
		VoltageL3, CurrentL3, CosphiL3,
		Cosphi,
		Sum, SumL1, SumL2, SumL3,
		Import, ImportL1, ImportL2, ImportL3,
		Export, ExportL1, ExportL2, ExportL3,
		ReactiveSum, ReactiveSumL1, ReactiveSumL2, ReactiveSumL3,
		ReactiveImport, ReactiveImportL1, ReactiveImportL2, ReactiveImportL3,
		ReactiveExport, ReactiveExportL1, ReactiveExportL2, ReactiveExportL3,
		SumT1, ImportT1, ExportT1, ReactiveSumT1, ReactiveImportT1, ReactiveExportT1,
		SumT2, ImportT2, ExportT2, ReactiveSumT2, ReactiveImportT2, ReactiveExportT2,
	} {
		res = append(res, p.snip32(op, 1))
	}

	// For Power values, we need to scale by 1000 (aka convert kW/kva -> W/va)
	for _, op := range []Measurement{
		PowerL1, ReactivePowerL1, ApparentPowerL1,
		PowerL2, ReactivePowerL2, ApparentPowerL2,
		PowerL3, ReactivePowerL3, ApparentPowerL3,
		Power, ReactivePower, ApparentPower,
	} {
		res = append(res, p.snip32(op, 0.001))
	}
	return res
}
