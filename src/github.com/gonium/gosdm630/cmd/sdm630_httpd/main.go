package main

import (
	"flag"
	"github.com/gonium/gosdm630"
)

var rtuDevice = flag.String("rtuDevice", "/dev/ttyUSB0", "Path to serial RTU device")
var verbose = flag.Bool("verbose", false, "Enables extensive logging")

func init() {
	flag.Parse()
}

func main() {
	var rc = make(sdm630.ReadingChannel)
	// Read the SDM630 every 5 seconds
	interval := 5
	qe := sdm630.NewQueryEngine(*rtuDevice, interval, *verbose, rc)
	go qe.Produce()
	mc := sdm630.NewMeasurementCache(rc, interval)
	go mc.ConsumeData()
	sdm630.Run_httpd(mc)
}
