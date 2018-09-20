package sdm630

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

type MeterType string
type MeterState uint8

const (
	METERSTATE_AVAILABLE   = iota // The device responds (initial state)
	METERSTATE_UNAVAILABLE        // The device does not respond
)

type Meter struct {
	Type          MeterType
	DeviceId      uint8
	Producer      Producer
	MeterReadings *MeterReadings
	state         MeterState
	mux           sync.Mutex // syncs the meter state variable
}

// Producer is the interface that produces query snips which represent
// modbus operations
type Producer interface {
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
	default:
		return nil, fmt.Errorf("Unknown meter type %s", typeid)
	}

	return NewMeter(MeterType(typeid), devid, p, timeToCacheReadings), nil
}

func NewMeter(
	typeid MeterType,
	devid uint8,
	producer Producer,
	timeToCacheReadings time.Duration,
) *Meter {
	r := NewMeterReadings(devid, timeToCacheReadings)
	return &Meter{
		Type:          typeid,
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

func (m *Meter) GetReadableState() string {
	var retval string
	switch m.GetState() {
	case METERSTATE_AVAILABLE:
		retval = "available"
	case METERSTATE_UNAVAILABLE:
		retval = "unavailable"
	default:
		log.Fatal("Unknown meter state, aborting.")
	}
	return retval
}

func (m *Meter) GetMeterType() MeterType {
	return m.Type
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
		ModbusDeviceId: devid,
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
		ModbusDeviceId: devid,
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
