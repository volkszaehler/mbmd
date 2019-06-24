package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/volkszaehler/mbmd/meters"
	"github.com/volkszaehler/mbmd/meters/bus"
	"github.com/volkszaehler/mbmd/meters/rs485"
	"github.com/volkszaehler/mbmd/meters/sunspec"
)

func createBusConnector(device string) (res bus.Bus) {
	// parse adapter string
	if tcp, _ := regexp.MatchString(":[0-9]+$", device); tcp {
		res = bus.NewTCP(device) // tcp connection
	} else {
		res = bus.NewRTU(device, 1) // serial connection
	}
	return res
}

func createDevices(connections []string, defaultDevice string) {
	devices := make(map[string]bus.Bus, 0)

	for _, deviceDef := range connections {
		deviceSplit := strings.Split(deviceDef, "@")
		if len(deviceSplit) == 0 || len(deviceSplit) > 2 {
			log.Fatalf("Cannot parse connect string %s. See -h for help.", deviceDef)
		}

		meterDef := deviceSplit[0]
		device := defaultDevice
		if len(deviceSplit) == 2 {
			device = deviceSplit[1]
		}

		if device == "" {
			log.Fatalf("Cannot parse connect string- missing physical device or connection for %s. See -h for help.", deviceDef)
		}

		busDevice := createBusConnector(device)
		devices[device] = busDevice

		meterSplit := strings.Split(meterDef, ":")
		if len(meterSplit) != 2 {
			log.Fatalf("Cannot parse device definition %s. See -h for help.", meterDef)
		}

		meterType, devID := meterSplit[0], meterSplit[1]
		id, err := strconv.Atoi(devID)
		if err != nil {
			log.Fatalf("Error parsing device id %s: %s. See -h for help.", meterDef, err.Error())
		}

		_ = meterType
		_ = id
		_ = busDevice

		fmt.Println(devices)

		var meter meters.Device
		if _, ok := busDevice.(*bus.TCP); ok {
			meter = sunspec.NewDevice()
		} else {
			meter, err = rs485.NewDevice(meterType)
			if err != nil {
				log.Fatalf("Error creating device %s: %s. See -h for help.", meterDef, err.Error())
			}
		}

		_ = meter

		// test
		log.Println(id)
		busDevice.Slave(uint8(id))
		if err := meter.Initialize(busDevice.(*bus.TCP).Client); err != nil {
			log.Fatal(err)
		}

		if res, err := meter.Query(busDevice.(*bus.TCP).Client); err != nil {
			log.Fatal(err)
		} else {
			log.Println(res)
		}

		if err := busDevice.Add(uint8(id), meter); err != nil {
			log.Fatal(err)
		}

		// meter, err := NewMeterByType(meterType, uint8(id))
		// if err != nil {
		// 	log.Fatalf("Unknown meter type %s for device %d. See -h for help.", metertype, id)
		// }
		// meters[uint8(id)] = meter
	}

	for _, b := range devices {
		println("Running...")
		b.Run()
	}
}

func main() {
	createDevices([]string{
		"sma:126@localhost:5061",
		// "sma:126@localhost:5062",
		// "sma:126@localhost:5063",
	}, "")
}
