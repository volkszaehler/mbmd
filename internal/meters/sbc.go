package meters

const (
	METERTYPE_SBC = "SBC"

	/***
	 * Opcodes for Saia Burgess ALE3
	 * http://datenblatt.stark-elektronik.de/saia_burgess/DE_DS_Energymeter-ALE3-with-Modbus.pdf
	 */
	OpCodeSaiaTotalImport   = 28 - 1 // double, scaler 100
	OpCodeSaiaPartialImport = 30 - 1 // double, scaler 100
	OpCodeSaiaTotalExport   = 32 - 1 // double, scaler 100
	OpCodeSaiaPartialExport = 34 - 1 // double, scaler 100

	OpCodeSaiaL1Voltage       = 36 - 1
	OpCodeSaiaL1Current       = 37 - 1 // scaler 10
	OpCodeSaiaL1Power         = 38 - 1 // scaler 100
	OpCodeSaiaL1ReactivePower = 39 - 1 // scaler 100
	OpCodeSaiaL1Cosphi        = 40 - 1 // scaler 100

	OpCodeSaiaL2Voltage       = 41 - 1
	OpCodeSaiaL2Current       = 42 - 1 // scaler 10
	OpCodeSaiaL2Power         = 43 - 1 // scaler 100
	OpCodeSaiaL2ReactivePower = 44 - 1 // scaler 100
	OpCodeSaiaL2Cosphi        = 45 - 1 // scaler 100

	OpCodeSaiaL3Voltage       = 46 - 1
	OpCodeSaiaL3Current       = 47 - 1 // scaler 10
	OpCodeSaiaL3Power         = 48 - 1 // scaler 100
	OpCodeSaiaL3ReactivePower = 49 - 1 // scaler 100
	OpCodeSaiaL3Cosphi        = 50 - 1 // scaler 100

	OpCodeSaiaTotalPower         = 51 - 1 // scaler 100
	OpCodeSaiaTotalReactivePower = 52 - 1 // scaler 100
)

type SBCProducer struct {
}

func NewSBCProducer() *SBCProducer {
	return &SBCProducer{}
}

func (p *SBCProducer) GetMeterType() string {
	return METERTYPE_SBC
}

func (p *SBCProducer) snip(opcode uint16, iec string, readlen uint16) Operation {
	return Operation{
		FuncCode: ReadHoldingReg,
		OpCode:   opcode,
		ReadLen:  readlen,
		IEC61850: iec,
	}
}

// op16 creates modbus operation for single register
func (p *SBCProducer) snip16(opcode uint16, iec string, scaler ...float64) Operation {
	snip := p.snip(opcode, iec, 1)

	snip.Transform = RTU16ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeRTU16ScaledIntToFloat64(scaler[0])
	}

	return snip
}

// op32 creates modbus operation for double register
func (p *SBCProducer) snip32(opcode uint16, iec string, scaler ...float64) Operation {
	snip := p.snip(opcode, iec, 2)

	snip.Transform = RTU32ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeRTU32ScaledIntToFloat64(scaler[0])
	}

	return snip
}

func (p *SBCProducer) Probe() Operation {
	return p.snip16(OpCodeSaiaL1Voltage, "VolLocPhsA")
}

func (p *SBCProducer) Produce() (res []Operation) {
	res = append(res, p.snip16(OpCodeSaiaL1Voltage, "VolLocPhsA"))
	res = append(res, p.snip16(OpCodeSaiaL2Voltage, "VolLocPhsB"))
	res = append(res, p.snip16(OpCodeSaiaL3Voltage, "VolLocPhsC"))

	res = append(res, p.snip16(OpCodeSaiaL1Current, "AmpLocPhsA", 10))
	res = append(res, p.snip16(OpCodeSaiaL2Current, "AmpLocPhsB", 10))
	res = append(res, p.snip16(OpCodeSaiaL3Current, "AmpLocPhsC", 10))

	res = append(res, p.snip16(OpCodeSaiaL1Power, "WLocPhsA", 100))
	res = append(res, p.snip16(OpCodeSaiaL2Power, "WLocPhsB", 100))
	res = append(res, p.snip16(OpCodeSaiaL3Power, "WLocPhsC", 100))

	res = append(res, p.snip16(OpCodeSaiaL1Cosphi, "AngLocPhsA", 100))
	res = append(res, p.snip16(OpCodeSaiaL2Cosphi, "AngLocPhsB", 100))
	res = append(res, p.snip16(OpCodeSaiaL3Cosphi, "AngLocPhsC", 100))

	// res = append(res, p.snip16(OpCodeSaiaTotalPower, "WLoc", 100))

	res = append(res, p.snip32(OpCodeSaiaTotalImport, "TotkWhImport", 100))
	res = append(res, p.snip32(OpCodeSaiaTotalExport, "TotkWhExport", 100))

	return res
}
