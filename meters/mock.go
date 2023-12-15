package meters

import (
	"errors"
	"math/rand"
	"time"

	"github.com/grid-x/modbus"
)

const (
	errorRate = 0
)

// Mock mocks a modbus connection
type Mock struct {
	address string
	Client  modbus.Client
}

// NewMock creates a mock modbus client
func NewMock(address string) Connection {
	client := NewMockClient(errorRate)

	b := &Mock{
		address: address,
		Client:  client,
	}

	return b
}

// String returns "simulate" as bus address
func (b *Mock) String() string {
	return "simulate"
}

// ModbusClient returns the mock modbus client
func (b *Mock) ModbusClient() modbus.Client {
	return b.Client
}

// Logger sets a logging instance for physical bus operations
func (b *Mock) Logger(l Logger) {
}

// Slave sets the modbus device id for the following operations
func (b *Mock) Slave(deviceID uint8) {
}

// Timeout sets the modbus timeout
func (b *Mock) Timeout(timeout time.Duration) time.Duration {
	return timeout
}

// ConnectDelay sets the the initial delay after connecting before starting communication
func (b *Mock) ConnectDelay(delay time.Duration) {
	// nop
}

// Close closes the modbus connection.
func (b *Mock) Close() {
}

// MockClient is a mock modbus client for testing that
// is able to simulate devices and errors
type MockClient struct {
	errorRate    int32
	responseTime time.Duration
}

// NewMockClient creates a mock modbus client
func NewMockClient(errorRate int32) *MockClient {
	return &MockClient{
		errorRate:    errorRate,
		responseTime: 10 * time.Millisecond,
	}
}

func (c *MockClient) fail() bool {
	if c.errorRate > 0 {
		return rand.Int31n(100) < c.errorRate
	}
	return false
}

func (c *MockClient) random(quantity uint16) (results []byte, err error) {
	bytes := make([]byte, quantity*2)
	//nolint:staticcheck // SA1019
	rand.Read(bytes)
	return bytes, nil
}

func (c *MockClient) read(quantity uint16) (results []byte, err error) {
	time.Sleep(c.responseTime)
	if c.fail() {
		return nil, errors.New("Failed")
	}
	return c.random(quantity)
}

// ReadInputRegisters implements modbus.Client
func (c *MockClient) ReadInputRegisters(address, quantity uint16) (results []byte, err error) {
	return c.read(quantity)
}

// ReadHoldingRegisters implements modbus.Client
func (c *MockClient) ReadHoldingRegisters(address, quantity uint16) (results []byte, err error) {
	return c.read(quantity)
}

// ReadCoils implements modbus.Client
func (c *MockClient) ReadCoils(address, quantity uint16) (results []byte, err error) {
	panic("Not implemented")
}

// ReadDiscreteInputs implements modbus.Client
func (c *MockClient) ReadDiscreteInputs(address, quantity uint16) (results []byte, err error) {
	panic("Not implemented")
}

// MaskWriteRegister implements modbus.Client
func (c *MockClient) MaskWriteRegister(address, andMask, orMask uint16) (results []byte, err error) {
	panic("Not implemented")
}

// ReadFIFOQueue implements modbus.Client
func (c *MockClient) ReadFIFOQueue(address uint16) (results []byte, err error) {
	panic("Not implemented")
}

// WriteSingleCoil implements modbus.Client
func (c *MockClient) WriteSingleCoil(address, value uint16) (results []byte, err error) {
	panic("Not implemented")
}

// WriteMultipleCoils implements modbus.Client
func (c *MockClient) WriteMultipleCoils(address, quantity uint16, value []byte) (results []byte, err error) {
	panic("Not implemented")
}

// WriteSingleRegister implements modbus.Client
func (c *MockClient) WriteSingleRegister(address, value uint16) (results []byte, err error) {
	panic("Not implemented")
}

// WriteMultipleRegisters implements modbus.Client
func (c *MockClient) WriteMultipleRegisters(address, quantity uint16, value []byte) (results []byte, err error) {
	panic("Not implemented")
}

// ReadWriteMultipleRegisters implements modbus.Client
func (c *MockClient) ReadWriteMultipleRegisters(readAddress, readQuantity, writeAddress, writeQuantity uint16, value []byte) (results []byte, err error) {
	panic("Not implemented")
}
