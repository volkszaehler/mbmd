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
	OpCodeL1Voltage     = 0x0000
	OpCodeL2Voltage     = 0x0002
	OpCodeL3Voltage     = 0x0004
	OpCodeL1Current     = 0x0006
	OpCodeL2Current     = 0x0008
	OpCodeL3Current     = 0x000A
	OpCodeL1Power       = 0x000C
	OpCodeL2Power       = 0x000E
	OpCodeL3Power       = 0x0010
	OpCodeL1PowerFactor = 0x001e
	OpCodeL2PowerFactor = 0x0020
	OpCodeL3PowerFactor = 0x0022

	MaxRetryCount = 3
)

type QueryEngine struct {
	client     modbus.Client
	interval   int
	handler    modbus.RTUClientHandler
	datastream ReadingChannel
}

func NewQueryEngine(
	rtuDevice string,
	interval int,
	verbose bool,
	channel ReadingChannel,
) *QueryEngine {
	// Modbus RTU/ASCII
	mbhandler := modbus.NewRTUClientHandler(rtuDevice)
	mbhandler.BaudRate = 9600
	mbhandler.DataBits = 8
	mbhandler.Parity = "N"
	mbhandler.StopBits = 1
	mbhandler.SlaveId = 1
	mbhandler.Timeout = 1000 * time.Millisecond
	if verbose {
		mbhandler.Logger = log.New(os.Stdout, "RTUClientHandler: ", log.LstdFlags)
		log.Printf("Connecting to RTU via %s\r\n", rtuDevice)
	}

	err := mbhandler.Connect()
	if err != nil {
		log.Fatal("Failed to connect: ", err)
	}

	mbclient := modbus.NewClient(mbhandler)

	return &QueryEngine{client: mbclient, interval: interval,
		handler: *mbhandler, datastream: channel}
}

func (q *QueryEngine) retrieveOpCode(opcode uint16) (retval float32,
	err error) {
	results, err := q.client.ReadInputRegisters(opcode, 2)
	if err == nil {
		retval = RtuToFloat32(results)
	}
	return retval, err
}

func (q *QueryEngine) queryOrFail(opcode uint16) (retval float32) {
	var err error
	tryCnt := 0
	for tryCnt = 0; tryCnt < MaxRetryCount; tryCnt++ {
		retval, err = q.retrieveOpCode(opcode)
		if err != nil {
			log.Printf("Closing broken handler, reconnecting attempt %d\r\n", tryCnt)
			// Note: Just close the handler here. If a new handler is manually
			// created it will create a resource leak (file descriptors). Just
			// close the handler, the modbus library will recreate one as
			// needed.
			q.handler.Close()
			time.Sleep(time.Duration(1) * time.Second)
		} else {
			break
		}
	}
	if tryCnt == MaxRetryCount {
		log.Fatal("Cannot query the sensor, reached maximum retry count.")
	}
	return retval
}

func (q *QueryEngine) Produce() {
	for {
		q.datastream <- Readings{
			Timestamp: time.Now(),
			L1Voltage: q.queryOrFail(OpCodeL1Voltage),
			L2Voltage: q.queryOrFail(OpCodeL2Voltage),
			L3Voltage: q.queryOrFail(OpCodeL3Voltage),
			L1Current: q.queryOrFail(OpCodeL1Current),
			L2Current: q.queryOrFail(OpCodeL2Current),
			L3Current: q.queryOrFail(OpCodeL3Current),
			L1Power:   q.queryOrFail(OpCodeL1Power),
			L2Power:   q.queryOrFail(OpCodeL2Power),
			L3Power:   q.queryOrFail(OpCodeL3Power),
			L1CosPhi:  q.queryOrFail(OpCodeL1PowerFactor),
			L2CosPhi:  q.queryOrFail(OpCodeL2PowerFactor),
			L3CosPhi:  q.queryOrFail(OpCodeL3PowerFactor),
		}
		time.Sleep(time.Duration(q.interval) * time.Second)
	}
	q.handler.Close()
}

func RtuToFloat32(b []byte) (f float32) {
	bits := binary.BigEndian.Uint32(b)
	f = math.Float32frombits(bits)
	return
}
