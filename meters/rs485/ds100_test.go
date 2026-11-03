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
