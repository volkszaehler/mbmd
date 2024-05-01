package meters

import (
	"time"

	"github.com/grid-x/modbus"
)

// RTUOverUDP is a RTU encoder over a TCP modbus connection
type RTUOverUDP struct {
	address string
	Client  modbus.Client
	Handler *modbus.RTUOverUDPClientHandler
	prevID  uint8
}

// NewRTUOverUDPClientHandler creates a RTU over TCP modbus handler
func NewRTUOverUDPClientHandler(device string) *modbus.RTUOverUDPClientHandler {
	return modbus.NewRTUOverUDPClientHandler(device)
}

// NewRTUOverUDP creates a TCP modbus client
func NewRTUOverUDP(address string) Connection {
	handler := NewRTUOverUDPClientHandler(address)
	client := modbus.NewClient(handler)

	b := &RTUOverUDP{
		address: address,
		Client:  client,
		Handler: handler,
	}

	return b
}

// String returns the bus connection address (TCP)
func (b *RTUOverUDP) String() string {
	return b.address
}

// ModbusClient returns the TCP modbus client
func (b *RTUOverUDP) ModbusClient() modbus.Client {
	return b.Client
}

// Logger sets a logging instance for physical bus operations
func (b *RTUOverUDP) Logger(l Logger) {
	b.Handler.Logger = l
}

// Slave sets the modbus device id for the following operations
func (b *RTUOverUDP) Slave(deviceID uint8) {
	// Some devices like SDM need to have a little pause between querying different device ids
	if b.prevID != 0 && deviceID != b.prevID {
		time.Sleep(time.Duration(100) * time.Millisecond)
	}

	b.prevID = deviceID
	b.Handler.SetSlave(deviceID)
}

// Timeout sets the modbus timeout
func (b *RTUOverUDP) Timeout(timeout time.Duration) time.Duration {
	return 0
}

// ConnectDelay sets the the initial delay after connecting before starting communication
func (b *RTUOverUDP) ConnectDelay(delay time.Duration) {
}

// Close closes the modbus connection.
// This forces the modbus client to reopen the connection before the next bus operations.
func (b *RTUOverUDP) Close() {
	b.Handler.Close()
}

// Clone clones the modbus connection.
func (b *RTUOverUDP) Clone() {
	b.Handler.Clone()
}
