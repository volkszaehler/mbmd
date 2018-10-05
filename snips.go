package sdm630

import (
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"time"

	. "github.com/gonium/gosdm630/internal/meters"
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
		IEC61850:    q.IEC61850,
		Description: GetIecDescription(q.IEC61850),
		Timestamp:   q.ReadTimestamp.UnixNano() / 1e6,
	})
}

type QuerySnipChannel chan QuerySnip

// QuerySnipBroadcaster acts as hub for broadcating QuerySnips
// to multiple recipients
type QuerySnipBroadcaster struct {
	in         QuerySnipChannel
	recipients []QuerySnipChannel
	mux        sync.Mutex // guard recipients
}

// NewQuerySnipBroadcaster creates QuerySnipBroadcaster
func NewQuerySnipBroadcaster(in QuerySnipChannel) *QuerySnipBroadcaster {
	return &QuerySnipBroadcaster{
		in:         in,
		recipients: make([]QuerySnipChannel, 0),
	}
}

// Run executes the broadcaster
func (b *QuerySnipBroadcaster) Run() {
	for {
		s := <-b.in
		b.mux.Lock()
		for _, recipient := range b.recipients {
			recipient <- s
		}
		b.mux.Unlock()
	}
}

// Attach creates and attaches a QuerySnipChannel to the broadcaster
func (b *QuerySnipBroadcaster) Attach() QuerySnipChannel {
	channel := make(QuerySnipChannel)

	b.mux.Lock()
	b.recipients = append(b.recipients, channel)
	b.mux.Unlock()

	return channel
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
