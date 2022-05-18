package rs485

import (
	"github.com/volkszaehler/mbmd/encoding"
)

// RTUTransform functions convert RTU bytes to meaningful data types.
type RTUTransform func([]byte) float64

// RTUIeee754ToFloat64 converts 32 bit IEEE 754 float readings
func RTUIeee754ToFloat64(b []byte) float64 {
	return float64(encoding.Float32(b))
}

// RTUIeee754ToFloat64Swapped converts 32 bit IEEE 754 float readings
func RTUIeee754ToFloat64Swapped(b []byte) float64 {
	return float64(encoding.Float32LswFirst(b))
}

// RTUFloat64ToFloat64 converts 64 bit float readings
func RTUFloat64ToFloat64(b []byte) float64 {
	return encoding.Float64(b)
}

// RTUUint16ToFloat64 converts 16 bit unsigned integer readings
func RTUUint16ToFloat64(b []byte) float64 {
	return float64(encoding.Uint16(b))
}

// RTUUint32ToFloat64 converts 32 bit unsigned integer readings
func RTUUint32ToFloat64(b []byte) float64 {
	return float64(encoding.Uint32(b))
}

// RTUUint32ToFloat64Swapped converts 32 bit unsigned integer readings with swapped word order
func RTUUint32ToFloat64Swapped(b []byte) float64 {
	return float64(encoding.Uint32LswFirst(b))
}

// RTUUint64ToFloat64 converts 64 bit unsigned integer readings
func RTUUint64ToFloat64(b []byte) float64 {
	return float64(encoding.Uint64(b))
}

// RTUInt16ToFloat64 converts 16 bit signed integer readings
func RTUInt16ToFloat64(b []byte) float64 {
	return float64(encoding.Int16(b))
}

// RTUInt32ToFloat64 converts 32 bit signed integer readings
func RTUInt32ToFloat64(b []byte) float64 {
	return float64(encoding.Int32(b))
}

// RTUInt32ToFloat64Swapped converts 32 bit unsigned integer readings with swapped word order
func RTUInt32ToFloat64Swapped(b []byte) float64 {
	return float64(encoding.Int32LswFirst(b))
}

// RTUInt64ToFloat64 converts 64 bit signed integer readings
func RTUInt64ToFloat64(b []byte) float64 {
	return float64(encoding.Int64(b))
}

// MakeScaledTransform creates an RTUTransform with applied scaler
func MakeScaledTransform(transform RTUTransform, scaler float64) RTUTransform {
	return RTUTransform(func(b []byte) float64 {
		unscaled := transform(b)
		f := unscaled / scaler
		return f
	})
}
