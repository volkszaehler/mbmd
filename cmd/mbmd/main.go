package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	latest "github.com/tcnksm/go-latest"
	. "github.com/volkszaehler/mbmd/meters"
	"github.com/volkszaehler/mbmd/meters/bus"
	_ "github.com/volkszaehler/mbmd/meters/impl"
	. "github.com/volkszaehler/mbmd/server"
	cli "gopkg.in/urfave/cli.v1"
)

const (
	DEFAULT_METER_STORE_SECONDS = 120 * time.Second
)

func checkVersion() {
	githubTag := &latest.GithubTag{
		Owner:      "volkszaehler",
		Repository: "mbmd",
	}

	if res, err := latest.Check(githubTag, Version); err == nil {
		if res.Outdated {
			log.Printf("updates available - please upgrade to ingress %s", res.Current)
		}
	}
}

func createMeters(deviceslice []string) map[uint8]*Meter {
	meters := make(map[uint8]*Meter)
	for _, meterdef := range deviceslice {
		splitdef := strings.Split(meterdef, ":")
		if len(splitdef) != 2 {
			log.Fatalf("Cannot parse device definition %s. See -h for help.", meterdef)
		}
		metertype, devid := splitdef[0], splitdef[1]
		id, err := strconv.Atoi(devid)
		if err != nil {
			log.Fatalf("Error parsing device id %s: %s. See -h for help.", meterdef, err.Error())
		}
		meter, err := NewMeterByType(metertype, uint8(id))
		if err != nil {
			log.Fatalf("Unknown meter type %s for device %d. See -h for help.", metertype, id)
		}
		meters[uint8(id)] = meter
	}
	return meters
}

func waitForSignal(signals ...os.Signal) {
	var wg sync.WaitGroup
	wg.Add(1)
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, signals...)
	go func() {
		<-channel
		wg.Done()
	}()
	wg.Wait()
}

func meterHelp() string {
	var s string
	for _, c := range []ConnectionType{RS485, TCP} {
		s += fmt.Sprintf("\n\t\t\t\t%s", c.String())
		// s += fmt.Sprintf("\n\t\t\t\t%s", strings.Repeat("-", len(c.String())))

		types := make([]string, 0)
		for t, f := range Producers {
			p := f()
			if c != p.ConnectionType() {
				continue
			}

			types = append(types, t)
		}

		sort.Strings(types)

		for _, t := range types {
			f := Producers[t]
			p := f()
			s += fmt.Sprintf("\n\t\t\t\t %-9s%s", t, p.Description())
		}
	}
	return s
}

func logMeterDetails(meters map[uint8]*Meter, qe *ModbusEngine) {
	for devid, meter := range meters {
		producer := meter.Producer
		log.Printf("#%d: %s (%s)", devid, producer.Description(), producer.ConnectionType())

		if sunspec, ok := producer.(SunSpecProducer); ok {
			op := sunspec.GetSunSpecCommonBlock()
			snip := QuerySnip{
				DeviceId:  meter.DeviceId,
				Operation: op,
			}

			if b, err := qe.Query(snip); err == nil {
				if descriptor, err := sunspec.DecodeSunSpecCommonBlock(b); err == nil {
					log.Printf("    Manufacturer: %s", descriptor.Manufacturer)
					log.Printf("    Model:        %s", descriptor.Model)
					log.Printf("    Options:      %s", descriptor.Options)
					log.Printf("    Version:      %s", descriptor.Version)
					log.Printf("    Serial:       %s", descriptor.Serial)
				} else {
					log.Println(err)
				}
			} else {
				log.Println(err)
			}
		}
	}
}


func main() {
	app := cli.NewApp()
	app.Name = "mbmd"
	app.Usage = "ModBus Measurement Daemon"
	app.Version = fmt.Sprintf("%s (https://github.com/volkszaehler/mbmd/commit/%s)", Version, Commit)
	app.HideVersion = true
	app.Flags = []cli.Flag{
		// general
		cli.StringFlag{
			Name:  "adapter, a",
			Value: "/dev/ttyUSB0",
			Usage: "MODBUS adapter - can be either serial RTU device (/dev/ttyUSB0) or TCP socket (localhost:502)",
		},
		cli.IntFlag{
			Name:  "comset, c",
			Value: ModbusComset9600_8N1,
			Usage: `Communication parameters:
			` + strconv.Itoa(ModbusComset2400_8N1) + `:  2400 baud, 8N1
			` + strconv.Itoa(ModbusComset9600_8N1) + `:  9600 baud, 8N1
			` + strconv.Itoa(ModbusComset19200_8N1) + `: 19200 baud, 8N1
			` + strconv.Itoa(ModbusComset2400_8E1) + `:  2400 baud, 8E1
			` + strconv.Itoa(ModbusComset9600_8E1) + `:  9600 baud, 8E1
			` + strconv.Itoa(ModbusComset19200_8E1) + `: 19200 baud, 8E1
			`,
		},
		cli.BoolFlag{
			Name:  "simulate, s",
			Usage: "Simulate MODBUS device for testing purposes",
		},
		cli.StringFlag{
			Name:  "devices, d",
			Value: "SDM:1",
			Usage: `MODBUS device type and ID to query, separated by comma.
			Valid types are:` + meterHelp() + `
			Example: -d JANITZA:1,SDM:22,DZG:23`,
		},
		cli.BoolFlag{
			Name:  "detect",
			Usage: "Detect MODBUS devices",
		},
		cli.StringFlag{
			Name:  "rate, r",
			Value: "1s",
			Usage: "Maximum update rate in seconds per message, 0 is unlimited",
			// Destination: &mqttRate,
		},
		cli.StringFlag{
			Name:  "idformat, f",
			Value: "Meter%d",
			Usage: `Meter id format. Determines meter id for REST and MQTT.
			Example: -f Meter%d
			The %d is replaced by the device ID`,
		},
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "print verbose messages",
		},

		// http api
		cli.StringFlag{
			Name:  "url, u",
			Value: "localhost:8080",
			Usage: "REST API url. Use 0.0.0.0:8080 to accept incoming connections.",
		},

		// mqtt api
		cli.StringFlag{
			Name:  "broker, b",
			Value: "",
			Usage: "MQTT: Broker URI. ex: tcp://10.10.1.1:1883",
			// Destination: &mqttBroker,
		},
		cli.StringFlag{
			Name:  "topic, t",
			Value: "mbmd",
			Usage: "MQTT: Base topic. Set empty to disable publishing.",
			// Destination: &mqttTopic,
		},
		cli.StringFlag{
			Name:  "user",
			Value: "",
			Usage: "MQTT: User (optional)",
			// Destination: &mqttUser,
		},
		cli.StringFlag{
			Name:  "password",
			Value: "",
			Usage: "MQTT: Password (optional)",
			// Destination: &mqttPassword,
		},
		cli.StringFlag{
			Name:  "clientid, i",
			Value: "mbmd",
			Usage: "MQTT: ClientID",
			// Destination: &mqttClientID,
		},
		cli.BoolFlag{
			Name:  "clean, l",
			Usage: "MQTT: Set Clean Session (default: false)",
			// Destination: &mqttCleanSession,
		},
		cli.IntFlag{
			Name:  "qos, q",
			Value: 0,
			Usage: "MQTT: Quality of Service 0,1,2",
			// Destination: &mqttQos,
		},
		cli.StringFlag{
			Name:  "homie",
			Value: "homie",
			Usage: "MQTT: Homie IoT discovery base topic (homieiot.github.io). Set empty to disable.",
		},
	}

	app.Action = func(c *cli.Context) {
		if c.NArg() > 0 {
			log.Fatalf("Unexpected arguments: %v", c.Args())
		}

		log.Printf("mbmd %s %s", Version, Commit)
		go checkVersion()

		// Set unique ID format
		UniqueIdFormat = c.String("idformat")

		// Parse the devices parameter
		meters := createMeters(strings.Split(c.String("devices"), ","))

		rate, err := time.ParseDuration(c.String("rate"))
		if err != nil {
			log.Fatalf("Invalid rate %s", err)
		}

		// create ModbusEngine with status
		status := NewStatus(meters)
		qe := NewModbusEngine(
			c.String("adapter"),
			c.Int("comset"),
			c.Bool("simulate"),
			c.Bool("verbose"),
			status,
		)

		// detect command
		if c.Bool("detect") {
			qe.Scan()
			return
		}

		// scheduler and meter data channel
		scheduler, snips := SetupScheduler(meters, qe)
		logMeterDetails(meters, qe)

		// tee that broadcasts meter messages to multiple recipients
		tee := NewQuerySnipBroadcaster(snips)
		go tee.Run()

		// websocket hub
		hub := NewSocketHub(status)
		tee.AttachRunner(hub.Run)

		// MQTT client
		if c.String("broker") != "" {
			mqtt := NewMqttClient(
				c.String("broker"),
				c.String("topic"),
				c.String("user"),
				c.String("password"),
				c.String("clientid"),
				c.Int("qos"),
				c.Bool("clean"),
				c.Bool("verbose"))

			// homie needs to scan the bust, start it first
			if c.String("homie") != "" {
				homieRunner := HomieRunner{MqttClient: mqtt}
				homieRunner.Register(c.String("homie"), meters, qe)
				tee.AttachRunner(homieRunner.Run)
			}

			// start "normal" mqtt handler after homie setup
			if c.String("topic") != "" {
				mqttRunner := MqttRunner{MqttClient: mqtt}
				tee.AttachRunner(mqttRunner.Run)
			}
		}

		// MeasurementCache for REST API
		mc := NewMeasurementCache(
			meters,
			scheduler,
			DEFAULT_METER_STORE_SECONDS,
			c.Bool("verbose"),
		)
		tee.AttachRunner(mc.Run)

		// start the scheduler
		ctx, cancelScheduler := context.WithCancel(context.Background())
		go scheduler.Run(ctx, rate)

		// handle os signals and gracefully exit Run methods
		go func() {
			waitForSignal(os.Interrupt, os.Kill)
			log.Println("Received signal - stopping")
			cancelScheduler() // cancel scheduler

			// wait for Run methods attached to tee to finish
			<-tee.Done()
			log.Println("Stopped")
			os.Exit(0)
		}()

		Run_httpd(
			mc,
			hub,
			status,
			c.String("url"),
		)
	}

	_ = app.Run(os.Args)
}
