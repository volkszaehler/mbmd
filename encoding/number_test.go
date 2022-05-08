package encoding

import "testing"

func TestUint32LswFirst(t *testing.T) {
	expect := uint32(0x03040102)
	out := uint32(Uint32LswFirst([]byte{1, 2, 3, 4}))
	if out != expect {
		t.Errorf("wanted: %08x, got %08x", expect, out)
	}
}
