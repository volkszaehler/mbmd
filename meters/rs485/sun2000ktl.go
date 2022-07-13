package rs485

import (
	. "github.com/volkszaehler/mbmd/meters"
)

func init() {
	Register("Sun2000", NewSun2000Producer)
}

const (
	METERTYPE_SUN2000 = "SUN2000"
)

type Sun2000Producer struct {
	Opcodes
}

func NewSun2000Producer() Producer {
	/**
	 * Opcodes as defined by Huawei Sun2000
	 * from: https://javierin.com/wp-content/uploads/sites/2/2021/09/Solar-Inverter-Modbus-Interface-Definitions.pdf
	 */
	ops := Opcodes{

		DCVoltageS1:          32016, // A maximum of 24 PV strings are supported. The number of PV
		DCCurrentS1:          32017, // strings read by the host is defined by the Number of PV strings signal.
		DCVoltageS2:          32018, // The voltage and current register addresses for each PV string are as follows:
		DCCurrentS2:          32019, // PV n voltage: 32014 + 2 n
		DCVoltageS3:          32020, // PV n current: 32015 + 2 n
		DCCurrentS3:          32021, // n indicates the PV string number, which ranges from 1 to 24.
		DCPower:              32064, // InputPower
		Power:                32080, // AC Power Output
		PeakPower:            32078, // Peak active power of current day
		VoltageL1:            32069, // Phase A voltage
		VoltageL2:            32070, // Phase B voltage
		VoltageL3:            32071, // Phase C voltage
		InsideTemperature:    32087, // Inside Temperature
		Efficiency:           32086, // Efficiancy Inverter
		GeneratedEnergyToday: 32114, // Daily energy yield
		GeneratedEnergyTotal: 32106, // Accumulated energy yield
		//VoltageL1: 37101, // Phase A voltage Smartmeter DTSU666-H
		//VoltageL2: 37103, // Phase B voltage Smartmeter DTSU666-H
		//VoltageL3: 37105, // Phase C voltage Smartmeter DTSU666-H
		CurrentL1:     37107, // Phase A current Smartmeter DTSU666-H
		CurrentL2:     37109, // Phase B current Smartmeter DTSU666-H
		CurrentL3:     37111, // Phase C current Smartmeter DTSU666-H
		ExportPower:   37113, // Active power Smartmeter DTSU666-H
		ExportPowerL1: 37132, // Phase A Active power Smartmeter DTSU666-H
		ExportPowerL2: 37134, // Phase B Active power Smartmeter DTSU666-H
		ExportPowerL3: 37136, // Phase C Active power Smartmeter DTSU666-H
		Frequency:     37118, // Grid frequency Smartmeter DTSU666-H

	}
	return &Sun2000Producer{Opcodes: ops}
}

func (p *Sun2000Producer) Type() string {
	return METERTYPE_SUN2000
}

func (p *Sun2000Producer) Description() string {
	return "Huawei Sun2000 Inverter"
}

// snip creates modbus operation
func (p *Sun2000Producer) snip(iec Measurement, readlen uint16) Operation {
	return Operation{
		FuncCode: ReadHoldingReg,
		OpCode:   p.Opcode(iec), // adjust according to docs
		ReadLen:  readlen,
		IEC61850: iec,
	}
}

// snip16 creates modbus operation for single register
func (p *Sun2000Producer) snip16(iec Measurement, scaler ...float64) Operation {
	snip := p.snip(iec, 2)

	snip.Transform = RTUUint16ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
	}

	return snip
}

// snip32 creates modbus operation for double register
func (p *Sun2000Producer) snip32(iec Measurement, scaler ...float64) Operation {
	snip := p.snip(iec, 4)

	snip.Transform = RTUInt32ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
	}

	return snip
}

func (p *Sun2000Producer) Probe() Operation {
	return p.snip32(VoltageL1, 10)
}

// Produce implements Producer interface
func (p *Sun2000Producer) Produce() (res []Operation) {

	for _, op := range []Measurement{
		DCCurrentS1, DCCurrentS2, DCCurrentS3, Frequency, Efficiency,
	} {
		res = append(res, p.snip16(op, 100))
	}

	for _, op := range []Measurement{
		DCVoltageS1, DCVoltageS2, DCVoltageS3, VoltageL1, VoltageL2, VoltageL3, InsideTemperature,
	} {
		res = append(res, p.snip16(op, 10))
	}

	for _, op := range []Measurement{
		GeneratedEnergyToday, GeneratedEnergyTotal,
	} {
		res = append(res, p.snip32(op, 100))
	}

	for _, op := range []Measurement{
		CurrentL1, CurrentL2, CurrentL3,
	} {
		res = append(res, p.snip32(op, 10))
	}

	for _, op := range []Measurement{
		Power, ExportPowerL1, ExportPowerL2, ExportPowerL3, ExportPower, DCPower, PeakPower,
	} {
		res = append(res, p.snip32(op, 1))
	}

	return res
}
