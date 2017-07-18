package sdm630

import (
	"sync"
)

type MeterType string
type MeterState uint8

const (
	METERTYPE_JANITZA = "JANITZA"
	METERTYPE_SDM     = "SDM"
)

const (
	METERSTATE_AVAILABLE   = iota // The device responds (initial state)
	METERSTATE_UNAVAILABLE        // The device does not respond
)

type Meter struct {
	Type      MeterType
	DeviceId  uint8
	Scheduler Scheduler
	state     MeterState
	mux       sync.Mutex // syncs the meter state variable
}

func NewMeter(
	typeid MeterType,
	devid uint8,
	scheduler Scheduler,
) *Meter {
	return &Meter{
		Type:      typeid,
		Scheduler: scheduler,
		DeviceId:  devid,
		state:     METERSTATE_AVAILABLE,
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
