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
	app.Version = sdm630.RELEASEVERSION
	app.HideVersion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "serialadapter, s",
			Value: "/dev/ttyUSB0",
			Usage: "path to serial RTU device",
		},
		cli.IntFlag{
			Name:  "comset, c",
			Value: sdm630.ModbusComset9600,
			Usage: `which communication parameter set to use. Valid sets are
		` + strconv.Itoa(sdm630.ModbusComset2400) + `:  2400 baud, 8N1
		` + strconv.Itoa(sdm630.ModbusComset9600) + `:  9600 baud, 8N1
		` + strconv.Itoa(sdm630.ModbusComset19200) + `: 19200 baud, 8N1
			`,
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
			Usage: `MODBUS device ID to query, separated by comma. 
			Example: -d 11,12,13`,
		},
		cli.StringFlag{
			Name: "unique_id_format, f",
			Value: "Instrument%d",
			Usage: `Unique ID format.
			Example: -f Instrument%d
			The %d is replaced by the device ID`,
		},
	}
	app.Action = func(c *cli.Context) {
		// Set unique ID format
		sdm630.UniqueIdFormat = c.String("unique_id_format")
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
		status := sdm630.NewStatus()

		// Create Channels that link the goroutines
		var scheduler2queryengine = make(sdm630.QuerySnipChannel)
		var queryengine2duplicator = make(sdm630.QuerySnipChannel)
		var duplicator2cache = make(sdm630.QuerySnipChannel)
		var duplicator2firehose = make(sdm630.QuerySnipChannel)

		scheduler := sdm630.NewRoundRobinScheduler(
			scheduler2queryengine,
			devids,
		)
		go scheduler.Produce()

		qe := sdm630.NewModbusEngine(
			c.String("serialadapter"),
			c.Int("comset"),
			c.Bool("verbose"),
			scheduler2queryengine,
			queryengine2duplicator,
			devids,
			status,
		)
		go qe.Transform()

		// This is the duplicator
		go func(in sdm630.QuerySnipChannel,
			out1 sdm630.QuerySnipChannel,
			out2 sdm630.QuerySnipChannel,
		) {
			for {
				snip := <-in
				out1 <- snip
				out2 <- snip
			}
		}(queryengine2duplicator, duplicator2cache, duplicator2firehose)

		firehose := sdm630.NewFirehose(duplicator2firehose,
			c.Bool("verbose"))
		go firehose.Run()

		mc := sdm630.NewMeasurementCache(
			duplicator2cache,
			120*time.Second, // TODO: How long to store data in the cache?.
			c.Bool("verbose"),
		)
		go mc.Consume()

		log.Printf("Starting API httpd at %s", c.String("url"))
		sdm630.Run_httpd(
			mc,
			firehose,
			status,
			c.String("url"),
		)
	}

	app.Run(os.Args)
}
