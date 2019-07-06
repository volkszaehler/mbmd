package main

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/volkszaehler/mbmd/meters"
	"github.com/volkszaehler/mbmd/meters/connection"
	"github.com/volkszaehler/mbmd/meters/rs485"
	"github.com/volkszaehler/mbmd/meters/sunspec"
)

var managers map[string]connection.Manager

func init() {
	managers = make(map[string]connection.Manager, 0)
}

func createConnection(device string) (res connection.Connection) {
	// parse adapter string
	if tcp, _ := regexp.MatchString(":[0-9]+$", device); tcp {
		res = connection.NewTCP(device) // tcp connection
	} else {
		res = connection.NewRTU(device, 1) // serial connection
	}
	return res
}

func createDevice(deviceDef string, defaultDevice string) {
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

	manager, ok := managers[device]
	if !ok {
		conn := createConnection(device)
		manager = connection.NewManager(conn)
		managers[device] = manager
	}

	meterSplit := strings.Split(meterDef, ":")
	if len(meterSplit) != 2 {
		log.Fatalf("Cannot parse device definition %s. See -h for help.", meterDef)
	}

	meterType, devID := meterSplit[0], meterSplit[1]
	id, err := strconv.Atoi(devID)
	if err != nil {
		log.Fatalf("Error parsing device id %s: %s. See -h for help.", meterDef, err.Error())
	}

	var meter meters.Device
	if _, ok := manager.Conn.(*connection.TCP); ok {
		meter = sunspec.NewDevice()
	} else {
		meter, err = rs485.NewDevice(meterType)
		if err != nil {
			log.Fatalf("Error creating device %s: %s. See -h for help.", meterDef, err.Error())
		}
	}

	manager.Add(uint8(id), meter)
}

func main() {
	devices := []string{
		"sma:126@localhost:5061",
		"sma:126@localhost:5062",
		"sma:126@localhost:5063",
	}

	for _, dev := range devices {
		createDevice(dev, "")
	}

	println("Init...")
	var wg sync.WaitGroup
	for _, m := range managers {
		wg.Add(1)

		go func(m connection.Manager) {
			m.All(func(id uint8, dev meters.Device) {
				if err := dev.Initialize(m.Conn.ModbusClient()); err != nil {
					log.Fatalf("initalizing %d at %s failed: %v", id, m.Conn, err)
				}
			})
			wg.Done()
		}(m)
	}
	wg.Wait()

	println("Found...")
	for _, m := range managers {
		wg.Add(1)

		go func(m connection.Manager) {
			m.All(func(id uint8, dev meters.Device) {
				desc := dev.Descriptor()
				log.Printf("%v", desc)
			})
			wg.Done()
		}(m)
	}
	wg.Wait()

	println("Probe...")
	for _, m := range managers {
		wg.Add(1)

		go func(m connection.Manager) {
			m.All(func(id uint8, dev meters.Device) {
				if val, err := dev.Probe(m.Conn.ModbusClient()); err != nil {
					log.Fatalf("probing %d at %s failed: %v", id, m.Conn, err)
				} else {
					log.Printf("%v", val)
				}
			})
			wg.Done()
		}(m)
	}
	wg.Wait()

	println("Running...")
	for _, m := range managers {
		m.Run()
	}
}
