package sdm630

import (
	"encoding/json"
	"io"
	"runtime"
	"time"
)

type MemoryStatus struct {
	Alloc      uint64
	TotalAlloc uint64
	HeapAlloc  uint64
	HeapSys    uint64
}

func CurrentMemoryStatus() MemoryStatus {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return MemoryStatus{
		Alloc:      mem.Alloc,
		TotalAlloc: mem.TotalAlloc,
		HeapAlloc:  mem.HeapAlloc,
		HeapSys:    mem.HeapSys,
	}
}

type Status struct {
	Memory           MemoryStatus
	Starttime        time.Time
	Uptime           float64
	ModbusReconnects uint64
}

func NewStatus() *Status {
	return &Status{
		Memory:           CurrentMemoryStatus(),
		Starttime:        time.Now(),
		Uptime:           1,
		ModbusReconnects: 0,
	}
}

func (s *Status) IncreaseModbusReconnectCounter() {
	s.ModbusReconnects = s.ModbusReconnects + 1
}

func (s *Status) Update() {
	s.Memory = CurrentMemoryStatus()
	s.Uptime = time.Since(s.Starttime).Seconds()
}

func (s *Status) UpdateAndJSON(w io.Writer) error {
	s.Update()
	return json.NewEncoder(w).Encode(s)
}
