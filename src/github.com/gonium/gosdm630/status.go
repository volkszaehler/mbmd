package sdm630

import (
	"encoding/json"
	"io"
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
	Starttime     time.Time
	UptimeSeconds float64
	Goroutines    int
	Memory        MemoryStatus
	Modbus        ModbusStatus
}

func NewStatus() *Status {
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
}

func (s *Status) UpdateAndJSON(w io.Writer) error {
	s.Update()
	return json.NewEncoder(w).Encode(s)
}
