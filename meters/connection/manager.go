package connection

import (
	"fmt"
	"log"

	"github.com/volkszaehler/mbmd/meters"
)

type Manager struct {
	Devices map[uint8]meters.Device
	Conn    Connection
}

// NewManager creates a new connection manager instance. connection managers operate devices on a connection instance
func NewManager(conn Connection) Manager {
	m := Manager{
		Devices: make(map[uint8]meters.Device, 1),
		Conn:    conn,
	}
	return m
}

func (m *Manager) Add(id uint8, device meters.Device) error {
	if _, ok := m.Devices[id]; ok {
		return fmt.Errorf("duplicate device id %d", id)
	}

	m.Devices[id] = device
	return nil
}

func (m *Manager) All(cb func(uint8, meters.Device)) {
	for id, device := range m.Devices {
		m.Conn.Slave(id)
		cb(id, device)
	}
}

func (m *Manager) Run() {
	for id, device := range m.Devices {
		m.Conn.Slave(id)

		if results, err := device.Query(m.Conn.ModbusClient()); err != nil {
			log.Fatal(err)
		} else {
			log.Println(results)
		}
	}
}
