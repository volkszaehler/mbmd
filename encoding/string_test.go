package encoding

import "testing"

func TestStringLswFirst(t *testing.T) {
	expect := "ABCD"
	out := StringLsbFirst([]byte("BADC"))
	if out != expect {
		t.Errorf("wanted: %s, got %s", expect, out)
	}
}
