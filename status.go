package sdm630

import (
	"runtime"
	"time"
)

type MemoryStatus struct {
	Alloc     uint64
	HeapAlloc uint64
}

type ModbusStatus struct {
	TotalModbusRequests        uint64
	ModbusRequestRatePerMinute float64
	TotalModbusErrors          uint64
	ModbusErrorRatePerMinute   float64
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
			TotalModbusRequests:        0,
			ModbusRequestRatePerMinute: 0,
			TotalModbusErrors:          0,
			ModbusErrorRatePerMinute:   0,
		},
		ConfiguredMeters: nil,
		metermap:         metermap,
	}
}

func (s *Status) IncreaseModbusRequestCounter() {
	s.Modbus.TotalModbusRequests = s.Modbus.TotalModbusRequests + 1
}

func (s *Status) IncreaseModbusReconnectCounter() {
	s.Modbus.TotalModbusErrors = s.Modbus.TotalModbusErrors + 1
}

func (s *Status) Update() {
	s.Memory = CurrentMemoryStatus()
	s.Goroutines = runtime.NumGoroutine()
	s.UptimeSeconds = time.Since(s.Starttime).Seconds()
	s.Modbus.ModbusErrorRatePerMinute =
		float64(s.Modbus.TotalModbusErrors) / (s.UptimeSeconds / 60)
	s.Modbus.ModbusRequestRatePerMinute =
		float64(s.Modbus.TotalModbusRequests) / (s.UptimeSeconds / 60)
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
