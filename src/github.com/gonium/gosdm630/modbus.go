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
	/***
	 * Opcodes as defined by Eastron.
	 * See http://bg-etech.de/download/manual/SDM630Register.pdf
	 * Please note that this is the superset of all SDM devices - some
	 * opcodes might not work on some devices.
	 */
	OpCodeSDML1Voltage   = 0x0000
	OpCodeSDML2Voltage   = 0x0002
	OpCodeSDML3Voltage   = 0x0004
	OpCodeSDML1Current   = 0x0006
	OpCodeSDML2Current   = 0x0008
	OpCodeSDML3Current   = 0x000A
	OpCodeSDML1Power     = 0x000C
	OpCodeSDML2Power     = 0x000E
	OpCodeSDML3Power     = 0x0010
	OpCodeSDML1Import    = 0x015a
	OpCodeSDML2Import    = 0x015c
	OpCodeSDML3Import    = 0x015e
	OpCodeSDMTotalImport = 0x0048
	OpCodeSDML1Export    = 0x0160
	OpCodeSDML2Export    = 0x0162
	OpCodeSDML3Export    = 0x0164
	OpCodeSDMTotalExport = 0x004a
	OpCodeSDML1Cosphi    = 0x001e
	OpCodeSDML2Cosphi    = 0x0020
	OpCodeSDML3Cosphi    = 0x0022
	//OpCodeL1THDCurrent         = 0x00F0
	//OpCodeL2THDCurrent         = 0x00F2
	//OpCodeL3THDCurrent         = 0x00F4
	//OpCodeAvgTHDCurrent        = 0x00Fa
	OpCodeSDML1THDVoltageNeutral  = 0x00ea
	OpCodeSDML2THDVoltageNeutral  = 0x00ec
	OpCodeSDML3THDVoltageNeutral  = 0x00ee
	OpCodeSDMAvgTHDVoltageNeutral = 0x00F8
	OpCodeSDMFrequency   = 0x0046

	/***
	 * Opcodes for Janitza B23.
	 * See https://www.janitza.de/betriebsanleitungen.html?file=files/download/manuals/current/B-Series/MID-Energy-Meters-Product-Manual.pdf
	 */
	OpCodeJanitzaL1Voltage   = 0x4A38
	OpCodeJanitzaL2Voltage   = 0x4A3A
	OpCodeJanitzaL3Voltage   = 0x4A3C
	OpCodeJanitzaL1Current   = 0x4A44
	OpCodeJanitzaL2Current   = 0x4A46
	OpCodeJanitzaL3Current   = 0x4A48
	OpCodeJanitzaL1Power     = 0x4A4C
	OpCodeJanitzaL2Power     = 0x4A4E
	OpCodeJanitzaL3Power     = 0x4A50
	OpCodeJanitzaL1Import    = 0x4A76
	OpCodeJanitzaL2Import    = 0x4A78
	OpCodeJanitzaL3Import    = 0x4A7A
	OpCodeJanitzaTotalImport = 0x4A7C
	OpCodeJanitzaL1Export    = 0x4A7E
	OpCodeJanitzaL2Export    = 0x4A80
	OpCodeJanitzaL3Export    = 0x4A82
	OpCodeJanitzaTotalExport = 0x4A84
	OpCodeJanitzaL1Cosphi    = 0x4A64
	OpCodeJanitzaL2Cosphi    = 0x4A66
	OpCodeJanitzaL3Cosphi    = 0x4A68

	/***
	 * Opcodes for DZG DVH4014.
	 * See "User Manual DVH4013", not public.
	 */
	OpCodeDZGTotalImportPower = 0x0000
	OpCodeDZGTotalExportPower = 0x0002
	OpCodeDZGL1Voltage        = 0x0004
	OpCodeDZGL2Voltage        = 0x0006
	OpCodeDZGL3Voltage        = 0x0008
	OpCodeDZGL1Current        = 0x000A
	OpCodeDZGL2Current        = 0x000C
	OpCodeDZGL3Current        = 0x000E
	OpCodeDZGL1Import         = 0x4020
	OpCodeDZGL2Import         = 0x4040
	OpCodeDZGL3Import         = 0x4060
	OpCodeDZGTotalImport      = 0x4000
	OpCodeDZGL1Export         = 0x4120
	OpCodeDZGL2Export         = 0x4140
	OpCodeDZGL3Export         = 0x4160
	OpCodeDZGTotalExport      = 0x4100
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

func (q *ModbusEngine) retrieveOpCode(deviceid uint8, funccode uint8,
	opcode uint16) (retval []byte, err error) {
	q.status.IncreaseModbusRequestCounter()
	// update the slave id in the handler
	q.handler.SlaveId = deviceid
	switch funccode {
	case ReadInputReg:
		retval, err = q.client.ReadInputRegisters(opcode, 2)
	case ReadHoldingReg:
		retval, err = q.client.ReadHoldingRegisters(opcode, 2)
	default:
		log.Fatalf("Unknown function code %d - cannot query device.",
			funccode)
	}
	if err != nil && q.verbose {
		log.Printf("Failed to retrieve opcode 0x%x, error was: %s\r\n", opcode, err.Error())
	}
	return retval, err
}

func (q *ModbusEngine) Transform(
	inputStream QuerySnipChannel,
	controlStream ControlSnipChannel,
	outputStream QuerySnipChannel,
) {
	var previousDeviceId uint8 = 0
	for {
		snip := <-inputStream
		// The SDM devices need to have a little pause between querying
		// different devices.
		if previousDeviceId != snip.DeviceId {
			time.Sleep(time.Duration(100) * time.Millisecond)
		}
		//if snip.OpCode == 0x00 {
		//	log.Printf("Skipping invalid Snip %+v", snip)
		//} else {
		//log.Printf("Executing Snip %+v", snip)
		previousDeviceId = snip.DeviceId
		//		value := q.queryOrFail(snip.DeviceId, snip.FuncCode, snip.OpCode, errorStream)

		var err error
		var reading []byte
		tryCnt := 0
		for tryCnt = 0; tryCnt < MaxRetryCount; tryCnt++ {
			reading, err = q.retrieveOpCode(snip.DeviceId, snip.FuncCode, snip.OpCode)
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
			//Now: convert bytes to value. Assume float64 conversion by
			//default.
			if snip.Transform != nil {
				snip.Value = snip.Transform(reading)
			} else {
				snip.Value = rtuToFloat64(reading)
			}
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
		BusId      uint8
		DeviceType MeterType
	}
	devicelist := make([]Device, 0)
	oldtimeout := q.handler.Timeout
	q.handler.Timeout = 50 * time.Millisecond
	log.Printf("Starting bus scan")
	// loop over all valid slave adresses
	var devid uint8
	for devid = 1; devid <= 247; devid++ {
		// Check if a SDM device responds: try to query L1 voltage
		voltage_L1, err := q.retrieveOpCode(devid, ReadInputReg, OpCodeSDML1Voltage)
		if err == nil {
			log.Printf("Device %d: SDM type device found, L1 voltage: %.2f\r\n", devid, rtuToFloat64(voltage_L1))
			dev := Device{
				BusId:      devid,
				DeviceType: METERTYPE_SDM,
			}
			devicelist = append(devicelist, dev)
		} else {
			// Check if a Janitza device responds: try to query L1 voltage
			voltage_L1, err := q.retrieveOpCode(devid, ReadHoldingReg, OpCodeJanitzaL1Voltage)
			if err == nil {
				log.Printf("Device %d: Janitza type device found, L1 voltage: %.2f\r\n", devid, rtuToFloat64(voltage_L1))
				dev := Device{
					BusId:      devid,
					DeviceType: METERTYPE_JANITZA,
				}
				devicelist = append(devicelist, dev)
			} else {
				// Check if a Janitza device responds: try to query L1 voltage
				voltage_L1, err := q.retrieveOpCode(devid, ReadHoldingReg, OpCodeDZGL1Voltage)
				if err == nil {
					log.Printf("Device %d: DZG type device found, L1 voltage: %.2f\r\n", devid, rtuToFloat64(voltage_L1))
					dev := Device{
						BusId:      devid,
						DeviceType: METERTYPE_DZG,
					}
					devicelist = append(devicelist, dev)
				} else {
					log.Printf("Device %d: n/a\r\n", devid)
				}
			}
		}
		// give the bus some time to recover before querying the next device
		time.Sleep(time.Duration(40) * time.Millisecond)
	}
	// restore timeout to old value
	q.handler.Timeout = oldtimeout
	log.Printf("Found %d active devices:\r\n", len(devicelist))
	for _, device := range devicelist {
		log.Printf("* slave address %d: type %s\r\n", device.BusId,
			device.DeviceType)
	}
	log.Println("WARNING: This lists only the devices that responded to " +
		"a known L1 voltage request. Devices with " +
		"different function code definitions might not be detected.")
}

// Transform functions to convert RTU bytes to meaningfull data types.
type RTUTransform func([]byte) float64

func rtuScaledIntToFloat64(b []byte, scalar float64) float64 {
	unscaled := float64(binary.BigEndian.Uint32(b))
	f := unscaled / scalar
	return float64(f)
}

func MkRTUScaledIntToFloat64(scalar float64) RTUTransform {
	return RTUTransform(func(b []byte) float64 {
		return rtuScaledIntToFloat64(b, scalar)
	})
}

func rtuToFloat64(b []byte) float64 {
	bits := binary.BigEndian.Uint32(b)
	f := math.Float32frombits(bits)
	return float64(f)
}
