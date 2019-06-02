package server

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"time"

	. "github.com/volkszaehler/mbmd/meters"
	"github.com/grid-x/modbus"
)

const (
	maxRetry = 3
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
	handler modbus.ClientHandler
	verbose bool
	status  *Status
}

// injectable logger for grid-x modbus implementation
type modbusLogger struct{}

func (l *modbusLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// NewRTUClientHandler creates a serial line  RTU/ASCII modbus handler
func NewRTUClientHandler(rtuDevice string, comset int, verbose bool) *modbus.RTUClientHandler {
	handler := modbus.NewRTUClientHandler(rtuDevice)

	handler.Parity = "N"
	handler.DataBits = 8
	handler.StopBits = 1

	switch comset {
	case ModbusComset2400_8N1:
		handler.BaudRate = 2400
	case ModbusComset9600_8N1:
		handler.BaudRate = 9600
	case ModbusComset19200_8N1:
		handler.BaudRate = 19200
	case ModbusComset2400_8E1:
		handler.BaudRate = 2400
		handler.Parity = "E"
	case ModbusComset9600_8E1:
		handler.BaudRate = 9600
		handler.Parity = "E"
	case ModbusComset19200_8E1:
		handler.BaudRate = 19200
		handler.Parity = "E"
	default:
		log.Fatal("Invalid communication set specified. See -h for help.")
	}

	handler.Timeout = 300 * time.Millisecond
	if verbose {
		logger := &modbusLogger{}
		handler.Logger = logger
		log.Printf("Connecting to RTU via %s, %d %d%s%d\r\n", rtuDevice,
			handler.BaudRate, handler.DataBits, handler.Parity,
			handler.StopBits)
	}

	return handler
}

// NewTCPClientHandler creates a TCP modbus handler
func NewTCPClientHandler(rtuDevice string, verbose bool) *modbus.TCPClientHandler {
	handler := modbus.NewTCPClientHandler(rtuDevice)
	handler.Timeout = 1 * time.Second
	handler.ProtocolRecoveryTimeout = 10 * time.Second
	handler.LinkRecoveryTimeout = 15 * time.Second

	if verbose {
		logger := &modbusLogger{}
		handler.Logger = logger
	}
	return handler
}

func NewModbusEngine(
	rtuDevice string,
	comset int,
	simulate bool,
	verbose bool,
	status *Status,
) *ModbusEngine {
	var handler modbus.ClientHandler
	var mbclient modbus.Client

	if simulate {
		log.Println("*** Simulation mode ***")
		mbclient = NewMockClient(20) // error rate for testing
	} else {
		// parse adapter string
		re := regexp.MustCompile(":[0-9]+$")
		if re.MatchString(rtuDevice) {
			// tcp connection
			handler = NewTCPClientHandler(rtuDevice, verbose)
		} else {
			// serial connection
			handler = NewRTUClientHandler(rtuDevice, comset, verbose)
		}

		mbclient = modbus.NewClient(handler)
		if err := handler.Connect(); err != nil {
			log.Fatal("Failed to connect: ", err)
		}
	}

	return &ModbusEngine{
		client:  mbclient,
		handler: handler,
		verbose: verbose,
		status:  status,
	}
}

func (q *ModbusEngine) setTimeout(timeout time.Duration) time.Duration {
	// update the slave id in the handler
	if handler, ok := q.handler.(*modbus.RTUClientHandler); ok {
		t := handler.Timeout
		handler.Timeout = timeout
		return t
	} else if handler, ok := q.handler.(*modbus.TCPClientHandler); ok {
		t := handler.Timeout
		handler.Timeout = timeout
		return t
	} else if handler != nil {
		log.Fatal("Unsupported modbus handler")
	}
	return 0
}

func (q *ModbusEngine) setSlave(deviceId uint8) {
	// update the slave id in the handler
	if handler, ok := q.handler.(*modbus.RTUClientHandler); ok {
		handler.SetSlave(deviceId)
	} else if handler, ok := q.handler.(*modbus.TCPClientHandler); ok {
		handler.SetSlave(deviceId)
	} else if handler != nil {
		log.Fatal("Unsupported modbus handler")
	}
}

// Reconnect refreshes underlying modbus TCP connection if connected via TCP
func (q *ModbusEngine) Reconnect() {
	if handler, ok := q.handler.(*modbus.TCPClientHandler); ok {
		handler.Close()
	}
}

// Query modbus registers
func (q *ModbusEngine) Query(snip QuerySnip) (retval []byte, err error) {
	q.setSlave(snip.DeviceId)
	q.status.IncreaseRequestCounter()

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
		log.Printf("Device %d failed to retrieve opcode 0x%x, error was: %s\r\n", snip.DeviceId, snip.OpCode, err.Error())
	}

	return retval, err
}

// Transform converts raw query result into one or more Readings
func (q *ModbusEngine) Transform(snip QuerySnip, bytes []byte) []QuerySnip {
	now := time.Now()

	if snip.Splitter != nil {
		// block reading - needs splitting
		snips := snip.Splitter(bytes)
		res := make([]QuerySnip, len(snips))

		for idx, sr := range snips {
			splitSnip := QuerySnip{
				DeviceId: snip.DeviceId,
				Operation: Operation{
					OpCode:   sr.OpCode,
					IEC61850: sr.IEC61850,
				},
				Value:         sr.Value,
				ReadTimestamp: now,
			}
			res[idx] = splitSnip
		}

		return res
	}

	// single reading
	if snip.Transform == nil {
		log.Fatalf("Snip transformation not defined: %v", snip)
	}

	// convert bytes to value
	snip.Value = snip.Transform(bytes)
	snip.ReadTimestamp = now

	// check if result is valid
	if math.IsNaN(snip.Value) {
		return []QuerySnip{}
	}

	return []QuerySnip{snip}
}

// Run consumes device operations and produces operation results
func (q *ModbusEngine) Run(
	snipChannel QuerySnipChannel,
	controlChannel ControlSnipChannel,
	outputChannel QuerySnipChannel,
) {
	defer close(outputChannel)
	defer close(controlChannel)

	var previousDeviceId uint8
	for snip := range snipChannel {
		// The SDM devices need to have a little pause between querying
		// different devices.
		if previousDeviceId != snip.DeviceId {
			time.Sleep(time.Duration(100) * time.Millisecond)
			previousDeviceId = snip.DeviceId
		}

		var bytes []byte
		var err error
		var controlSnip ControlSnip

		for retry := 0; retry < maxRetry; retry++ {
			bytes, err = q.Query(snip)
			if err == nil {
				break
			}

			q.status.IncreaseReconnectCounter()
			log.Printf("Device %d failed to respond (%d/%d)",
				snip.DeviceId, retry+1, maxRetry)
			time.Sleep(time.Duration(100) * time.Millisecond)
		}

		if err == nil {
			// success
			snips := q.Transform(snip, bytes)
			for _, snip := range snips {
				if q.verbose {
					log.Printf("Device %d - %s: %.2f\n", snip.DeviceId, snip.IEC61850.String(), snip.Value)
				}
				outputChannel <- snip
			}

			// signal ok
			controlSnip = ControlSnip{
				Type:     CONTROLSNIP_OK,
				Message:  "OK",
				DeviceId: snip.DeviceId,
			}
		} else {
			// signal error
			controlSnip = ControlSnip{
				Type:     CONTROLSNIP_ERROR,
				Message:  fmt.Sprintf("Device %d did not respond.", snip.DeviceId),
				DeviceId: snip.DeviceId,
			}
		}

		controlChannel <- controlSnip
	}
}

func (q *ModbusEngine) Scan() {
	type DeviceInfo struct {
		DeviceId  uint8
		MeterType string
	}

	var deviceId uint8
	deviceList := make([]DeviceInfo, 0)
	oldtimeout := q.setTimeout(50 * time.Millisecond)
	log.Printf("Starting bus scan")

SCAN:
	// loop over all valid slave adresses
	for deviceId = 1; deviceId <= 247; deviceId++ {
		// give the bus some time to recover before querying the next device
		time.Sleep(time.Duration(40) * time.Millisecond)

		for _, factory := range Producers {
			producer := factory()
			operation := producer.Probe()
			snip := NewQuerySnip(deviceId, operation)

			value, err := q.Query(snip)
			if err == nil {
				log.Printf("Device %d: %s type device found, %s: %.2f\r\n",
					deviceId,
					producer.Type(),
					snip.IEC61850,
					snip.Transform(value))
				dev := DeviceInfo{
					DeviceId:  deviceId,
					MeterType: producer.Type(),
				}
				deviceList = append(deviceList, dev)
				continue SCAN
			}
		}

		log.Printf("Device %d: n/a\r\n", deviceId)
	}

	// restore timeout to old value
	q.setTimeout(oldtimeout)

	log.Printf("Found %d active devices:\r\n", len(deviceList))
	for _, device := range deviceList {
		log.Printf("* slave address %d: type %s\r\n", device.DeviceId,
			device.MeterType)
	}
	log.Println("WARNING: This lists only the devices that responded to " +
		"a known probe request. Devices with different " +
		"function code definitions might not be detected.")
}
