package meters

import (
	"log"
	"strings"
	"time"

	"github.com/grid-x/modbus"
)

// RTU is an RTU modbus connection
type RTU struct {
	device  string
	Client  modbus.Client
	Handler *modbus.RTUClientHandler
	prevID  uint8
}

// NewClientHandler creates a serial line RTU modbus handler
func NewClientHandler(device string, baudrate int, comset string) *modbus.RTUClientHandler {
	handler := modbus.NewRTUClientHandler(device)

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

// NewRTU creates a RTU modbus client
func NewRTU(device string, baudrate int, comset string) Connection {
	handler := NewClientHandler(device, baudrate, comset)
	client := modbus.NewClient(handler)

	b := &RTU{
		device:  device,
		Client:  client,
		Handler: handler,
	}

	return b
}

// String returns the bus device
func (b *RTU) String() string {
	return b.device
}

// ModbusClient returns the RTU modbus client
func (b *RTU) ModbusClient() modbus.Client {
	return b.Client
}

// Logger sets a logging instance for physical bus operations
func (b *RTU) Logger(l Logger) {
	b.Handler.Logger = l
}

// Slave sets the modbus device id for the following operations
func (b *RTU) Slave(deviceID uint8) {
	// Some devices like SDM need to have a little pause between querying different device ids
	if b.prevID == 0 || deviceID != b.prevID {
		b.prevID = deviceID
		b.Handler.SetSlave(deviceID)
		time.Sleep(200 * time.Millisecond)
	}
}

// Timeout sets the modbus timeout
func (b *RTU) Timeout(timeout time.Duration) time.Duration {
	t := b.Handler.Timeout
	b.Handler.Timeout = timeout
	return t
}

// ConnectDelay sets the the initial delay after connecting before starting communication
func (b *RTU) ConnectDelay(delay time.Duration) {
	// nop
}

// Close closes the modbus connection.
// This forces the modbus client to reopen the connection before the next bus operations.
func (b *RTU) Close() {
	b.Handler.Close()
}
