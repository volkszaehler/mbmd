/*
Package serial provides a cross-platform serial reader and writer.
*/
package serial

import (
	"errors"
	"io"
	"time"
)

var (
	// ErrTimeout is occurred when timing out.
	ErrTimeout = errors.New("serial: timeout")
)

// Config is common configuration for serial port.
type Config struct {
	// Device path (/dev/ttyS0)
	Address string
	// Baud rate (default 19200)
	BaudRate int
	// Data bits: 5, 6, 7 or 8 (default 8)
	DataBits int
	// Stop bits: 1 or 2 (default 1)
	StopBits int
	// Parity: N - None, E - Even, O - Odd (default E)
	// (The use of no parity requires 2 stop bits.)
	Parity string
	// Read (Write) timeout.
	Timeout time.Duration
}

// Port is the interface for controlling serial port.
type Port interface {
	io.ReadWriteCloser
	// Connect connects to the serial port.
	Open(*Config) error
}

// Open opens a serial port.
func Open(c *Config) (p Port, err error) {
	p = New()
	err = p.Open(c)
	return
}
