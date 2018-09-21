package sdm630

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/goburrow/modbus"
	. "github.com/gonium/gosdm630/internal/meters"
)

const (
	MaxRetryCount = 5
)

const (
	_ = iota
	ModbusComset2400_8N1
	ModbusComset9600_8N1
	ModbusComset19200_8N1
	ModbusComset2400_8E1
	ModbusComset9600_8E1
	ModbusComset19200_8E1
)

type ModbusEngine struct {
	client  modbus.Client
	handler *modbus.RTUClientHandler
	verbose bool
	status  *Status
}

func NewRTUClient(rtuDevice string, comset int, verbose bool) *modbus.RTUClientHandler {
	// Modbus RTU/ASCII
	rtuclient := modbus.NewRTUClientHandler(rtuDevice)

	rtuclient.Parity = "N"
	rtuclient.DataBits = 8
	rtuclient.StopBits = 1

	switch comset {
	case ModbusComset2400_8N1:
		rtuclient.BaudRate = 2400
	case ModbusComset9600_8N1:
		rtuclient.BaudRate = 9600
	case ModbusComset19200_8N1:
		rtuclient.BaudRate = 19200
	case ModbusComset2400_8E1:
		rtuclient.BaudRate = 2400
		rtuclient.Parity = "E"
	case ModbusComset9600_8E1:
		rtuclient.BaudRate = 9600
		rtuclient.Parity = "E"
	case ModbusComset19200_8E1:
		rtuclient.BaudRate = 19200
		rtuclient.Parity = "E"
	default:
		log.Fatal("Invalid communication set specified. See -h for help.")
	}

	rtuclient.Timeout = 300 * time.Millisecond
	if verbose {
		rtuclient.Logger = log.New(os.Stdout, "RTUClientHandler: ", log.LstdFlags)
		log.Printf("Connecting to RTU via %s, %d %d%s%d\r\n", rtuDevice,
			rtuclient.BaudRate, rtuclient.DataBits, rtuclient.Parity,
			rtuclient.StopBits)
	}

	err := rtuclient.Connect()
	if err != nil {
		log.Fatal("Failed to connect: ", err)
	}
	defer rtuclient.Close()

	return rtuclient
}

func NewModbusEngine(
	rtuDevice string,
	comset int,
	verbose bool,
	status *Status,
) *ModbusEngine {
	var rtuclient *modbus.RTUClientHandler
	var mbclient modbus.Client

	if rtuDevice == "simulation" {
		rtuclient = &modbus.RTUClientHandler{}
		mbclient = NewMockClient(0)
	} else {
		rtuclient = NewRTUClient(rtuDevice, comset, verbose)
		mbclient = modbus.NewClient(rtuclient)
	}

	return &ModbusEngine{
		client:  mbclient,
		handler: rtuclient,
		verbose: verbose,
		status:  status,
	}
}

func (q *ModbusEngine) query(snip QuerySnip) (retval []byte, err error) {
	q.status.IncreaseModbusRequestCounter()

	// update the slave id in the handler
	q.handler.SlaveId = snip.DeviceId

	if snip.ReadLen <= 0 {
		log.Fatalf("Invalid meter operation %v.", snip)
	}

	switch snip.FuncCode {
	case ReadInputReg:
		retval, err = q.client.ReadInputRegisters(snip.OpCode, snip.ReadLen)
	case ReadHoldingReg:
		retval, err = q.client.ReadHoldingRegisters(snip.OpCode, snip.ReadLen)
	default:
		log.Fatalf("Unknown function code %d - cannot query device.",
			snip.FuncCode)
	}

	if err != nil && q.verbose {
		log.Printf("Failed to retrieve opcode 0x%x, error was: %s\r\n", snip.OpCode, err.Error())
	}

	return retval, err
}

func (q *ModbusEngine) Transform(
	inputStream QuerySnipChannel,
	controlStream ControlSnipChannel,
	outputStream QuerySnipChannel,
) {
	var previousDeviceId uint8
	for {
	PROCESS_READINGS:
		snip := <-inputStream
		// The SDM devices need to have a little pause between querying
		// different devices.
		if previousDeviceId != snip.DeviceId {
			time.Sleep(time.Duration(100) * time.Millisecond)
		}
		previousDeviceId = snip.DeviceId

		for retryCount := 0; retryCount < MaxRetryCount; retryCount++ {
			reading, err := q.query(snip)
			if err == nil {
				// convert bytes to value
				snip.Value = snip.Transform(reading)
				snip.ReadTimestamp = time.Now()
				outputStream <- snip

				// signal ok
				successSnip := ControlSnip{
					Type:     CONTROLSNIP_OK,
					Message:  "OK",
					DeviceId: snip.DeviceId,
				}
				controlStream <- successSnip

				goto PROCESS_READINGS
			} else {
				q.status.IncreaseModbusReconnectCounter()
				log.Printf("Device %d failed to respond - retry attempt %d of %d",
					snip.DeviceId, retryCount+1, MaxRetryCount)
				time.Sleep(time.Duration(100) * time.Millisecond)
			}
		}

		// signal error
		errorSnip := ControlSnip{
			Type:     CONTROLSNIP_ERROR,
			Message:  fmt.Sprintf("Device %d did not respond.", snip.DeviceId),
			DeviceId: snip.DeviceId,
		}
		controlStream <- errorSnip
	}
}

func (q *ModbusEngine) Scan() {
	type DeviceInfo struct {
		DeviceId  uint8
		MeterType string
	}

	var deviceId uint8
	deviceList := make([]DeviceInfo, 0)
	oldtimeout := q.handler.Timeout
	q.handler.Timeout = 50 * time.Millisecond
	log.Printf("Starting bus scan")

	producers := []Producer{
		NewSDMProducer(),
		NewJanitzaProducer(),
		NewDZGProducer(),
	}

SCAN:
	// loop over all valid slave adresses
	for deviceId = 1; deviceId <= 247; deviceId++ {
		// give the bus some time to recover before querying the next device
		time.Sleep(time.Duration(40) * time.Millisecond)

		for _, producer := range producers {
			operation := producer.Probe()
			snip := NewQuerySnip(deviceId, operation)

			value, err := q.query(snip)
			if err == nil {
				log.Printf("Device %d: %s type device found, %s: %.2f\r\n",
					deviceId,
					producer.GetMeterType(),
					GetIecDescription(snip.IEC61850),
					snip.Transform(value))
				dev := DeviceInfo{
					DeviceId:  deviceId,
					MeterType: producer.GetMeterType(),
				}
				deviceList = append(deviceList, dev)
				continue SCAN
			}
		}

		log.Printf("Device %d: n/a\r\n", deviceId)
	}

	// restore timeout to old value
	q.handler.Timeout = oldtimeout
	log.Printf("Found %d active devices:\r\n", len(deviceList))
	for _, device := range deviceList {
		log.Printf("* slave address %d: type %s\r\n", device.DeviceId,
			device.MeterType)
	}
	log.Println("WARNING: This lists only the devices that responded to " +
		"a known probe request. Devices with different " +
		"function code definitions might not be detected.")
}
