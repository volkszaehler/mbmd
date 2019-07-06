package server

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/volkszaehler/mbmd/meters"
)

type Cache struct {
	items   map[string]CacheItem
	maxAge  time.Duration
	verbose bool
}

type CacheItem struct {
	// *meters.Device
	*MeterReadings
}

func NewCache(
	devices map[string]meters.Device,
	scheduler *MeterScheduler,
	maxAge time.Duration,
	isVerbose bool,
) *Cache {
	items := make(map[string]CacheItem)

	for deviceID := range devices {
		items[deviceID] = CacheItem{
			NewMeterReadings(deviceID, maxAge),
		}
	}

	cache := &Cache{
		items:   items,
		maxAge:  maxAge,
		verbose: isVerbose,
	}

	scheduler.SetCache(cache)
	return cache
}

// Run consumes meter readings into snip cache
func (mc *Cache) Run(in QuerySnipChannel) {
	for snip := range in {
		devid := snip.Device
		// Search corresponding meter
		if ci, ok := mc.items[devid]; ok {
			// add the snip to the cache
			ci.AddSnip(snip)
		} else {
			log.Fatalf("Snip for unknown meter received - this should not happen (%v).", snip)
		}
	}
}

// Purge removes accumulated data for specified device
func (mc *Cache) Purge(deviceID string) error {
	if meter, ok := mc.items[deviceID]; ok {
		meter.Purge(deviceID)
		return nil
	}

	return fmt.Errorf("Device with id %d does not exist.", deviceID)
}

func (mc *Cache) SortedIDs() []byte {
	var keys ByteSlice
	for k := range mc.items {
		keys = append(keys, k)
	}
	sort.Sort(keys)
	return keys
}

func (mc *Cache) GetCurrent(deviceID string) (*Readings, error) {
	if meter, ok := mc.items[deviceID]; ok {
		if meter.State() == AVAILABLE {
			return &meter.Current, nil
		}
		return nil, fmt.Errorf("Device %d is not available.", deviceID)
	}
	return nil, fmt.Errorf("Device %d does not exist.", deviceID)
}

func (mc *Cache) GetMinuteAvg(deviceID string) (*Readings, error) {
	if meter, ok := mc.items[deviceID]; ok {
		if meter.State() == AVAILABLE {
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
		return nil, fmt.Errorf("Device %d is not available.", deviceID)
	}
	return nil, fmt.Errorf("Device %d does not exist.", deviceID)
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

func NewMeterReadings(deviceID string, maxAge time.Duration) *MeterReadings {
	res := &MeterReadings{
		Historic: ReadingSlice{},
		Current: Readings{
			DeviceID: deviceID,
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

func (mr *MeterReadings) Purge(deviceID string) {
	mr.mux.Lock()
	defer mr.mux.Unlock()

	mr.Historic = ReadingSlice{}
	mr.Current = Readings{
		DeviceID: deviceID,
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
