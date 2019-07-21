package server

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"
)

// Cache caches and aggregates meter reasings
type Cache struct {
	items   map[string]*MeterReadings
	maxAge  time.Duration
	status  *Status
	verbose bool
}

// NewCache creates new meter reading cache
func NewCache(maxAge time.Duration, status *Status, verbose bool) *Cache {
	items := make(map[string]*MeterReadings)

	cache := &Cache{
		items:   items,
		maxAge:  maxAge,
		status:  status,
		verbose: verbose,
	}

	return cache
}

// Run consumes meter readings into snip cache
func (mc *Cache) Run(in <-chan QuerySnip) {
	for snip := range in {
		uniqueID := snip.Device

		// Search corresponding meter
		ci, ok := mc.items[uniqueID]
		if !ok {
			ci = NewMeterReadings(uniqueID, mc.maxAge)
			mc.items[uniqueID] = ci
		}
		ci.Add(snip)
	}
}

// Purge removes accumulated data for specified device
func (mc *Cache) Purge(device string) error {
	if meter, ok := mc.items[device]; ok {
		meter.Purge(device)
		return nil
	}

	return fmt.Errorf("device with id %s does not exist", device)
}

// SortedIDs returns the sorted list of cache ids
func (mc *Cache) SortedIDs() []string {
	var keys []string
	for k := range mc.items {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// GetCurrent returns the last meter reading
func (mc *Cache) GetCurrent(device string) (*Readings, error) {
	if readings, ok := mc.items[device]; ok {
		if mc.status.Online(device) {
			readings.Lock()
			defer readings.Unlock()

			// return a copy
			copy := readings.Current
			return &copy, nil
		}

		return nil, fmt.Errorf("device %d is not available", device)
	}

	return nil, fmt.Errorf("device %d does not exist", device)
}

// GetAverage returns averages meter readings
func (mc *Cache) GetAverage(device string) (*Readings, error) {
	if readings, ok := mc.items[device]; ok {
		readings.Lock()
		defer readings.Unlock()

		if mc.status.Online(device) {
			measurements := readings.Historic
			lastminute := measurements.NotOlderThan(time.Now().Add(-1 * time.Minute))

			res, err := lastminute.Average()
			if err != nil {
				return nil, err
			}

			if mc.verbose {
				log.Printf("averaging over %d measurements:\r\n%s\r\n", len(measurements), res.String())
			}
			return res, nil
		}

		return nil, fmt.Errorf("device %d is not available", device)
	}

	return nil, fmt.Errorf("device %d does not exist", device)
}

// ByteSlice attaches the methods of sort.Interface to []byte, sorting in increasing order.
type ByteSlice []byte

func (s ByteSlice) Len() int           { return len(s) }
func (s ByteSlice) Less(i, j int) bool { return s[i] < s[j] }
func (s ByteSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// MeterReadings holds entire sets of current and recent meter readings
type MeterReadings struct {
	sync.Mutex
	Historic ReadingSlice
	Current  Readings
}

// NewMeterReadings container for current and recent meter readings
func NewMeterReadings(device string, maxAge time.Duration) *MeterReadings {
	res := &MeterReadings{
		Historic: ReadingSlice{},
		Current: Readings{},
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

// Purge clears meter readings for specified device
func (mr *MeterReadings) Purge(device string) {
	mr.Lock()
	defer mr.Unlock()

	mr.Historic = ReadingSlice{}
	mr.Current = Readings{}
}

// Add adds a meter reading for specified device
func (mr *MeterReadings) Add(snip QuerySnip) {
	mr.Lock()
	defer mr.Unlock()

	// 1. Merge the snip to the last values.
	reading := mr.Current
	reading.MergeSnip(snip)

	// 2. store it
	mr.Current = reading
	mr.Historic = append(mr.Historic, reading)
}
