package rs485

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/volkszaehler/mbmd/meters"
)

func TestDS100SignedSum(t *testing.T) {
	ops := make(map[Measurement]Operation)
	for _, op := range NewDS100Producer().Produce() {
		ops[op.IEC61850] = op
	}

	// -1.00 kWh, e.g. combination code Import - Export with more export than import
	negative := []byte{0xff, 0xff, 0xff, 0x9c}
	for _, iec := range []Measurement{
		Sum, SumL1, SumL2, SumL3,
		ReactiveSum, ReactiveSumL1, ReactiveSumL2, ReactiveSumL3,
	} {
		assert.InDelta(t, -1.0, ops[iec].Transform(negative), 1e-9, iec.String())
	}

	// 1234.56 kWh reads the same as before
	positive := []byte{0x00, 0x01, 0xe2, 0x40}
	assert.InDelta(t, 1234.56, ops[Sum].Transform(positive), 1e-9)

	// import and export counters cannot go negative and stay unsigned
	for _, iec := range []Measurement{Import, ReactiveImport, ReactiveExport} {
		assert.InDelta(t, 42949671.96, ops[iec].Transform(negative), 1e-6, iec.String())
	}
}

func TestDS100EnergyRegisters(t *testing.T) {
	ops := make(map[Measurement]Operation)
	for _, op := range NewDS100Producer().Produce() {
		ops[op.IEC61850] = op
	}

	// addresses from the DS100 manual, section 3.2 electric quantity parameters
	for iec, opcode := range map[Measurement]uint16{
		Import:   0x010E,
		Export:   0x0118,
		Sum:      0x0122,
		SumL1:    0x0500,
		ImportL1: 0x050A,
		ExportL1: 0x0514,
		SumL2:    0x0564,
		ImportL2: 0x056E,
		ExportL2: 0x0578,
		SumL3:    0x05C8,
		ImportL3: 0x05D2,
		ExportL3: 0x05DC,

		ReactiveImport:   0x012C,
		ReactiveExport:   0x0136,
		ReactiveSum:      0x0140,
		ReactiveSumL1:    0x051E,
		ReactiveImportL1: 0x0528,
		ReactiveExportL1: 0x0532,
		ReactiveSumL2:    0x0582,
		ReactiveImportL2: 0x058C,
		ReactiveExportL2: 0x0596,
		ReactiveSumL3:    0x05E6,
		ReactiveImportL3: 0x05F0,
		ReactiveExportL3: 0x05FA,
	} {
		assert.Equal(t, opcode, ops[iec].OpCode, iec.String())
		assert.Equal(t, uint16(2), ops[iec].ReadLen, iec.String())
	}
}

func TestDS100Opcodes(t *testing.T) {
	p := NewDS100Producer().(*DS100Producer)

	produced := make(map[Measurement]bool)
	for _, op := range p.Produce() {
		iec := op.IEC61850

		// a measurement without opcode would silently read register 0
		opcode, ok := p.Opcodes[iec]
		assert.True(t, ok, "%s is read but has no opcode", iec)
		assert.Equal(t, opcode, op.OpCode, iec.String())

		assert.False(t, produced[iec], "%s is read twice", iec)
		produced[iec] = true

		// frequency and power factor are single registers, all other values double registers
		readlen := uint16(2)
		if op.OpCode >= 0x0432 && op.OpCode <= 0x0439 {
			readlen = 1
		}
		assert.Equal(t, readlen, op.ReadLen, iec.String())
	}

	addresses := make(map[uint16]Measurement)
	for iec, opcode := range p.Opcodes {
		assert.True(t, produced[iec], "%s has an opcode but is never read", iec)

		if other, ok := addresses[opcode]; ok {
			assert.Failf(t, "duplicate opcode", "%s and %s both use 0x%04X", iec, other, opcode)
		}
		addresses[opcode] = iec
	}
}
