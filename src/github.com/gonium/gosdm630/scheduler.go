package sdm630

type Scheduler interface {
	Run()
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

func (s *RoundRobinScheduler) Run() {
	for _, devid := range s.devids {
		s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1Voltage}
		s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2Voltage}
		s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3Voltage}

		s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1Current}
		s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2Current}
		s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3Current}

		s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1Cosphi}
		s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2Cosphi}
		s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3Cosphi}

		s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1Import}
		s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2Import}
		s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3Import}

		s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1Export}
		s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2Export}
		s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3Export}
	}
}
