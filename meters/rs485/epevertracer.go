package rs485

import (
	. "github.com/volkszaehler/mbmd/meters"
)

func init() {
	Register("Epever", NewEpeverProducer)
}

const (
	METERTYPE_EPEVER = "EPEVER"
)

type EpeverProducer struct {
	Opcodes
}

func NewEpeverProducer() Producer {
	/**
	 * Opcodes as defined by EPSolar
	 * from: http://www.solar-elektro.cz/data/dokumenty/1733_modbus_protocol.pdf
	 */
	ops := Opcodes{

		DCVoltage:            0x3100, // Solar charge controller--PV array voltage
		DCCurrent:            0x3101, // Solar charge controller--PV array current
		DCPower:              0x3102, // Solar charge controller--PV array power
		BatteryVoltage:       0x3104, // Batterie voltage
		BatteryCurrent:       0x3105, // Batterie charging current
		BatteryPower:         0x3106, // Batterie charging power
		LoadVoltage:          0x310C, // Load voltage
		LoadCurrent:          0x310D, // Load  current
		LoadPower:            0x310E, // Load  power
		BatteryTemperature:   0x3110, // Battery Temperature
		InsideTemperature:    0x3111, // Temperature inside case
		MaxPVVoltage:         0x3300, // maximum input voltage (PV) today
		MinPVVoltage:         0x3301, // minimum input voltage (PV) today
		MaxBattVoltage:       0x3302, // maximum battery voltage (PV) today
		MinBattVoltage:       0x3303, // minimum battery voltage (PV) today
		ConsumedEnergyToday:  0x3304, // Consumed Energy Today
		ConsumedEnergyMonth:  0x3306, // Consumed Energy Month
		ConsumedEnergyYear:   0x3308, // Consumed Energy Year
		ConsumedEnergyTotal:  0x330A, // Consumed Energy Totaly
		GeneratedEnergyToday: 0x330C, // Generated Energy Today
		GeneratedEnergyMonth: 0x330E, // Generated Energy Month
		GeneratedEnergyYear:  0x3310, // Generated Energy Year
		GeneratedEnergyTotal: 0x3312, // Generated Energy Totaly
	}
	return &EpeverProducer{Opcodes: ops}
}

func (p *EpeverProducer) Type() string {
	return METERTYPE_EPEVER
}

func (p *EpeverProducer) Description() string {
	return "Epever Tracer"
}

// snip creates modbus operation
func (p *EpeverProducer) snip(iec Measurement, readlen uint16) Operation {
	return Operation{
		FuncCode: ReadInputReg,
		OpCode:   p.Opcode(iec), // adjust according to docs
		ReadLen:  readlen,
		IEC61850: iec,
	}
}

// snip16 creates modbus operation for single register
func (p *EpeverProducer) snip16(iec Measurement, scaler ...float64) Operation {
	snip := p.snip(iec, 2)

	snip.Transform = RTUUint16ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
	}

	return snip
}

// snip32 creates modbus operation for double register
func (p *EpeverProducer) snip32(iec Measurement, scaler ...float64) Operation {
	snip := p.snip(iec, 4)

	snip.Transform = RTUInt32ToFloat64Swapped // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
	}

	return snip
}

func (p *EpeverProducer) Probe() Operation {
	return p.snip16(BatteryVoltage, 100)
}

// Produce implements Producer interface
func (p *EpeverProducer) Produce() (res []Operation) {

	for _, op := range []Measurement{
		DCVoltage, DCCurrent, BatteryVoltage, BatteryCurrent, LoadVoltage, LoadCurrent, BatteryTemperature, InsideTemperature, MaxPVVoltage, MinPVVoltage, MaxBattVoltage, MinBattVoltage,
	} {
		res = append(res, p.snip16(op, 100))
	}

	for _, op := range []Measurement{
		DCPower, BatteryPower, LoadPower, ConsumedEnergyToday, ConsumedEnergyMonth, ConsumedEnergyYear, ConsumedEnergyTotal, GeneratedEnergyToday, GeneratedEnergyMonth, GeneratedEnergyYear, GeneratedEnergyTotal,
	} {
		res = append(res, p.snip32(op, 100))
	}

	return res
}
