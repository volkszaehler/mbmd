package connection

import (
	"time"

	"github.com/volkszaehler/mbmd/meters"
)

// Connection encapsulates a physical modbus connection, either RTU or TCP
type Connection interface {
	// ModbusClient returns the underlying modbus client
	ModbusClient() meters.ModbusClient

	// Slave sets the modbus device id
	Slave(deviceID uint8)

	// Timeout sets the modbus timeout
	Timeout(timeout time.Duration) time.Duration

	// Close closes the modbus connection.
	// This forces the modbus client to reopen the connection before the next bus operations.
	Close()

	// Logger sets a logging instance for physical bus operations
	Logger(l Logger)

	// String returns the bus device/address
	String() string
}
