package rs485

import (
	"math"

	. "github.com/volkszaehler/mbmd/meters"
)

func init() {
	Register("SEMTR", NewSEMTRProducer)
}

type SEMTRProducer struct {
	Opcodes
}

func NewSEMTRProducer() Producer {
	/**
	 * Opcodes as defined by SolarEdge SE-MTR-3Y
	 * reverse engineered from: https://github.com/nmakel/solaredge_meterproxy/blob/master/semp-rtu.py
	 * Uses the modbus RTU protocol over RS485.
	 *
	 * These are only necessary, if you'd like to connect the power meter
	 * directly via rs485 (e.g. if you have no inverter from solaredge).
	 * Otherwise, the values can be accessed readily through the network
	 * connection of the inverter via the sunspec protocol.
	 */
	ops := Opcodes{
		Sum:      0x03E8, // total active energy
		Import:   0x03EA, // imported active energy
		SumT1:    0x03EC, // total active energy non-reset
		ImportT1: 0x03EE, // imported active energy non-reset
		Power:    0x03F0, // total power
		PowerL1:  0x03F2, // power l1
		PowerL2:  0x03F4, // power l2
		PowerL3:  0x03F6, // power l3
		//		VoltageLN:		0x03F8,	// l-n voltage
		VoltageL1: 0x03FA, // l1-n voltage
		VoltageL2: 0x03FC, // l2-n voltage
		VoltageL3: 0x03FE, // l3-n voltage
		/*	        VoltageLL:		0x0400,	// l-l voltage
		VoltageL12:		0x0402,	// l1-l2 voltage
		VoltageL23:		0x0404,	// l2-l3 voltage
		VoltageL31:		0x0406,	// l3-l1 voltage */
		Frequency: 0x0408, // line frequency

		SumL1:         0x044C, // total active energy l1
		SumL2:         0x044E, // total active energy l2
		SumL3:         0x0450, // total active energy l3
		ImportL1:      0x0452, // imported active energy l1
		ImportL2:      0x0454, // imported active energy l2
		ImportL3:      0x0456, // imported active energy l3
		Export:        0x0458, // total exported active energy
		ExportT1:      0x045A, // total exported active energy non-reset
		ExportL1:      0x045C, // exported energy l1
		ExportL2:      0x045E, // exported energy l2
		ExportL3:      0x0460, // exported energy l3
		ReactiveSum:   0x0462, // total reactive energy
		ReactiveSumL1: 0x0464, // reactive energy l1
		ReactiveSumL2: 0x0468, // reactive energy l2
		ReactiveSumL3: 0x046A, // reactive energy l3
		//		EnergyApparent:		0x046C, // total apparent energy
		//		EnergyApparentL1:	0x046E, // apparent energy l1
		//		EnergyApparentL2:	0x0470, // apparent energy l2
		//		EnergyApparentL3:	0x0472, // apparent energy l3
		Cosphi:          0x0472, // power factor
		CosphiL1:        0x0474, // power factor l1
		CosphiL2:        0x0476, // power factor l2
		CosphiL3:        0x0478, // power factor l3
		ReactivePower:   0x047A, // total reactive power
		ReactivePowerL1: 0x047C, // reactive power l1
		ReactivePowerL2: 0x047e, // reactive power l2
		ReactivePowerL3: 0x0480, // reactive power l3
		ApparentPower:   0x0482, // total apparent power
		ApparentPowerL1: 0x0484, // apparent power l1
		ApparentPowerL2: 0x0486, // apparent power l2
		ApparentPowerL3: 0x0488, // apparent power l3
		CurrentL1:       0x048A, // current l1
		CurrentL2:       0x048C, // current l2
		CurrentL3:       0x048E, // current l3
		//		PowerDemand:		0x0490, // demand power
		//		MinimumPowerDemand:	0x0492, // minimum demand power
		//		MaximumPowerDemand:	0x0494, // maximum demand power
		//		ApparentPowerDemand:	0x0496, // apparent demand power
		//		PowerDemandL1:		0x0498, // demand power l1
		//		PowerDemandL2:		0x049A, // demand power l2
		//		PowerDemandL3:		0x049C, // demand power l3
	}
	return &SEMTRProducer{Opcodes: ops}
}

func (p *SEMTRProducer) Description() string {
	return "SolarEdge SE-MTR-3Y"
}

// RTUIeee754SolaredgeToFloat64 converts 32 bit IEEE 754 solar edge big endian float readings
// The wire protocol seems to have some strange byte ordering (?)
func RTUIeee754SolaredgeToFloat64(b []byte) float64 {
	_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
	bits := uint32(b[1]) | uint32(b[0])<<8 | uint32(b[3])<<16 | uint32(b[2])<<24
	f := math.Float32frombits(bits)
	return float64(f)
}

func (p *SEMTRProducer) snip(iec Measurement) Operation {
	operation := Operation{
		FuncCode:  ReadInputReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   2,
		IEC61850:  iec,
		Transform: RTUIeee754SolaredgeToFloat64,
	}
	return operation
}

func (p *SEMTRProducer) Probe() Operation {
	return p.snip(VoltageL1)
}

func (p *SEMTRProducer) Produce() (res []Operation) {
	for op := range p.Opcodes {
		res = append(res, p.snip(op))
	}

	return res
}
