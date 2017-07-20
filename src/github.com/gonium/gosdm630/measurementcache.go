package sdm630

import (
	"fmt"
	"log"
	"sort"
	"time"
)

type MeasurementCache struct {
	datastream     QuerySnipChannel
	meters         map[uint8]*Meter
	secondsToStore time.Duration
	verbose        bool
}

func NewMeasurementCache(meters map[uint8]*Meter, ds QuerySnipChannel, secondsToStore time.Duration, isVerbose bool) *MeasurementCache {
	return &MeasurementCache{
		datastream:     ds,
		meters:         meters,
		secondsToStore: secondsToStore,
		verbose:        isVerbose,
	}
}

func (mc *MeasurementCache) Consume() {
	for {
		snip := <-mc.datastream
		devid := snip.DeviceId
		// Search corresponding meter
		if meter, ok := mc.meters[devid]; ok {
			// add the snip to the meter's cache
			meter.AddSnip(snip)
			if mc.verbose {
				log.Printf("%s\r\n", meter.MeterReadings.lastreading.String())
			}
		} else {
			log.Fatal("Snip for unknown meter received - this should not happen.")
		}
	}
}

func (mc *MeasurementCache) GetSortedIDs() []byte {
	var keys ByteSlice
	for k, _ := range mc.meters {
		keys = append(keys, k)
	}
	sort.Sort(keys)
	return keys
}

func (mc *MeasurementCache) GetLast(id byte) (*Readings, error) {
	if meter, ok := mc.meters[id]; ok {
		return &meter.MeterReadings.lastreading, nil
	} else {
		return nil, fmt.Errorf("No device with id %d available.", id)
	}
}

func (mc *MeasurementCache) GetMinuteAvg(id byte) (Readings, error) {
	if _, ok := mc.meters[id]; !ok {
		return Readings{}, fmt.Errorf("No device with id %d available.", id)
	}
	measurements := mc.meters[id].MeterReadings.lastminutereadings
	lastminute := measurements.NotOlderThan(time.Now().Add(-1 *
		time.Minute))
	avg := Readings{UniqueId: fmt.Sprintf(UniqueIdFormat, id), ModbusDeviceId: id}
	for _, r := range lastminute {
		var err error
		avg, err = r.add(&avg)
		if err != nil {
			return avg, err
		}
	}
	retval := avg.divide(float64(len(lastminute)))
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
