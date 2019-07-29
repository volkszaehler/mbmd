package meters

import (
	"fmt"
)

// Manager handles devices attached to a connection
type Manager struct {
	devices map[uint8]Device
	Conn    Connection
}

// NewManager creates a new connection manager instance. connection managers operate devices on a connection instance
func NewManager(conn Connection) Manager {
	m := Manager{
		devices: make(map[uint8]Device, 1),
		Conn:    conn,
	}
	return m
}

// Add adds device to the device manager at specified device id
func (m *Manager) Add(id uint8, device Device) error {
	if _, ok := m.devices[id]; ok {
		return fmt.Errorf("duplicate device id %d", id)
	}

	m.devices[id] = device
	return nil
}

// All iterates over all devices and executes the callback per device.
// Before the callback, the slave id is set on the underlying connection if access is true.
func (m *Manager) All(access bool, cb func(uint8, Device)) {
	for id, device := range m.devices {
		if access {
			m.Conn.Slave(id)
		}
		cb(id, device)
	}
}
