package meters

import (
	"sync"
)

const (
	ReadInputReg   = 4
	ReadHoldingReg = 3
)

type MeterState int

const (
	ONLINE  MeterState = iota // The device responds (initial state)
	OFFLINE                   // The device does not respond
)

func (ms MeterState) String() string {
	if ms == ONLINE {
		return "available"
	}
	return "unavailable"
}

type Meter struct {
	DeviceID   uint8
	Device     Device
	state      MeterState
	sync.Mutex // syncs the meter state variable
}

func NewMeter(devID uint8, device Device) *Meter {
	return &Meter{
		DeviceID: devID,
		Device:   device,
		state:    ONLINE,
	}
}

func (m *Meter) SetState(state MeterState) {
	m.Lock()
	defer m.Unlock()
	m.state = state
}

func (m *Meter) State() MeterState {
	m.Lock()
	defer m.Unlock()
	return m.state
}
