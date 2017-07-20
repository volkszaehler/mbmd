package sdm630

import (
	"fmt"
	"sync"
	"time"
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
	Type          MeterType
	DeviceId      uint8
	Scheduler     Scheduler
	MeterReadings *MeterReadings
	state         MeterState
	mux           sync.Mutex // syncs the meter state variable
}

func NewMeter(
	typeid MeterType,
	devid uint8,
	scheduler Scheduler,
) *Meter {
	r := NewMeterReadings(devid, DEFAULT_METER_STORE_SECONDS)
	return &Meter{
		Type:          typeid,
		Scheduler:     scheduler,
		DeviceId:      devid,
		MeterReadings: r,
		state:         METERSTATE_AVAILABLE,
	}
}

func (m *Meter) UpdateState(newstate MeterState) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.state = newstate
	if newstate == METERSTATE_UNAVAILABLE {
		m.MeterReadings.Purge(m.DeviceId)
	}
}

func (m *Meter) GetState() MeterState {
	m.mux.Lock()
	defer m.mux.Unlock()
	return m.state
}

func (m *Meter) AddSnip(snip QuerySnip) {
	m.MeterReadings.AddSnip(snip)
}

type MeterReadings struct {
	lastminutereadings ReadingSlice
	lastreading        Readings
}

func NewMeterReadings(devid uint8, secondsToStore time.Duration) (retval *MeterReadings) {
	reading := Readings{
		UniqueId:       fmt.Sprintf(UniqueIdFormat, devid),
		ModbusDeviceId: devid,
	}
	retval = &MeterReadings{
		lastminutereadings: ReadingSlice{},
		lastreading:        reading,
	}
	go func() {
		for {
			time.Sleep(time.Minute * 1)
			//before := len(retval.lastminutereadings)
			retval.lastminutereadings =
				retval.lastminutereadings.NotOlderThan(time.Now().Add(-1 *
					secondsToStore))
			//after := len(retval.lastminutereadings)
			//if isVerbose {
			//	log.Printf("Cache cleanup: Before %d, after %d", before, after)
			//}
		}
	}()
	return retval
}

func (mr *MeterReadings) Purge(devid uint8) {
	mr.lastminutereadings = ReadingSlice{}
	mr.lastreading = Readings{
		UniqueId:       fmt.Sprintf(UniqueIdFormat, devid),
		ModbusDeviceId: devid,
	}
}

func (mr *MeterReadings) AddSnip(snip QuerySnip) {
	// 1. Merge the snip to the last values.
	reading := mr.lastreading
	reading.MergeSnip(snip)
	// 2. store it
	mr.lastreading = reading
	mr.lastminutereadings = append(mr.lastminutereadings, reading)
}
