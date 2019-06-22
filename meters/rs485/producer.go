package rs485

import (
	// "fmt"
	"log"
	// "strings"

	 "github.com/volkszaehler/mbmd/meters"
)

type Operation struct {
	FuncCode  uint8
	OpCode    uint16
	ReadLen   uint16
	IEC61850  meters.Measurement
	Transform RTUTransform `json:"-"`
}

// type MeterState uint8

// const (
// 	AVAILABLE   MeterState = iota // The device responds (initial state)
// 	UNAVAILABLE                   // The device does not respond
// )

// func (ms MeterState) String() string {
// 	if ms == AVAILABLE {
// 		return "available"
// 	} else {
// 		return "unavailable"
// 	}
// }

// type Meter struct {
// 	DeviceId uint8
// 	Producer Producer
// 	state    MeterState
// 	mux      sync.Mutex // syncs the meter state variable
// }

// Producer is the interface that produces query snips which represent
// modbus operations
type Producer interface {
	Type() string
	Description() string
	Produce() []Operation
	Probe() Operation
}


// Opcodes map measurements to phyiscal registers
type Opcodes map[meters.Measurement]uint16

// Opcode returns physical register for measurement type
func (o *Opcodes) Opcode(iec meters.Measurement) uint16 {
	if opcode, ok := (*o)[iec]; ok {
		return opcode
	}

	log.Fatalf("Undefined opcode for measurement %s", iec.String())
	return 0
}

// NewProducerByType meter factory
// func NewProducerByType(typeid string, devid uint8) (*Meter, error) {
// 	typeid = strings.ToUpper(typeid)

// 	f, ok := Producers[typeid]
// 	if !ok {
// 		return nil, fmt.Errorf("Unknown meter type %s", typeid)
// 	}

// 	return NewMeter(devid, f()), nil
// }

// func NewMeter(devid uint8, producer Producer) *Meter {
// 	return &Meter{
// 		Producer: producer,
// 		DeviceId: devid,
// 		state:    AVAILABLE,
// 	}
// }

// func (m *Meter) SetState(newstate MeterState) {
// 	m.mux.Lock()
// 	defer m.mux.Unlock()
// 	m.state = newstate
// }

// func (m *Meter) State() MeterState {
// 	m.mux.Lock()
// 	defer m.mux.Unlock()
// 	return m.state
// }
