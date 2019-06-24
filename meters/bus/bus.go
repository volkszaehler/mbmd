package bus

import (
	"time"

	"github.com/volkszaehler/mbmd/meters"
)

// Bus encapsulates a physical modbus connection, either RTU or TCP
type Bus interface {
	// physical access

	// Slave sets the modbus device id
	Slave(deviceID uint8)

	// Timeout sets the modbus timeout
	Timeout(timeout time.Duration) time.Duration

	// Reconnect closes the modbus connection.
	// This forces the modbus client to reopen the connection before the next bus operations.
	Reconnect()

	// Logger sets a logging instance for physical bus operations
	Logger(l Logger)

	// String returns the bus device/address
	String() string

	// device management

	// Add adds a device connected to a specific id
	Add(id uint8, device meters.Device) error

	// Run queries all devices connected to the bus
	Run()
}
