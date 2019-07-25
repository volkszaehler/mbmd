package server

import (
	"errors"
	"fmt"
	"math"
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

func (r *Readings) fp2f(key meters.Measurement) float64 {
	if v, ok := r.Values[key]; ok {
		return v
	}
	return math.NaN()
}

func (r *Readings) String() string {
	fmtString := "" +
		"L1: %.1fV %.2fA %.0fW %.2fcos | " +
		"L2: %.1fV %.2fA %.0fW %.2fcos | " +
		"L3: %.1fV %.2fA %.0fW %.2fcos | " +
		"%.1fHz"
	return fmt.Sprintf(fmtString,
		r.fp2f(meters.VoltageL1),
		r.fp2f(meters.CurrentL1),
		r.fp2f(meters.PowerL1),
		r.fp2f(meters.CosphiL1),
		r.fp2f(meters.VoltageL2),
		r.fp2f(meters.CurrentL2),
		r.fp2f(meters.PowerL2),
		r.fp2f(meters.CosphiL2),
		r.fp2f(meters.VoltageL3),
		r.fp2f(meters.CurrentL3),
		r.fp2f(meters.PowerL3),
		r.fp2f(meters.CosphiL3),
		r.fp2f(meters.Frequency),
	)
}

// After is true if the reading is older than the given timestamp.
func (r *Readings) After(ts time.Time) (retval bool) {
	return r.Timestamp.After(ts)
}

// Add two readings. The individual values are added except for
// time- the latter of the two times is copied over to the result
func (r *Readings) add(rhs *Readings) (*Readings, error) {
	res := &Readings{
		Timestamp: r.Timestamp,
		Values:    make(map[meters.Measurement]float64),
	}

	for m, rhsv := range rhs.Values {
		if lhsv, ok := r.Values[m]; ok {
			res.Values[m] = lhsv + rhsv
		}
	}

	if r.Timestamp.Before(rhs.Timestamp) {
		res.Timestamp = rhs.Timestamp
	}

	return res, nil
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
func (r *Readings) Clone() Readings {
	r.Lock()
	defer r.Unlock()

	res := Readings{
		Timestamp: r.Timestamp,
		Values:    make(map[meters.Measurement]float64, len(r.Values)),
	}

	for k, v := range r.Values {
		res.Values[k] = v
	}

	return res
}

// ReadingSlice is a type alias for a slice of readings.
type ReadingSlice []Readings

// After creates a new ReadingSlice of latest data
func (r ReadingSlice) After(ts time.Time) ReadingSlice {
	res := ReadingSlice{}
	for _, reading := range r {
		if reading.After(ts) {
			res = append(res, reading)
		}
	}
	return res
}

// Average calculates average across a ReadingSlice.
// It is assumed that each set of readings is fully populated.
func (r *ReadingSlice) Average() (avg *Readings, err error) {
	for idx, r := range *r {
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

	if len(*r) == 0 {
		return nil, errors.New("readings empty")
	}

	return avg.divide(float64(len(*r))), nil
}
