package connection

import (
	"fmt"

	"github.com/volkszaehler/mbmd/meters"
)

// Manager handles devices attached to a connection
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

// Add adds device to the device manager at specified device id
func (m *Manager) Add(id uint8, device meters.Device) error {
	if _, ok := m.Devices[id]; ok {
		return fmt.Errorf("duplicate device id %d", id)
	}

	m.Devices[id] = device
	return nil
}

// All iterates over all devices and executes the callback per device.
// Before the callback, the slave id is set on the underlying connection.
func (m *Manager) All(cb func(uint8, meters.Device)) {
	for id, device := range m.Devices {
		m.Conn.Slave(id)
		cb(id, device)
	}
}
