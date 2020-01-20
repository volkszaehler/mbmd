package rs485

import (
	"fmt"
	"time"

	"github.com/grid-x/modbus"
	"github.com/pkg/errors"
	"github.com/volkszaehler/mbmd/meters"
)

const (
	ReadHoldingReg = 3
	ReadInputReg   = 4
)

type rs485 struct {
	producer Producer
	ops      chan Operation
	inflight Operation
}

// NewDevice creates a device who's type must exist in the producer registry
func NewDevice(typeid string) (meters.Device, error) {
	if factory, ok := Producers[typeid]; ok {
		device := &rs485{
			producer: factory(),
			ops:      make(chan Operation),
		}

		// ringbuffer of device operations
		go func(d *rs485) {
			for {
				for _, op := range d.producer.Produce() {
					d.ops <- op
				}
			}
		}(device)

		return device, nil
	}

	return nil, fmt.Errorf("unknown meter type %s", typeid)
}

// Initialize prepares the device for usage. Any setup or initilization should be done here.
func (d *rs485) Initialize(client modbus.Client) error {
	return nil
}

// Descriptor returns the device descriptor. Since this method doe not have bus access the descriptor should be preared
// during initilization.
func (d *rs485) Descriptor() meters.DeviceDescriptor {
	return meters.DeviceDescriptor{
		Manufacturer: d.producer.Type(),
		Model:        d.producer.Description(),
	}
}

func (d *rs485) rawQuery(client modbus.Client, op Operation) (bytes []byte, err error) {
	if op.ReadLen == 0 {
		return bytes, fmt.Errorf("invalid meter operation %v", op)
	}

	switch op.FuncCode {
	case ReadHoldingReg:
		bytes, err = client.ReadHoldingRegisters(op.OpCode, op.ReadLen)
	case ReadInputReg:
		bytes, err = client.ReadInputRegisters(op.OpCode, op.ReadLen)
	default:
		return bytes, fmt.Errorf("unknown function code %d", op.FuncCode)
	}

	if err != nil {
		return bytes, errors.Wrap(err, "read failed")
	}

	return bytes, nil
}

func (d *rs485) query(client modbus.Client, op Operation) (res meters.MeasurementResult, err error) {
	if op.Transform == nil {
		return res, fmt.Errorf("transformation not defined: %v", op)
	}

	bytes, err := d.rawQuery(client, op)
	if err != nil {
		return res, err
	}

	res = meters.MeasurementResult{
		Measurement: op.IEC61850,
		Value:       op.Transform(bytes),
		Timestamp:   time.Now(),
	}

	return res, nil
}

// Probe is called by the handler after preparing the bus by setting the device id
func (d *rs485) Probe(client modbus.Client) (res bool, err error) {
	op := d.producer.Probe()

	// use specific identificator for devices that are able to recognize
	// themselves reliably
	if idf, ok := d.producer.(Identificator); ok {
		bytes, err := d.rawQuery(client, op)
		if err != nil {
			return false, err
		}

		match := idf.Identify(bytes)
		return match, nil
	}

	// use default validator looking for 110/230V
	measurement, err := d.query(client, op)
	if err != nil {
		return false, err
	}

	v := validator{[]float64{110, 230}}
	match := v.validate(measurement.Value)

	return match, nil
}

// Query is called by the handler after preparing the bus by setting the device id and waiting for rate limit
func (d *rs485) Query(client modbus.Client) (res []meters.MeasurementResult, err error) {
	res = make([]meters.MeasurementResult, 0)

	// Query loop will try to read all operations in a single run. It will
	// always start with the current inflight operation. If an error is encountered,
	// the partial results are returned. The loop is terminated after as many
	// operations have been executed as the producer provides in a single run.
	// In case of a flakey connection this guarantees that all registers are
	// read at an equal rate.
	for range d.producer.Produce() {
		// get next inflight
		if d.inflight.FuncCode == 0 {
			d.inflight = <-d.ops
		}

		m, err := d.query(client, d.inflight)
		if err != nil {
			return res, err
		}

		// mark inflight operation as completed
		d.inflight.FuncCode = 0

		res = append(res, m)
	}

	return res, nil
}
