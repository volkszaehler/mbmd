package server

import (
	"encoding/json"
	"math"
	"runtime"
	"sync"
	"time"
)

// MemoryStatus represents daemon memory allocation
type MemoryStatus struct {
	Alloc     uint64
	HeapAlloc uint64
}

// ModbusStatus represents device request and error status
type ModbusStatus struct {
	Requests          uint64
	RequestsPerMinute float64
	Errors            uint64
	ErrorsPerMinute   float64
}

// DeviceStatus represents a devices runtime status
type DeviceStatus struct {
	Device string
	Type   string
	Model  string
	Serial string
	Online bool
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

// Status represents the daemon and device status.
// It is updated when marshaled to JSON
type Status struct {
	qe         DeviceInfo
	StartTime  time.Time
	UpTime     float64
	Goroutines int
	Memory     MemoryStatus
	Meters     []DeviceStatus
	mu         sync.RWMutex
	meterMap   map[string]DeviceStatus
}

// NewStatus creates status cache that collects device status from control channel.
// It needs to be Update()d in order to refresh its data for consumption
func NewStatus(qe DeviceInfo, control <-chan ControlSnip) *Status {
	s := &Status{
		qe:         qe,
		Memory:     memoryStatus(),
		Goroutines: runtime.NumGoroutine(),
		StartTime:  time.Now(),
		UpTime:     1,
		meterMap:   make(map[string]DeviceStatus),
	}

	go func() {
		for c := range control {
			s.mu.Lock()

			minutes := s.UpTime / 60

			desc := s.qe.DeviceDescriptorByID(c.Device)
			s.meterMap[c.Device] = DeviceStatus{
				Device: c.Device,
				Type:   desc.Manufacturer,
				Model:  desc.Model,
				Serial: desc.Serial,
				Online: c.Status.Online,
				ModbusStatus: ModbusStatus{
					Requests:          c.Status.Requests,
					Errors:            c.Status.Errors,
					ErrorsPerMinute:   math.Round(float64(c.Status.Errors)/minutes*1000) / 1000,
					RequestsPerMinute: math.Round(float64(c.Status.Requests)/minutes*1000) / 1000,
				},
			}

			s.mu.Unlock()
		}
	}()

	return s
}

// Online returns device's online status or false if the device does not exist
func (s *Status) Online(device string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if ds, ok := s.meterMap[device]; ok {
		return ds.Online
	}

	return false
}

// Update status
func (s *Status) update() {
	s.Memory = memoryStatus()
	s.Goroutines = runtime.NumGoroutine()
	s.UpTime = time.Since(s.StartTime).Seconds()

	s.Meters = make([]DeviceStatus, 0, len(s.meterMap))
	for _, ms := range s.meterMap {
		s.Meters = append(s.Meters, ms)
	}
}

// MarshalJSON will syncronize access to the status object
// see http://choly.ca/post/go-json-marshalling/ for avoiding infinite loop
func (s *Status) MarshalJSON() ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.update()

	type Alias Status
	return json.Marshal(&struct{ *Alias }{Alias: (*Alias)(s)})
}
