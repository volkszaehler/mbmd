package rs485

import (
	"encoding/binary"
	"fmt"

	"github.com/grid-x/modbus"
	. "github.com/volkszaehler/mbmd/meters"
)

func init() {
	Register(NewJanitzaProducer)
}

const (
	METERTYPE_JANITZA = "JANITZA"
)

type JanitzaProducer struct {
	Opcodes
}

func NewJanitzaProducer() Producer {
	/**
	 * Opcodes for Janitza B23.
	 * See https://www.janitza.de/betriebsanleitungen.html?file=files/download/manuals/current/B-Series/MID-Energy-Meters-Product-Manual.pdf
	 */
	ops := Opcodes{
		VoltageL1: 0x4A38,
		VoltageL2: 0x4A3A,
		VoltageL3: 0x4A3C,
		CurrentL1: 0x4A44,
		CurrentL2: 0x4A46,
		CurrentL3: 0x4A48,
		PowerL1:   0x4A4C,
		PowerL2:   0x4A4E,
		PowerL3:   0x4A50,
		ImportL1:  0x4A76,
		ImportL2:  0x4A78,
		ImportL3:  0x4A7A,
		Import:    0x4A7C,
		ExportL1:  0x4A7E,
		ExportL2:  0x4A80,
		ExportL3:  0x4A82,
		Export:    0x4A84,
		CosphiL1:  0x4A64,
		CosphiL2:  0x4A66,
		CosphiL3:  0x4A68,
	}
	return &JanitzaProducer{Opcodes: ops}
}

// Type implements Producer interface
func (p *JanitzaProducer) Type() string {
	return METERTYPE_JANITZA
}

// Description implements Producer interface
func (p *JanitzaProducer) Description() string {
	return "Janitza MID B-Series meters"
}

func (p *JanitzaProducer) Initialize(client modbus.Client, descriptor *DeviceDescriptor) error {
	// serial
	if bytes, err := client.ReadHoldingRegisters(0x8900, 2); err == nil {
		descriptor.Serial = fmt.Sprintf("%4d", binary.BigEndian.Uint32(bytes))
	}
	// firmware
	if bytes, err := client.ReadHoldingRegisters(0x8908, 8); err == nil {
		descriptor.Version = string(bytes)
	}
	// type
	if bytes, err := client.ReadHoldingRegisters(0x8960, 6); err == nil {
		descriptor.Model = string(bytes)
	}

	// assume success
	return nil
}

func (p *JanitzaProducer) snip(iec Measurement) Operation {
	snip := Operation{
		FuncCode:  readHoldingReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   2,
		IEC61850:  iec,
		Transform: RTUIeee754ToFloat64,
	}
	return snip
}

// Probe implements Producer interface
func (p *JanitzaProducer) Probe() Operation {
	return p.snip(VoltageL1)
}

// Produce implements Producer interface
func (p *JanitzaProducer) Produce() (res []Operation) {
	for op := range p.Opcodes {
		res = append(res, p.snip(op))
	}

	return res
}
