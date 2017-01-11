package sdm630

import (
	"fmt"
	"github.com/zfjagann/golang-ring"
	"log"
	"sort"
)

type DeviceReadings struct {
	lastminutebuffer *ring.Ring
	lastreading      Readings
}

func NewDeviceReadings(interval int) (retval *DeviceReadings) {
	r := &ring.Ring{}
	r.SetCapacity(60 / interval)
	retval = &DeviceReadings{
		lastminutebuffer: r,
	}
	return retval
}

type MeasurementCache struct {
	datastream          ReadingChannel
	deviceReadings      map[uint8]*DeviceReadings
	secsBetweenReadings int
	verbose             bool
}

func NewMeasurementCache(ds ReadingChannel, interval int, isVerbose bool) *MeasurementCache {
	return &MeasurementCache{
		datastream:          ds,
		deviceReadings:      make(map[uint8]*DeviceReadings),
		secsBetweenReadings: interval,
		verbose:             isVerbose,
	}
}

func (mc *MeasurementCache) ConsumeData() {
	for {
		reading := <-mc.datastream
		devid := reading.ModbusDeviceId
		if devreading, ok := mc.deviceReadings[devid]; ok {
			// The device has already a DeviceReadings object
			devreading.lastreading = reading
			devreading.lastminutebuffer.Enqueue(reading)
		} else {
			// create a new DeviceReadings object
			mc.deviceReadings[devid] = NewDeviceReadings(mc.secsBetweenReadings)
			devreading = mc.deviceReadings[devid]
			devreading.lastreading = reading
			devreading.lastminutebuffer.Enqueue(reading)
		}
		if mc.verbose {
			devreading := mc.deviceReadings[devid]
			log.Printf("%s\r\n", devreading.lastreading.String())
		}
	}
}

func (mc *MeasurementCache) GetSortedIDs() []byte {
	var keys ByteSlice
	for k := range mc.deviceReadings {
		keys = append(keys, k)
	}
	sort.Sort(keys)
	return keys
}

func (mc *MeasurementCache) GetLast(id byte) (*Readings, error) {
	if r, ok := mc.deviceReadings[id]; ok {
		return &r.lastreading, nil
	} else {
		return nil, fmt.Errorf("No reading with id %d available.", id)
	}
}

func (mc *MeasurementCache) GetMinuteAvg() Readings {
	measurements := mc.deviceReadings[1].lastminutebuffer.Values()
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

// Helper for dealing with Modbus device ids (bytes).
// ByteSlice attaches the methods of sort.Interface to []byte, sorting in increasing order.
type ByteSlice []byte

func (s ByteSlice) Len() int           { return len(s) }
func (s ByteSlice) Less(i, j int) bool { return s[i] < s[j] }
func (s ByteSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
