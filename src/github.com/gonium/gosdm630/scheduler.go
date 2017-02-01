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
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1Voltage,
				Value: math.NaN(), Unit: "L1 Voltage (V)"}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2Voltage,
				Value: math.NaN(), Unit: "L2 Voltage (V)"}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3Voltage,
				Value: math.NaN(), Unit: "L3 Voltage (V)"}

			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1Current,
				Value: math.NaN(), Unit: "L1 Current (A)"}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2Current,
				Value: math.NaN(), Unit: "L2 Current (A)"}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3Current,
				Value: math.NaN(), Unit: "L3 Current (A)"}

			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1Power, Value: math.NaN(), Unit: "L1 Power (W)"}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2Power, Value: math.NaN(), Unit: "L2 Power (W)"}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3Power, Value: math.NaN(), Unit: "L3 Power (W)"}

			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1Cosphi, Value: math.NaN(), Unit: "L1 Cosphi"}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2Cosphi, Value: math.NaN(), Unit: "L2 Cosphi"}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3Cosphi, Value: math.NaN(), Unit: "L3 Cosphi"}

			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1Import, Value: math.NaN(), Unit: "L1 Import (kWh)"}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2Import, Value: math.NaN(), Unit: "L2 Import (kWh)"}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3Import, Value: math.NaN(), Unit: "L3 Import (kWh)"}

			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1Export, Value: math.NaN(), Unit: "L1 Export (kWh)"}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2Export, Value: math.NaN(), Unit: "L2 Export (kWh)"}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3Export, Value: math.NaN(), Unit: "L3 Export (kWh)"}
		}
	}
}
