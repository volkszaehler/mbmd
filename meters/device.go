package meters

import (
	"github.com/grid-x/modbus"
)

type DeviceDescriptor struct {
	Manufacturer string
	Model        string
	Options      string
	Version      string
	Serial       string
}

type Device interface {
	// Initialize prepares the device for usage. Any setup or initilization should be done here.
	Initialize(client modbus.Client) error

	// Descriptor returns the device descriptor. Since this method does not have
	// bus access the descriptor should be preared during initilization.
	Descriptor() DeviceDescriptor

	// Probe is called by the BusManager after preparing the bus by setting the device id and waiting for rate limit
	Probe(client modbus.Client) (MeasurementResult, error)

	// Query is called by the BusManager after preparing the bus by setting the device id and waiting for rate limit
	Query(client modbus.Client) ([]MeasurementResult, error)
}
