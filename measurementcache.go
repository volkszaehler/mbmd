package sdm630

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	. "github.com/gonium/gosdm630/internal/meters"
)

type MeasurementCache struct {
	meters  map[uint8]MeasurementCacheItem
	maxAge  time.Duration
	verbose bool
}

type MeasurementCacheItem struct {
	*Meter
	*MeterReadings
}

func NewMeasurementCache(
	meters map[uint8]*Meter,
	scheduler *MeterScheduler,
	maxAge time.Duration,
	isVerbose bool,
) *MeasurementCache {
	items := make(map[uint8]MeasurementCacheItem)

	for _, meter := range meters {
		items[meter.DeviceId] = MeasurementCacheItem{
			meter,
			NewMeterReadings(meter.DeviceId, maxAge),
		}
	}

	cache := &MeasurementCache{
		meters:  items,
		maxAge:  maxAge,
		verbose: isVerbose,
	}

	scheduler.SetCache(cache)
	return cache
}

// Run consumes meter readings into snip cache
func (mc *MeasurementCache) Run(in QuerySnipChannel) {
	for snip := range in {
		devid := snip.DeviceId
		// Search corresponding meter
		if meter, ok := mc.meters[devid]; ok {
			// add the snip to the cache
			meter.AddSnip(snip)
			if mc.verbose {
				// log.Printf("%s\r\n", meter.Current.String())
			}
		} else {
			log.Fatalf("Snip for unknown meter received - this should not happen (%v).", snip)
		}
	}
}

// Purge removes accumulated data for specified device
func (mc *MeasurementCache) Purge(deviceId byte) error {
	if meter, ok := mc.meters[deviceId]; ok {
		meter.Purge(deviceId)
		return nil
	}

	return fmt.Errorf("Device with id %d does not exist.", deviceId)
}

func (mc *MeasurementCache) GetSortedIDs() []byte {
	var keys ByteSlice
	for k, _ := range mc.meters {
		keys = append(keys, k)
	}
	sort.Sort(keys)
	return keys
}

func (mc *MeasurementCache) GetCurrent(deviceId byte) (*Readings, error) {
	if meter, ok := mc.meters[deviceId]; ok {
		if meter.GetState() == AVAILABLE {
			return &meter.Current, nil
		}
		return nil, fmt.Errorf("Device %d is not available.", deviceId)
	}
	return nil, fmt.Errorf("Device %d does not exist.", deviceId)
}

func (mc *MeasurementCache) GetMinuteAvg(deviceId byte) (*Readings, error) {
	if meter, ok := mc.meters[deviceId]; ok {
		if meter.GetState() == AVAILABLE {
			measurements := meter.Historic
			lastminute := measurements.NotOlderThan(time.Now().Add(-1 * time.Minute))

			res, err := lastminute.Average()
			if err != nil {
				return nil, err
			}

			if mc.verbose {
				log.Printf("Averaging over %d measurements:\r\n%s\r\n",
					len(measurements), res.String())
			}
			return res, nil
		}
		return nil, fmt.Errorf("Device %d is not available.", deviceId)
	}
	return nil, fmt.Errorf("Device %d does not exist.", deviceId)
}

// ByteSlice attaches the methods of sort.Interface to []byte, sorting in increasing order.
type ByteSlice []byte

func (s ByteSlice) Len() int           { return len(s) }
func (s ByteSlice) Less(i, j int) bool { return s[i] < s[j] }
func (s ByteSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

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
