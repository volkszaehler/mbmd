package main

import (
	"github.com/gonium/gosdm630"
	"gopkg.in/urfave/cli.v1"
	"os"
	"strconv"
)

func main() {
	app := cli.NewApp()
	app.Name = "sdm_detect"
	app.Usage = "Attempts to detect available SDM devices."
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
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "print verbose messages",
		},
	}
	app.Action = func(c *cli.Context) {
		status := sdm630.NewStatus(nil)
		qe := sdm630.NewModbusEngine(
			c.String("serialadapter"),
			c.Int("comset"),
			c.Bool("verbose"),
			status,
		)

		qe.Scan()
	}

	app.Run(os.Args)
}
