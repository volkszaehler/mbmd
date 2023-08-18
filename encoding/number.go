package encoding

import (
	"encoding/binary"
	"math"
)

// Uint16 decodes bytes as uint16 in network byte order (big endian)
func Uint16(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

// Int16 decodes bytes as int16 in network byte order (big endian)
func Int16(b []byte) int16 {
	return int16(binary.BigEndian.Uint16(b))
}

// Uint32 decodes bytes as uint32 in network byte order (big endian)
func Uint32(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

// Int32 decodes bytes as int32 in network byte order (big endian)
func Int32(b []byte) int32 {
	return int32(binary.BigEndian.Uint32(b))
}

// Uint64 decodes bytes as uint64 in network byte order (big endian)
func Uint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

// Int64 decodes bytes as int64 in network byte order (big endian)
func Int64(b []byte) int64 {
	return int64(binary.BigEndian.Uint64(b))
}

// Float32 decodes bytes as float32 in network byte order (big endian)
func Float32(b []byte) float32 {
	return math.Float32frombits(Uint32(b))
}

// Float64 decodes bytes as float64 in network byte order (big endian)
func Float64(b []byte) float64 {
	return math.Float64frombits(Uint64(b))
}

// Uint32LswFirst decodes bytes as uint32 in network byte order (big endian) with least significant word first
func Uint32LswFirst(b []byte) uint32 {
	_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
	return uint32(b[2])<<24 | uint32(b[3])<<16 | uint32(b[0])<<8 | uint32(b[1])
}

// Int32LswFirst decodes bytes as int32 in network byte order (big endian) with least significant word first
func Int32LswFirst(b []byte) int32 {
	_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
	return int32(uint32(b[2])<<24 | uint32(b[3])<<16 | uint32(b[0])<<8 | uint32(b[1]))
}

// Float32LswFirst decodes bytes as float32 in network byte order (big endian) with least significant word first
func Float32LswFirst(b []byte) float32 {
	return math.Float32frombits(Uint32LswFirst(b))
}

// Uint64LswFirst decodes bytes as uint64 in network byte order (big endian) with least significant word first
func Uint64LswFirst(b []byte) uint64 {
	_ = b[7] // bounds check hint to compiler; see golang.org/issue/14808
	return uint64(b[6])<<56 | uint64(b[7])<<48 | uint64(b[4])<<40 | uint64(b[5])<<32 | uint64(b[2])<<24 | uint64(b[3])<<16 | uint64(b[0])<<8 | uint64(b[1])
}

// Int64LswFirst decodes bytes as int64 in network byte order (big endian) with least significant word first
func Int64LswFirst(b []byte) int64 {
	_ = b[7] // bounds check hint to compiler; see golang.org/issue/14808
	return int64(b[6])<<56 | int64(b[7])<<48 | int64(b[4])<<40 | int64(b[5])<<32 | int64(b[2])<<24 | int64(b[3])<<16 | int64(b[0])<<8 | int64(b[1])
}

// Float64LswFirst decodes bytes as float64 in network byte order (big endian) with least significant word first
func Float64LswFirst(b []byte) float64 {
	return math.Float64frombits(Uint64LswFirst(b))
}
