package sdm630

import (
	"fmt"
	"log"
	"sort"
	"time"

	. "github.com/gonium/gosdm630/internal/meters"
)

type MeasurementCache struct {
	in      QuerySnipChannel
	items   map[uint8]MeasurementCacheItem
	maxAge  time.Duration
	verbose bool
}

type MeasurementCacheItem struct {
	meter    *Meter
	readings *MeterReadings
}

func NewMeasurementCache(
	meters map[uint8]*Meter,
	inChannel QuerySnipChannel,
	scheduler *MeterScheduler,
	maxAge time.Duration,
	isVerbose bool,
) *MeasurementCache {
	items := make(map[uint8]MeasurementCacheItem)

	for _, meter := range meters {
		items[meter.DeviceId] = MeasurementCacheItem{
			meter:    meter,
			readings: NewMeterReadings(meter.DeviceId, maxAge),
		}
	}

	cache := &MeasurementCache{
		in:      inChannel,
		items:   items,
		maxAge:  maxAge,
		verbose: isVerbose,
	}

	scheduler.SetCache(cache)
	return cache
}

func (mc *MeasurementCache) Consume() {
	for {
		snip := <-mc.in
		devid := snip.DeviceId
		// Search corresponding meter
		if item, ok := mc.items[devid]; ok {
			// add the snip to the cache
			item.readings.AddSnip(snip)
			if mc.verbose {
				log.Printf("%s\r\n", item.readings.Current)
			}
		} else {
			log.Fatal("Snip for unknown meter received - this should not happen.")
		}
	}
}

func (mc *MeasurementCache) Purge(deviceId byte) error {
	if item, ok := mc.items[deviceId]; ok {
		item.readings.Purge(deviceId)
		return nil
	} else {
		return fmt.Errorf("No device with id %d available.", deviceId)
	}
}

func (mc *MeasurementCache) GetSortedIDs() []byte {
	var keys ByteSlice
	for k, _ := range mc.items {
		keys = append(keys, k)
	}
	sort.Sort(keys)
	return keys
}

func (mc *MeasurementCache) GetLast(deviceId byte) (*Readings, error) {
	if item, ok := mc.items[deviceId]; ok {
		if item.meter.GetState() == AVAILABLE {
			return &item.readings.Current, nil
		} else {
			return nil, fmt.Errorf("Meter %d is not available.", deviceId)
		}
	} else {
		return nil, fmt.Errorf("No device with id %d available.", deviceId)
	}
}

func average(readings ReadingSlice) (*Readings, error) {
	var avg *Readings
	var err error

	for idx, r := range readings {
		if idx == 0 {
			// This is the first element - initialize our accumulator
			avg = &r
		} else {
			avg, err = r.add(avg)
			if err != nil {
				return nil, err
			}
		}
	}

	res := avg.divide(float64(len(readings)))
	return res, nil
}

func (mc *MeasurementCache) GetMinuteAvg(deviceId byte) (*Readings, error) {
	if item, ok := mc.items[deviceId]; ok {
		if item.meter.GetState() == AVAILABLE {
			measurements := item.readings.Historic
			lastminute := measurements.NotOlderThan(time.Now().Add(-1 * time.Minute))

			res, err := average(lastminute)
			if err != nil {
				return nil, err
			}

			if mc.verbose {
				log.Printf("Averaging over %d measurements:\r\n%s\r\n",
					len(measurements), res.String())
			}
			return res, nil
		}
		return nil, fmt.Errorf("Meter %d is not available.", deviceId)
	} else {
		return nil, fmt.Errorf("No device with id %d available.", deviceId)
	}
}

// ByteSlice attaches the methods of sort.Interface to []byte, sorting in increasing order.
type ByteSlice []byte

func (s ByteSlice) Len() int           { return len(s) }
func (s ByteSlice) Less(i, j int) bool { return s[i] < s[j] }
func (s ByteSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
