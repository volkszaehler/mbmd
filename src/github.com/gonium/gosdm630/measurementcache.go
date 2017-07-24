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
				log.Printf("%s\r\n", meter.MeterReadings.Lastreading.String())
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
		if meter.GetState() == METERSTATE_AVAILABLE {
			return &meter.MeterReadings.Lastreading, nil
		} else {
			return nil, fmt.Errorf("Meter %d is not available.", id)
		}
	} else {
		return nil, fmt.Errorf("No device with id %d available.", id)
	}
}

func (mc *MeasurementCache) GetMinuteAvg(id byte) (*Readings, error) {
	if meter, ok := mc.meters[id]; !ok {
		return nil, fmt.Errorf("No device with id %d available.", id)
	} else {
		if meter.GetState() == METERSTATE_AVAILABLE {
			measurements := meter.MeterReadings.Lastminutereadings
			lastminute := measurements.NotOlderThan(time.Now().Add(-1 *
				time.Minute))
			var err error
			var avg Readings
			for idx, r := range lastminute {
				if idx == 0 {
					// This is the first element - initialize our accumulator
					avg = r
				} else {
					avg, err = r.add(&avg)
					if err != nil {
						return nil, err
					}
				}
			}
			retval := avg.divide(float64(len(lastminute)))
			if mc.verbose {
				log.Printf("Averaging over %d measurements:\r\n%s\r\n",
					len(measurements), retval.String())
			}
			return &retval, nil
		} else { // !METERSTATE_AVAILABLE
			return nil, fmt.Errorf("Meter %d is not available.", id)
		}
	}
}

// Helper for dealing with Modbus device ids (bytes).
// ByteSlice attaches the methods of sort.Interface to []byte, sorting in increasing order.
type ByteSlice []byte

func (s ByteSlice) Len() int           { return len(s) }
func (s ByteSlice) Less(i, j int) bool { return s[i] < s[j] }
func (s ByteSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
