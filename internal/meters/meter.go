package sdm630

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

type MeterState uint8

const (
	METERSTATE_AVAILABLE   MeterState = iota // The device responds (initial state)
	METERSTATE_UNAVAILABLE        // The device does not respond
)

func (ms MeterState) String() string {
	if ms == METERSTATE_AVAILABLE {
		return "available"
	} else {
		return "unavailable"
	}
}

type Meter struct {
	DeviceId      uint8
	Producer      Producer
	MeterReadings *MeterReadings
	state         MeterState
	mux           sync.Mutex // syncs the meter state variable
}

// Producer is the interface that produces query snips which represent
// modbus operations
type Producer interface {
	GetMeterType() string
	Produce(devid uint8) []QuerySnip
	Probe(devid uint8) QuerySnip
}

func NewMeterByType(
	typeid string,
	devid uint8,
	timeToCacheReadings time.Duration,
) (*Meter, error) {
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

	return NewMeter(devid, p, timeToCacheReadings), nil
}

func NewMeter(
	devid uint8,
	producer Producer,
	timeToCacheReadings time.Duration,
) *Meter {
	r := NewMeterReadings(devid, timeToCacheReadings)
	return &Meter{
		Producer:      producer,
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
	Lastminutereadings ReadingSlice
	Lastreading        Readings
}

func NewMeterReadings(devid uint8, secondsToStore time.Duration) (retval *MeterReadings) {
	reading := Readings{
		UniqueId:       fmt.Sprintf(UniqueIdFormat, devid),
		DeviceId: devid,
	}
	retval = &MeterReadings{
		Lastminutereadings: ReadingSlice{},
		Lastreading:        reading,
	}
	go func() {
		for {
			time.Sleep(secondsToStore)
			//before := len(retval.lastminutereadings)
			retval.Lastminutereadings =
				retval.Lastminutereadings.NotOlderThan(time.Now().Add(-1 *
					secondsToStore))
			//after := len(retval.lastminutereadings)
			//fmt.Printf("Cache cleanup: Before %d, after %d\r\n", before, after)
		}
	}()
	return retval
}

func (mr *MeterReadings) Purge(devid uint8) {
	mr.Lastminutereadings = ReadingSlice{}
	mr.Lastreading = Readings{
		UniqueId:       fmt.Sprintf(UniqueIdFormat, devid),
		DeviceId: devid,
	}
}

func (mr *MeterReadings) AddSnip(snip QuerySnip) {
	// 1. Merge the snip to the last values.
	reading := mr.Lastreading
	reading.MergeSnip(snip)
	// 2. store it
	mr.Lastreading = reading
	mr.Lastminutereadings = append(mr.Lastminutereadings, reading)
}
