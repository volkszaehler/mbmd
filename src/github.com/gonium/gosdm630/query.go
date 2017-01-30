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
	client       modbus.Client
	handler      *modbus.RTUClientHandler
	inputStream  QuerySnipChannel
	outputStream QuerySnipChannel
	devids       []uint8
	verbose      bool
	status       *Status
}

func NewQueryEngine(
	rtuDevice string,
	verbose bool,
	inputChannel QuerySnipChannel,
	outputChannel QuerySnipChannel,
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

	return &QueryEngine{
		client: mbclient, handler: rtuclient,
		inputStream: inputChannel, outputStream: outputChannel,
		devids: devids, verbose: verbose,
		status: status,
	}
}

func (q *QueryEngine) retrieveOpCode(opcode uint16) (retval float64,
	err error) {
	q.status.IncreaseModbusRequestCounter()
	results, err := q.client.ReadInputRegisters(opcode, 2)
	if err == nil {
		retval = rtuToFloat64(results)
	} else if q.verbose {
		log.Printf("Failed to retrieve opcode 0x%x, error was: %s\r\n", opcode, err.Error())
	}
	return retval, err
}

func (q *QueryEngine) queryOrFail(opcode uint16) (retval float64) {
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

func (q *QueryEngine) Transform() {
	var previousDeviceId uint8 = 0
	for {
		snip := <-q.inputStream
		q.handler.SlaveId = snip.DeviceId
		// apparently the turnaround timeout must be respected
		// See http://www.modbus.org/docs/Modbus_over_serial_line_V1_02.pdf
		// 3.5 chars at 9600 Baud take 36 ms
		if previousDeviceId != snip.DeviceId {
			time.Sleep(time.Duration(40) * time.Millisecond)
		}
		previousDeviceId = snip.DeviceId
		value := q.queryOrFail(snip.OpCode)
		snip.Value = value
		snip.ReadTimestamp = time.Now()
		q.outputStream <- snip
	}
	q.handler.Close()
}

func rtuToFloat64(b []byte) float64 {
	bits := binary.BigEndian.Uint32(b)
	f := math.Float32frombits(bits)
	return float64(f)
}
