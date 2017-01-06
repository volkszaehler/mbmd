package main

import (
	"github.com/codegangsta/cli"
	"github.com/gonium/gosdm630"
	"log"
	"os"
	"strconv"
	"strings"
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
		cli.IntFlag{
			Name:  "interval, i",
			Value: 10,
			Usage: "seconds between getting new values from the SDM630",
		},
		cli.StringFlag{
			Name:  "device_list, d",
			Value: "1",
			Usage: "MODBUS device ID to query",
		},
	}
	app.Action = func(c *cli.Context) {
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
		log.Println("Will query MODBUS IDs", devids)
		qe := sdm630.NewQueryEngine(
			c.String("serialadapter"),
			c.Int("interval"),
			c.Bool("verbose"),
			rc,
			devids,
		)
		go qe.Produce()
		mc := sdm630.NewMeasurementCache(
			rc,
			c.Int("interval"),
			c.Bool("verbose"),
		)
		go mc.ConsumeData()
		log.Println("Starting API httpd.")
		sdm630.Run_httpd(mc, c.String("url"))
	}

	app.Run(os.Args)
}
