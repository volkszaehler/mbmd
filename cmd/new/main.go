package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/volkszaehler/mbmd/meters"
	"github.com/volkszaehler/mbmd/meters/connection"
	"github.com/volkszaehler/mbmd/meters/rs485"
	"github.com/volkszaehler/mbmd/meters/sunspec"
	"github.com/volkszaehler/mbmd/server"
)

var managers map[string]connection.Manager

func init() {
	managers = make(map[string]connection.Manager, 0)
}

// createConnection parses adapter string to create TCP or RTU connection
func createConnection(device string) (res connection.Connection) {
	if tcp, _ := regexp.MatchString(":[0-9]+$", device); tcp {
		res = connection.NewTCP(device) // tcp connection
	} else {
		res = connection.NewRTU(device, 1) // serial connection
	}
	return res
}

// createDevice creates new device and adds it to the
func createDevice(deviceDef string, defaultDevice string) {
	deviceSplit := strings.Split(deviceDef, "@")
	if len(deviceSplit) == 0 || len(deviceSplit) > 2 {
		log.Fatalf("Cannot parse connect string %s. See -h for help.", deviceDef)
	}

	meterDef := deviceSplit[0]
	connSpec := defaultDevice
	if len(deviceSplit) == 2 {
		connSpec = deviceSplit[1]
	}

	if connSpec == "" {
		log.Fatalf("Cannot parse connect string- missing physical device or connection for %s. See -h for help.", deviceDef)
	}

	manager, ok := managers[connSpec]
	if !ok {
		conn := createConnection(connSpec)
		manager = connection.NewManager(conn)
		managers[connSpec] = manager
	}

	meterSplit := strings.Split(meterDef, ":")
	if len(meterSplit) != 2 {
		log.Fatalf("Cannot parse device definition: %s. See -h for help.", meterDef)
	}

	meterType, devID := meterSplit[0], meterSplit[1]
	if len(strings.TrimSpace(meterType)) == 0 {
		log.Fatalf("Cannot parse device definition- meter type empty: %s. See -h for help.", meterDef)
	}

	id, err := strconv.Atoi(devID)
	if err != nil {
		log.Fatalf("Error parsing device id %s: %v. See -h for help.", devID, err)
	}

	var meter meters.Device
	if _, ok := manager.Conn.(*connection.TCP); ok {
		meter = sunspec.NewDevice()
	} else {
		meter, err = rs485.NewDevice(meterType)
		if err != nil {
			log.Fatalf("Error creating device %s: %v. See -h for help.", meterDef, err)
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

	// status := make(map[string]server.MeterStatus)
	qe := server.NewQueryEngine(managers)

	cc := make(chan server.ControlSnip)
	rc := make(chan server.QuerySnip)

	// tee that broadcasts meter messages to multiple recipients
	tee := server.NewQuerySnipBroadcaster(rc)
	go tee.Run()

	status := server.NewStatus(cc)

	// websocket hub
	hub := server.NewSocketHub(status) // status
	tee.AttachRunner(hub.Run)

	cache := server.NewCache(time.Second, true) // measurement cache
	tee.AttachRunner(cache.Run)

	httpd := &server.Httpd{}
	go httpd.Run(cache, hub, nil, "127.1:8080")

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan bool)
	go func() {
		qe.Run(ctx, cc, rc)
		println("qe done")
		done <- true
		return
	}()

	// send signals to exit channel
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, os.Kill)

	timer := time.After(5000 * time.Second)

loop:
	for {
		time.Sleep(50 * time.Millisecond)
		select {
		// case <-cc:
		// case <-rc:
		case <-exit:
			println("exit")
			cancel()
		case <-timer:
			println("timer")
			cancel()
		case <-done:
			// stop processing only if sender has closed
			break loop
		}
	}

	println("main loop done")
}
