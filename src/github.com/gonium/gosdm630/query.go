package sdm630

import (
	"encoding/binary"
	"github.com/goburrow/modbus"
	"log"
	"math"
	"os"
	"time"
)

// See http://bg-etech.de/download/manual/SDM630Register.pdf
const (
	OpCodeL1Voltage     = 0x0000
	OpCodeL2Voltage     = 0x0002
	OpCodeL3Voltage     = 0x0004
	OpCodeL1Current     = 0x0006
	OpCodeL2Current     = 0x0008
	OpCodeL3Current     = 0x000A
	OpCodeL1Power       = 0x000C
	OpCodeL2Power       = 0x000E
	OpCodeL3Power       = 0x0010
	OpCodeL1Import      = 0x015a
	OpCodeL2Import      = 0x015c
	OpCodeL3Import      = 0x015e
	OpCodeL1Export      = 0x0160
	OpCodeL2Export      = 0x0162
	OpCodeL3Export      = 0x0164
	OpCodeL1PowerFactor = 0x001e
	OpCodeL2PowerFactor = 0x0020
	OpCodeL3PowerFactor = 0x0022

	MaxRetryCount = 3
)

type QueryEngine struct {
	client     modbus.Client
	interval   int
	handler    *modbus.RTUClientHandler
	datastream ReadingChannel
	devids     []uint8
	verbose    bool
}

func NewQueryEngine(
	rtuDevice string,
	interval int,
	verbose bool,
	channel ReadingChannel,
	devids []uint8,
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

	return &QueryEngine{client: mbclient, interval: interval,
		handler: rtuclient, datastream: channel,
		devids: devids, verbose: verbose}
}

func (q *QueryEngine) retrieveOpCode(opcode uint16) (retval float32,
	err error) {
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
					L1: q.queryOrFail(OpCodeL1PowerFactor),
					L2: q.queryOrFail(OpCodeL2PowerFactor),
					L3: q.queryOrFail(OpCodeL3PowerFactor),
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
		if q.interval > 0 {
			time.Sleep(time.Duration(q.interval) * time.Second)
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
