package sdm630

import (
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"time"

	. "github.com/gonium/gosdm630/meters"
)

// QuerySnip represents modbus query operations
type QuerySnip struct {
	DeviceId      uint8
	Operation     `json:"-"`
	Value         float64
	ReadTimestamp time.Time
}

func NewQuerySnip(deviceId uint8, operation Operation) QuerySnip {
	snip := QuerySnip{
		DeviceId:  deviceId,
		Operation: operation,
		Value:     math.NaN(),
	}
	return snip
}

// String representation
func (q *QuerySnip) String() string {
	return fmt.Sprintf("DevID: %d, FunCode: %d, Opcode: %x, IEC: %s, Value: %.3f",
		q.DeviceId, q.FuncCode, q.OpCode, q.IEC61850, q.Value)
}

// MarshalJSON converts QuerySnip to json, replacing ReadTimestamp with unix time representation
func (q *QuerySnip) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		DeviceId    uint8
		Value       float64
		IEC61850    string
		Description string
		Timestamp   int64
	}{
		DeviceId:    q.DeviceId,
		Value:       q.Value,
		IEC61850:    q.IEC61850.String(),
		Description: q.IEC61850.Description(),
		Timestamp:   q.ReadTimestamp.UnixNano() / 1e6,
	})
}

type QuerySnipChannel chan QuerySnip

// QuerySnipBroadcaster acts as hub for broadcating QuerySnips
// to multiple recipients
type QuerySnipBroadcaster struct {
	in         QuerySnipChannel
	recipients []QuerySnipChannel
	done       chan bool
	mux        sync.Mutex // guard recipients
	wg         sync.WaitGroup
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
		b.mux.Lock()
		for _, recipient := range b.recipients {
			recipient <- s
		}
		b.mux.Unlock()
	}
	b.stop()
}

// Done returns a channel signalling when broadcasting has stopped
func (b *QuerySnipBroadcaster) Done() <-chan bool {
	return b.done
}

// stop closes broadcast receiver channels and waits for run methods to finish
func (b *QuerySnipBroadcaster) stop() {
	b.mux.Lock()
	defer b.mux.Unlock()
	for _, recipient := range b.recipients {
		close(recipient)
	}
	b.wg.Wait()
	b.done <- true
}

// attach creates and attaches a QuerySnipChannel to the broadcaster
func (b *QuerySnipBroadcaster) attach() QuerySnipChannel {
	channel := make(QuerySnipChannel)

	b.mux.Lock()
	b.recipients = append(b.recipients, channel)
	b.mux.Unlock()

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
	Type     ControlSnipType
	Message  string
	DeviceId uint8
}

type ControlSnipType uint8

const (
	CONTROLSNIP_OK ControlSnipType = iota
	CONTROLSNIP_ERROR
)

type ControlSnipChannel chan ControlSnip
