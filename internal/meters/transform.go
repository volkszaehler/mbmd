package meters

import (
	"encoding/binary"
	"math"
)

// RTUTransform functions convert RTU bytes to meaningful data types.
type RTUTransform func([]byte) float64

// RTU32ToFloat64 converts 32 bit readings
func RTU32ToFloat64(b []byte) float64 {
	bits := binary.BigEndian.Uint32(b)
	f := math.Float32frombits(bits)
	return float64(f)
}

// RTU16ToFloat64 converts 16 bit readings
func RTU16ToFloat64(b []byte) float64 {
	u := binary.BigEndian.Uint16(b)
	return float64(u)
}

func rtuScaledInt32ToFloat64(b []byte, scalar float64) float64 {
	unscaled := float64(binary.BigEndian.Uint32(b))
	f := unscaled / scalar
	return float64(f)
}

// MakeRTU32ScaledIntToFloat64 creates a 32 bit scaled reading transform
func MakeRTU32ScaledIntToFloat64(scalar float64) RTUTransform {
	return RTUTransform(func(b []byte) float64 {
		return rtuScaledInt32ToFloat64(b, scalar)
	})
}

func rtuScaledInt16ToFloat64(b []byte, scalar float64) float64 {
	unscaled := float64(binary.BigEndian.Uint16(b))
	f := unscaled / scalar
	return float64(f)
}

// MakeRTU16ScaledIntToFloat64 creates a 16 bit scaled reading transform
func MakeRTU16ScaledIntToFloat64(scalar float64) RTUTransform {
	return RTUTransform(func(b []byte) float64 {
		return rtuScaledInt16ToFloat64(b, scalar)
	})
}
