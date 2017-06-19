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
	return &RoundRobinScheduler{out: out,
		devids: devices,
	}
}

func (s *RoundRobinScheduler) Produce() {
	for {
		for _, devid := range s.devids {
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1Voltage,
				Value: math.NaN(), Description: "L1 Voltage (V)", IEC61850: "VolLocPhsA"}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2Voltage,
				Value: math.NaN(), Description: "L2 Voltage (V)", IEC61850: "VolLocPhsB"}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3Voltage,
				Value: math.NaN(), Description: "L3 Voltage (V)", IEC61850: "VolLocPhsC"}

			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1Current,
				Value: math.NaN(), Description: "L1 Current (A)", IEC61850: "AmpLocPhsA"}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2Current,
				Value: math.NaN(), Description: "L2 Current (A)", IEC61850: "AmpLocPhsB"}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3Current,
				Value: math.NaN(), Description: "L3 Current (A)", IEC61850: "AmpLocPhsC"}

			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1Power, Value: math.NaN(), Description: "L1 Power (W)", IEC61850: "WLocPhsA"}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2Power, Value: math.NaN(), Description: "L2 Power (W)", IEC61850: "WLocPhsB"}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3Power, Value: math.NaN(), Description: "L3 Power (W)", IEC61850: "WLocPhsC"}

			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1Cosphi, Value: math.NaN(), Description: "L1 Cosphi", IEC61850: "AngLocPhsA"}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2Cosphi, Value: math.NaN(), Description: "L2 Cosphi", IEC61850: "AngLocPhsB"}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3Cosphi, Value: math.NaN(), Description: "L3 Cosphi", IEC61850: "AngLocPhsC"}

			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1Import, Value: math.NaN(), Description: "L1 Import (kWh)", IEC61850: "TotkWhImportPhsA"}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2Import, Value: math.NaN(), Description: "L2 Import (kWh)",
				IEC61850: "TotkWhImportPhsB"}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3Import, Value: math.NaN(), Description: "L3 Import (kWh)", IEC61850: "TotkWhImportPhsC"}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeTotalImport, Value: math.NaN(), Description: "Total Import (kWh)", IEC61850: "TotkWhImport"}

			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1Export, Value: math.NaN(), Description: "L1 Export (kWh)", IEC61850: "TotkWhExportPhsA"}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2Export, Value: math.NaN(), Description: "L2 Export (kWh)", IEC61850: "TotkWhExportPhsB"}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3Export, Value: math.NaN(), Description: "L3 Export (kWh)", IEC61850: "TotkWhExportPhsC"}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeTotalExport, Value: math.NaN(), Description: "Total Export (kWh)", IEC61850: "TotkWhExport"}

			//	s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1THDCurrent, Value: math.NaN(), Description: "L1 Current THD (%)", IEC61850: "ThdAPhsA"}
			//	s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2THDCurrent, Value: math.NaN(), Description: "L2 Current THD (%)", IEC61850: "ThdAPhsB"}
			//	s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3THDCurrent, Value: math.NaN(), Description: "L3 Current THD (%)", IEC61850: "ThdAPhsC"}
			//	s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeAvgTHDCurrent, Value: math.NaN(), Description: "Average current to neutral THD (%)", IEC61850: "ThdAmp"}

			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1THDVoltageNeutral, Value: math.NaN(), Description: "L1 Voltage to neutral THD (%)", IEC61850: "ThdVolPhsA"}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2THDVoltageNeutral, Value: math.NaN(), Description: "L2 Voltage to neutral THD (%)", IEC61850: "ThdVolPhsB"}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3THDVoltageNeutral, Value: math.NaN(), Description: "L3 Voltage to neutral THD (%)", IEC61850: "ThdVolPhsC"}
			s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeAvgTHDVoltageNeutral, Value: math.NaN(), Description: "Average voltage to neutral THD (%)", IEC61850: "ThdVol"}
		}
	}
}
