package server

import (
	"encoding/json"
	"runtime"
	"sync"
	"time"
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

type MeterStatus struct {
	Device string
	Type   string
	Status string
	ModbusStatus
}

func memoryStatus() MemoryStatus {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return MemoryStatus{
		Alloc:     mem.Alloc,
		HeapAlloc: mem.HeapAlloc,
	}
}

type Status struct {
	sync.Mutex
	StartTime  time.Time
	UpTime     float64
	Goroutines int
	Memory     MemoryStatus
	Meters     []MeterStatus
	meterMap   map[string]MeterStatus
}

// NewStatus creates status cache that collects device status from control channel.
// It needs to be Update()d in order to refresh its data for consumption
func NewStatus(control ControlSnipChannel) *Status {
	s := &Status{
		Memory:     memoryStatus(),
		Goroutines: runtime.NumGoroutine(),
		StartTime:  time.Now(),
		UpTime:     1,
		meterMap:   make(map[string]MeterStatus),
	}

	go func() {
		for c := range control {
			s.Lock()

			var status string
			if c.Status.Online {
				status = "online"
			} else {
				status = "offline"
			}

			mbs := ModbusStatus{
				Requests:          c.Status.Requests,
				Errors:            c.Status.Errors,
				ErrorsPerMinute:   float64(c.Status.Errors) / (s.UpTime / 60),
				RequestsPerMinute: float64(c.Status.Requests) / (s.UpTime / 60),
			}

			ms := MeterStatus{
				Device: c.Device,
				// Type:         dev.Descriptor().Manufacturer,
				Status:       status,
				ModbusStatus: mbs,
			}
			s.meterMap[c.Device] = ms

			s.Unlock()
		}
	}()

	return s
}

// Update status
func (s *Status) update() {
	s.Memory = memoryStatus()
	s.Goroutines = runtime.NumGoroutine()
	s.UpTime = time.Since(s.StartTime).Seconds()

	s.Meters = make([]MeterStatus, 0)
	for _, ms := range s.meterMap {
		s.Meters = append(s.Meters, ms)
	}
}

// MarshalJSON will syncronize access to the status object
// see http://choly.ca/post/go-json-marshalling/ for avoiding infinite loop
func (s *Status) MarshalJSON() ([]byte, error) {
	s.Lock()
	defer s.Unlock()

	s.update()

	type Alias Status
	return json.Marshal(&struct{ *Alias }{Alias: (*Alias)(s)})
}
