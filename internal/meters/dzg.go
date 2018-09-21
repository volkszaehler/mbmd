package meters

const (
	METERTYPE_DZG = "DZG"

	/***
	 * Opcodes for DZG DVH4014.
	 * See "User Manual DVH4013", not public.
	 */
	OpCodeDZGTotalImportPower = 0x0000
	OpCodeDZGTotalExportPower = 0x0002
	OpCodeDZGL1Voltage        = 0x0004
	OpCodeDZGL2Voltage        = 0x0006
	OpCodeDZGL3Voltage        = 0x0008
	OpCodeDZGL1Current        = 0x000A
	OpCodeDZGL2Current        = 0x000C
	OpCodeDZGL3Current        = 0x000E
	OpCodeDZGL1Import         = 0x4020
	OpCodeDZGL2Import         = 0x4040
	OpCodeDZGL3Import         = 0x4060
	OpCodeDZGTotalImport      = 0x4000
	OpCodeDZGL1Export         = 0x4120
	OpCodeDZGL2Export         = 0x4140
	OpCodeDZGL3Export         = 0x4160
	OpCodeDZGTotalExport      = 0x4100
)

type DZGProducer struct {
}

func NewDZGProducer() *DZGProducer {
	return &DZGProducer{}
}

func (p *DZGProducer) GetMeterType() string {
	return METERTYPE_DZG
}

func (p *DZGProducer) snip(opcode uint16, iec string, scaler ...float64) Operation {
	transform := RTU32ToFloat64 // default conversion
	if len(scaler) > 0 {
		transform = MakeRTU32ScaledIntToFloat64(scaler[0])
	}

	return Operation{
		FuncCode:  ReadHoldingReg,
		OpCode:    opcode,
		ReadLen:   2,
		IEC61850:  iec,
		Transform: transform,
	}
}

func (p *DZGProducer) Probe() Operation {
	return p.snip(OpCodeDZGL1Voltage, "VolLocPhsA", 100)
}

func (p *DZGProducer) Produce() (res []Operation) {
	res = append(res, p.snip(OpCodeDZGL1Voltage, "VolLocPhsA", 100))
	res = append(res, p.snip(OpCodeDZGL2Voltage, "VolLocPhsB", 100))
	res = append(res, p.snip(OpCodeDZGL3Voltage, "VolLocPhsC", 100))

	res = append(res, p.snip(OpCodeDZGL1Current, "AmpLocPhsA", 1000))
	res = append(res, p.snip(OpCodeDZGL2Current, "AmpLocPhsB", 1000))
	res = append(res, p.snip(OpCodeDZGL3Current, "AmpLocPhsC", 1000))

	res = append(res, p.snip(OpCodeDZGL1Import, "TotkWhImportPhsA"))
	res = append(res, p.snip(OpCodeDZGL2Import, "TotkWhImportPhsB"))
	res = append(res, p.snip(OpCodeDZGL3Import, "TotkWhImportPhsC"))
	res = append(res, p.snip(OpCodeDZGTotalImport, "TotkWhImport", 1000))

	res = append(res, p.snip(OpCodeDZGL1Export, "TotkWhExportPhsA"))
	res = append(res, p.snip(OpCodeDZGL2Export, "TotkWhExportPhsB"))
	res = append(res, p.snip(OpCodeDZGL3Export, "TotkWhExportPhsC"))
	res = append(res, p.snip(OpCodeDZGTotalExport, "TotkWhExport", 1000))

	return res
}
