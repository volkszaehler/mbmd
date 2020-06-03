package rs485

import (
	"encoding/binary"
	"math"
)

// BigEndianUint32Swapped converts bytes to uint32 wrapped as uint64 with swapped word order.
// To use the result as int32 value make sure to convert to uint32 first before converting to int32.
func BigEndianUint32Swapped(b []byte) uint64 {
	_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
	return uint64(b[3])<<16 | uint64(b[2])<<24 | uint64(b[1]) | uint64(b[0])<<8
}

// RTUTransform functions convert RTU bytes to meaningful data types.
type RTUTransform func([]byte) float64

// RTUIeee754ToFloat64 converts 32 bit IEEE 754 float readings
func RTUIeee754ToFloat64(b []byte) float64 {
	bits := binary.BigEndian.Uint32(b)
	f := math.Float32frombits(bits)
	return float64(f)
}

// RTUIeee754ToFloat64Swapped converts 32 bit IEEE 754 float readings
func RTUIeee754ToFloat64Swapped(b []byte) float64 {
	bits := uint32(BigEndianUint32Swapped(b))
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

// RTUUint32ToFloat64Swapped converts 32 bit unsigned integer readings with swapped word order
func RTUUint32ToFloat64Swapped(b []byte) float64 {
	u := uint32(BigEndianUint32Swapped(b))
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

// RTUInt32ToFloat64Swapped converts 32 bit unsigned integer readings with swapped word order
func RTUInt32ToFloat64Swapped(b []byte) float64 {
	u := int32(BigEndianUint32Swapped(b))
	return float64(u)
}

// RTUInt64ToFloat64 converts 64 bit signed integer readings
func RTUInt64ToFloat64(b []byte) float64 {
	u := int64(binary.BigEndian.Uint64(b))
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
