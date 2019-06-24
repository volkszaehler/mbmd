package bus

import (
	"time"

	"github.com/grid-x/modbus"
)

type TCP struct {
	manager
	address string
	Client  modbus.Client
	Handler *modbus.TCPClientHandler
}

// NewTCPClientHandler creates a TCO modbus handler
func NewTCPClientHandler(device string) *modbus.TCPClientHandler {
	handler := modbus.NewTCPClientHandler(device)

	handler.Timeout = 1 * time.Second
	handler.ProtocolRecoveryTimeout = 10 * time.Second
	handler.LinkRecoveryTimeout = 15 * time.Second

	// if verbose {
	// 	logger := &ModbusLogger{}
	// 	handler.Logger = logger
	// }

	return handler
}

func NewTCP(address string) Bus {
	handler := NewTCPClientHandler(address)
	client := modbus.NewClient(handler)

	b := &TCP{
		// manager: NewManager(),
		address: address,
		Client:  client,
		Handler: handler,
	}

	b.manager = NewManager(b)

	return b
}

func (b *TCP) String() string {
	return b.address
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

// Reconnect refreshes underlying modbus TCP connection by closing it
// and thus forcing the client to reopen
func (b *TCP) Reconnect() {
	b.Handler.Close()
}
