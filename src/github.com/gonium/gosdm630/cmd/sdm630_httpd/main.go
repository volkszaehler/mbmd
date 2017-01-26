package main

import (
	"github.com/gonium/gosdm630"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "sdm630_httpd"
	app.Usage = "SDM630 power measurements via HTTP."
	app.Version = "0.2.0"
	app.HideVersion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "serialadapter, s",
			Value: "/dev/ttyUSB0",
			Usage: "path to serial RTU device",
		},
		cli.StringFlag{
			Name:  "url, u",
			Value: ":8080",
			Usage: "the URL the server should respond on",
		},
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "print verbose messages",
		},
		cli.StringFlag{
			Name:  "device_list, d",
			Value: "1",
			Usage: "MODBUS device ID to query",
		},
	}
	app.Action = func(c *cli.Context) {
		status := sdm630.NewStatus()
		var rc = make(sdm630.ReadingChannel)
		// Parse the device_list parameter
		deviceslice := strings.Split(c.String("device_list"), ",")
		devids := make([]uint8, 0, len(deviceslice))
		for _, devid := range deviceslice {
			id, err := strconv.Atoi(devid)
			if err != nil {
				log.Fatalf("Error parsing device id %s: %s", devid, err.Error())
			}
			devids = append(devids, uint8(id))
		}
		qe := sdm630.NewQueryEngine(
			c.String("serialadapter"),
			c.Bool("verbose"),
			rc,
			devids,
			status,
		)
		go qe.Produce()
		mc := sdm630.NewMeasurementCache(
			rc,
			120*time.Second, // TODO: How long to store data in the cache?.
			c.Bool("verbose"),
		)
		go mc.ConsumeData()
		log.Printf("Starting API httpd at %s", c.String("url"))
		sdm630.Run_httpd(mc, status, c.String("url"))
	}

	app.Run(os.Args)
}
