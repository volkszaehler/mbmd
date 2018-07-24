package sdm630

import (
	"log"
	"math"
	"time"
)

type MeterScheduler struct {
	out     QuerySnipChannel
	control ControlSnipChannel
	meters  map[uint8]*Meter
}

func NewMeterScheduler(
	out QuerySnipChannel,
	control ControlSnipChannel,
	devices map[uint8]*Meter,
) *MeterScheduler {
	return &MeterScheduler{
		out:     out,
		meters:  devices,
		control: control,
	}
}

// SetupScheduler creates a scheduler and its wiring
func SetupScheduler(meters map[uint8]*Meter, qe *ModbusEngine) (*MeterScheduler, QuerySnipChannel) {
	// Create Channels that link the goroutines
	var scheduler2queryengine = make(QuerySnipChannel)
	var queryengine2scheduler = make(ControlSnipChannel)
	var queryengine2tee = make(QuerySnipChannel)

	scheduler := NewMeterScheduler(
		scheduler2queryengine,
		queryengine2scheduler,
		meters,
	)

	go qe.Transform(
		scheduler2queryengine, // input
		queryengine2scheduler, // error
		queryengine2tee,       // output
	)

	return scheduler, queryengine2tee
}

func (q *MeterScheduler) produceSnips(out QuerySnipChannel) {
	for {
		for _, meter := range q.meters {
			sniplist := meter.Scheduler.Produce(meter.DeviceId)
			for _, snip := range sniplist {
				// Check if meter is still valid
				if meter.GetState() != METERSTATE_UNAVAILABLE {
					q.out <- snip
				}
			}
		}
	}
}

func (q *MeterScheduler) supervisor() {
	for {
		for _, meter := range q.meters {
			if meter.GetState() == METERSTATE_UNAVAILABLE {
				log.Printf("Attempting to ping unavailable meter %d", meter.DeviceId)
				// inject probe snip - the re-enabling logic is in Run()
				q.out <- meter.Scheduler.GetProbeSnip(meter.DeviceId)
			}
		}
		time.Sleep(15 * time.Minute)
	}
}

func (q *MeterScheduler) Run() {
	source := make(QuerySnipChannel)
	go q.supervisor()
	go q.produceSnips(source)
	for {
		select {
		case snip := <-source:
			q.out <- snip
		case controlSnip := <-q.control:
			switch controlSnip.Type {
			case CONTROLSNIP_ERROR:
				log.Printf("Failure - deactivating meter %d: %s",
					controlSnip.DeviceId, controlSnip.Message)
				// search meter and deactivate it...
				meter, ok := q.meters[controlSnip.DeviceId]
				if !ok {
					log.Fatal("Internal device id mismatch - this should not happen!")
				} else {
					meter.UpdateState(METERSTATE_UNAVAILABLE)
				}
			case CONTROLSNIP_OK:
				// search meter and reactivate it...
				meter, ok := q.meters[controlSnip.DeviceId]
				if !ok {
					log.Fatal("Internal device id mismatch - this should not happen!")
				} else {
					if meter.GetState() != METERSTATE_AVAILABLE {
						log.Printf("Re-activating meter %d", controlSnip.DeviceId)
						meter.UpdateState(METERSTATE_AVAILABLE)
					}
				}
			default:
				log.Fatal("Received unknown control snip - something weird happened.")
			}
		}
	}
}

// This is the interface each scheduler must implement.
type Scheduler interface {
	Produce(devid uint8) []QuerySnip
	GetProbeSnip(devid uint8) QuerySnip
}

// ####################################################################
// Round-Robin Scheduler for the Eastron SDM Devices
// ####################################################################
type SDMRoundRobinScheduler struct {
}

func NewSDMRoundRobinScheduler() *SDMRoundRobinScheduler {
	return &SDMRoundRobinScheduler{}
}

func (s *SDMRoundRobinScheduler) GetProbeSnip(devid uint8) (retval QuerySnip) {
	retval = QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
		OpCode: OpCodeSDML1Voltage, Value: math.NaN(), Description: "L1 Voltage (V)", IEC61850: "VolLocPhsA"}
	return retval
}

func (s *SDMRoundRobinScheduler) Produce(devid uint8) (retval []QuerySnip) {
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadInputReg,
		OpCode: OpCodeSDML1Voltage, Value: math.NaN(), Description: "L1 Voltage (V)", IEC61850: "VolLocPhsA"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadInputReg, OpCode: OpCodeSDML2Voltage, Value: math.NaN(),
		Description: "L2 Voltage (V)", IEC61850: "VolLocPhsB"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadInputReg, OpCode: OpCodeSDML3Voltage, Value: math.NaN(),
		Description: "L3 Voltage (V)", IEC61850: "VolLocPhsC"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadInputReg, OpCode: OpCodeSDML1Current, Value: math.NaN(),
		Description: "L1 Current (A)", IEC61850: "AmpLocPhsA"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadInputReg, OpCode: OpCodeSDML2Current, Value: math.NaN(),
		Description: "L2 Current (A)", IEC61850: "AmpLocPhsB"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadInputReg, OpCode: OpCodeSDML3Current, Value: math.NaN(),
		Description: "L3 Current (A)", IEC61850: "AmpLocPhsC"})

	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadInputReg, OpCode: OpCodeSDML1Power, Value: math.NaN(),
		Description: "L1 Power (W)", IEC61850: "WLocPhsA"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadInputReg, OpCode: OpCodeSDML2Power, Value: math.NaN(),
		Description: "L2 Power (W)", IEC61850: "WLocPhsB"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadInputReg, OpCode: OpCodeSDML3Power, Value: math.NaN(),
		Description: "L3 Power (W)", IEC61850: "WLocPhsC"})

	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadInputReg, OpCode: OpCodeSDML1Cosphi, Value: math.NaN(),
		Description: "L1 Cosphi", IEC61850: "AngLocPhsA"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadInputReg, OpCode: OpCodeSDML2Cosphi, Value: math.NaN(),
		Description: "L2 Cosphi", IEC61850: "AngLocPhsB"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadInputReg, OpCode: OpCodeSDML3Cosphi, Value: math.NaN(),
		Description: "L3 Cosphi", IEC61850: "AngLocPhsC"})

	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadInputReg, OpCode: OpCodeSDML1Import, Value: math.NaN(),
		Description: "L1 Import (kWh)", IEC61850: "TotkWhImportPhsA"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadInputReg, OpCode: OpCodeSDML2Import, Value: math.NaN(),
		Description: "L2 Import (kWh)", IEC61850: "TotkWhImportPhsB"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadInputReg, OpCode: OpCodeSDML3Import, Value: math.NaN(),
		Description: "L3 Import (kWh)", IEC61850: "TotkWhImportPhsC"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadInputReg, OpCode: OpCodeSDMTotalImport, Value: math.NaN(),
		Description: "Total Import (kWh)", IEC61850: "TotkWhImport"})

	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadInputReg, OpCode: OpCodeSDML1Export, Value: math.NaN(),
		Description: "L1 Export (kWh)", IEC61850: "TotkWhExportPhsA"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadInputReg, OpCode: OpCodeSDML2Export, Value: math.NaN(),
		Description: "L2 Export (kWh)", IEC61850: "TotkWhExportPhsB"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadInputReg, OpCode: OpCodeSDML3Export, Value: math.NaN(),
		Description: "L3 Export (kWh)", IEC61850: "TotkWhExportPhsC"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadInputReg, OpCode: OpCodeSDMTotalExport, Value: math.NaN(),
		Description: "Total Export (kWh)", IEC61850: "TotkWhExport"})

	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadInputReg, OpCode: OpCodeSDML1THDVoltageNeutral, Value: math.NaN(),
		Description: "L1 Voltage to neutral THD (%)", IEC61850: "ThdVolPhsA"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadInputReg, OpCode: OpCodeSDML2THDVoltageNeutral, Value: math.NaN(),
		Description: "L2 Voltage to neutral THD (%)", IEC61850: "ThdVolPhsB"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadInputReg, OpCode: OpCodeSDML3THDVoltageNeutral, Value: math.NaN(),
		Description: "L3 Voltage to neutral THD (%)", IEC61850: "ThdVolPhsC"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadInputReg, OpCode: OpCodeSDMAvgTHDVoltageNeutral, Value: math.NaN(),
		Description: "Average voltage to neutral THD (%)", IEC61850: "ThdVol"})

	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadInputReg, OpCode: OpCodeSDMFrequency, Value: math.NaN(),
		Description: "Frequency of supply voltages", IEC61850: "Freq"})

	return retval
}

// ####################################################################
// Round-Robin Scheduler for the Janitza B23 DIN-Rail meters
// ####################################################################

type JanitzaRoundRobinScheduler struct {
}

func NewJanitzaRoundRobinScheduler() *JanitzaRoundRobinScheduler {
	return &JanitzaRoundRobinScheduler{}
}

func (s *JanitzaRoundRobinScheduler) GetProbeSnip(devid uint8) (retval QuerySnip) {
	retval = QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeJanitzaL1Voltage, Value: math.NaN(), Description: "L1 Voltage (V)", IEC61850: "VolLocPhsA"}
	return retval
}

func (s *JanitzaRoundRobinScheduler) Produce(devid uint8) (retval []QuerySnip) {
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeJanitzaL1Voltage, Value: math.NaN(), Description: "L1 Voltage (V)", IEC61850: "VolLocPhsA"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeJanitzaL2Voltage, Value: math.NaN(), Description: "L2 Voltage (V)", IEC61850: "VolLocPhsB"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeJanitzaL3Voltage, Value: math.NaN(), Description: "L3 Voltage (V)", IEC61850: "VolLocPhsC"})

	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg, OpCode: OpCodeJanitzaL1Current, Value: math.NaN(),
		Description: "L1 Current (A)", IEC61850: "AmpLocPhsA"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg, OpCode: OpCodeJanitzaL2Current, Value: math.NaN(),
		Description: "L2 Current (A)", IEC61850: "AmpLocPhsB"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg, OpCode: OpCodeJanitzaL3Current, Value: math.NaN(),
		Description: "L3 Current (A)", IEC61850: "AmpLocPhsC"})

	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg, OpCode: OpCodeJanitzaL1Power, Value: math.NaN(),
		Description: "L1 Power (W)", IEC61850: "WLocPhsA"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg, OpCode: OpCodeJanitzaL2Power, Value: math.NaN(),
		Description: "L2 Power (W)", IEC61850: "WLocPhsB"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg, OpCode: OpCodeJanitzaL3Power, Value: math.NaN(),
		Description: "L3 Power (W)", IEC61850: "WLocPhsC"})

	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg, OpCode: OpCodeJanitzaL1Cosphi, Value: math.NaN(),
		Description: "L1 Cosphi", IEC61850: "AngLocPhsA"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg, OpCode: OpCodeJanitzaL2Cosphi, Value: math.NaN(),
		Description: "L2 Cosphi", IEC61850: "AngLocPhsB"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg, OpCode: OpCodeJanitzaL3Cosphi, Value: math.NaN(),
		Description: "L3 Cosphi", IEC61850: "AngLocPhsC"})

	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg, OpCode: OpCodeJanitzaL1Import, Value: math.NaN(),
		Description: "L1 Import (kWh)", IEC61850: "TotkWhImportPhsA"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg, OpCode: OpCodeJanitzaL2Import, Value: math.NaN(),
		Description: "L2 Import (kWh)", IEC61850: "TotkWhImportPhsB"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg, OpCode: OpCodeJanitzaL3Import, Value: math.NaN(),
		Description: "L3 Import (kWh)", IEC61850: "TotkWhImportPhsC"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg, OpCode: OpCodeJanitzaTotalImport, Value: math.NaN(),
		Description: "Total Import (kWh)", IEC61850: "TotkWhImport"})

	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg, OpCode: OpCodeJanitzaL1Export, Value: math.NaN(),
		Description: "L1 Export (kWh)", IEC61850: "TotkWhExportPhsA"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg, OpCode: OpCodeJanitzaL2Export, Value: math.NaN(),
		Description: "L2 Export (kWh)", IEC61850: "TotkWhExportPhsB"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg, OpCode: OpCodeJanitzaL3Export, Value: math.NaN(),
		Description: "L3 Export (kWh)", IEC61850: "TotkWhExportPhsC"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg, OpCode: OpCodeJanitzaTotalExport, Value: math.NaN(),
		Description: "Total Export (kWh)", IEC61850: "TotkWhExport"})
	return retval
}

// ####################################################################
// Round-Robin Scheduler for the Eastron SDM Devices
// ####################################################################
type DZGRoundRobinScheduler struct {
}

func NewDZGRoundRobinScheduler() *DZGRoundRobinScheduler {
	return &DZGRoundRobinScheduler{}
}

func (s *DZGRoundRobinScheduler) GetProbeSnip(devid uint8) (retval QuerySnip) {
	retval = QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeDZGL1Voltage, Value: math.NaN(),
		Description: "L1 Voltage (V)", IEC61850: "VolLocPhsA",
		Transform: MkRTUScaledIntToFloat64(100)}
	return retval
}

func (s *DZGRoundRobinScheduler) Produce(devid uint8) (retval []QuerySnip) {
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeDZGL1Voltage, Value: math.NaN(), Description: "L1 Voltage (V)", IEC61850: "VolLocPhsA", Transform: MkRTUScaledIntToFloat64(100)})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeDZGL2Voltage, Value: math.NaN(), Description: "L2 Voltage (V)", IEC61850: "VolLocPhsB", Transform: MkRTUScaledIntToFloat64(100)})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeDZGL3Voltage, Value: math.NaN(), Description: "L3 Voltage (V)", IEC61850: "VolLocPhsC", Transform: MkRTUScaledIntToFloat64(100)})

	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg, OpCode: OpCodeDZGL1Current, Value: math.NaN(),
		Description: "L1 Current (A)", IEC61850: "AmpLocPhsA", Transform: MkRTUScaledIntToFloat64(1000)})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg, OpCode: OpCodeDZGL2Current, Value: math.NaN(),
		Description: "L2 Current (A)", IEC61850: "AmpLocPhsB", Transform: MkRTUScaledIntToFloat64(1000)})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg, OpCode: OpCodeDZGL3Current, Value: math.NaN(),
		Description: "L3 Current (A)", IEC61850: "AmpLocPhsC", Transform: MkRTUScaledIntToFloat64(1000)})

	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeDZGL1Import, Value: math.NaN(),
		Description: "L1 Import (kWh)", IEC61850: "TotkWhImportPhsA"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeDZGL2Import, Value: math.NaN(),
		Description: "L2 Import (kWh)", IEC61850: "TotkWhImportPhsB"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeDZGL3Import, Value: math.NaN(),
		Description: "L3 Import (kWh)", IEC61850: "TotkWhImportPhsC"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeDZGTotalImport, Value: math.NaN(),
		Description: "Total Import", IEC61850: "TotkWhImport",
		Transform: MkRTUScaledIntToFloat64(1000)})

	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeDZGL1Export, Value: math.NaN(),
		Description: "L1 Export (kWh)", IEC61850: "TotkWhExportPhsA"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeDZGL2Export, Value: math.NaN(),
		Description: "L2 Export (kWh)", IEC61850: "TotkWhExportPhsB"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeDZGL3Export, Value: math.NaN(),
		Description: "L3 Export (kWh)", IEC61850: "TotkWhExportPhsC"})
	retval = append(retval, QuerySnip{DeviceId: devid, FuncCode: ReadHoldingReg,
		OpCode: OpCodeDZGTotalExport, Value: math.NaN(),
		Description: "Total Export", IEC61850: "TotkWhExport",
		Transform: MkRTUScaledIntToFloat64(1000)})

	return retval
}
