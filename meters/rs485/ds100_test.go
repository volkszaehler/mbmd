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
	for _, iec := range []Measurement{Sum, SumL1, SumL2, SumL3} {
		assert.InDelta(t, -1.0, ops[iec].Transform(negative), 1e-9, iec.String())
	}

	// 1234.56 kWh reads the same as before
	positive := []byte{0x00, 0x01, 0xe2, 0x40}
	assert.InDelta(t, 1234.56, ops[Sum].Transform(positive), 1e-9)

	// import and export counters cannot go negative and stay unsigned
	assert.InDelta(t, 42949671.96, ops[Import].Transform(negative), 1e-6)
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
	} {
		assert.Equal(t, opcode, ops[iec].OpCode, iec.String())
		assert.Equal(t, uint16(2), ops[iec].ReadLen, iec.String())
	}
}
