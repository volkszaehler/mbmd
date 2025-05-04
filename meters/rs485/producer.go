package rs485

import (
	"fmt"

	"github.com/grid-x/modbus"
	"github.com/volkszaehler/mbmd/meters"
)

// Operation describes a physical bus operation and its result
type Operation struct {
	FuncCode  uint8
	OpCode    uint16
	ReadLen   uint16
	IEC61850  meters.Measurement
	Transform RTUTransform
}

// Producer is the interface that produces query snips which represent
// modbus operations
type Producer interface {
	// Type returns device description, typically static
	Description() string

	// Produce creates a slice of possible device operations
	Produce() []Operation

	// Produce creates a device operation suited to detect the device during
	// scanning, typically a L1 voltage read operation
	Probe() Operation

	// Allow a device to do special things for initialization,
	// like detecting a subtype for the EM24 from Carlo Gavazzi
	Initialize(client modbus.Client)
}

// Opcodes map measurements to physical registers
type Opcodes map[meters.Measurement]uint16

// Opcode returns physical register for measurement type
func (o *Opcodes) Opcode(iec meters.Measurement) uint16 {
	if opcode, ok := (*o)[iec]; ok {
		return opcode
	}

	panic(fmt.Sprintf("Undefined opcode for measurement %s", iec.String()))
}
