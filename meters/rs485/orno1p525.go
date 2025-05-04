package rs485

import (
	"github.com/grid-x/modbus"
	. "github.com/volkszaehler/mbmd/meters"
)

func init() {
	Register("ORNO1p525", NewORNO1P525Producer)
}

type ORNO1P525Producer struct {
	Opcodes
}

func NewORNO1P525Producer() Producer {
	/***
	 * Opcodes for ORNO WE-525 and WE-526
	 * https://files.orno.pl/support/Others/ORNO/ORWE525_5908254827846/OR-WE-525_rejestry.pdf
	 */
	ops := Opcodes{
		Frequency: 0x10A, // 16 bit, 0.1Hz

		Voltage:       0x100, // 32 bit, 0.001V
		Current:       0x102, // 32 bit, 0.001A
		Power:         0x104, // 32 bit, 0W
		ReactivePower: 0x108, // 32 bit, 0var
		ApparentPower: 0x106, // 32 bit, 0va
		Cosphi:        0x10B, // 16 bit, 0.001

		Import:   0x010E, //32 Bit, 0.00kWh
		ImportT1: 0x0110, //32 Bit, 0.00kWh
		ImportT2: 0x0112, //32 Bit, 0.00kWh
		//ImportT3: 0x0114, //32 Bit, 0.00kWh  // currently not supported
		//ImportT4: 0x0116, //32 Bit, 0.00kWh  // currently not supported

		Export:   0x0118, //32 Bit, 0.00kWh
		ExportT1: 0x011A, //32 Bit, 0.00kWh
		ExportT2: 0x011C, //32 Bit, 0.00kWh
		//ExportT3: 0x011E, //32 Bit, 0.00kWh  // currently not supported
		//ExportT4: 0x0120, //32 Bit, 0.00kWh  // currently not supported

		Sum:   0x0122, //32 Bit, 0.01kWh
		SumT1: 0x0124, //32 Bit, 0.01kWh
		SumT2: 0x0126, //32 Bit, 0.01kWh
		//SumT3:   0x0128, //32 Bit, 0.01kWh  // currently not supported
		//SumT4:   0x012A, //32 Bit, 0.01kWh  // currently not supported

		ReactiveImport:   0x012C, //32 Bit, 0.00kWh
		ReactiveImportT1: 0x012E, //32 Bit, 0.00kWh
		ReactiveImportT2: 0x0130, //32 Bit, 0.00kWh
		//ReactiveImportT3: 0x0132, //32 Bit, 0.00kWh  // currently not supported
		//ReactiveImportT4: 0x0134, //32 Bit, 0.00kWh  // currently not supported

		ReactiveExport:   0x0136, //32 Bit, 0.00kWh
		ReactiveExportT1: 0x0138, //32 Bit, 0.00kWh
		ReactiveExportT2: 0x013A, //32 Bit, 0.00kWh
		//ReactiveExportT3: 0x013C, //32 Bit, 0.00kWh  // currently not supported
		//ReactiveExportT4: 0x013E, //32 Bit, 0.00kWh  // currently not supported

		ReactiveSum:   0x0140, //32 Bit, 0.01kvarh
		ReactiveSumT1: 0x0142, //32 Bit, 0.01kvarh
		ReactiveSumT2: 0x0144, //32 Bit, 0.01kvarh
		//ReactiveSumT3: 0x0146, //32 Bit, 0.01kvarh  // currently not supported
		//ReactiveSumT4: 0x0148, //32 Bit, 0.01kvarh  // currently not supported
	}

	return &ORNO1P525Producer{Opcodes: ops}
}

// Initialize implements Producer interface
func (p *ORNO1P525Producer) Initialize(client modbus.Client) {

}

// Description implements Producer interface
func (p *ORNO1P525Producer) Description() string {
	return "ORNO WE-525 & WE-526"
}

// snip creates modbus operation
func (p *ORNO1P525Producer) snip(iec Measurement, readlen uint16) Operation {
	return Operation{
		FuncCode: ReadInputReg,
		OpCode:   p.Opcode(iec), // adjust according to docs
		ReadLen:  readlen,
		IEC61850: iec,
	}
}

// snip16 creates modbus operation for single register
func (p *ORNO1P525Producer) snip16(iec Measurement, scaler ...float64) Operation {
	snip := p.snip(iec, 1)

	snip.Transform = RTUInt16ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
	}

	return snip
}

// snip32 creates modbus operation for double register
func (p *ORNO1P525Producer) snip32(iec Measurement, scaler ...float64) Operation {
	snip := p.snip(iec, 2)

	snip.Transform = RTUInt32ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
	}

	return snip
}

func (p *ORNO1P525Producer) Probe() Operation {
	return p.snip32(Voltage, 100)
}

// Produce implements Producer interface
func (p *ORNO1P525Producer) Produce() (res []Operation) {

	for _, op := range []Measurement{
		Frequency,
	} {
		res = append(res, p.snip16(op, 10))
	}

	for _, op := range []Measurement{
		Current, Voltage,
	} {
		res = append(res, p.snip32(op, 1000))
	}

	for _, op := range []Measurement{
		Power, ReactivePower, ApparentPower,
	} {
		res = append(res, p.snip32(op, 1))
	}

	for _, op := range []Measurement{
		Cosphi,
	} {
		res = append(res, p.snip16(op, 1000))
	}

	for _, op := range []Measurement{
		Sum, Import, Export, ReactiveSum, ReactiveImport, ReactiveExport,
		SumT1, ImportT1, ExportT1, ReactiveSumT1, ReactiveImportT1, ReactiveExportT1,
		SumT2, ImportT2, ExportT2, ReactiveSumT2, ReactiveImportT2, ReactiveExportT2,
		//SumT3, ImportT3, ExportT3, ReactiveSumT3, ReactiveImportT3, ReactiveExportT3,  // currently not supported
		//SumT4, ImportT4, ExportT4, ReactiveSumT4, ReactiveImportT4, ReactiveExportT4,  // currently not supported
	} {
		res = append(res, p.snip32(op, 100))
	}
	return res
}
