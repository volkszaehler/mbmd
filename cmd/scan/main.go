package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	sunspec "github.com/crabmusket/gosunspec"
	bus "github.com/crabmusket/gosunspec/modbus"
	_ "github.com/crabmusket/gosunspec/models" // import models
	"github.com/crabmusket/gosunspec/smdx"

	"github.com/grid-x/modbus"
)

const (
	base = 40000
)

func scaledValue(p sunspec.Point) float64 {
	f := p.ScaledValue()

	switch p.Type() {
	case "int16":
		if p.Value() == int16(math.MinInt16) {
			f = math.NaN()
		}
	case "int32":
		if p.Value() == int32(math.MinInt32) {
			f = math.NaN()
		}
	case "int64":
		if p.Value() == int64(math.MinInt64) {
			f = math.NaN()
		}
	case "uint16":
		if p.Value() == uint16(math.MaxUint16) {
			f = math.NaN()
		}
	case "uint32":
		if p.Value() == uint32(math.MaxUint32) {
			f = math.NaN()
		}
	case "uint64":
		if p.Value() == uint64(math.MaxUint64) {
			f = math.NaN()
		}
	}

	return f
}

func pf(format string, v ...interface{}) {
	format = strings.TrimSuffix(format, "\n") + "\n"
	fmt.Printf(format, v...)
}

// injectable logger for grid-x modbus implementation
type modbusLogger struct{}

func (l *modbusLogger) Printf(format string, v ...interface{}) {
	pf(format, v...)
}

func doModels(d sunspec.Device, cb func(m sunspec.Model)) {
	modelIds := []sunspec.ModelId{1, 101, 103}
	models := d.Collect(sunspec.OneOfSeveralModelIds(modelIds))

	for _, model := range models {
		cb(model)
	}
}

func scanSunspec(client modbus.Client) {
	in, err := bus.Open(client)
	if err != nil {
		log.Fatal(err)
	}

	in.Do(func(d sunspec.Device) {
		d.Do(func(m sunspec.Model) {
			// doModels(d, func(m sunspec.Model) {
			pf("--------- Model %d %s ---------", m.Id(), modelName(m))

			printModel(smdx.GetModel(uint16(m.Id())))
			pf("-- Data --")

			m.Do(func(b sunspec.Block) {
				err = b.Read()
				if err != nil {
					log.Fatal(err)
				}

				b.Do(func(p sunspec.Point) {
					t := p.Type()[0:3]
					v := ""
					if t == "int" || t == "uin" || t == "acc" {
						v = fmt.Sprintf("%.2f", scaledValue(p))
					}
					pf("%10s %-18s %10v %10s", p.Type(), p.Id(), p.Value(), v)
				})
			})
		})
	})
}

func modelName(m sunspec.Model) string {
	model := smdx.GetModel(uint16(m.Id()))
	if model == nil {
		return ""
	}
	return model.Name
}

func printModel(m *smdx.ModelElement) {
	pf("-- Definition --")
	// pf("----")
	// pf("Model:  %d - %s", m.Id, m.Name)
	pf("Length: %d (0x%02x words, 0x%02x bytes)", m.Length, m.Length, 2*m.Length)
	pf("Blocks: %d", len(m.Blocks))

	for i, b := range m.Blocks {
		pf("-- block #%d - %s", i, b.Name)
		pf("Length: %d", b.Length)

		for _, p := range b.Points {
			u := ""
			if p.Units != "" {
				u = p.Units
			}
			pf("%4d %4d %12s %-8s %s", p.Offset, p.Length, p.Id, u, p.Type)
		}
	}
}

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

	mbl := &modbusLogger{}
	handler.Logger = mbl

	handler.SetSlave(byte(deviceID))

	// b, err := client.ReadHoldingRegisters(40072, 7)
	// if err != nil {
	// 	panic(err)
	// }
	// pf("% x", b)
	// os.Exit(0)

	scanSunspec(client)
	// scanCustom(client)
}
