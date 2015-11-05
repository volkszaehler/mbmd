package main

import (
	"github.com/codegangsta/cli"
	"github.com/gonium/gosdm630"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "sdm630_httpd"
	app.Usage = "SDM630 power measurements via HTTP."
	app.Version = "0.1.0"
	app.HideVersion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "device, d",
			Value: "/dev/ttyUSB0",
			Usage: "path to serial RTU device",
		},
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "print verbose messages",
		},
	}
	app.Action = func(c *cli.Context) {
		var rc = make(sdm630.ReadingChannel)
		// Read the SDM630 every 5 seconds
		interval := 5
		qe := sdm630.NewQueryEngine(c.String("device"), interval,
			c.Bool("verbose"), rc)
		go qe.Produce()
		mc := sdm630.NewMeasurementCache(rc, interval)
		go mc.ConsumeData()
		sdm630.Run_httpd(mc)
	}

	app.Run(os.Args)
}
