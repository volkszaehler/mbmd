package sdm630

import (
	"github.com/zfjagann/golang-ring"
	"log"
)

type MeasurementCache struct {
	datastream          ReadingChannel
	secsBetweenReadings int
	lastminutebuffer    *ring.Ring
	lastreading         Readings
	verbose             bool
}

func NewMeasurementCache(ds ReadingChannel, interval int, isVerbose bool) *MeasurementCache {
	r := &ring.Ring{}
	r.SetCapacity(60 / interval)
	return &MeasurementCache{
		datastream:          ds,
		secsBetweenReadings: interval,
		lastminutebuffer:    r,
		verbose:             isVerbose,
	}
}

func (mc *MeasurementCache) ConsumeData() {
	for {
		reading := <-mc.datastream
		mc.lastreading = reading
		mc.lastminutebuffer.Enqueue(reading)
		if mc.verbose {
			log.Printf("%s\r\n", &mc.lastreading)
		}
	}
}

func (mc *MeasurementCache) GetLast() Readings {
	return mc.lastreading
}

func (mc *MeasurementCache) GetMinuteAvg() Readings {
	measurements := mc.lastminutebuffer.Values()
	var avg Readings
	for _, m := range measurements {
		r, _ := m.(Readings)
		avg = r.add(&avg)
	}
	if mc.verbose {
		log.Printf("%s\r\n", avg.String())
	}
	return avg.divide(float32(len(measurements)))
}
