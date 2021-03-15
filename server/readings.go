package server

import (
	"fmt"
	"sync"
	"time"

	"github.com/volkszaehler/mbmd/meters"
)

// Readings combines readings of all measurements into one data structure
type Readings struct {
	sync.Mutex
	Timestamp time.Time
	Values    map[meters.Measurement]float64
}

func (r *Readings) f2s(key meters.Measurement, digits int) string {
	if v, ok := r.Values[key]; ok {
		format := fmt.Sprintf("%%.%df", digits)
		return fmt.Sprintf(format, v)
	}
	return "0.0"
}

func (r *Readings) String() string {
	fmtString := "" +
		"L1: %sV %sA %sW %scos | " +
		"L2: %sV %sA %sW %scos | " +
		"L3: %sV %sA %sW %scos | " +
		"%sHz"
	return fmt.Sprintf(fmtString,
		r.f2s(meters.VoltageL1, 0),
		r.f2s(meters.CurrentL1, 1),
		r.f2s(meters.PowerL1, 0),
		r.f2s(meters.CosphiL1, 2),
		r.f2s(meters.VoltageL2, 0),
		r.f2s(meters.CurrentL2, 1),
		r.f2s(meters.PowerL2, 0),
		r.f2s(meters.CosphiL2, 2),
		r.f2s(meters.VoltageL3, 0),
		r.f2s(meters.CurrentL3, 1),
		r.f2s(meters.PowerL3, 0),
		r.f2s(meters.CosphiL3, 2),
		r.f2s(meters.Frequency, 0),
	)
}

// Add adds the values represented by the QuerySnip to the
// Readings and updates the current time stamp
func (r *Readings) Add(q QuerySnip) {
	r.Lock()
	defer r.Unlock()

	r.Timestamp = q.Timestamp

	if r.Values == nil {
		r.Values = make(map[meters.Measurement]float64)
	}

	r.Values[q.Measurement] = q.Value
}

// Clone clones a Readings including its values map
func (r *Readings) Clone() *Readings {
	r.Lock()
	defer r.Unlock()

	res := Readings{
		Timestamp: r.Timestamp,
		Values:    make(map[meters.Measurement]float64, len(r.Values)),
	}

	for k, v := range r.Values {
		res.Values[k] = v
	}

	return &res
}

// MeterReadings holds entire sets of current and recent meter readings for a single device
type MeterReadings struct {
	sync.Mutex
	Current  Readings
	Historic []*Readings
}

// NewMeterReadings container for current and recent meter readings
func NewMeterReadings(maxAge time.Duration) *MeterReadings {
	res := &MeterReadings{
		Current:  Readings{},
		Historic: make([]*Readings, 0),
	}

	// housekeeping
	go func(mr *MeterReadings) {
		for {
			time.Sleep(maxAge)
			mr.TrimBefore(time.Now().Add(-1 * maxAge))
		}
	}(res)

	return res
}

// Add adds a meter reading for specified device
func (mr *MeterReadings) Add(snip QuerySnip) {
	mr.Lock()
	defer mr.Unlock()

	mr.Current.Add(snip)
	mr.Historic = append(mr.Historic, mr.Current.Clone())
}

// Average averages historic readings after given timestamp
func (mr *MeterReadings) Average(timestamp time.Time) *Readings {
	mr.Lock()
	defer mr.Unlock()

	mcv := make(map[meters.Measurement]struct {
		count int
		sum   float64
	})

	for _, r := range mr.Historic {
		if r.Timestamp.Before(timestamp) {
			continue
		}

		// calculate sum and count per measurement
		for k, v := range r.Values {
			if m, ok := mcv[k]; ok {
				mcv[k] = struct {
					count int
					sum   float64
				}{m.count + 1, m.sum + v}
			} else {
				mcv[k] = struct {
					count int
					sum   float64
				}{1, v}
			}
		}
	}

	res := Readings{
		Timestamp: mr.Current.Timestamp,
		Values:    make(map[meters.Measurement]float64, len(mcv)),
	}
	for m, cv := range mcv {
		res.Values[m] = cv.sum / float64(cv.count)
	}

	return &res
}

// TrimBefore removes historic readings older than timestamp
func (mr *MeterReadings) TrimBefore(timestamp time.Time) {
	mr.Lock()
	defer mr.Unlock()

	for idx, r := range mr.Historic {
		// trim everything before first recent timestamp
		if r.Timestamp.After(timestamp) {
			mr.Historic = mr.Historic[idx : len(mr.Historic)-1]
			return
		}
	}
}

// Purge clears meter readings
func (mr *MeterReadings) Purge() {
	mr.Lock()
	defer mr.Unlock()

	mr.Current = Readings{}
	mr.Historic = make([]*Readings, 0)
}
