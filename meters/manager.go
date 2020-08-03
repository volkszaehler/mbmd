package meters

type device struct {
	id  uint8
	dev Device
}

// Manager handles devices attached to a connection
type Manager struct {
	devices []device
	Conn    Connection
}

// NewManager creates a new connection manager instance. connection managers operate devices on a connection instance
func NewManager(conn Connection) Manager {
	m := Manager{
		devices: make([]device, 0),
		Conn:    conn,
	}
	return m
}

// Add adds device to the device manager at specified device id
func (m *Manager) Add(id uint8, dev Device) error {
	device := device{
		id:  id,
		dev: dev,
	}

	m.devices = append(m.devices, device)
	return nil
}

// Count returns the number of devices attached to the connection
func (m *Manager) Count() int {
	return len(m.devices)
}

// All iterates over all devices and executes the callback per device.
// Before the callback, the slave id is set on the underlying connection if access is true.
func (m *Manager) All(cb func(uint8, Device)) {
	for _, device := range m.devices {
		cb(device.id, device.dev)
	}
}
