package impl

import . "github.com/gonium/gosdm630/meters"

func init() {
	Register(NewSEProducer)
}

const (
	METERTYPE_SE = "SE"
)

type SEProducer struct {
	SunSpecCore
}

func NewSEProducer() Producer {
	/***
	 * Opcodes for SunSpec-compatible Inverters like SolarEdge
	 * https://www.solaredge.com/sites/default/files/sunspec-implementation-technical-note.pdf
	 */
	ops := Opcodes{
		Current:   72,
		CurrentL1: 73,
		CurrentL2: 74,
		CurrentL3: 75, // + scaler

		VoltageL1: 80,
		VoltageL2: 81,
		VoltageL3: 82, // + scaler

		Power:         84, // + scaler
		ApparentPower: 88, // + scaler
		ReactivePower: 90, // + scaler
		Export:        94, // + scaler

		Cosphi:    92, // + scaler
		Frequency: 86, // + scaler

		DCCurrent: 97,  // + scaler
		DCVoltage: 99,  // + scaler
		DCPower:   101, // + scaler

		HeatSinkTemp: 104, // + scaler
	}
	return &SEProducer{
		SunSpecCore{ops},
	}
}

func (p *SEProducer) Type() string {
	return METERTYPE_SE
}

func (p *SEProducer) Description() string {
	return "SolarEdge SunSpec-compatible inverters (e.g. SolarEdge 9k)"
}

func (p *SEProducer) Probe() Operation {
	return p.snip16uint(VoltageL1, 10)
}

func (p *SEProducer) Produce() (res []Operation) {
	res = []Operation{
		// uint16
		p.scaleSnip16(p.mkSplitUint16, VoltageL1, VoltageL2, VoltageL3),
		p.scaleSnip16(p.mkSplitUint16, Current, CurrentL1, CurrentL2, CurrentL3),

		p.scaleSnip16(p.mkSplitUint16, Frequency),
		p.scaleSnip16(p.mkSplitUint16, DCCurrent),
		p.scaleSnip16(p.mkSplitUint16, DCVoltage),

		// int16
		p.scaleSnip16(p.mkSplitInt16, Cosphi),
		p.scaleSnip16(p.mkSplitInt16, Power),
		p.scaleSnip16(p.mkSplitInt16, DCPower),
		p.scaleSnip16(p.mkSplitInt16, HeatSinkTemp),

		// uint32
		p.scaleSnip32(p.mkSplitUint32, Export),
	}

	return res
}
