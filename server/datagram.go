package server

import (
	"errors"
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

// After is true if the reading is older than the given timestamp.
func (r *Readings) After(ts time.Time) (retval bool) {
	return r.Timestamp.After(ts)
}

// Add two readings. The individual values are added except for
// time- the latter of the two times is copied over to the result.
// If the right-hand side value does not exist in left-hand side, it
// is ignored, so order matters.
func (r *Readings) add(rhs *Readings) *Readings {
	res := &Readings{
		Timestamp: r.Timestamp,
		Values:    make(map[meters.Measurement]float64),
	}

	for m, rhsv := range rhs.Values {
		lhsv := r.Values[m]
		res.Values[m] = lhsv + rhsv
	}

	if r.Timestamp.Before(rhs.Timestamp) {
		res.Timestamp = rhs.Timestamp
	}

	return res
}

// Divide a reading by an integer. The individual values are divided except
// for time which is copied over to the result
func (r *Readings) divide(scaler float64) *Readings {
	res := &Readings{
		Timestamp: r.Timestamp,
		Values:    make(map[meters.Measurement]float64),
	}

	for m, v := range r.Values {
		res.Values[m] = v / scaler
	}

	return res
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

// ReadingSlice is a type alias for a slice of readings.
type ReadingSlice []*Readings

// After creates a new ReadingSlice of latest data
func (rs ReadingSlice) After(ts time.Time) ReadingSlice {
	res := ReadingSlice{}
	for _, reading := range rs {
		if reading.After(ts) {
			res = append(res, reading)
		}
	}
	return res
}

// Average calculates average across a ReadingSlice.
// It is assumed that each set of readings is fully populated.
func (rs *ReadingSlice) Average() (avg *Readings, err error) {
	for idx := len(*rs) - 1; idx >= 0; idx-- {
		r := (*rs)[idx]

		// This is the first element - initialize our accumulator
		if avg == nil {
			avg = r
		} else {
			avg = r.add(avg)
		}
	}

	if avg == nil {
		return nil, errors.New("readings empty")
	}

	return avg.divide(float64(len(*rs))), nil
}
