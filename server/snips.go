package server

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/volkszaehler/mbmd/meters"
)

// ControlSnip wraps device status information
type ControlSnip struct {
	Device string
	Status RuntimeInfo
}

// QuerySnip wraps query results
type QuerySnip struct {
	Device string
	meters.MeasurementResult
}

// String representation
func (q *QuerySnip) String() string {
	return fmt.Sprintf("Dev: %s, IEC: %s, Value: %.3f", q.Device, q.Measurement.String(), q.Value)
}

// MarshalJSON converts QuerySnip to json, replacing Timestamp with unix time representation
func (q *QuerySnip) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Device      string
		Value       float64
		IEC61850    string
		Description string
		Timestamp   int64
	}{
		Device:      q.Device,
		Value:       q.Value,
		IEC61850:    q.Measurement.String(),
		Description: q.Measurement.Description(),
		Timestamp:   q.Timestamp.UnixNano() / 1e6,
	})
}

// NewSnipRunner adapts a chan QuerySnip to chan interface
func NewSnipRunner(run func(c <-chan QuerySnip)) func(c <-chan interface{}) {
	return func(c <-chan interface{}) {
		out := make(chan QuerySnip)

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			run(out)
			wg.Done()
		}()

		for x := range c {
			if snip, ok := x.(QuerySnip); ok {
				out <- snip
			} else {
				panic("runner: unexpected type")
			}
		}

		close(out)
		wg.Wait()
	}
}

// NewControlRunner adapts a chan ControlSnip to chan interface
func NewControlRunner(run func(c <-chan ControlSnip)) func(c <-chan interface{}) {
	return func(c <-chan interface{}) {
		out := make(chan ControlSnip)

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			run(out)
			wg.Done()
		}()

		for x := range c {
			if snip, ok := x.(ControlSnip); ok {
				out <- snip
			} else {
				panic("runner: unexpected type")
			}
		}

		close(out)
		wg.Wait()
	}
}

// FromSnipChannel adapts a chan QuerySnip to chan interface
func FromSnipChannel(in <-chan QuerySnip) <-chan interface{} {
	out := make(chan interface{})

	go func(in <-chan QuerySnip, out chan<- interface{}) {
		for snip := range in {
			out <- snip
		}
		close(out)
	}(in, out)

	return out
}

// FromControlChannel adapts a chan ControlSnip to chan interface
func FromControlChannel(in <-chan ControlSnip) <-chan interface{} {
	out := make(chan interface{})

	go func(in <-chan ControlSnip, out chan<- interface{}) {
		for snip := range in {
			out <- snip
		}
		close(out)
	}(in, out)

	return out
}

// ToControlChannel adapts a chan interface to chan ControlSnip
func ToControlChannel(in <-chan interface{}) <-chan ControlSnip {
	out := make(chan ControlSnip)

	go func(in <-chan interface{}, out chan<- ControlSnip) {
		for x := range in {
			if snip, ok := x.(ControlSnip); ok {
				out <- snip
			} else {
				panic("runner: unexpected type")
			}
		}
		close(out)
	}(in, out)

	return out
}
