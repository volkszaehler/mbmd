package bus

import (
	"fmt"
	"log"

	"github.com/volkszaehler/mbmd/meters"
)

type manager struct {
	devices map[uint8]meters.Device
	bus     Bus
}

// NewManager creates a new bus manager instance. Bus managers operate devices on a bus instance
func NewManager(bus Bus) manager {
	m := manager{
		bus:     bus,
		devices: make(map[uint8]meters.Device, 1),
	}
	return m
}

func (m *manager) Add(id uint8, device meters.Device) error {
	if _, ok := m.devices[id]; ok {
		return fmt.Errorf("duplicate device id %d", id)
	}

	m.devices[id] = device
	return nil
}

func (m *manager) Run() {
	for id, device := range m.devices {
		m.bus.Slave(id)
		if results, err := device.Query(m.bus.(*TCP).Client); err != nil {
			log.Fatal(err)
		} else {
			log.Println(results)
		}
	}
}
