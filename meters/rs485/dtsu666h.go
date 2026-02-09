package rs485

import (
	. "github.com/volkszaehler/mbmd/meters"
)

func init() {
	Register("DTSU666H", NewDTSU666HProducer)
}

type DTSU666HProducer struct {
	Opcodes
}

func NewDTSU666HProducer() Producer {
	/***
	 * Opcodes as defined by Chint DTSU666-H (Huawei firmware).
	 * Based on https://github.com/salakrzy/DTSU666_CHINT_to_HUAWEI_translator
	 * Huawei firmware uses different register addresses starting at 0x0836 (2102).
	 */
	ops := Opcodes{
		CurrentL1:       2102, // Phase current A (IEEE754 float32, A, scale: 1000) - Pos 0
		CurrentL2:       2104, // Phase current B (IEEE754 float32, A, scale: 1000) - Pos 1
		CurrentL3:       2106, // Phase current C (IEEE754 float32, A, scale: 1000) - Pos 2
		VoltageL1:       2110, // Phase voltage A (IEEE754 float32, V, scale: 10) - Pos 4
		VoltageL2:       2112, // Phase voltage B (IEEE754 float32, V, scale: 10) - Pos 5
		VoltageL3:       2114, // Phase voltage C (IEEE754 float32, V, scale: 10) - Pos 6
		Frequency:       2124, // Grid frequency (IEEE754 float32, Hz, scale: 100) - Pos 11
		Power:           2126, // Active power total (IEEE754 float32, kW, scale: 10) - Pos 12
		PowerL1:         2128, // Active power A (IEEE754 float32, kW, scale: 10) - Pos 13
		PowerL2:         2130, // Active power B (IEEE754 float32, kW, scale: 10) - Pos 14
		PowerL3:         2132, // Active power C (IEEE754 float32, kW, scale: 10) - Pos 15
		ReactivePower:   2134, // Reactive power total (IEEE754 float32, kVar, scale: 10) - Pos 16
		ReactivePowerL1: 2136, // Reactive power A (IEEE754 float32, kVar, scale: 10) - Pos 17
		ReactivePowerL2: 2138, // Reactive power B (IEEE754 float32, kVar, scale: 10) - Pos 18
		ReactivePowerL3: 2140, // Reactive power C (IEEE754 float32, kVar, scale: 10) - Pos 19
		Cosphi:          2144, // Power factor total (IEEE754 float32, scale: 10) - Pos 24
		CosphiL1:        2146, // Power factor A (IEEE754 float32, scale: 10) - Pos 25
		CosphiL2:        2148, // Power factor B (IEEE754 float32, scale: 10) - Pos 26
		CosphiL3:        2150, // Power factor C (IEEE754 float32, scale: 10) - Pos 27
		Import:          2166, // Total positive active energy (IEEE754 float32, kWh, scale: 1) - Pos 32
		Export:          2180, // Total negative active energy (IEEE754 float32, kWh, scale: 1) - Pos 36
	}

	return &DTSU666HProducer{Opcodes: ops}
}

func (p *DTSU666HProducer) Description() string {
	return "Chint DTSU666-H (Huawei)"
}

func (p *DTSU666HProducer) snip(iec Measurement, scaler ...float64) Operation {
	operation := Operation{
		FuncCode:  ReadHoldingReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   2,
		IEC61850:  iec,
		Transform: RTUIeee754ToFloat64,
	}

	if len(scaler) > 0 {
		operation.Transform = MakeScaledTransform(operation.Transform, scaler[0])
	}

	return operation
}

func (p *DTSU666HProducer) Probe() Operation {
	return p.snip(VoltageL1, 10)
}

func (p *DTSU666HProducer) Produce() (res []Operation) {
	for _, op := range []Measurement{
		VoltageL1, VoltageL2, VoltageL3,
		Power, PowerL1, PowerL2, PowerL3,
		ReactivePower, ReactivePowerL1, ReactivePowerL2, ReactivePowerL3,
		Cosphi, CosphiL1, CosphiL2, CosphiL3,
	} {
		res = append(res, p.snip(op, 10))
	}

	for _, op := range []Measurement{
		CurrentL1, CurrentL2, CurrentL3,
	} {
		res = append(res, p.snip(op, 1000))
	}

	for _, op := range []Measurement{
		Frequency,
	} {
		res = append(res, p.snip(op, 100))
	}

	for _, op := range []Measurement{
		Import, Export,
	} {
		res = append(res, p.snip(op, 1))
	}

	return res
}
