package sdm630

import (
	"fmt"
	"time"
)

type MeterReadings struct {
	Historic ReadingSlice
	Current  Readings
}

func NewMeterReadings(devid uint8, maxAge time.Duration) *MeterReadings {
	res := &MeterReadings{
		Historic: ReadingSlice{},
		Current: Readings{
			UniqueId: fmt.Sprintf(UniqueIdFormat, devid),
			DeviceId: devid,
		},
	}
	go func() {
		for {
			time.Sleep(maxAge)
			res.Historic = res.Historic.NotOlderThan(
				time.Now().Add(-1 * maxAge))
		}
	}()
	return res
}

func (mr *MeterReadings) Purge(devid uint8) {
	mr.Historic = ReadingSlice{}
	mr.Current = Readings{
		UniqueId: fmt.Sprintf(UniqueIdFormat, devid),
		DeviceId: devid,
	}
}

func (mr *MeterReadings) AddSnip(snip QuerySnip) {
	// 1. Merge the snip to the last values.
	reading := mr.Current
	reading.MergeSnip(snip)
	// 2. store it
	mr.Current = reading
	mr.Historic = append(mr.Historic, reading)
}
