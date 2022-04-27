package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register("DZG", NewDZGProducer)
}

type DZGProducer struct {
	Opcodes
}

func NewDZGProducer() Producer {
	/**
	 * Opcodes for DZG DVH4014.
	 * https://www.dzg.de/fileadmin/dzg/content/downloads/produkte-zaehler/dvh4013/Communication-Protocol_DVH4013.pdf
	 */
	ops := Opcodes{
		ImportPower: 0x0000, // 0x0000 - parameters
		ExportPower: 0x0002,
		VoltageL1:   0x0004,
		VoltageL2:   0x0006,
		VoltageL3:   0x0008,
		CurrentL1:   0x000A,
		CurrentL2:   0x000C,
		CurrentL3:   0x000E,
		Cosphi:      0x0010,
		Frequency:   0x0012,
		Import:      0x4000, // 0x4000 - energy
		ImportL1:    0x4020,
		ImportL2:    0x4040,
		ImportL3:    0x4060,
		Export:      0x4100, // 0x.1.. - reverse
		ExportL1:    0x4120,
		ExportL2:    0x4140,
		ExportL3:    0x4160,
		// ImportPower:   0x8000, // 0x8000 - demand(power)
		// ImportPowerL1: 0x8020,
		// ImportPowerL2: 0x8040,
		// ImportPowerL3: 0x8060,
		// ExportPower:   0x8100, // 0x.1.. - reverse
		// ExportPowerL1: 0x8120,
		// ExportPowerL2: 0x8140,
		// ExportPowerL3: 0x8160,
		// ImportPower:0x0014, // exception '2' (illegal data address)
		// ExportPower:0x0016, // exception '2' (illegal data address)
	}
	return &DZGProducer{Opcodes: ops}
}

// Description implements Producer interface
func (p *DZGProducer) Description() string {
	return "DZG Metering GmbH DVH4013 meters"
}

func (p *DZGProducer) snip(iec Measurement, scaler ...float64) Operation {
	transform := RTUUint32ToFloat64 // default conversion
	if len(scaler) > 0 {
		transform = MakeScaledTransform(transform, scaler[0])
	}

	snip := Operation{
		FuncCode:  ReadHoldingReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   2,
		IEC61850:  iec,
		Transform: transform,
	}
	return snip
}

// Probe implements Producer interface
func (p *DZGProducer) Probe() Operation {
	return p.snip(VoltageL1, 100)
}

// Produce implements Producer interface
func (p *DZGProducer) Produce() (res []Operation) {
	for _, op := range []Measurement{
		VoltageL1, VoltageL2, VoltageL3,
	} {
		res = append(res, p.snip(op, 100))
	}

	for _, op := range []Measurement{
		CurrentL1, CurrentL2, CurrentL3,
		Cosphi, Frequency,
	} {
		res = append(res, p.snip(op, 1000))
	}

	for _, op := range []Measurement{
		ImportPower, ExportPower,
	} {
		res = append(res, p.snip(op, 10)) // W
	}

	// these are "maximum" values, apparently retrieving "current" does not work
	// for _, op := range []Measurement{
	// 	ImportPower, ImportPowerL1, ImportPowerL2, ImportPowerL3,
	// 	ExportPower, ExportPowerL1, ExportPowerL2, ExportPowerL3,
	// } {
	// 	res = append(res, p.snip(op, 10000)) // factor 10000 = kW according to docs, but is W?
	// }

	for _, op := range []Measurement{
		Import, ImportL1, ImportL2, ImportL3,
		Export, ExportL1, ExportL2, ExportL3,
	} {
		res = append(res, p.snip(op, 1000)) // factor 1000 = kWh
	}

	return res
}
