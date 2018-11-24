package meters

import (
	"encoding/binary"
	"math"
)

const (
	// MODBUS protocol address (base 0)
	sunspecBase = 40000
)

// RTUUint16ToFloat64WithNaN converts 16 bit unsigned integer readings
// If byte sequence is 0xffff, NaN is returned for compatibility with SunSpec/SE 1-phase inverters
func RTUUint16ToFloat64WithNaN(b []byte) float64 {
	u := binary.BigEndian.Uint16(b)
	if u == 0xffff {
		return math.NaN()
	}
	return float64(u)
}
