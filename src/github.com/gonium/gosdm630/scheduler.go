package sdm630

import (
	"math"
)

type Scheduler interface {
	Produce()
}

type RoundRobinScheduler struct {
	out    QuerySnipChannel
	devids []uint8
}

func NewRoundRobinScheduler(
	out QuerySnipChannel,
	devices []uint8,
) *RoundRobinScheduler {
	return &RoundRobinScheduler{
		out:    out,
		devids: devices,
	}
}

func (s *RoundRobinScheduler) Produce() {
	for {
		for _, devid := range s.devids {
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1Voltage, Value: math.NaN()}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2Voltage, Value: math.NaN()}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3Voltage, Value: math.NaN()}

			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1Current, Value: math.NaN()}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2Current, Value: math.NaN()}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3Current, Value: math.NaN()}

			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1Power, Value: math.NaN()}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2Power, Value: math.NaN()}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3Power, Value: math.NaN()}

			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1Cosphi, Value: math.NaN()}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2Cosphi, Value: math.NaN()}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3Cosphi, Value: math.NaN()}

			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1Import, Value: math.NaN()}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2Import, Value: math.NaN()}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3Import, Value: math.NaN()}

			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1Export, Value: math.NaN()}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2Export, Value: math.NaN()}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3Export, Value: math.NaN()}
		}
	}
}
