package meters

import (
	"github.com/grid-x/modbus"
)

// DeviceDescriptor describes a device
type DeviceDescriptor struct {
	Type         string
	Manufacturer string
	Model        string
	Options      string
	Version      string
	Serial       string
	SubDevice    int
}

// Device is a modbus device that can be described, probed and queried
type Device interface {
	// Initialize prepares the device for usage. Any setup or initialization should be done here.
	// It requires that the client has the correct device id applied.
	Initialize(client modbus.Client) error

	// Descriptor returns the device descriptor. Since this method does not have
	// bus access the descriptor should be prepared during initialization.
	Descriptor() DeviceDescriptor

	// Probe tests if a basic register, typically VoltageL1, can be read.
	// It requires that the client has the correct device id applied.
	Probe(client modbus.Client) (MeasurementResult, error)

	// Query retrieves all registers that the device supports.
	// It requires that the client has the correct device id applied.
	Query(client modbus.Client) ([]MeasurementResult, error)
}
