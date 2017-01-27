package sdm630

import (
	"fmt"
	"log"
	"sort"
	"time"
)

type DeviceReadings struct {
	lastminutereadings ReadingSlice
	lastreading        Readings
}

func NewDeviceReadings(secondsToStore time.Duration, isVerbose bool) (retval *DeviceReadings) {
	retval = &DeviceReadings{
		lastminutereadings: ReadingSlice{},
	}
	go func() {
		for {
			time.Sleep(time.Minute * 1)
			before := len(retval.lastminutereadings)
			retval.lastminutereadings =
				retval.lastminutereadings.NotOlderThan(time.Now().Add(-1 *
					secondsToStore))
			after := len(retval.lastminutereadings)
			if isVerbose {
				log.Printf("Cache cleanup: Before %d, after %d", before, after)
			}
		}
	}()
	return retval
}

type MeasurementCache struct {
	datastream     ReadingChannel
	deviceReadings map[uint8]*DeviceReadings
	secondsToStore time.Duration
	verbose        bool
}

func NewMeasurementCache(ds ReadingChannel, secondsToStore time.Duration, isVerbose bool) *MeasurementCache {
	return &MeasurementCache{
		datastream:     ds,
		deviceReadings: make(map[uint8]*DeviceReadings),
		secondsToStore: secondsToStore,
		verbose:        isVerbose,
	}
}

func (mc *MeasurementCache) Consume() {
	for {
		reading := <-mc.datastream
		devid := reading.ModbusDeviceId
		if devreading, ok := mc.deviceReadings[devid]; ok {
			// The device has already a DeviceReadings object
			devreading.lastreading = reading
			devreading.lastminutereadings = append(devreading.lastminutereadings, reading)
		} else {
			// create a new DeviceReadings object
			mc.deviceReadings[devid] = NewDeviceReadings(mc.secondsToStore,
				mc.verbose)
			devreading = mc.deviceReadings[devid]
			devreading.lastreading = reading
			devreading.lastminutereadings = append(devreading.lastminutereadings, reading)
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
		return nil, fmt.Errorf("No device with id %d available.", id)
	}
}

func (mc *MeasurementCache) GetMinuteAvg(id byte) (Readings, error) {
	if _, ok := mc.deviceReadings[id]; !ok {
		return Readings{}, fmt.Errorf("No device with id %d available.", id)
	}
	measurements := mc.deviceReadings[id].lastminutereadings
	lastminute := measurements.NotOlderThan(time.Now().Add(-1 *
		time.Minute))
	avg := Readings{ModbusDeviceId: id}
	for _, r := range lastminute {
		var err error
		avg, err = r.add(&avg)
		if err != nil {
			return avg, err
		}
	}
	retval := avg.divide(float32(len(lastminute)))
	if mc.verbose {
		log.Printf("Averaging over %d measurements:\r\n%s\r\n",
			len(measurements), retval.String())
	}
	return retval, nil
}

// Helper for dealing with Modbus device ids (bytes).
// ByteSlice attaches the methods of sort.Interface to []byte, sorting in increasing order.
type ByteSlice []byte

func (s ByteSlice) Len() int           { return len(s) }
func (s ByteSlice) Less(i, j int) bool { return s[i] < s[j] }
func (s ByteSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
