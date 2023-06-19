package meters

import (
	"log"
	"strings"
	"time"

	"github.com/grid-x/modbus"
)

// ASCII is an ASCII modbus connection
type ASCII struct {
	device  string
	Client  modbus.Client
	Handler *modbus.ASCIIClientHandler
	prevID  uint8
}

// NewASCIIClientHandler creates a serial line ASCII modbus handler
func NewASCIIClientHandler(device string, baudrate int, comset string) *modbus.ASCIIClientHandler {
	handler := modbus.NewASCIIClientHandler(device)

	handler.BaudRate = baudrate
	handler.DataBits = 8
	handler.StopBits = 1

	switch strings.ToUpper(comset) {
	case "8N1":
		handler.Parity = "N"
	case "8E1":
		handler.Parity = "E"
	default:
		log.Fatalf("Invalid communication set specified: %s. See -h for help.", comset)
	}

	return handler
}

// NewASCII creates a RTU modbus client
func NewASCII(device string, baudrate int, comset string) Connection {
	handler := NewASCIIClientHandler(device, baudrate, comset)
	client := modbus.NewClient(handler)

	b := &ASCII{
		device:  device,
		Client:  client,
		Handler: handler,
	}

	return b
}

// String returns the bus device
func (b *ASCII) String() string {
	return b.device
}

// ModbusClient returns the RTU modbus client
func (b *ASCII) ModbusClient() modbus.Client {
	return b.Client
}

// Logger sets a logging instance for physical bus operations
func (b *ASCII) Logger(l Logger) {
	b.Handler.Logger = l
}

// Slave sets the modbus device id for the following operations
func (b *ASCII) Slave(deviceID uint8) {
	// Some devices like SDM need to have a little pause between querying different device ids
	if b.prevID != 0 && deviceID != b.prevID {
		time.Sleep(time.Duration(100) * time.Millisecond)
		b.prevID = deviceID
	}

	b.Handler.SetSlave(deviceID)
}

// Timeout sets the modbus timeout
func (b *ASCII) Timeout(timeout time.Duration) time.Duration {
	t := b.Handler.Timeout
	b.Handler.Timeout = timeout
	return t
}

// ConnectDelay sets the the initial delay after connecting before starting communication
func (b *ASCII) ConnectDelay(delay time.Duration) {
	// nop
}

// Close closes the modbus connection.
// This forces the modbus client to reopen the connection before the next bus operations.
func (b *ASCII) Close() {
	b.Handler.Close()
}
