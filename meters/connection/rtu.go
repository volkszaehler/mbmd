package connection

import (
	"log"
	"time"

	"github.com/grid-x/modbus"
	"github.com/volkszaehler/mbmd/meters"
)

const (
	_ = iota
	Comset2400_8N1
	Comset9600_8N1
	Comset19200_8N1
	Comset2400_8E1
	Comset9600_8E1
	Comset19200_8E1
)

type RTU struct {
	device  string
	Client  meters.Client
	Handler *modbus.RTUClientHandler
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

func (b *RTU) String() string {
	return b.device
}

func (b *RTU) ModbusClient() meters.Client {
	return b.Client
}

func (b *RTU) Logger(l Logger) {
	b.Handler.Logger = l
}

func (b *RTU) Slave(deviceID uint8) {
	b.Handler.SetSlave(deviceID)
}

func (b *RTU) Timeout(timeout time.Duration) time.Duration {
	t := b.Handler.Timeout
	b.Handler.Timeout = timeout
	return t
}

func (b *RTU) Reconnect() {

}
