package rs485

import (
	"encoding/binary"
	"math"
)

// RTUTransform functions convert RTU bytes to meaningful data types.
type RTUTransform func([]byte) float64

// RTUIeee754ToFloat64 converts 32 bit IEEE 754 float readings
func RTUIeee754ToFloat64(b []byte) float64 {
	bits := binary.BigEndian.Uint32(b)
	f := math.Float32frombits(bits)
	return float64(f)
}

// RTUUint16ToFloat64 converts 16 bit unsigned integer readings
func RTUUint16ToFloat64(b []byte) float64 {
	u := binary.BigEndian.Uint16(b)
	return float64(u)
}

// RTUUint32ToFloat64 converts 32 bit unsigned integer readings
func RTUUint32ToFloat64(b []byte) float64 {
	u := binary.BigEndian.Uint32(b)
	return float64(u)
}

// RTUUint64ToFloat64 converts 64 bit unsigned integer readings
func RTUUint64ToFloat64(b []byte) float64 {
	u := binary.BigEndian.Uint64(b)
	return float64(u)
}

// RTUInt16ToFloat64 converts 16 bit signed integer readings
func RTUInt16ToFloat64(b []byte) float64 {
	u := int16(binary.BigEndian.Uint16(b))
	return float64(u)
}

// RTUInt32ToFloat64 converts 32 bit signed integer readings
func RTUInt32ToFloat64(b []byte) float64 {
	u := int32(binary.BigEndian.Uint32(b))
	return float64(u)
}

// MakeScaledTransform creates an RTUTransform with applied scaler
func MakeScaledTransform(transform RTUTransform, scaler float64) RTUTransform {
	return RTUTransform(func(b []byte) float64 {
		unscaled := transform(b)
		f := unscaled / scaler
		return f
	})
}
