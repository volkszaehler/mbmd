package sdm630

import (
	"math"
)

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

func (p *DZGProducer) snip(devid uint8, opcode uint16, iec string, scaler ...float64) QuerySnip {
	transform := RTU32ToFloat64 // default conversion
	if len(scaler) > 0 {
		transform = MakeRTU32ScaledIntToFloat64(scaler[0])
	}

	snip := QuerySnip{
		DeviceId:  devid,
		FuncCode:  ReadHoldingReg,
		OpCode:    opcode,
		ReadLen:   2,
		Value:     math.NaN(),
		IEC61850:  iec,
		Transform: transform,
	}
	return snip
}

func (p *DZGProducer) Probe(devid uint8) QuerySnip {
	return p.snip(devid, OpCodeDZGL1Voltage, "VolLocPhsA", 100)
}

func (p *DZGProducer) Produce(devid uint8) (res []QuerySnip) {
	res = append(res, p.snip(devid, OpCodeDZGL1Voltage, "VolLocPhsA", 100))
	res = append(res, p.snip(devid, OpCodeDZGL2Voltage, "VolLocPhsB", 100))
	res = append(res, p.snip(devid, OpCodeDZGL3Voltage, "VolLocPhsC", 100))

	res = append(res, p.snip(devid, OpCodeDZGL1Current, "AmpLocPhsA", 1000))
	res = append(res, p.snip(devid, OpCodeDZGL2Current, "AmpLocPhsB", 1000))
	res = append(res, p.snip(devid, OpCodeDZGL3Current, "AmpLocPhsC", 1000))

	res = append(res, p.snip(devid, OpCodeDZGL1Import, "TotkWhImportPhsA"))
	res = append(res, p.snip(devid, OpCodeDZGL2Import, "TotkWhImportPhsB"))
	res = append(res, p.snip(devid, OpCodeDZGL3Import, "TotkWhImportPhsC"))
	res = append(res, p.snip(devid, OpCodeDZGTotalImport, "TotkWhImport", 1000))

	res = append(res, p.snip(devid, OpCodeDZGL1Export, "TotkWhExportPhsA"))
	res = append(res, p.snip(devid, OpCodeDZGL2Export, "TotkWhExportPhsB"))
	res = append(res, p.snip(devid, OpCodeDZGL3Export, "TotkWhExportPhsC"))
	res = append(res, p.snip(devid, OpCodeDZGTotalExport, "TotkWhExport", 1000))

	return res
}
