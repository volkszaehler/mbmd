package rs485

import (
	"fmt"

	"github.com/grid-x/modbus"

	"github.com/volkszaehler/mbmd/meters"
)

type rs485 struct {
	producer Producer
}

func NewDevice(typeid string) (meters.Device, error) {
	if factory, ok := producers[typeid]; ok {
		device := &rs485{
			producer: factory(),
		}
		return device, nil
	}

	return nil, fmt.Errorf("Unknown meter type %s", typeid)
}

// Initialize prepares the device for usage. Any setup or initilization should be done here.
func (d *rs485) Initialize(client modbus.Client) error {
	return nil
}

// Descriptor returns the device descriptor. Since this method doe not have bus access the descriptor should be preared
// during initilization.
func (d *rs485) Descriptor() meters.DeviceDescriptor {
	return meters.DeviceDescriptor{
		Manufacturer: d.producer.Description(),
	}
}

// Probe is called by the BusManager after preparing the bus by setting the device id and waiting for rate limit
func (d *rs485) Probe(client modbus.Client) (meters.MeasurementResult, error) {
	res := meters.MeasurementResult{}
	return res, nil
}

// Query is called by the BusManager after preparing the bus by setting the device id and waiting for rate limit
func (d *rs485) Query(client modbus.Client) ([]meters.MeasurementResult, error) {
	res := make([]meters.MeasurementResult, 0)
	return res, nil
}
