package sdm630

import (
	"encoding/binary"
	"github.com/goburrow/modbus"
	"log"
	"math"
	"os"
	"time"
)

const (
	MaxRetryCount = 5
)

type QueryEngine struct {
	client     modbus.Client
	handler    *modbus.RTUClientHandler
	datastream ReadingChannel
	devids     []uint8
	verbose    bool
	status     *Status
}

func NewQueryEngine(
	rtuDevice string,
	verbose bool,
	channel ReadingChannel,
	devids []uint8,
	status *Status,
) *QueryEngine {
	// Modbus RTU/ASCII
	rtuclient := modbus.NewRTUClientHandler(rtuDevice)
	rtuclient.BaudRate = 9600
	rtuclient.DataBits = 8
	rtuclient.Parity = "N"
	rtuclient.StopBits = 1
	// TODO: Add support for more than one slave ID.
	rtuclient.SlaveId = devids[0]
	rtuclient.Timeout = 1000 * time.Millisecond
	if verbose {
		rtuclient.Logger = log.New(os.Stdout, "RTUClientHandler: ", log.LstdFlags)
		log.Printf("Connecting to RTU via %s\r\n", rtuDevice)
	}

	err := rtuclient.Connect()
	if err != nil {
		log.Fatal("Failed to connect: ", err)
	}
	defer rtuclient.Close()

	mbclient := modbus.NewClient(rtuclient)

	return &QueryEngine{client: mbclient,
		handler: rtuclient, datastream: channel,
		devids: devids, verbose: verbose,
		status: status,
	}
}

func (q *QueryEngine) retrieveOpCode(opcode uint16) (retval float32,
	err error) {
	q.status.IncreaseModbusRequestCounter()
	results, err := q.client.ReadInputRegisters(opcode, 2)
	if err == nil {
		retval = rtuToFload32(results)
	} else if q.verbose {
		log.Printf("Failed to retrieve opcode 0x%x, error was: %s\r\n", opcode, err.Error())
	}
	return retval, err
}

func (q *QueryEngine) queryOrFail(opcode uint16) (retval float32) {
	var err error
	tryCnt := 0
	for tryCnt = 0; tryCnt < MaxRetryCount; tryCnt++ {
		retval, err = q.retrieveOpCode(opcode)
		if err != nil {
			q.status.IncreaseModbusReconnectCounter()
			log.Printf("Failed to retrieve opcode - retry attempt %d of %d\r\n", tryCnt+1,
				MaxRetryCount)
			time.Sleep(time.Duration(100) * time.Millisecond)
		} else {
			break
		}
	}
	if tryCnt == MaxRetryCount {
		log.Fatal("Cannot query the sensor, reached maximum retry count. Abort.")
	}
	return retval
}

func (q *QueryEngine) Produce() {
	for {
		//start := time.Now()
		for _, devid := range q.devids {
			// Set the current device id as "slave id" in the modbus rtu
			// library. Somewhat ugly...
			q.handler.SlaveId = devid
			timestamp := time.Now()
			q.datastream <- Readings{
				Timestamp:      timestamp,
				Unix:           timestamp.Unix(),
				ModbusDeviceId: devid,
				Voltage: ThreePhaseReadings{
					L1: q.queryOrFail(OpCodeL1Voltage),
					L2: q.queryOrFail(OpCodeL2Voltage),
					L3: q.queryOrFail(OpCodeL3Voltage),
				},
				Current: ThreePhaseReadings{
					L1: q.queryOrFail(OpCodeL1Current),
					L2: q.queryOrFail(OpCodeL2Current),
					L3: q.queryOrFail(OpCodeL3Current),
				},
				Power: ThreePhaseReadings{
					L1: q.queryOrFail(OpCodeL1Power),
					L2: q.queryOrFail(OpCodeL2Power),
					L3: q.queryOrFail(OpCodeL3Power),
				},
				Cosphi: ThreePhaseReadings{
					L1: q.queryOrFail(OpCodeL1Cosphi),
					L2: q.queryOrFail(OpCodeL2Cosphi),
					L3: q.queryOrFail(OpCodeL3Cosphi),
				},
				Import: ThreePhaseReadings{
					L1: q.queryOrFail(OpCodeL1Import),
					L2: q.queryOrFail(OpCodeL2Import),
					L3: q.queryOrFail(OpCodeL3Import),
				},
				Export: ThreePhaseReadings{
					L1: q.queryOrFail(OpCodeL1Export),
					L2: q.queryOrFail(OpCodeL2Export),
					L3: q.queryOrFail(OpCodeL3Export),
				},
			}
			time.Sleep(20 * time.Millisecond)
		}
		//elapsed := time.Since(start)
		//log.Printf("Reading all values took %s", elapsed)
	}
	q.handler.Close()
}

func rtuToFload32(b []byte) (f float32) {
	bits := binary.BigEndian.Uint32(b)
	f = math.Float32frombits(bits)
	return
}
