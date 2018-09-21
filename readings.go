package sdm630

import (
	"fmt"
	"sync"
	"time"
)

type MeterReadings struct {
	Historic ReadingSlice
	Current  Readings
	mux      sync.Mutex
}

func NewMeterReadings(devid uint8, maxAge time.Duration) *MeterReadings {
	res := &MeterReadings{
		Historic: ReadingSlice{},
		Current: Readings{
			UniqueId: fmt.Sprintf(UniqueIdFormat, devid),
			DeviceId: devid,
		},
	}
	go func(mr *MeterReadings) {
		for {
			time.Sleep(maxAge)
			mr.mux.Lock()
			mr.Historic = mr.Historic.NotOlderThan(time.Now().Add(-1 * maxAge))
			mr.mux.Unlock()
		}
	}(res)
	return res
}

func (mr *MeterReadings) Purge(devid uint8) {
	mr.mux.Lock()
	defer mr.mux.Unlock()

	mr.Historic = ReadingSlice{}
	mr.Current = Readings{
		UniqueId: fmt.Sprintf(UniqueIdFormat, devid),
		DeviceId: devid,
	}
}

func (mr *MeterReadings) AddSnip(snip QuerySnip) {
	mr.mux.Lock()
	defer mr.mux.Unlock()

	// 1. Merge the snip to the last values.
	reading := mr.Current
	reading.MergeSnip(snip)

	// 2. store it
	mr.Current = reading
	mr.Historic = append(mr.Historic, reading)
}
