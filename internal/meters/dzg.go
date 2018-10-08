package meters

const (
	METERTYPE_DZG = "DZG"
)

type DZGProducer struct {
	MeasurementMapping
}

func NewDZGProducer() *DZGProducer {
	/**
	 * Opcodes for DZG DVH4014.
	 * https://www.dzg.de/fileadmin/dzg/content/downloads/produkte-zaehler/dvh4013/Communication-Protocol_DVH4013.pdf
	 */
	ops := Measurements{
		ActivePower:   0x0000, // 0x0 instant values and parameters
		ReactivePower: 0x0002,
		VoltageL1:     0x0004,
		VoltageL2:     0x0006,
		VoltageL3:     0x0008,
		CurrentL1:     0x000A,
		CurrentL2:     0x000C,
		CurrentL3:     0x000E,
		Cosphi:        0x0010, // DVH4013
		Frequency:     0x0012, // DVH4013
		// Import:        0x0014, // DVH4013
		// Export:        0x0016, // DVH4013
		Import:   0x4000, // 0x4 energy
		ImportL1: 0x4020,
		ImportL2: 0x4040,
		ImportL3: 0x4060,
		Export:   0x4100,
		ExportL1: 0x4120,
		ExportL2: 0x4140,
		ExportL3: 0x4160,
		// 0x8 max demand
	}
	return &DZGProducer{
		MeasurementMapping{ops},
	}
}

func (p *DZGProducer) GetMeterType() string {
	return METERTYPE_DZG
}

func (p *DZGProducer) snip(iec Measurement, scaler ...float64) Operation {
	transform := RTUUint32ToFloat64 // default conversion
	if len(scaler) > 0 {
		transform = MakeRTUScaledUint32ToFloat64(scaler[0])
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

func (p *DZGProducer) Probe() Operation {
	return p.snip(VoltageL1, 100)
}

func (p *DZGProducer) Produce() (res []Operation) {
	for _, op := range []Measurement{
		VoltageL1, VoltageL2, VoltageL3,
	} {
		res = append(res, p.snip(op, 100))
	}

	for _, op := range []Measurement{
		CurrentL1, CurrentL2, CurrentL3,
		Import, Export, Cosphi,
	} {
		res = append(res, p.snip(op, 1000))
	}

	for _, op := range []Measurement{
		ImportL1, ImportL2, ImportL3,
		ExportL1, ExportL2, ExportL3,
	} {
		res = append(res, p.snip(op))
	}

	return res
}
