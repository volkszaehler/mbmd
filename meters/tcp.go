package meters

import (
	"time"

	"github.com/grid-x/modbus"
)

// TCP is a TCP modbus connection
type TCP struct {
	address string
	Client  modbus.Client
	Handler *modbus.TCPClientHandler
}

// NewTCPClientHandler creates a TCP modbus handler
func NewTCPClientHandler(device string) *modbus.TCPClientHandler {
	handler := modbus.NewTCPClientHandler(device)

	// set default timings
	handler.Timeout = 1 * time.Second
	handler.IdleTimeout = 5 * time.Second
	handler.ProtocolRecoveryTimeout = 10 * time.Second
	handler.LinkRecoveryTimeout = 15 * time.Second

	return handler
}

// NewTCP creates a TCP modbus client
func NewTCP(address string) Connection {
	handler := NewTCPClientHandler(address)
	client := modbus.NewClient(handler)

	b := &TCP{
		address: address,
		Client:  client,
		Handler: handler,
	}

	return b
}

// String returns the bus connection address (TCP)
func (b *TCP) String() string {
	return b.address
}

// ModbusClient returns the TCP modbus client
func (b *TCP) ModbusClient() modbus.Client {
	return b.Client
}

// Logger sets a logging instance for physical bus operations
func (b *TCP) Logger(l Logger) {
	b.Handler.Logger = l
}

// Slave sets the modbus device id for the following operations
func (b *TCP) Slave(deviceID uint8) {
	b.Handler.SetSlave(deviceID)
}

// Timeout sets the modbus timeout
func (b *TCP) Timeout(timeout time.Duration) time.Duration {
	t := b.Handler.Timeout
	b.Handler.Timeout = timeout
	return t
}

// ConnectDelay sets the the initial delay after connecting before starting communication
func (b *TCP) ConnectDelay(delay time.Duration) {
	b.Handler.ConnectDelay = delay
}

// Close closes the modbus connection.
// This forces the modbus client to reopen the connection before the next bus operations.
func (b *TCP) Close() {
	b.Handler.Close()
}
