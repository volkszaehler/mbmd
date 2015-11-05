package serial

import (
	"os"
	"testing"
)

const (
	// socat -d -d pty,raw,echo=0 pty,raw,echo=0
	pty1 = "/dev/ttys009"
	pty2 = "/dev/ttys010"
)

func TestReadWrite(t *testing.T) {
	checkPty(t)

	config1 := Config{Address: pty1}
	port1, err := Open(&config1)
	if err != nil {
		t.Fatal(err)
	}
	defer port1.Close()

	config2 := Config{
		Address:  pty2,
		BaudRate: 57600,
		DataBits: 7,
		Parity:   "N",
		StopBits: 2,
	}
	port2, err := Open(&config2)
	if err != nil {
		t.Fatal(err)
	}
	defer port2.Close()

	message := "test serial"
	n, err := port1.Write([]byte(message))
	if err != nil {
		t.Fatal(err)
	}
	if n != len(message) {
		t.Fatalf("unexpected write length %v", n)
	}
	var buf [16]byte
	n, err = port2.Read(buf[:])
	if err != nil {
		t.Fatal(err)
	}
	if string(buf[:n]) != message {
		t.Fatalf("unexpected response %q (len: %d)", buf[:n], n)
	}
}

func checkPty(t *testing.T) {
	for _, p := range [...]string{pty1, pty2} {
		if _, err := os.Stat(p); err != nil {
			t.Skipf("%v does not exist", p)
		}
	}
}
