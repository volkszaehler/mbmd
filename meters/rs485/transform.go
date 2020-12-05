package rs485

import (
	"encoding/binary"
	"math"

	"github.com/volkszaehler/mbmd/encoding"
)

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
	bits := encoding.BigEndianUint32Swapped(b)
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
	u := uint32(encoding.BigEndianUint32Swapped(b))
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
	u := int32(encoding.BigEndianUint32Swapped(b))
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
