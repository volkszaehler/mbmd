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
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1Voltage, Value: float32(math.NaN())}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2Voltage, Value: float32(math.NaN())}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3Voltage, Value: float32(math.NaN())}

			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1Current, Value: float32(math.NaN())}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2Current, Value: float32(math.NaN())}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3Current, Value: float32(math.NaN())}

			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1Cosphi, Value: float32(math.NaN())}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2Cosphi, Value: float32(math.NaN())}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3Cosphi, Value: float32(math.NaN())}

			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1Import, Value: float32(math.NaN())}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2Import, Value: float32(math.NaN())}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3Import, Value: float32(math.NaN())}

			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1Export, Value: float32(math.NaN())}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2Export, Value: float32(math.NaN())}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3Export, Value: float32(math.NaN())}
		}
	}
}
