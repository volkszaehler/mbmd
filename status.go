package sdm630

import (
	"encoding/json"
	"runtime"
	"sync"
	"time"

	. "github.com/gonium/gosdm630/internal/meters"
)

type MemoryStatus struct {
	Alloc     uint64
	HeapAlloc uint64
}

type ModbusStatus struct {
	Requests          uint64
	RequestsPerMinute float64
	Errors            uint64
	ErrorsPerMinute   float64
}

func CurrentMemoryStatus() MemoryStatus {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return MemoryStatus{
		Alloc:     mem.Alloc,
		HeapAlloc: mem.HeapAlloc,
	}
}

type Status struct {
	Starttime        time.Time
	UptimeSeconds    float64
	Goroutines       int
	Memory           MemoryStatus
	Modbus           ModbusStatus
	ConfiguredMeters []MeterStatus
	metermap         map[uint8]*Meter
	mux              sync.RWMutex `json:"-"`
}

type MeterStatus struct {
	Id     uint8
	Type   string
	Status string
}

func NewStatus(metermap map[uint8]*Meter) *Status {
	return &Status{
		Memory:        CurrentMemoryStatus(),
		Starttime:     time.Now(),
		Goroutines:    runtime.NumGoroutine(),
		UptimeSeconds: 1,
		Modbus: ModbusStatus{
			Requests:          0,
			RequestsPerMinute: 0,
			Errors:            0,
			ErrorsPerMinute:   0,
		},
		ConfiguredMeters: nil,
		metermap:         metermap,
	}
}

func (s *Status) IncreaseRequestCounter() {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.Modbus.Requests++
}

func (s *Status) IncreaseReconnectCounter() {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.Modbus.Errors++
}

func (s *Status) Update() {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.Memory = CurrentMemoryStatus()
	s.Goroutines = runtime.NumGoroutine()
	s.UptimeSeconds = time.Since(s.Starttime).Seconds()
	s.Modbus.ErrorsPerMinute = float64(s.Modbus.Errors) / (s.UptimeSeconds / 60)
	s.Modbus.RequestsPerMinute = float64(s.Modbus.Requests) / (s.UptimeSeconds / 60)

	var confmeters []MeterStatus
	for id, meter := range s.metermap {
		ms := MeterStatus{
			Id:     id,
			Type:   meter.Producer.GetMeterType(),
			Status: meter.GetState().String(),
		}

		confmeters = append(confmeters, ms)
	}
	s.ConfiguredMeters = confmeters
}

// MarshalJSON will syncronize access to the status object
// see http://choly.ca/post/go-json-marshalling/ for avoiding infinite loop
func (s *Status) MarshalJSON() ([]byte, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	type Alias Status
	return json.Marshal(&struct{ *Alias }{Alias: (*Alias)(s)})
}
