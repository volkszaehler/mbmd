package meters

const (
	METERTYPE_JANITZA = "JANITZA"
)

type JanitzaProducer struct {
	MeasurementMapping
}

func NewJanitzaProducer() *JanitzaProducer {
	/**
	 * Opcodes for Janitza B23.
	 * See https://www.janitza.de/betriebsanleitungen.html?file=files/download/manuals/current/B-Series/MID-Energy-Meters-Product-Manual.pdf
	 */
	ops := Measurements{
		VoltageL1: 0x4A38,
		VoltageL2: 0x4A3A,
		VoltageL3: 0x4A3C,
		CurrentL1: 0x4A44,
		CurrentL2: 0x4A46,
		CurrentL3: 0x4A48,
		PowerL1:   0x4A4C,
		PowerL2:   0x4A4E,
		PowerL3:   0x4A50,
		ImportL1:  0x4A76,
		ImportL2:  0x4A78,
		ImportL3:  0x4A7A,
		Import:    0x4A7C,
		ExportL1:  0x4A7E,
		ExportL2:  0x4A80,
		ExportL3:  0x4A82,
		Export:    0x4A84,
		CosphiL1:  0x4A64,
		CosphiL2:  0x4A66,
		CosphiL3:  0x4A68,
	}
	return &JanitzaProducer{
		MeasurementMapping{ops},
	}
}

func (p *JanitzaProducer) GetMeterType() string {
	return METERTYPE_JANITZA
}

func (p *JanitzaProducer) snip(iec Measurement) Operation {
	snip := Operation{
		FuncCode:  ReadHoldingReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   2,
		IEC61850:  iec,
		Transform: RTU32ToFloat64,
	}
	return snip
}

func (p *JanitzaProducer) Probe() Operation {
	return p.snip(VoltageL1)
}

func (p *JanitzaProducer) Produce() (res []Operation) {
	for op := range p.ops {
		res = append(res, p.snip(op))
	}

	return res
}
