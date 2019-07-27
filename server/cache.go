package server

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// Cache caches and aggregates meter reasings
type Cache struct {
	sync.Mutex
	readings map[string]*MeterReadings
	maxAge   time.Duration
	status   *Status
	verbose  bool
}

// NewCache creates new meter reading cache
func NewCache(maxAge time.Duration, status *Status, verbose bool) *Cache {
	readings := make(map[string]*MeterReadings)

	cache := &Cache{
		readings: readings,
		maxAge:   maxAge,
		status:   status,
		verbose:  verbose,
	}

	return cache
}

// Run consumes meter readings into snip cache
func (mc *Cache) Run(in <-chan QuerySnip) {
	for snip := range in {
		uniqueID := snip.Device

		// Search corresponding meter
		readings, ok := mc.readings[uniqueID]
		if !ok {
			readings = NewMeterReadings(mc.maxAge)
			mc.Lock()
			mc.readings[uniqueID] = readings
			mc.Unlock()
		}
		readings.Add(snip)
	}
}

// SortedIDs returns the sorted list of cache ids
func (mc *Cache) SortedIDs() []string {
	mc.Lock()
	defer mc.Unlock()

	var keys []string
	for k := range mc.readings {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Current returns the latest set of meter reading
func (mc *Cache) Current(device string) (res *Readings, err error) {
	mc.Lock()
	defer mc.Unlock()

	if readings, ok := mc.readings[device]; ok {
		if mc.status.Online(device) {
			// return a copy
			return readings.Current.Clone(), nil
		}

		return res, fmt.Errorf("device %s is not available", device)
	}

	return res, fmt.Errorf("device %s does not exist", device)
}

// Average returns averaged sets of meter readings
func (mc *Cache) Average(device string) (*Readings, error) {
	mc.Lock()
	defer mc.Unlock()

	if readings, ok := mc.readings[device]; ok {
		if mc.status.Online(device) {
			res := readings.Average(time.Now().Add(-1 * time.Minute))
			return res, nil
		}

		return nil, fmt.Errorf("device %s is not available", device)
	}

	return nil, fmt.Errorf("device %s does not exist", device)
}

// Purge removes accumulated data for specified device
func (mc *Cache) Purge(device string) error {
	mc.Lock()
	defer mc.Unlock()

	if readings, ok := mc.readings[device]; ok {
		readings.Purge()
		return nil
	}

	return fmt.Errorf("device with id %s does not exist", device)
}
