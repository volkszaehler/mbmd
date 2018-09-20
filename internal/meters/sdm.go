package sdm630

import "math"

const (
	METERTYPE_SDM = "SDM"

	/***
	 * Opcodes as defined by Eastron.
	 * See http://bg-etech.de/download/manual/SDM630Register.pdf
	 * Please note that this is the superset of all SDM devices - some
	 * opcodes might not work on some devicep.
	 */
	OpCodeSDML1Voltage   = 0x0000
	OpCodeSDML2Voltage   = 0x0002
	OpCodeSDML3Voltage   = 0x0004
	OpCodeSDML1Current   = 0x0006
	OpCodeSDML2Current   = 0x0008
	OpCodeSDML3Current   = 0x000A
	OpCodeSDML1Power     = 0x000C
	OpCodeSDML2Power     = 0x000E
	OpCodeSDML3Power     = 0x0010
	OpCodeSDML1Import    = 0x015a
	OpCodeSDML2Import    = 0x015c
	OpCodeSDML3Import    = 0x015e
	OpCodeSDMTotalImport = 0x0048
	OpCodeSDML1Export    = 0x0160
	OpCodeSDML2Export    = 0x0162
	OpCodeSDML3Export    = 0x0164
	OpCodeSDMTotalExport = 0x004a
	OpCodeSDML1Cosphi    = 0x001e
	OpCodeSDML2Cosphi    = 0x0020
	OpCodeSDML3Cosphi    = 0x0022
	//OpCodeL1THDCurrent         = 0x00F0
	//OpCodeL2THDCurrent         = 0x00F2
	//OpCodeL3THDCurrent         = 0x00F4
	//OpCodeAvgTHDCurrent        = 0x00Fa
	OpCodeSDML1THDVoltageNeutral  = 0x00ea
	OpCodeSDML2THDVoltageNeutral  = 0x00ec
	OpCodeSDML3THDVoltageNeutral  = 0x00ee
	OpCodeSDMAvgTHDVoltageNeutral = 0x00F8
	OpCodeSDMFrequency            = 0x0046
)

type SDMProducer struct {
}

func NewSDMProducer() *SDMProducer {
	return &SDMProducer{}
}

func (p *SDMProducer) GetMeterType() string {
	return METERTYPE_SDM
}

func (p *SDMProducer) snip(devid uint8, opcode uint16, iec string) QuerySnip {
	snip := QuerySnip{
		DeviceId:  devid,
		FuncCode:  ReadInputReg,
		OpCode:    opcode,
		ReadLen:   2,
		Value:     math.NaN(),
		IEC61850:  iec,
		Transform: RTU32ToFloat64,
	}
	return snip
}

func (p *SDMProducer) Probe(devid uint8) QuerySnip {
	return p.snip(devid, OpCodeSDML1Voltage, "VolLocPhsA")
}

func (p *SDMProducer) Produce(devid uint8) (res []QuerySnip) {
	res = append(res, p.snip(devid, OpCodeSDML1Voltage, "VolLocPhsA"))
	res = append(res, p.snip(devid, OpCodeSDML2Voltage, "VolLocPhsB"))
	res = append(res, p.snip(devid, OpCodeSDML3Voltage, "VolLocPhsC"))
	res = append(res, p.snip(devid, OpCodeSDML1Current, "AmpLocPhsA"))
	res = append(res, p.snip(devid, OpCodeSDML2Current, "AmpLocPhsB"))
	res = append(res, p.snip(devid, OpCodeSDML3Current, "AmpLocPhsC"))

	res = append(res, p.snip(devid, OpCodeSDML1Power, "WLocPhsA"))
	res = append(res, p.snip(devid, OpCodeSDML2Power, "WLocPhsB"))
	res = append(res, p.snip(devid, OpCodeSDML3Power, "WLocPhsC"))

	res = append(res, p.snip(devid, OpCodeSDML1Cosphi, "AngLocPhsA"))
	res = append(res, p.snip(devid, OpCodeSDML2Cosphi, "AngLocPhsB"))
	res = append(res, p.snip(devid, OpCodeSDML3Cosphi, "AngLocPhsC"))

	res = append(res, p.snip(devid, OpCodeSDML1Import, "TotkWhImportPhsA"))
	res = append(res, p.snip(devid, OpCodeSDML2Import, "TotkWhImportPhsB"))
	res = append(res, p.snip(devid, OpCodeSDML3Import, "TotkWhImportPhsC"))
	res = append(res, p.snip(devid, OpCodeSDMTotalImport, "TotkWhImport"))

	res = append(res, p.snip(devid, OpCodeSDML1Export, "TotkWhExportPhsA"))
	res = append(res, p.snip(devid, OpCodeSDML2Export, "TotkWhExportPhsB"))
	res = append(res, p.snip(devid, OpCodeSDML3Export, "TotkWhExportPhsC"))
	res = append(res, p.snip(devid, OpCodeSDMTotalExport, "TotkWhExport"))

	res = append(res, p.snip(devid, OpCodeSDML1THDVoltageNeutral, "ThdVolPhsA"))
	res = append(res, p.snip(devid, OpCodeSDML2THDVoltageNeutral, "ThdVolPhsB"))
	res = append(res, p.snip(devid, OpCodeSDML3THDVoltageNeutral, "ThdVolPhsC"))
	res = append(res, p.snip(devid, OpCodeSDMAvgTHDVoltageNeutral, "ThdVol"))

	res = append(res, p.snip(devid, OpCodeSDMFrequency, "Freq"))

	return res
}
