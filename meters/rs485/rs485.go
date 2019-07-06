package rs485

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/volkszaehler/mbmd/meters"
)

const (
	ReadHoldingReg = 3
	ReadInputReg   = 4
)

type rs485 struct {
	producer Producer
}

// NewDevice creates a device who's type must exist in the producer registry
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
func (d *rs485) Initialize(client meters.Client) error {
	return nil
}

// Descriptor returns the device descriptor. Since this method doe not have bus access the descriptor should be preared
// during initilization.
func (d *rs485) Descriptor() meters.DeviceDescriptor {
	return meters.DeviceDescriptor{
		Manufacturer: d.producer.Description(),
	}
}

func (d *rs485) query(client meters.Client, op Operation) (res meters.MeasurementResult, err error) {
	var bytes []byte

	if op.ReadLen <= 0 {
		return res, fmt.Errorf("invalid meter operation %v", op)
	}

	if op.Transform == nil {
		return res, fmt.Errorf("transformation not defined: %v", op)
	}

	switch op.FuncCode {
	case ReadHoldingReg:
		bytes, err = client.ReadHoldingRegisters(op.OpCode, op.ReadLen)
	case ReadInputReg:
		bytes, err = client.ReadInputRegisters(op.OpCode, op.ReadLen)
	default:
		return res, fmt.Errorf("unknown function code %d", op.FuncCode)
	}

	if err != nil {
		return res, errors.Wrap(err, "read failed")
	}

	res = meters.MeasurementResult{
		Measurement: op.IEC61850,
		Value:       op.Transform(bytes),
		Timestamp:   time.Now(),
	}

	return res, nil
}

// Probe is called by the BusManager after preparing the bus by setting the device id and waiting for rate limit
func (d *rs485) Probe(client meters.Client) (res meters.MeasurementResult, err error) {
	op := d.producer.Probe()

	res, err = d.query(client, op)
	if err != nil {
		return res, err
	}

	return res, nil
}

// Query is called by the BusManager after preparing the bus by setting the device id and waiting for rate limit
func (d *rs485) Query(client meters.Client) (res []meters.MeasurementResult, err error) {
	res = make([]meters.MeasurementResult, 0)

	for _, op := range d.producer.Produce() {
		m, err := d.query(client, op)
		if err != nil {
			return res, err
		}

		res = append(res, m)
	}

	return res, nil
}
