package meters

import (
	"time"

	"github.com/grid-x/modbus"
)

// TCP is a TCP modbus connection
type TCP struct {
	Client  modbus.Client
	Handler *modbus.TCPClientHandler
}

// NewTCPClientHandler creates a TCP modbus handler
func NewTCPClientHandler(device string) *modbus.TCPClientHandler {
	handler := modbus.NewTCPClientHandler(device)

	// set default timings
	handler.ProtocolRecoveryTimeout = 10 * time.Second
	handler.LinkRecoveryTimeout = 15 * time.Second

	return handler
}

var _ Connection = (*TCP)(nil)

// NewTCP creates a TCP modbus client
func NewTCP(address string) *TCP {
	handler := NewTCPClientHandler(address)
	client := modbus.NewClient(handler)

	b := &TCP{
		Client:  client,
		Handler: handler,
	}

	// TODO prometheus: TCPConnectionCreated

	return b
}

// String returns the bus connection address (TCP)
func (b *TCP) String() string {
	return b.Handler.Address
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
	// TODO prometheus: TCPConnectionClosed
}

// Clone clones the modbus connection.
func (b *TCP) Clone(deviceID byte) Connection {
	handler := b.Handler.Clone()
	handler.SetSlave(deviceID)

	return &TCP{
		Client:  modbus.NewClient(handler),
		Handler: handler,
	}
}
