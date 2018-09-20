package meters

const (
	METERTYPE_JANITZA = "JANITZA"

	/***
	 * Opcodes for Janitza B23.
	 * See https://www.janitza.de/betriebsanleitungen.html?file=files/download/manuals/current/B-Series/MID-Energy-Meters-Product-Manual.pdf
	 */
	OpCodeJanitzaL1Voltage   = 0x4A38
	OpCodeJanitzaL2Voltage   = 0x4A3A
	OpCodeJanitzaL3Voltage   = 0x4A3C
	OpCodeJanitzaL1Current   = 0x4A44
	OpCodeJanitzaL2Current   = 0x4A46
	OpCodeJanitzaL3Current   = 0x4A48
	OpCodeJanitzaL1Power     = 0x4A4C
	OpCodeJanitzaL2Power     = 0x4A4E
	OpCodeJanitzaL3Power     = 0x4A50
	OpCodeJanitzaL1Import    = 0x4A76
	OpCodeJanitzaL2Import    = 0x4A78
	OpCodeJanitzaL3Import    = 0x4A7A
	OpCodeJanitzaTotalImport = 0x4A7C
	OpCodeJanitzaL1Export    = 0x4A7E
	OpCodeJanitzaL2Export    = 0x4A80
	OpCodeJanitzaL3Export    = 0x4A82
	OpCodeJanitzaTotalExport = 0x4A84
	OpCodeJanitzaL1Cosphi    = 0x4A64
	OpCodeJanitzaL2Cosphi    = 0x4A66
	OpCodeJanitzaL3Cosphi    = 0x4A68
)

type JanitzaProducer struct {
}

func NewJanitzaProducer() *JanitzaProducer {
	return &JanitzaProducer{}
}

func (p *JanitzaProducer) GetMeterType() string {
	return METERTYPE_JANITZA
}

func (p *JanitzaProducer) op(opcode uint16, iec string) Operation {
	op := Operation{
		FuncCode:  ReadHoldingReg,
		OpCode:    opcode,
		ReadLen:   2,
		IEC61850:  iec,
		Transform: RTU32ToFloat64,
	}
	return op
}

func (p *JanitzaProducer) Probe() Operation {
	return p.op(OpCodeJanitzaL1Voltage, "VolLocPhsA")
}

func (p *JanitzaProducer) Produce() (res []Operation) {
	res = append(res, p.op(OpCodeJanitzaL1Voltage, "VolLocPhsA"))
	res = append(res, p.op(OpCodeJanitzaL2Voltage, "VolLocPhsB"))
	res = append(res, p.op(OpCodeJanitzaL3Voltage, "VolLocPhsC"))

	res = append(res, p.op(OpCodeJanitzaL1Current, "AmpLocPhsA"))
	res = append(res, p.op(OpCodeJanitzaL2Current, "AmpLocPhsB"))
	res = append(res, p.op(OpCodeJanitzaL3Current, "AmpLocPhsC"))

	res = append(res, p.op(OpCodeJanitzaL1Power, "WLocPhsA"))
	res = append(res, p.op(OpCodeJanitzaL2Power, "WLocPhsB"))
	res = append(res, p.op(OpCodeJanitzaL3Power, "WLocPhsC"))

	res = append(res, p.op(OpCodeJanitzaL1Cosphi, "AngLocPhsA"))
	res = append(res, p.op(OpCodeJanitzaL2Cosphi, "AngLocPhsB"))
	res = append(res, p.op(OpCodeJanitzaL3Cosphi, "AngLocPhsC"))

	res = append(res, p.op(OpCodeJanitzaL1Import, "TotkWhImportPhsA"))
	res = append(res, p.op(OpCodeJanitzaL2Import, "TotkWhImportPhsB"))
	res = append(res, p.op(OpCodeJanitzaL3Import, "TotkWhImportPhsC"))
	res = append(res, p.op(OpCodeJanitzaTotalImport, "TotkWhImport"))

	res = append(res, p.op(OpCodeJanitzaL1Export, "TotkWhExportPhsA"))
	res = append(res, p.op(OpCodeJanitzaL2Export, "TotkWhExportPhsB"))
	res = append(res, p.op(OpCodeJanitzaL3Export, "TotkWhExportPhsC"))
	res = append(res, p.op(OpCodeJanitzaTotalExport, "TotkWhExport"))
	return res
}
