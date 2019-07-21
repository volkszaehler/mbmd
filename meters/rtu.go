package meters

import (
	"log"
	"time"

	"github.com/grid-x/modbus"
)

const (
	_ = iota
	// Comset2400_8N1 is a 2400baud 8bit uneven parity connection with 1 stop bit
	Comset2400_8N1
	// Comset9600_8N1 is a 9600baud 8bit uneven parity connection with 1 stop bit
	Comset9600_8N1
	// Comset19200_8N1 is a 19200baud 8bit uneven parity connection with 1 stop bit
	Comset19200_8N1
	// Comset2400_8E1 is a 2400baud 8bit even parity connection with 1 stop bit
	Comset2400_8E1
	// Comset9600_8E1 is a 9600baud 8bit even parity connection with 1 stop bit
	Comset9600_8E1
	// Comset19200_8E1 is a 19200baud 8bit even parity connection with 1 stop bit
	Comset19200_8E1
)

// RTU is an RTU modbus connection
type RTU struct {
	device  string
	Client  modbus.Client
	Handler *modbus.RTUClientHandler
	prevID  uint8
}

// NewClientHandler creates a serial line RTU modbus handler
func NewClientHandler(device string, comset int) *modbus.RTUClientHandler {
	handler := modbus.NewRTUClientHandler(device)

	handler.Parity = "N"
	handler.DataBits = 8
	handler.StopBits = 1

	switch comset {
	case Comset2400_8N1:
		handler.BaudRate = 2400
	case Comset9600_8N1:
		handler.BaudRate = 9600
	case Comset19200_8N1:
		handler.BaudRate = 19200
	case Comset2400_8E1:
		handler.BaudRate = 2400
		handler.Parity = "E"
	case Comset9600_8E1:
		handler.BaudRate = 9600
		handler.Parity = "E"
	case Comset19200_8E1:
		handler.BaudRate = 19200
		handler.Parity = "E"
	default:
		log.Fatal("Invalid communication set specified. See -h for help.")
	}

	handler.Timeout = 300 * time.Millisecond

	// if verbose {
	// 	logger := &ModbusLogger{}
	// 	handler.Logger = logger
	// 	log.Printf("Connecting to RTU via %s, %d %d%s%d\r\n", device,
	// 		handler.BaudRate, handler.DataBits, handler.Parity,
	// 		handler.StopBits)
	// }

	return handler
}

// NewRTU creates a RTU modbus client
func NewRTU(device string, comset int) Connection {
	handler := NewClientHandler(device, comset)
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
	if b.prevID != 0 && deviceID != b.prevID {
		time.Sleep(time.Duration(100) * time.Millisecond)
		b.prevID = deviceID
	}

	b.Handler.SetSlave(deviceID)
}

// Timeout sets the modbus timeout
func (b *RTU) Timeout(timeout time.Duration) time.Duration {
	t := b.Handler.Timeout
	b.Handler.Timeout = timeout
	return t
}

// Close closes the modbus connection.
// This forces the modbus client to reopen the connection before the next bus operations.
func (b *RTU) Close() {
	b.Handler.Close()
}
