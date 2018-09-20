package sdm630

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/goburrow/modbus"
)

const (
	MaxRetryCount  = 5
	ReadInputReg   = 4
	ReadHoldingReg = 3
)

const (
	ModbusComset2400_8N1  = 1
	ModbusComset9600_8N1  = 2
	ModbusComset19200_8N1 = 3
	ModbusComset2400_8E1  = 4
	ModbusComset9600_8E1  = 5
	ModbusComset19200_8E1 = 6
)

type ModbusEngine struct {
	client  modbus.Client
	handler *modbus.RTUClientHandler
	verbose bool
	status  *Status
}

func NewModbusEngine(
	rtuDevice string,
	comset int,
	verbose bool,
	status *Status,
) *ModbusEngine {
	// Modbus RTU/ASCII
	rtuclient := modbus.NewRTUClientHandler(rtuDevice)
	switch comset {
	case ModbusComset2400_8N1:
		rtuclient.BaudRate = 2400
		rtuclient.DataBits = 8
		rtuclient.Parity = "N"
		rtuclient.StopBits = 1
	case ModbusComset9600_8N1:
		rtuclient.BaudRate = 9600
		rtuclient.DataBits = 8
		rtuclient.Parity = "N"
		rtuclient.StopBits = 1
	case ModbusComset19200_8N1:
		rtuclient.BaudRate = 19200
		rtuclient.DataBits = 8
		rtuclient.Parity = "N"
		rtuclient.StopBits = 1
	case ModbusComset2400_8E1:
		rtuclient.BaudRate = 2400
		rtuclient.DataBits = 8
		rtuclient.Parity = "E"
		rtuclient.StopBits = 1
	case ModbusComset9600_8E1:
		rtuclient.BaudRate = 9600
		rtuclient.DataBits = 8
		rtuclient.Parity = "E"
		rtuclient.StopBits = 1
	case ModbusComset19200_8E1:
		rtuclient.BaudRate = 19200
		rtuclient.DataBits = 8
		rtuclient.Parity = "E"
		rtuclient.StopBits = 1
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

	mbclient := modbus.NewClient(rtuclient)

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
		snip := <-inputStream
		// The SDM devices need to have a little pause between querying
		// different devices.
		if previousDeviceId != snip.DeviceId {
			time.Sleep(time.Duration(100) * time.Millisecond)
		}
		previousDeviceId = snip.DeviceId

		var err error
		var reading []byte

		tryCnt := 0
		for tryCnt = 0; tryCnt < MaxRetryCount; tryCnt++ {
			reading, err = q.query(snip)
			if err != nil {
				q.status.IncreaseModbusReconnectCounter()
				log.Printf("Device %d failed to respond - retry attempt %d of %d",
					snip.DeviceId, tryCnt+1, MaxRetryCount)
				time.Sleep(time.Duration(100) * time.Millisecond)
			} else {
				break
			}
		}

		if tryCnt == MaxRetryCount {
			errorSnip := ControlSnip{
				Type:     CONTROLSNIP_ERROR,
				Message:  fmt.Sprintf("Device %d did not respond.", snip.DeviceId),
				DeviceId: snip.DeviceId,
			}
			controlStream <- errorSnip
		} else {
			// convert bytes to value
			snip.Value = snip.Transform(reading)
			snip.ReadTimestamp = time.Now()
			outputStream <- snip

			successSnip := ControlSnip{
				Type:     CONTROLSNIP_OK,
				Message:  "OK",
				DeviceId: snip.DeviceId,
			}
			controlStream <- successSnip
		}
	}
}

func (q *ModbusEngine) Scan() {
	type Device struct {
		DeviceId   uint8
		DeviceType MeterType
	}

	devicelist := make([]Device, 0)
	oldtimeout := q.handler.Timeout
	q.handler.Timeout = 50 * time.Millisecond
	log.Printf("Starting bus scan")

	probe := func(meterType MeterType, snip QuerySnip) bool {
		value, err := q.query(snip)
		if err == nil {
			log.Printf("Device %d: %s type device found, %s: %.2f\r\n",
				snip.DeviceId,
				meterType,
				GetIecDescription(snip.IEC61850),
				snip.Transform(value))
			dev := Device{
				DeviceId:   snip.DeviceId,
				DeviceType: meterType,
			}
			devicelist = append(devicelist, dev)
			return true
		}
		return false
	}

	// loop over all valid slave adresses
	var devid uint8
	for devid = 1; devid <= 247; devid++ {
		if probe(METERTYPE_SDM, NewSDMProducer().Probe(devid)) {
			continue
		}
		if probe(METERTYPE_JANITZA, NewJanitzaProducer().Probe(devid)) {
			continue
		}
		if probe(METERTYPE_DZG, NewDZGProducer().Probe(devid)) {
			continue
		}

		log.Printf("Device %d: n/a\r\n", devid)

		// give the bus some time to recover before querying the next device
		time.Sleep(time.Duration(40) * time.Millisecond)
	}

	// restore timeout to old value
	q.handler.Timeout = oldtimeout
	log.Printf("Found %d active devices:\r\n", len(devicelist))
	for _, device := range devicelist {
		log.Printf("* slave address %d: type %s\r\n", device.DeviceId,
			device.DeviceType)
	}
	log.Println("WARNING: This lists only the devices that responded to " +
		"a known probe request. Devices with different " +
		"function code definitions might not be detected.")
}

// RTUTransform functions convert RTU bytes to meaningful data types.
type RTUTransform func([]byte) float64

// RTU32ToFloat64 converts 32 bit readings
func RTU32ToFloat64(b []byte) float64 {
	bits := binary.BigEndian.Uint32(b)
	f := math.Float32frombits(bits)
	return float64(f)
}

// RTU16ToFloat64 converts 16 bit readings
func RTU16ToFloat64(b []byte) float64 {
	u := binary.BigEndian.Uint16(b)
	return float64(u)
}

func rtuScaledInt32ToFloat64(b []byte, scalar float64) float64 {
	unscaled := float64(binary.BigEndian.Uint32(b))
	f := unscaled / scalar
	return float64(f)
}

// MakeRTU32ScaledIntToFloat64 creates a 32 bit scaled reading transform
func MakeRTU32ScaledIntToFloat64(scalar float64) RTUTransform {
	return RTUTransform(func(b []byte) float64 {
		return rtuScaledInt32ToFloat64(b, scalar)
	})
}

func rtuScaledInt16ToFloat64(b []byte, scalar float64) float64 {
	unscaled := float64(binary.BigEndian.Uint16(b))
	f := unscaled / scalar
	return float64(f)
}

// MakeRTU16ScaledIntToFloat64 creates a 16 bit scaled reading transform
func MakeRTU16ScaledIntToFloat64(scalar float64) RTUTransform {
	return RTUTransform(func(b []byte) float64 {
		return rtuScaledInt16ToFloat64(b, scalar)
	})
}
