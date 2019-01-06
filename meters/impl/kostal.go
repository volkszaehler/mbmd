package impl

import . "github.com/gonium/gosdm630/meters"

func init() {
	Register(NewKostalProducer)
}

const (
	METERTYPE_KOSTAL = "KOSTAL"
)

type KostalProducer struct {
	SunSpecCore
}

func NewKostalProducer() Producer {
	/***
	 * Opcodes for SunSpec-compatible Inverters from Kostal
	 * https://www.kostal-solar-electric.com/de-de/download/-/media/document%20library%20folder%20-%20kse/2018/08/30/08/53/ba_kostal_interface_modbus-tcp_sunspec.pdf
	 */
	ops := Opcodes{
		DCPower: 100, // + scaler
		/*
			HomeOwnBatteryPower: 106, // + scaler
			HomeOwnGridPower:    108, // + scaler
			HomeOwnPVPower:      116, // + scaler

			TotalHomePVConsumption:      110, // + scaler
			TotalHomeGridConsumption:    112, // + scaler
			TotalHomeBatteryConsumption: 114, // + scaler
			TotalHomeConsumption:        118, // + scaler

			EVUPowerLimit:            122, // + scaler
			TotalHomeConsumptionRate: 124, // + scaler
		*/
		Cosphi:    150, // + scaler
		Frequency: 152, // + scaler

		CurrentL1: 154, // + scaler
		PowerL1:   156, // + scaler
		VoltageL1: 158, // + scaler

		CurrentL2: 160, // + scaler
		PowerL2:   162, // + scaler
		VoltageL2: 164, // + scaler

		CurrentL3: 166, // + scaler
		PowerL3:   168, // + scaler
		VoltageL3: 170, // + scaler

		Power:         172, // + scaler
		ReactivePower: 174, // + scaler
		ApparentPower: 176, // + scaler
		/*
			BatteryVoltage: 216, // + scaler

			CurrentDC1: 258, // + scaler
			PowerDC1:   260, // + scaler
			VoltageDC1: 266, // + scaler

			CurrentDC2: 268, // + scaler
			PowerDC2:   270, // + scaler
			VoltageDC2: 276, // + scaler

			CurrentDC3: 278, // + scaler
			PowerDC3:   280, // + scaler
			VoltageDC3: 286, // + scaler

			TotalYield:   320, // + scaler
			DailyYield:   322, // + scaler
			YearlyYield:  324, // + scaler
			MonthlyYield: 326, // + scaler
		*/
	}
	return &KostalProducer{
		SunSpecCore{ops},
	}
}

func (p *KostalProducer) Type() string {
	return METERTYPE_KOSTAL
}

func (p *KostalProducer) Description() string {
	return "Kostal SunSpec-compatible inverters (e.g. Pico IQ) (experimental)"
}

func (p *KostalProducer) Probe() Operation {
	return p.snip16uint(VoltageL1, 10)
}

func (p *KostalProducer) Produce() (res []Operation) {
	res = []Operation{
		// int16
		p.scaleSnip16(p.mkSplitInt16, CurrentL1),
		p.scaleSnip16(p.mkSplitInt16, PowerL1),
		p.scaleSnip16(p.mkSplitInt16, VoltageL1),

		p.scaleSnip16(p.mkSplitInt16, CurrentL2),
		p.scaleSnip16(p.mkSplitInt16, PowerL2),
		p.scaleSnip16(p.mkSplitInt16, VoltageL2),

		p.scaleSnip16(p.mkSplitInt16, CurrentL3),
		p.scaleSnip16(p.mkSplitInt16, PowerL3),
		p.scaleSnip16(p.mkSplitInt16, VoltageL3),

		p.scaleSnip16(p.mkSplitInt16, Power),
		p.scaleSnip16(p.mkSplitInt16, DCPower),

		p.scaleSnip16(p.mkSplitInt16, Cosphi),
		p.scaleSnip16(p.mkSplitInt16, Frequency),
	}

	return res
}
