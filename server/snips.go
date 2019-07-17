package server

import (
	"encoding/json"
	"fmt"
	"sync"

	. "github.com/volkszaehler/mbmd/meters"
)

// QuerySnip represents modbus query operations
type QuerySnip struct {
	Device string
	MeasurementResult
}

// String representation
func (q *QuerySnip) String() string {
	return fmt.Sprintf("Dev: %s, IEC: %s, Value: %.3f",
		q.Device, q.Measurement.String(), q.Value)
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

type QuerySnipChannel chan QuerySnip

// QuerySnipBroadcaster acts as hub for broadcating QuerySnips
// to multiple recipients
type QuerySnipBroadcaster struct {
	sync.Mutex // guard recipients
	wg         sync.WaitGroup
	in         QuerySnipChannel
	recipients []QuerySnipChannel
	done       chan bool
}

// NewQuerySnipBroadcaster creates QuerySnipBroadcaster
func NewQuerySnipBroadcaster(in QuerySnipChannel) *QuerySnipBroadcaster {
	return &QuerySnipBroadcaster{
		in:         in,
		recipients: make([]QuerySnipChannel, 0),
		done:       make(chan bool),
	}
}

// Run executes the broadcaster
func (b *QuerySnipBroadcaster) Run() {
	for s := range b.in {
		b.Lock()
		for _, recipient := range b.recipients {
			recipient <- s
		}
		b.Unlock()
	}
	b.stop()
}

// Done returns a channel signalling when broadcasting has stopped
func (b *QuerySnipBroadcaster) Done() <-chan bool {
	return b.done
}

// stop closes broadcast receiver channels and waits for run methods to finish
func (b *QuerySnipBroadcaster) stop() {
	b.Lock()
	defer b.Unlock()
	for _, recipient := range b.recipients {
		close(recipient)
	}
	b.wg.Wait()
	b.done <- true
}

// attach creates and attaches a QuerySnipChannel to the broadcaster
func (b *QuerySnipBroadcaster) attach() QuerySnipChannel {
	channel := make(QuerySnipChannel)

	b.Lock()
	b.recipients = append(b.recipients, channel)
	b.Unlock()

	return channel
}

// AttachRunner attaches a Run method as broadcast receiver and adds it
// to the waitgroup
func (b *QuerySnipBroadcaster) AttachRunner(runner func(QuerySnipChannel)) {
	b.wg.Add(1)
	go func() {
		ch := b.attach()
		runner(ch)
		b.wg.Done()
	}()
}

// ControlSnip wraps control information like query success or failure.
type ControlSnip struct {
	Device  string
	Result  ControlSnipType
	Message string
}

type ControlSnipType uint8

const (
	ok ControlSnipType = iota
	failure
)

type ControlSnipChannel chan ControlSnip
