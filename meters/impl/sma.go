package impl

import . "github.com/gonium/gosdm630/meters"

func init() {
	Register(NewSMAProducer)
}

const (
	METERTYPE_SMA = "SMA"
)

type SMAProducer struct {
	SunSpecCore
}

func NewSMAProducer() Producer {
	/***
	 * Opcodes for SMA SunSpec-compatible Inverters
	 * https://www.sma.de/fileadmin/content/landingpages/pl/FAQ/SunSpec_Modbus-TI-en-15.pdf
	 */
	ops := Opcodes{
		Current:   188, // uint16
		CurrentL1: 189,
		CurrentL2: 190,
		CurrentL3: 191, // + scaler

		VoltageL1: 196, // uint16
		VoltageL2: 197,
		VoltageL3: 198, // + scaler

		Power: 200, // int16 + scaler
		// ApparentPower: 204, // int16 + scaler
		// ReactivePower: 206, // int16 + scaler
		Export: 210, // uint32 + scaler

		Cosphi:    208, // int16 + scaler
		Frequency: 202, // uint16 + scaler

		DCPower: 217, // int16 + scaler

		// DC block with global scale factors
		// DCCurrent1: 641,  // uint16
		// DCVoltage1: 642,  // uint16
		// DCPower1:   643, // uint16
		// DCCurrent2: 661,  // uint16
		// DCVoltage2: 662,  // uint16
		// DCPower2:   663, // uint16

		HeatSinkTemp: 219, // int16 + scaler
	}
	return &SMAProducer{
		SunSpecCore{ops},
	}
}

func (p *SMAProducer) Type() string {
	return METERTYPE_SMA
}

func (p *SMAProducer) Description() string {
	return "SMA SunSpec-compatible inverters (e.g. Sunny Boy or Tripower) (experimental)"
}

func (p *SMAProducer) Probe() Operation {
	return p.snip16uint(VoltageL1, 10)
}

func (p *SMAProducer) Produce() (res []Operation) {
	res = []Operation{
		// uint16
		p.scaleSnip16(p.mkSplitUint16, VoltageL1, VoltageL2, VoltageL3),
		p.scaleSnip16(p.mkSplitUint16, Frequency),

		// int16
		p.scaleSnip16(p.mkSplitInt16, Current, CurrentL1, CurrentL2, CurrentL3),
		p.scaleSnip16(p.mkSplitInt16, Cosphi),
		p.scaleSnip16(p.mkSplitInt16, Power),
		p.scaleSnip16(p.mkSplitInt16, DCPower),
		p.scaleSnip16(p.mkSplitInt16, HeatSinkTemp),

		// uint32
		p.scaleSnip32(p.mkSplitUint32, Export),
	}

	return res
}
