package sdm630

import (
	"runtime"
	"time"

	. "github.com/gonium/gosdm630/internal/meters"
)

type MemoryStatus struct {
	Alloc     uint64
	HeapAlloc uint64
}

type ModbusStatus struct {
	TotalRequests        uint64
	RequestRatePerMinute float64
	TotalErrors          uint64
	ErrorRatePerMinute   float64
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
			TotalRequests:        0,
			RequestRatePerMinute: 0,
			TotalErrors:          0,
			ErrorRatePerMinute:   0,
		},
		ConfiguredMeters: nil,
		metermap:         metermap,
	}
}

func (s *Status) IncreaseRequestCounter() {
	s.Modbus.TotalRequests = s.Modbus.TotalRequests + 1
}

func (s *Status) IncreaseReconnectCounter() {
	s.Modbus.TotalErrors = s.Modbus.TotalErrors + 1
}

func (s *Status) Update() {
	s.Memory = CurrentMemoryStatus()
	s.Goroutines = runtime.NumGoroutine()
	s.UptimeSeconds = time.Since(s.Starttime).Seconds()
	s.Modbus.ErrorRatePerMinute =
		float64(s.Modbus.TotalErrors) / (s.UptimeSeconds / 60)
	s.Modbus.RequestRatePerMinute =
		float64(s.Modbus.TotalRequests) / (s.UptimeSeconds / 60)
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
