package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/viper"
	latest "github.com/tcnksm/go-latest"
	cli "gopkg.in/urfave/cli.v1"

	"github.com/volkszaehler/mbmd"
	"github.com/volkszaehler/mbmd/meters"
	"github.com/volkszaehler/mbmd/server"
)

const (
	cacheDuration = 1 * time.Minute
	defaultConfig = "mbmd.yaml"
)

func checkVersion() {
	githubTag := &latest.GithubTag{
		Owner:      "volkszaehler",
		Repository: "mbmd",
	}

	if res, err := latest.Check(githubTag, server.Version); err == nil {
		if res.Outdated {
			log.Printf("updates available - please upgrade to ingress %s", res.Current)
		}
	}
}

func meterHelp() string {
	return "FOO"
}

// func meterHelp() string {
// 	var s string
// 	for _, c := range []ConnectionType{RS485, TCP} {
// 		s += fmt.Sprintf("\n\t\t\t\t%s", c.String())
// 		// s += fmt.Sprintf("\n\t\t\t\t%s", strings.Repeat("-", len(c.String())))

// 		types := make([]string, 0)
// 		for t, f := range Producers {
// 			p := f()
// 			if c != p.ConnectionType() {
// 				continue
// 			}

// 			types = append(types, t)
// 		}

// 		sort.Strings(types)

// 		for _, t := range types {
// 			f := Producers[t]
// 			p := f()
// 			s += fmt.Sprintf("\n\t\t\t\t %-9s%s", t, p.Description())
// 		}
// 	}
// 	return s
// }

func main() {
	app := cli.NewApp()
	app.Name = "mbmd"
	app.Usage = "ModBus Measurement Daemon"
	app.Version = fmt.Sprintf("%s (https://github.com/volkszaehler/mbmd/commit/%s)", server.Version, server.Commit)
	app.HideVersion = true
	app.Flags = []cli.Flag{
		// general
		cli.StringFlag{
			Name:  "adapter, a",
			Value: "/dev/ttyUSB0",
			Usage: "ModBus adapter - can be either serial RTU device (/dev/ttyUSB0) or TCP socket (localhost:502)",
		},
		cli.IntFlag{
			Name:  "comset",
			Value: meters.Comset9600_8N1,
			Usage: `Communication parameters:
			` + strconv.Itoa(meters.Comset2400_8N1) + `:  2400 baud, 8N1
			` + strconv.Itoa(meters.Comset9600_8N1) + `:  9600 baud, 8N1
			` + strconv.Itoa(meters.Comset19200_8N1) + `: 19200 baud, 8N1
			` + strconv.Itoa(meters.Comset2400_8E1) + `:  2400 baud, 8E1
			` + strconv.Itoa(meters.Comset9600_8E1) + `:  9600 baud, 8E1
			` + strconv.Itoa(meters.Comset19200_8E1) + `: 19200 baud, 8E1
			`,
		},
		cli.StringFlag{
			Name:  "config, c",
			Value: defaultConfig,
			Usage: `Configuration file`,
		},
		cli.StringFlag{
			Name:  "devices, d",
			Value: "",
			Usage: `MODBUS device type and ID to query, separated by comma.
			Append device or address separated by @.
			Valid types are:` + meterHelp() + `
			Example: -d SDM:22,SMA:126@localhost:502`,
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
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "print verbose messages",
		},

		// http api
		cli.StringFlag{
			Name:  "url, u",
			Value: ":8080",
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

		log.Printf("mbmd %s %s", server.Version, server.Commit)
		go checkVersion()

		// rate, err := time.ParseDuration(c.String("rate"))
		// if err != nil {
		// 	log.Fatalf("Invalid rate %s", err)
		// }

		// detect command
		// if c.Bool("detect") {
		// 	qe.Scan()
		// 	return
		// }

		var conf mbmd.Config
		viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")

		if configFile := c.String("config"); configFile != defaultConfig {
			log.Println("SetConfigFile")
			viper.SetConfigFile(configFile)
		} else {
			viper.SetConfigName("mbmd")  // name of config file (without extension)
			viper.AddConfigPath("/etc")  // path to look for the config file in
			viper.AddConfigPath("$HOME") // call multiple times to add many search paths
			viper.AddConfigPath(".")     // optionally look for config in the working directory
		}

		confHandler := mbmd.NewDeviceConfigHandler()
		if err := viper.ReadInConfig(); err != nil { // handle errors reading the config file
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				log.Fatal(err)
			}
		} else {
			// config file found
			log.Printf("using %s", viper.ConfigFileUsed())
			if err := viper.Unmarshal(&conf); err != nil {
				log.Fatalf("failed parsing config file: %v", err)
			}

			confHandler.DefaultDevice = conf.Default.Adapter
			for _, dev := range conf.Devices {
				confHandler.CreateDevice(dev)
			}
		}

		// remaining command line options
		confHandler.DefaultDevice = c.String("adapter")
		for _, dev := range strings.Split(c.String("devices"), ",") {
			if dev != "" {
				confHandler.CreateDeviceFromSpec(dev)
			}
		}

		// query engine
		qe := server.NewQueryEngine(confHandler.Managers)

		// result channels
		rc := make(chan server.QuerySnip)
		cc := make(chan server.ControlSnip)

		// tee that broadcasts meter messages to multiple recipients
		tee := server.NewQuerySnipBroadcaster(rc)
		go tee.Run()

		// status cache
		status := server.NewStatus(cc)

		// websocket hub
		hub := server.NewSocketHub(status)
		tee.AttachRunner(hub.Run)

		// measurement cache for REST api
		cache := server.NewCache(cacheDuration, status, c.Bool("verbose"))
		tee.AttachRunner(cache.Run)

		httpd := server.NewHttpd(qe)
		go httpd.Run(cache, hub, status, c.String("url"))

		// MQTT client
		if c.String("broker") != "" {
			mqtt := server.NewMqttClient(
				c.String("broker"),
				c.String("topic"),
				c.String("user"),
				c.String("password"),
				c.String("clientid"),
				c.Int("qos"),
				c.Bool("clean"),
				c.Bool("verbose"),
			)

			// homie needs to scan the bus, start it first
			if c.String("homie") != "" {
				homieRunner := server.NewHomieRunner(mqtt, qe, c.String("homie"))
				// homieRunner.Register(c.String("homie"), meters, qe)
				tee.AttachRunner(homieRunner.Run)
			}

			// start "normal" mqtt handler after homie setup
			if c.String("topic") != "" {
				mqttRunner := server.MqttRunner{MqttClient: mqtt}
				tee.AttachRunner(mqttRunner.Run)
			}
		}

		ctx, cancel := context.WithCancel(context.Background())
		go qe.Run(ctx, cc, rc)

		// wait for signal on exit channel and cancel context
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		<-exit
		log.Println("received signal - stopping")
		cancel()

		// wait for Run methods attached to tee to finish
		<-tee.Done()
		log.Println("stopped")
	}

	_ = app.Run(os.Args)
}
