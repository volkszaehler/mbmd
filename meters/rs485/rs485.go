package rs485

import (
	"fmt"
	"strings"
	"time"

	"github.com/grid-x/modbus"
	"github.com/volkszaehler/mbmd/meters"
)

const (
	ReadHoldingReg = 3
	ReadInputReg   = 4
)

// RS485 implements meters.Device
type RS485 struct {
	typ      string
	producer Producer
	ops      chan Operation
	inflight Operation
}

// NewDevice creates a device who's type must exist in the producer registry
func NewDevice(typ string) (*RS485, error) {
	for t, factory := range Producers {
		if strings.EqualFold(t, typ) {
			device := &RS485{
				typ:      typ,
				producer: factory(),
			}
			return device, nil
		}
	}

	return nil, fmt.Errorf("unknown meter type: %s", typ)
}

// Initialize prepares the device for usage. Any setup or initialization should be done here.
func (d *RS485) Initialize(client modbus.Client) error {
	return nil
}

// Producer returns the underlying producer. The producer can be used to understand which operations the device supports.
func (d *RS485) Producer() Producer {
	return d.producer
}

// Descriptor returns the device descriptor. Since this method does not have bus access the descriptor should be
// prepared during initialization.
func (d *RS485) Descriptor() meters.DeviceDescriptor {
	return meters.DeviceDescriptor{
		Type:         d.typ,
		Manufacturer: d.typ,
		Model:        d.producer.Description(),
	}
}

// Probe is called by the handler after preparing the bus by setting the device id
func (d *RS485) Probe(client modbus.Client) (res meters.MeasurementResult, err error) {
	op := d.producer.Probe()

	// check for empty op in case Probe isn't supported
	if op.FuncCode == 0 {
		return res, fmt.Errorf("meter type %s doesn't support Probe", d.producer.Description())
	}

	res, err = d.QueryOp(client, op)
	if err != nil {
		return res, err
	}

	return res, nil
}

// QueryOp executes a single query operation on the bus
func (d *RS485) QueryOp(client modbus.Client, op Operation) (res meters.MeasurementResult, err error) {
	var bytes []byte

	if op.ReadLen == 0 {
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
		return res, fmt.Errorf("read failed: %v", err)
	}

	res = meters.MeasurementResult{
		Measurement: op.IEC61850,
		Value:       op.Transform(bytes),
		Timestamp:   time.Now(),
	}

	return res, nil
}

// Query is called by the handler after preparing the bus by setting the device id and waiting for rate limit
func (d *RS485) Query(client modbus.Client) (res []meters.MeasurementResult, err error) {
	res = make([]meters.MeasurementResult, 0)

	if d.ops == nil {
		d.ops = make(chan Operation)

		// ringbuffer of device operations
		go func(d *RS485) {
			for {
				for _, op := range d.producer.Produce() {
					d.ops <- op
				}
			}
		}(d)
	}

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

		m, err := d.QueryOp(client, d.inflight)
		if err != nil {
			return res, err
		}

		// mark inflight operation as completed
		d.inflight.FuncCode = 0

		res = append(res, m)
	}

	return res, nil
}
