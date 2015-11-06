package main

import (
	"github.com/codegangsta/cli"
	"github.com/gonium/gosdm630"
	"log"
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
	}
	app.Action = func(c *cli.Context) {
		// Check the interval - only values between 5 and 20 are
		// useful.
		if c.Int("interval") < 5 || c.Int("interval") > 20 {
			log.Fatal("the interval must be between 5 and 20 seconds.")
		}
		var rc = make(sdm630.ReadingChannel)
		qe := sdm630.NewQueryEngine(
			c.String("device"),
			c.Int("interval"),
			c.Bool("verbose"),
			rc,
		)
		go qe.Produce()
		mc := sdm630.NewMeasurementCache(rc, c.Int("interval"))
		go mc.ConsumeData()
		sdm630.Run_httpd(mc, c.String("url"))
	}

	app.Run(os.Args)
}
