package meters

const (
	METERTYPE_SUN = "SUN"

	// MODBUS protocol address (base 0)
	Base = 40000 - 1

	/***
	 * Opcodes for SunSpec- compatible Inverters like SolarEdge
	 * https://www.solaredge.com/sites/default/files/sunspec-implementation-technical-note.pdf
	 */
	OpCodeSunSpecInverterL1Current = 73
	OpCodeSunSpecInverterL2Current = 74
	OpCodeSunSpecInverterL3Current = 75

	OpCodeSunSpecInverterL1Voltage = 80
	OpCodeSunSpecInverterL2Voltage = 81
	OpCodeSunSpecInverterL3Voltage = 82

	OpCodeSunSpecInverterFrequency = 86

	// // power
	// q(base+84, 1)
	// q(base+85, 1) // SF

	// // PF
	// q(base+92, 1)
	// q(base+93, 1) // SF

	// q(base+1, 2)
	// q(base+3, 1)
	// q(base+5, 16)
	// q(base+21, 16)
	// q(base+45, 8)
	// q(base+53, 16)

	// q(base+69, 1)

	// q(base+70, 1)
	// q(base+71, 1)

	// // strom
	// q(base+72, 1) // total

	// // power
	// q(base+84, 1)
	// q(base+85, 1) // SF

	// // freq
	// q(base+86, 1)
	// q(base+87, 1) // SF

	// // apparent
	// q(base+88, 1)
	// q(base+89, 1) // SF

	// // reactive
	// q(base+90, 1)
	// q(base+91, 1) // SF

	// // PF
	// q(base+92, 1)
	// q(base+93, 1) // SF

	// // energy
	// q(base+94, 2)

	// q(base+97, 1)  // DC current + SF
	// q(base+99, 1)  // DC voltage + SF
	// q(base+101, 1) // DC power + SF

	// // hreatsink temp
	// q(base+104, 1)
	// q(base+105, 1) // SF

	// q(base+108, 1) // status
	// q(base+109, 1) // vendor status

	// // meter
	// q(base+121, 1)
)

type SUNProducer struct {
}

func NewSUNProducer() *SUNProducer {
	return &SUNProducer{}
}

func (p *SUNProducer) GetMeterType() string {
	return METERTYPE_SUN
}

// func (p *SUNProducer) snip(opcode uint16, iec string) Operation {
// 	return Operation{
// 		FuncCode:  ReadHoldingReg,
// 		OpCode:    Base + opcode,
// 		ReadLen:   1,
// 		IEC61850:  iec,
// 		Transform: MakeRTU16ScaledIntToFloat64(10),
// 	}
// }

func (p *SUNProducer) snip(opcode uint16, iec string, readlen uint16) Operation {
	return Operation{
		FuncCode: ReadHoldingReg,
		OpCode:   Base + opcode,
		ReadLen:  readlen,
		IEC61850: iec,
	}
}

// op16 creates modbus operation for single register
func (p *SUNProducer) snip16(opcode uint16, iec string, scaler ...float64) Operation {
	snip := p.snip(opcode, iec, 1)

	snip.Transform = RTU16ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeRTU16ScaledIntToFloat64(scaler[0])
	}

	return snip
}

func (p *SUNProducer) Probe() Operation {
	return p.snip(OpCodeSunSpecInverterL1Voltage, "VolLocPhsA", 10)
}

func (p *SUNProducer) Produce() (res []Operation) {
	res = append(res, p.snip16(OpCodeSunSpecInverterL1Voltage, "VolLocPhsA", 10))
	res = append(res, p.snip16(OpCodeSunSpecInverterL2Voltage, "VolLocPhsB", 10))
	res = append(res, p.snip16(OpCodeSunSpecInverterL3Voltage, "VolLocPhsC", 10))

	res = append(res, p.snip16(OpCodeSunSpecInverterL1Current, "AmpLocPhsA", 100))
	res = append(res, p.snip16(OpCodeSunSpecInverterL2Current, "AmpLocPhsB", 100))
	res = append(res, p.snip16(OpCodeSunSpecInverterL3Current, "AmpLocPhsC", 100))

	res = append(res, p.snip16(OpCodeSunSpecInverterFrequency, "Freq", 100))

	return res
}
