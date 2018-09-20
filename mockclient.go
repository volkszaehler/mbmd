package sdm630

import (
	"errors"
	"math/rand"
	"time"
)

type MockClient struct {
	errorRate    int32
	responseTime time.Duration
}

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

func (c *MockClient) ReadInputRegisters(address, quantity uint16) (results []byte, err error) {
	return c.read(quantity)
}

func (c *MockClient) ReadHoldingRegisters(address, quantity uint16) (results []byte, err error) {
	return c.read(quantity)
}

func (c *MockClient) ReadCoils(address, quantity uint16) (results []byte, err error) {
	panic("Not implemented")
}

func (c *MockClient) ReadDiscreteInputs(address, quantity uint16) (results []byte, err error) {
	panic("Not implemented")
}

func (c *MockClient) MaskWriteRegister(address, andMask, orMask uint16) (results []byte, err error) {
	panic("Not implemented")
}

func (c *MockClient) ReadFIFOQueue(address uint16) (results []byte, err error) {
	panic("Not implemented")
}

func (c *MockClient) WriteSingleCoil(address, value uint16) (results []byte, err error) {
	panic("Not implemented")
}

func (c *MockClient) WriteMultipleCoils(address, quantity uint16, value []byte) (results []byte, err error) {
	panic("Not implemented")
}

func (c *MockClient) WriteSingleRegister(address, value uint16) (results []byte, err error) {
	panic("Not implemented")
}

func (c *MockClient) WriteMultipleRegisters(address, quantity uint16, value []byte) (results []byte, err error) {
	panic("Not implemented")
}

func (c *MockClient) ReadWriteMultipleRegisters(readAddress, readQuantity, writeAddress, writeQuantity uint16, value []byte) (results []byte, err error) {
	panic("Not implemented")
}
