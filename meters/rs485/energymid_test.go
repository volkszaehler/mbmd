package rs485

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/volkszaehler/mbmd/meters"
)

func TestEnergyMIDMantissaTransform(t *testing.T) {
	// voltage 2309, exponent -1
	b := make([]byte, 2*9)
	b[0], b[1] = 0x09, 0x05
	b[len(b)-1] = 0xff
	assert.InDelta(t, 230.9, energyMIDMantissaTransform(b), 1e-9)

	// power -1234, exponent 1
	b = make([]byte, 2*10)
	b[0], b[1] = 0xfb, 0x2e
	b[len(b)-1] = 0x01
	assert.InDelta(t, -12340, energyMIDMantissaTransform(b), 1e-9)

	// undefined
	b[0], b[1] = 0x80, 0x00
	assert.True(t, math.IsNaN(energyMIDMantissaTransform(b)))
}

func TestEnergyMIDEnergyTransform(t *testing.T) {
	// 4561240 Wh, factor 1
	b := make([]byte, 2*10)
	b[0], b[1], b[2], b[3] = 0x00, 0x45, 0x99, 0x58
	b[len(b)-1] = 0x01
	assert.InDelta(t, 4561.24, energyMIDEnergyTransform(b), 1e-9)

	// factor 10
	b[len(b)-1] = 0x0a
	assert.InDelta(t, 45612.4, energyMIDEnergyTransform(b), 1e-9)
}

func TestEnergyMIDProduce(t *testing.T) {
	p := NewEnergyMIDProducer()

	ops := make(map[Measurement]Operation)
	for _, op := range p.Produce() {
		assert.Equal(t, uint8(ReadInputReg), op.FuncCode)
		ops[op.IEC61850] = op
	}

	for iec, expect := range map[Measurement]struct{ opcode, readlen uint16 }{
		VoltageL1:     {4, 9},
		CurrentL1:     {100, 9},
		PowerL1:       {200, 13},
		Power:         {203, 10},
		ReactivePower: {207, 6},
		Import:        {300, 10},
		Export:        {302, 8},
		Frequency:     {11, 1},
	} {
		assert.Equal(t, expect.opcode, ops[iec].OpCode, iec.String())
		assert.Equal(t, expect.readlen, ops[iec].ReadLen, iec.String())
	}
}
