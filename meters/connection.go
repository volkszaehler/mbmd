package meters

import (
	"time"

	"github.com/grid-x/modbus"
)

// Connection encapsulates a physical modbus connection, either RTU or TCP
type Connection interface {
	// ModbusClient returns the underlying modbus client
	ModbusClient() modbus.Client

	// Slave sets the modbus device id for the following operations
	Slave(deviceID uint8)

	// Timeout sets the modbus timeout
	Timeout(timeout time.Duration) time.Duration

	// ConnectDelay sets the the initial delay after connecting before starting communication
	ConnectDelay(delay time.Duration)

	// Close closes the modbus connection.
	// This forces the modbus client to reopen the connection before the next bus operations.
	Close()

	// Logger sets a logging instance for physical bus operations
	Logger(l Logger)

	// String returns the bus device (RTU) or bus connection address (TCP)
	String() string
}
