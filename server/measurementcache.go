package server

import (
	"fmt"
	"sort"
	"sync"
	"time"
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

func NewCache(maxAge time.Duration, verbose bool) *Cache {
	items := make(map[string]CacheItem)

	cache := &Cache{
		items:   items,
		maxAge:  maxAge,
		verbose: verbose,
	}

	return cache
}

// Run consumes meter readings into snip cache
func (mc *Cache) Run(in QuerySnipChannel) {
	for snip := range in {
		uniqueID := snip.Device

		ci, ok := mc.items[uniqueID]
		// Search corresponding meter
		if !ok {
			ci = CacheItem{
				NewMeterReadings(uniqueID, mc.maxAge),
			}
			mc.items[uniqueID] = ci
		}
		ci.AddSnip(snip)
	}
}

// Purge removes accumulated data for specified device
func (mc *Cache) Purge(device string) error {
	if meter, ok := mc.items[device]; ok {
		meter.Purge(device)
		return nil
	}

	return fmt.Errorf("Device with id %d does not exist.", device)
}

func (mc *Cache) SortedIDs() []string {
	var keys []string
	for k := range mc.items {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (mc *Cache) GetCurrent(device string) (*Readings, error) {
	// if meter, ok := mc.items[device]; ok {
	// 	if meter.State() == AVAILABLE {
	// 		return &meter.Current, nil
	// 	}
	// 	return nil, fmt.Errorf("Device %d is not available.", device)
	// }
	return nil, fmt.Errorf("Device %d does not exist.", device)
}

func (mc *Cache) GetMinuteAvg(device string) (*Readings, error) {
	// if meter, ok := mc.items[device]; ok {
	// 	if meter.State() == AVAILABLE {
	// 		measurements := meter.Historic
	// 		lastminute := measurements.NotOlderThan(time.Now().Add(-1 * time.Minute))

	// 		res, err := lastminute.Average()
	// 		if err != nil {
	// 			return nil, err
	// 		}

	// 		if mc.verbose {
	// 			log.Printf("Averaging over %d measurements:\r\n%s\r\n",
	// 				len(measurements), res.String())
	// 		}
	// 		return res, nil
	// 	}
	// 	return nil, fmt.Errorf("Device %d is not available.", device)
	// }
	return nil, fmt.Errorf("Device %d does not exist.", device)
}

// ByteSlice attaches the methods of sort.Interface to []byte, sorting in increasing order.
type ByteSlice []byte

func (s ByteSlice) Len() int           { return len(s) }
func (s ByteSlice) Less(i, j int) bool { return s[i] < s[j] }
func (s ByteSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type MeterReadings struct {
	sync.Mutex
	Historic ReadingSlice
	Current  Readings
}

func NewMeterReadings(device string, maxAge time.Duration) *MeterReadings {
	res := &MeterReadings{
		Historic: ReadingSlice{},
		Current: Readings{
			Device: device,
		},
	}
	go func(mr *MeterReadings) {
		for {
			time.Sleep(maxAge)
			mr.Lock()
			mr.Historic = mr.Historic.NotOlderThan(time.Now().Add(-1 * maxAge))
			mr.Unlock()
		}
	}(res)
	return res
}

func (mr *MeterReadings) Purge(device string) {
	mr.Lock()
	defer mr.Unlock()

	mr.Historic = ReadingSlice{}
	mr.Current = Readings{
		Device: device,
	}
}

func (mr *MeterReadings) AddSnip(snip QuerySnip) {
	mr.Lock()
	defer mr.Unlock()

	// 1. Merge the snip to the last values.
	reading := mr.Current
	reading.MergeSnip(snip)

	// 2. store it
	mr.Current = reading
	mr.Historic = append(mr.Historic, reading)
}
