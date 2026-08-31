package rs485

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProducers(t *testing.T) {
	for name, fun := range Producers {
		p := fun()
		assert.Equal(t, len(p.Measurements()), len(p.Produce()), name)
	}
}
