package encoding

import (
	"encoding/binary"
	"math"
)

// Uint16 decodes bytes as uint16 in network byte order (big endian)
func Uint16(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

func PutUint16(b []byte, v uint16) {
	binary.BigEndian.PutUint16(b, v)
}

// Int16 decodes bytes as int16 in network byte order (big endian)
func Int16(b []byte) int16 {
	return int16(binary.BigEndian.Uint16(b))
}

func PutInt16(b []byte, v int16) {
	binary.BigEndian.PutUint16(b, uint16(v))
}

// Uint32 decodes bytes as uint32 in network byte order (big endian)
func Uint32(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

func PutUint32(b []byte, v uint32) {
	binary.BigEndian.PutUint32(b, v)
}

// Int32 decodes bytes as int32 in network byte order (big endian)
func Int32(b []byte) int32 {
	return int32(binary.BigEndian.Uint32(b))
}

func PutInt32(b []byte, v int32) {
	binary.BigEndian.PutUint32(b, uint32(v))
}

// Uint64 decodes bytes as uint64 in network byte order (big endian)
func Uint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func PutUint64(b []byte, v uint64) {
	binary.BigEndian.PutUint64(b, v)
}

// Int64 decodes bytes as int64 in network byte order (big endian)
func Int64(b []byte) int64 {
	return int64(binary.BigEndian.Uint64(b))
}

func PutInt64(b []byte, v int64) {
	binary.BigEndian.PutUint64(b, uint64(v))
}

// Float32 decodes bytes as float32 in network byte order (big endian)
func Float32(b []byte) float32 {
	return math.Float32frombits(Uint32(b))
}

func PutFloat32(b []byte, v float32) {
	binary.BigEndian.PutUint32(b, math.Float32bits(v))
}

// Float64 decodes bytes as float64 in network byte order (big endian)
func Float64(b []byte) float64 {
	return math.Float64frombits(Uint64(b))
}

func PutFloat64(b []byte, v float64) {
	binary.BigEndian.PutUint64(b, math.Float64bits(v))
}

// Uint32LswFirst decodes bytes as uint32 in network byte order (big endian) with least significant word first
func Uint32LswFirst(b []byte) uint32 {
	_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
	return uint32(b[2])<<24 | uint32(b[3])<<16 | uint32(b[0])<<8 | uint32(b[1])
}

func PutUint32LswFirst(b []byte, v uint32) {
	binary.BigEndian.PutUint32(b, v<<16|v>>16)
}

// Int32LswFirst decodes bytes as int32 in network byte order (big endian) with least significant word first
func Int32LswFirst(b []byte) int32 {
	_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
	return int32(uint32(b[2])<<24 | uint32(b[3])<<16 | uint32(b[0])<<8 | uint32(b[1]))
}

func PutInt32LswFirst(b []byte, v int32) {
	PutUint32LswFirst(b, uint32(v))
}

// Float32LswFirst decodes bytes as float32 in network byte order (big endian) with least significant word first
func Float32LswFirst(b []byte) float32 {
	return math.Float32frombits(Uint32LswFirst(b))
}

func PutFloat32LswFirst(b []byte, v float32) {
	PutUint32LswFirst(b, math.Float32bits(v))
}

// Uint64LswFirst decodes bytes as uint64 in network byte order (big endian) with least significant word first
func Uint64LswFirst(b []byte) uint64 {
	_ = b[7] // bounds check hint to compiler; see golang.org/issue/14808
	return uint64(b[6])<<56 | uint64(b[7])<<48 | uint64(b[4])<<40 | uint64(b[5])<<32 | uint64(b[2])<<24 | uint64(b[3])<<16 | uint64(b[0])<<8 | uint64(b[1])
}

func PutUint64LswFirst(b []byte, v uint64) {
	PutUint32LswFirst(b, uint32(v))
	PutUint32LswFirst(b[4:], uint32(v>>32))
}

// Int64LswFirst decodes bytes as int64 in network byte order (big endian) with least significant word first
func Int64LswFirst(b []byte) int64 {
	_ = b[7] // bounds check hint to compiler; see golang.org/issue/14808
	return int64(b[6])<<56 | int64(b[7])<<48 | int64(b[4])<<40 | int64(b[5])<<32 | int64(b[2])<<24 | int64(b[3])<<16 | int64(b[0])<<8 | int64(b[1])
}

func PutInt64LswFirst(b []byte, v int64) {
	PutUint64LswFirst(b, uint64(v))
}

// Float64LswFirst decodes bytes as float64 in network byte order (big endian) with least significant word first
func Float64LswFirst(b []byte) float64 {
	return math.Float64frombits(Uint64LswFirst(b))
}

func PutFloat64LswFirst(b []byte, v float64) {
	PutUint64LswFirst(b, math.Float64bits(v))
}
