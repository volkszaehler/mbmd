package meters

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

const (
	ReadInputReg   = 4
	ReadHoldingReg = 3
)

type Operation struct {
	FuncCode  uint8
	OpCode    uint16
	ReadLen   uint16
	IEC61850  Measurement
	Splitter  Splitter     `json:"-"`
	Transform RTUTransform `json:"-"`
}

type SplitResult struct {
	OpCode   uint16
	IEC61850 Measurement
	Value    float64
}

type Splitter func(b []byte) []SplitResult

type MeterState uint8

const (
	AVAILABLE   MeterState = iota // The device responds (initial state)
	UNAVAILABLE                   // The device does not respond
)

func (ms MeterState) String() string {
	if ms == AVAILABLE {
		return "available"
	} else {
		return "unavailable"
	}
}

type Meter struct {
	DeviceId uint8
	Producer Producer
	state    MeterState
	mux      sync.Mutex // syncs the meter state variable
}

// Producer is the interface that produces query snips which represent
// modbus operations
type Producer interface {
	GetMeterType() string
	Produce() []Operation
	Probe() Operation
}

// MeasurementMapping maps measurements to phyiscal registers
type MeasurementMapping struct {
	ops Measurements
}

// Opcode returns physical register for measurement type
func (o *MeasurementMapping) Opcode(iec Measurement) uint16 {
	if opcode, ok := o.ops[iec]; ok {
		return opcode
	}

	log.Fatalf("Undefined opcode for measurement %s", iec.String())
	return 0
}

// NewMeterByType meter factory
func NewMeterByType(typeid string, devid uint8) (*Meter, error) {
	var p Producer
	typeid = strings.ToUpper(typeid)

	switch typeid {
	case METERTYPE_ABB:
		p = NewABBProducer()
	case METERTYPE_SDM:
		p = NewSDMProducer()
	case METERTYPE_INEPRO:
		p = NewIneproProducer()
	case METERTYPE_JANITZA:
		p = NewJanitzaProducer()
	case METERTYPE_DZG:
		p = NewDZGProducer()
	case METERTYPE_SBC:
		p = NewSBCProducer()
	case METERTYPE_SE:
		p = NewSEProducer()
	case METERTYPE_SMA:
		p = NewSMAProducer()
	case METERTYPE_KOSTAL:
		p = NewKostalProducer()
	default:
		return nil, fmt.Errorf("Unknown meter type %s", typeid)
	}

	return NewMeter(devid, p), nil
}

func NewMeter(devid uint8, producer Producer) *Meter {
	return &Meter{
		Producer: producer,
		DeviceId: devid,
		state:    AVAILABLE,
	}
}

func (m *Meter) SetState(newstate MeterState) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.state = newstate
}

func (m *Meter) State() MeterState {
	m.mux.Lock()
	defer m.mux.Unlock()
	return m.state
}
