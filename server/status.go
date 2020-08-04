package server

import (
	"encoding/json"
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
	sync.Mutex
	qe         DeviceInfo
	StartTime  time.Time
	UpTime     float64
	Goroutines int
	Memory     MemoryStatus
	Meters     []DeviceStatus
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
			s.Lock()

			minutes := s.UpTime / 60
			mbs := ModbusStatus{
				Requests:          c.Status.Requests,
				Errors:            c.Status.Errors,
				ErrorsPerMinute:   float64(c.Status.Errors) / minutes,
				RequestsPerMinute: float64(c.Status.Requests) / minutes,
			}

			desc := s.qe.DeviceDescriptorByID(c.Device)

			ds := DeviceStatus{
				Device:       c.Device,
				Type:         desc.Manufacturer,
				Online:       c.Status.Online,
				ModbusStatus: mbs,
			}
			s.meterMap[c.Device] = ds

			s.Unlock()
		}
	}()

	return s
}

// Online returns device's online status or false if the device does not exist
func (s *Status) Online(device string) bool {
	s.Lock()
	defer s.Unlock()

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

	s.Meters = make([]DeviceStatus, 0)
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
