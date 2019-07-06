package connection

import (
	"time"

	"github.com/grid-x/modbus"
	"github.com/volkszaehler/mbmd/meters"
)

type TCP struct {
	address string
	Client  meters.ModbusClient
	Handler *modbus.TCPClientHandler
}

// NewTCPClientHandler creates a TCO modbus handler
func NewTCPClientHandler(device string) *modbus.TCPClientHandler {
	handler := modbus.NewTCPClientHandler(device)

	// set default timings
	handler.Timeout = 1 * time.Second
	handler.ProtocolRecoveryTimeout = 10 * time.Second
	handler.LinkRecoveryTimeout = 15 * time.Second

	// if verbose {
	// 	logger := &ModbusLogger{}
	// 	handler.Logger = logger
	// }

	return handler
}

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

func (b *TCP) String() string {
	return b.address
}

func (b *TCP) ModbusClient() meters.ModbusClient {
	return b.Client
}

func (b *TCP) Logger(l Logger) {
	b.Handler.Logger = l
}

func (b *TCP) Slave(deviceID uint8) {
	b.Handler.SetSlave(deviceID)
}

func (b *TCP) Timeout(timeout time.Duration) time.Duration {
	t := b.Handler.Timeout
	b.Handler.Timeout = timeout
	return t
}

// Close closes the modbus connection.
// This forces the modbus client to reopen the connection before the next bus operations.
func (b *TCP) Close() {
	b.Handler.Close()
}
