package sdm630

import (
	"math"
)

// This is the interface each scheduler must implement.
type Scheduler interface {
	Produce()
}

// ####################################################################
// Round-Robin Scheduler for the Eastron SDM Devices
// ####################################################################
type SDMRoundRobinScheduler struct {
	out    QuerySnipChannel
	devids []uint8
}

func NewSDMRoundRobinScheduler(
	out QuerySnipChannel,
	devices []uint8,
) *SDMRoundRobinScheduler {
	return &SDMRoundRobinScheduler{out: out,
		devids: devices,
	}
}

func (s *SDMRoundRobinScheduler) Produce() {
	for {
		for _, devid := range s.devids {
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg, OpCode: OpCodeSDML1Voltage, Value: math.NaN(), Description: "L1 Voltage (V)", IEC61850: "VolLocPhsA"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
				OpCode: OpCodeSDML2Voltage, Value: math.NaN(), Description: "L2 Voltage (V)", IEC61850: "VolLocPhsB"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
				OpCode: OpCodeSDML3Voltage, Value: math.NaN(), Description: "L3 Voltage (V)", IEC61850: "VolLocPhsC"}

			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
				OpCode: OpCodeSDML1Current, Value: math.NaN(), Description: "L1 Current (A)", IEC61850: "AmpLocPhsA"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
				OpCode: OpCodeSDML2Current, Value: math.NaN(), Description: "L2 Current (A)", IEC61850: "AmpLocPhsB"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
				OpCode: OpCodeSDML3Current, Value: math.NaN(), Description: "L3 Current (A)", IEC61850: "AmpLocPhsC"}

			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
				OpCode: OpCodeSDML1Power, Value: math.NaN(), Description: "L1 Power (W)", IEC61850: "WLocPhsA"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
				OpCode: OpCodeSDML2Power, Value: math.NaN(), Description: "L2 Power (W)", IEC61850: "WLocPhsB"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
				OpCode: OpCodeSDML3Power, Value: math.NaN(), Description: "L3 Power (W)", IEC61850: "WLocPhsC"}

			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
				OpCode: OpCodeSDML1Cosphi, Value: math.NaN(), Description: "L1 Cosphi", IEC61850: "AngLocPhsA"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
				OpCode: OpCodeSDML2Cosphi, Value: math.NaN(), Description: "L2 Cosphi", IEC61850: "AngLocPhsB"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
				OpCode: OpCodeSDML3Cosphi, Value: math.NaN(), Description: "L3 Cosphi", IEC61850: "AngLocPhsC"}

			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
				OpCode: OpCodeSDML1THDVoltageNeutral, Value: math.NaN(), Description: "L1 Voltage to neutral THD (%)", IEC61850: "ThdVolPhsA"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
				OpCode: OpCodeSDML2THDVoltageNeutral, Value: math.NaN(), Description: "L2 Voltage to neutral THD (%)", IEC61850: "ThdVolPhsB"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
				OpCode: OpCodeSDML3THDVoltageNeutral, Value: math.NaN(), Description: "L3 Voltage to neutral THD (%)", IEC61850: "ThdVolPhsC"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
				OpCode: OpCodeSDMAvgTHDVoltageNeutral, Value: math.NaN(), Description: "Average voltage to neutral THD (%)", IEC61850: "ThdVol"}

			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
				OpCode: OpCodeSDML1Import, Value: math.NaN(), Description: "L1 Import (kWh)", IEC61850: "TotkWhImportPhsA"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
				OpCode: OpCodeSDML2Import, Value: math.NaN(), Description: "L2 Import (kWh)", IEC61850: "TotkWhImportPhsB"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
				OpCode: OpCodeSDML3Import, Value: math.NaN(), Description: "L3 Import (kWh)", IEC61850: "TotkWhImportPhsC"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
				OpCode: OpCodeSDMTotalImport, Value: math.NaN(), Description: "Total Import (kWh)", IEC61850: "TotkWhImport"}

			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
				OpCode: OpCodeSDML1Export, Value: math.NaN(), Description: "L1 Export (kWh)", IEC61850: "TotkWhExportPhsA"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
				OpCode: OpCodeSDML2Export, Value: math.NaN(), Description: "L2 Export (kWh)", IEC61850: "TotkWhExportPhsB"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
				OpCode: OpCodeSDML3Export, Value: math.NaN(), Description: "L3 Export (kWh)", IEC61850: "TotkWhExportPhsC"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
				OpCode: OpCodeSDMTotalExport, Value: math.NaN(), Description: "Total Export (kWh)", IEC61850: "TotkWhExport"}

			//	s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1THDCurrent, Value: math.NaN(), Description: "L1 Current THD (%)", IEC61850: "ThdAPhsA"}
			//	s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2THDCurrent, Value: math.NaN(), Description: "L2 Current THD (%)", IEC61850: "ThdAPhsB"}
			//	s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3THDCurrent, Value: math.NaN(), Description: "L3 Current THD (%)", IEC61850: "ThdAPhsC"}
			//	s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeAvgTHDCurrent, Value: math.NaN(), Description: "Average current to neutral THD (%)", IEC61850: "ThdAmp"}
		}
	}
}

// ####################################################################
// Round-Robin Scheduler for the Janitza B23 DIN-Rail meters
// ####################################################################

type JanitzaRoundRobinScheduler struct {
	out    QuerySnipChannel
	devids []uint8
}

func NewJanitzaRoundRobinScheduler(
	out QuerySnipChannel,
	devices []uint8,
) *JanitzaRoundRobinScheduler {
	return &JanitzaRoundRobinScheduler{
		out:    out,
		devids: devices,
	}
}

func (s *JanitzaRoundRobinScheduler) Produce() {
	for {
		for _, devid := range s.devids {
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg, OpCode: OpCodeJanitzaL1Voltage, Value: math.NaN(), Description: "L1 Voltage (V)", IEC61850: "VolLocPhsA"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
				OpCode: OpCodeJanitzaL2Voltage, Value: math.NaN(), Description: "L2 Voltage (V)", IEC61850: "VolLocPhsB"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
				OpCode: OpCodeJanitzaL3Voltage, Value: math.NaN(), Description: "L3 Voltage (V)", IEC61850: "VolLocPhsC"}

			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg, OpCode: OpCodeJanitzaL1Current, Value: math.NaN(), Description: "L1 Current (A)", IEC61850: "AmpLocPhsA"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
				OpCode: OpCodeJanitzaL2Current, Value: math.NaN(), Description: "L2 Current (A)", IEC61850: "AmpLocPhsB"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
				OpCode: OpCodeJanitzaL3Current, Value: math.NaN(), Description: "L3 Current (A)", IEC61850: "AmpLocPhsC"}

			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
				OpCode: OpCodeJanitzaL1Power, Value: math.NaN(), Description: "L1 Power (W)", IEC61850: "WLocPhsA"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
				OpCode: OpCodeJanitzaL2Power, Value: math.NaN(), Description: "L2 Power (W)", IEC61850: "WLocPhsB"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
				OpCode: OpCodeJanitzaL3Power, Value: math.NaN(), Description: "L3 Power (W)", IEC61850: "WLocPhsC"}

			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
				OpCode: OpCodeJanitzaL1Cosphi, Value: math.NaN(), Description: "L1 Cosphi", IEC61850: "AngLocPhsA"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
				OpCode: OpCodeJanitzaL2Cosphi, Value: math.NaN(), Description: "L2 Cosphi", IEC61850: "AngLocPhsB"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
				OpCode: OpCodeJanitzaL3Cosphi, Value: math.NaN(), Description: "L3 Cosphi", IEC61850: "AngLocPhsC"}

			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
				OpCode: OpCodeJanitzaL1Import, Value: math.NaN(), Description: "L1 Import (kWh)", IEC61850: "TotkWhImportPhsA"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
				OpCode: OpCodeJanitzaL2Import, Value: math.NaN(), Description: "L2 Import (kWh)", IEC61850: "TotkWhImportPhsB"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
				OpCode: OpCodeJanitzaL3Import, Value: math.NaN(), Description: "L3 Import (kWh)", IEC61850: "TotkWhImportPhsC"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
				OpCode: OpCodeJanitzaTotalImport, Value: math.NaN(), Description: "Total Import (kWh)", IEC61850: "TotkWhImport"}

			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
				OpCode: OpCodeJanitzaL1Export, Value: math.NaN(), Description: "L1 Export (kWh)", IEC61850: "TotkWhExportPhsA"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
				OpCode: OpCodeJanitzaL2Export, Value: math.NaN(), Description: "L2 Export (kWh)", IEC61850: "TotkWhExportPhsB"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
				OpCode: OpCodeJanitzaL3Export, Value: math.NaN(), Description: "L3 Export (kWh)", IEC61850: "TotkWhExportPhsC"}
			s.out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
				OpCode: OpCodeJanitzaTotalExport, Value: math.NaN(), Description: "Total Export (kWh)", IEC61850: "TotkWhExport"}

		}
	}
}
