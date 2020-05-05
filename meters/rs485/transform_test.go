package rs485

import "testing"

func TestBigEndianUint32Swapped(t *testing.T) {
	expect := uint32(0x03040102)
	out := uint32(BigEndianUint32Swapped([]byte{1, 2, 3, 4}))
	if out != expect {
		t.Errorf("wanted: %08x, got %08x", expect, out)
	}
}
