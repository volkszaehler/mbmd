package rs485

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/grid-x/modbus"
	"github.com/pkg/errors"
	"github.com/volkszaehler/mbmd/meters"
)

const (
	// modbus operation types
	readHoldingReg = 3
	readInputReg   = 4
)

// modbusClient is the minimal interface that is usable by the initializer interface.
// It is used to keep the producers free of modbus implementation dependencies.
type modbusClient interface {
	ReadHoldingRegisters(address, quantity uint16) (results []byte, err error)
	ReadInputRegisters(address, quantity uint16) (results []byte, err error)
}

// initializer can be implemented by producers to perform bus operations for device discovery
type initializer interface {
	// Initialize prepares the device for usage. Any setup or initialization should be done here.
	// It requires that the client has the correct device id applied.
	Initialize(client modbusClient, descriptor *meters.DeviceDescriptor) error
}

// MID meters initialization method used by Janitza and ABB
func initializeMID(client modbusClient, descriptor *meters.DeviceDescriptor) error {
	// serial
	if bytes, err := client.ReadHoldingRegisters(0x8900, 2); err == nil {
		descriptor.Serial = fmt.Sprintf("%4x", binary.BigEndian.Uint32(bytes))
	}
	// firmware
	if bytes, err := client.ReadHoldingRegisters(0x8908, 8); err == nil {
		descriptor.Version = string(bytes)
	}
	// type
	if bytes, err := client.ReadHoldingRegisters(0x8960, 6); err == nil {
		descriptor.Model = string(bytes)
	}

	// assume success
	return nil
}

type rs485 struct {
	producer   Producer
	descriptor meters.DeviceDescriptor
	ops        chan Operation
	inflight   Operation
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

// Initialize prepares the device for usage. Any setup or initialization should be done here.
func (d *rs485) Initialize(client modbus.Client) error {
	d.descriptor = meters.DeviceDescriptor{
		Manufacturer: d.producer.Type(),
		Model:        d.producer.Description(),
	}

	// does device support initializing itself?
	if p, ok := d.producer.(initializer); ok {
		return p.Initialize(client, &d.descriptor)
	}

	return nil
}

// Descriptor returns the device descriptor. Since this method doe not have bus access the descriptor should be preared
// during initialization.
func (d *rs485) Descriptor() meters.DeviceDescriptor {
	return d.descriptor
}

func (d *rs485) query(client modbus.Client, op Operation) (res meters.MeasurementResult, err error) {
	var bytes []byte

	if op.ReadLen == 0 {
		return res, fmt.Errorf("invalid meter operation %v", op)
	}

	if op.Transform == nil {
		return res, fmt.Errorf("transformation not defined: %v", op)
	}

	switch op.FuncCode {
	case readHoldingReg:
		bytes, err = client.ReadHoldingRegisters(op.OpCode, op.ReadLen)
	case readInputReg:
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

// Probe is called by the handler after preparing the bus by setting the device id
func (d *rs485) Probe(client modbus.Client) (res meters.MeasurementResult, err error) {
	op := d.producer.Probe()

	res, err = d.query(client, op)
	if err != nil {
		return res, err
	}

	return res, nil
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
