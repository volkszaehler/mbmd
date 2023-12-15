package encoding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	u16 = uint16(0x1234)
	i16 = int16(0x1234)
	u32 = uint32(0x12345678)
	i32 = int32(0x12345678)
	u64 = uint64(0x1234567812345678)
	i64 = int64(0x1234567812345678)
)

func TestUint16(t *testing.T) {
	b := make([]byte, 2)

	PutUint16(b, u16)
	assert.Equal(t, u16, Uint16(b))
}

func TestInt16(t *testing.T) {
	b := make([]byte, 2)

	PutInt16(b, i16)
	assert.Equal(t, i16, Int16(b))
}

func TestUint32(t *testing.T) {
	b := make([]byte, 4)

	PutUint32(b, u32)
	assert.Equal(t, u32, Uint32(b))
}

func TestInt32(t *testing.T) {
	b := make([]byte, 4)

	PutInt32(b, i32)
	assert.Equal(t, i32, Int32(b))
}

func TestUint64(t *testing.T) {
	b := make([]byte, 8)

	PutUint64(b, u64)
	assert.Equal(t, u64, Uint64(b))
}

func TestInt64(t *testing.T) {
	b := make([]byte, 8)

	PutInt64(b, i64)
	assert.Equal(t, i64, Int64(b))
}

func TestFloat32(t *testing.T) {
	v := float32(1)
	b := make([]byte, 4)

	PutFloat32(b, v)
	assert.Equal(t, v, Float32(b))
}

func TestFloat64(t *testing.T) {
	v := float64(1)
	b := make([]byte, 8)

	PutFloat64(b, v)
	assert.Equal(t, v, Float64(b))
}

func TestUint32LswFirst(t *testing.T) {
	b := make([]byte, 4)

	PutUint32LswFirst(b, u32)
	assert.Equal(t, u32, Uint32LswFirst(b))
}

func TestInt32LswFirst(t *testing.T) {
	b := make([]byte, 4)

	PutInt32LswFirst(b, i32)
	assert.Equal(t, i32, Int32LswFirst(b))
}

func TestFloat32LswFirst(t *testing.T) {
	v := float32(1)
	b := make([]byte, 4)

	PutFloat32LswFirst(b, v)
	assert.Equal(t, v, Float32LswFirst(b))
}

func TestUint64LswFirst(t *testing.T) {
	v := uint64(1)
	b := make([]byte, 8)

	PutUint64LswFirst(b, v)
	assert.Equal(t, v, Uint64LswFirst(b))
}

func TestInt64LswFirst(t *testing.T) {
	v := int64(1)
	b := make([]byte, 8)

	PutInt64LswFirst(b, v)
	assert.Equal(t, v, Int64LswFirst(b))
}

func TestFloat64LswFirst(t *testing.T) {
	v := float64(1)
	b := make([]byte, 8)

	PutFloat64LswFirst(b, v)
	assert.Equal(t, v, Float64LswFirst(b))
}
