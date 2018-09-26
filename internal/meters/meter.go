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
	IEC61850  string
	Transform RTUTransform
}

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

// NewMeterByType meter factory
func NewMeterByType(typeid string, devid uint8) (*Meter, error) {
	var p Producer
	typeid = strings.ToUpper(typeid)

	switch typeid {
	case METERTYPE_SDM:
		p = NewSDMProducer()
	case METERTYPE_JANITZA:
		p = NewJanitzaProducer()
	case METERTYPE_DZG:
		log.Println(`WARNING: The DZG DVH 4013 does not report the same
		measurements as the other meters. Only limited functionality is
		implemented.`)
		p = NewDZGProducer()
	case METERTYPE_SBC:
		log.Println(`WARNING: The SBC ALE3 does not report the same
		measurements as the other meters. Only limited functionality is
		implemented.`)
		p = NewSBCProducer()
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

func (m *Meter) UpdateState(newstate MeterState) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.state = newstate
}

func (m *Meter) GetState() MeterState {
	m.mux.Lock()
	defer m.mux.Unlock()
	return m.state
}
