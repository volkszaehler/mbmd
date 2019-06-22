package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/volkszaehler/mbmd/meters/sunspec"
	"github.com/volkszaehler/mbmd/server"

	"github.com/grid-x/modbus"
)

const (
	base = 40000
)

func main() {
	addr := os.Args[1]

	deviceID, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	handler := modbus.NewTCPClientHandler(addr)
	client := modbus.NewClient(handler)
	if err := handler.Connect(); err != nil {
		log.Fatal("Failed to connect: ", err)
	}

	mbl := &server.ModbusLogger{}
	_ = mbl
	// handler.Logger = mbl

	handler.SetSlave(byte(deviceID))

	device := sunspec.NewDevice()
	device.Initialize(client)

	fmt.Printf("%+v", device.Descriptor())
	if results, err := device.Query(client); err != nil {
		panic(err)
	} else {
		for _, r := range results {
			fmt.Printf("%s: %v\n", r.Measurement, r.Value)
		}
	}
}
