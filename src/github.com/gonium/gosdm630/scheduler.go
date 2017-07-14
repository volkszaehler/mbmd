package sdm630

import (
	"math"
)

type QueryScheduler struct {
	out    QuerySnipChannel
	meters []*Meter
}

func NewQueryScheduler(
	out QuerySnipChannel,
	devices []*Meter,
) *QueryScheduler {
	return &QueryScheduler{
		out:    out,
		meters: devices,
	}
}

func (q *QueryScheduler) Run() {
	for {
		for _, meter := range q.meters {
			// TODO: Implement state of meter, skip/probe defective ones.
			// TODO: The probe function can also be used by the sdm-detect program.
			meter.Scheduler.Produce(q.out, meter.DevId)
		}
	}
}

// This is the interface each scheduler must implement.
type Scheduler interface {
	Produce(out QuerySnipChannel, devid uint8)
}

// ####################################################################
// Round-Robin Scheduler for the Eastron SDM Devices
// ####################################################################
type SDMRoundRobinScheduler struct {
}

func NewSDMRoundRobinScheduler() *SDMRoundRobinScheduler {
	return &SDMRoundRobinScheduler{}
}

func (s *SDMRoundRobinScheduler) Produce(out QuerySnipChannel, devid uint8) {
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
		OpCode: OpCodeSDML1Voltage, Value: math.NaN(), Description: "L1 Voltage (V)", IEC61850: "VolLocPhsA"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
		OpCode: OpCodeSDML2Voltage, Value: math.NaN(), Description: "L2 Voltage (V)", IEC61850: "VolLocPhsB"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
		OpCode: OpCodeSDML3Voltage, Value: math.NaN(), Description: "L3 Voltage (V)", IEC61850: "VolLocPhsC"}

	out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
		OpCode: OpCodeSDML1Current, Value: math.NaN(), Description: "L1 Current (A)", IEC61850: "AmpLocPhsA"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
		OpCode: OpCodeSDML2Current, Value: math.NaN(), Description: "L2 Current (A)", IEC61850: "AmpLocPhsB"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
		OpCode: OpCodeSDML3Current, Value: math.NaN(), Description: "L3 Current (A)", IEC61850: "AmpLocPhsC"}

	out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
		OpCode: OpCodeSDML1Power, Value: math.NaN(), Description: "L1 Power (W)", IEC61850: "WLocPhsA"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
		OpCode: OpCodeSDML2Power, Value: math.NaN(), Description: "L2 Power (W)", IEC61850: "WLocPhsB"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
		OpCode: OpCodeSDML3Power, Value: math.NaN(), Description: "L3 Power (W)", IEC61850: "WLocPhsC"}

	out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
		OpCode: OpCodeSDML1Cosphi, Value: math.NaN(), Description: "L1 Cosphi", IEC61850: "AngLocPhsA"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
		OpCode: OpCodeSDML2Cosphi, Value: math.NaN(), Description: "L2 Cosphi", IEC61850: "AngLocPhsB"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
		OpCode: OpCodeSDML3Cosphi, Value: math.NaN(), Description: "L3 Cosphi", IEC61850: "AngLocPhsC"}

	out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
		OpCode: OpCodeSDML1THDVoltageNeutral, Value: math.NaN(), Description: "L1 Voltage to neutral THD (%)", IEC61850: "ThdVolPhsA"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
		OpCode: OpCodeSDML2THDVoltageNeutral, Value: math.NaN(), Description: "L2 Voltage to neutral THD (%)", IEC61850: "ThdVolPhsB"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
		OpCode: OpCodeSDML3THDVoltageNeutral, Value: math.NaN(), Description: "L3 Voltage to neutral THD (%)", IEC61850: "ThdVolPhsC"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
		OpCode: OpCodeSDMAvgTHDVoltageNeutral, Value: math.NaN(), Description: "Average voltage to neutral THD (%)", IEC61850: "ThdVol"}

	out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
		OpCode: OpCodeSDML1Import, Value: math.NaN(), Description: "L1 Import (kWh)", IEC61850: "TotkWhImportPhsA"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
		OpCode: OpCodeSDML2Import, Value: math.NaN(), Description: "L2 Import (kWh)", IEC61850: "TotkWhImportPhsB"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
		OpCode: OpCodeSDML3Import, Value: math.NaN(), Description: "L3 Import (kWh)", IEC61850: "TotkWhImportPhsC"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
		OpCode: OpCodeSDMTotalImport, Value: math.NaN(), Description: "Total Import (kWh)", IEC61850: "TotkWhImport"}

	out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
		OpCode: OpCodeSDML1Export, Value: math.NaN(), Description: "L1 Export (kWh)", IEC61850: "TotkWhExportPhsA"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
		OpCode: OpCodeSDML2Export, Value: math.NaN(), Description: "L2 Export (kWh)", IEC61850: "TotkWhExportPhsB"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
		OpCode: OpCodeSDML3Export, Value: math.NaN(), Description: "L3 Export (kWh)", IEC61850: "TotkWhExportPhsC"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
		OpCode: OpCodeSDMTotalExport, Value: math.NaN(), Description: "Total Export (kWh)", IEC61850: "TotkWhExport"}

	//	s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL1THDCurrent, Value: math.NaN(), Description: "L1 Current THD (%)", IEC61850: "ThdAPhsA"}
	//	s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL2THDCurrent, Value: math.NaN(), Description: "L2 Current THD (%)", IEC61850: "ThdAPhsB"}
	//	s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeL3THDCurrent, Value: math.NaN(), Description: "L3 Current THD (%)", IEC61850: "ThdAPhsC"}
	//	s.out <- QuerySnip{DeviceId: devid, OpCode: OpCodeAvgTHDCurrent, Value: math.NaN(), Description: "Average current to neutral THD (%)", IEC61850: "ThdAmp"}
}

// ####################################################################
// Round-Robin Scheduler for the Janitza B23 DIN-Rail meters
// ####################################################################

type JanitzaRoundRobinScheduler struct {
}

func NewJanitzaRoundRobinScheduler() *JanitzaRoundRobinScheduler {
	return &JanitzaRoundRobinScheduler{}
}

func (s *JanitzaRoundRobinScheduler) Produce(out QuerySnipChannel, devid uint8) {
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeJanitzaL1Voltage, Value: math.NaN(), Description: "L1 Voltage (V)", IEC61850: "VolLocPhsA"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeJanitzaL2Voltage, Value: math.NaN(), Description: "L2 Voltage (V)", IEC61850: "VolLocPhsB"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeJanitzaL3Voltage, Value: math.NaN(), Description: "L3 Voltage (V)", IEC61850: "VolLocPhsC"}

	out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg, OpCode: OpCodeJanitzaL1Current, Value: math.NaN(), Description: "L1 Current (A)", IEC61850: "AmpLocPhsA"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeJanitzaL2Current, Value: math.NaN(), Description: "L2 Current (A)", IEC61850: "AmpLocPhsB"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeJanitzaL3Current, Value: math.NaN(), Description: "L3 Current (A)", IEC61850: "AmpLocPhsC"}

	out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeJanitzaL1Power, Value: math.NaN(), Description: "L1 Power (W)", IEC61850: "WLocPhsA"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeJanitzaL2Power, Value: math.NaN(), Description: "L2 Power (W)", IEC61850: "WLocPhsB"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeJanitzaL3Power, Value: math.NaN(), Description: "L3 Power (W)", IEC61850: "WLocPhsC"}

	out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeJanitzaL1Cosphi, Value: math.NaN(), Description: "L1 Cosphi", IEC61850: "AngLocPhsA"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeJanitzaL2Cosphi, Value: math.NaN(), Description: "L2 Cosphi", IEC61850: "AngLocPhsB"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeJanitzaL3Cosphi, Value: math.NaN(), Description: "L3 Cosphi", IEC61850: "AngLocPhsC"}

	out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeJanitzaL1Import, Value: math.NaN(), Description: "L1 Import (kWh)", IEC61850: "TotkWhImportPhsA"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeJanitzaL2Import, Value: math.NaN(), Description: "L2 Import (kWh)", IEC61850: "TotkWhImportPhsB"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeJanitzaL3Import, Value: math.NaN(), Description: "L3 Import (kWh)", IEC61850: "TotkWhImportPhsC"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeJanitzaTotalImport, Value: math.NaN(), Description: "Total Import (kWh)", IEC61850: "TotkWhImport"}

	out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeJanitzaL1Export, Value: math.NaN(), Description: "L1 Export (kWh)", IEC61850: "TotkWhExportPhsA"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeJanitzaL2Export, Value: math.NaN(), Description: "L2 Export (kWh)", IEC61850: "TotkWhExportPhsB"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeJanitzaL3Export, Value: math.NaN(), Description: "L3 Export (kWh)", IEC61850: "TotkWhExportPhsC"}
	out <- QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeJanitzaTotalExport, Value: math.NaN(), Description: "Total Export (kWh)", IEC61850: "TotkWhExport"}
}
