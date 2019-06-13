package main

import (
	"encoding/binary"
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

	custom "github.com/volkszaehler/mbmd/meters"
	customImpl "github.com/volkszaehler/mbmd/meters/impl"

	"github.com/grid-x/modbus"
)

const (
	base = 40000
)

func scaledValue(p sunspec.Point) float64 {
	f := p.ScaledValue()

	switch p.Type() {
	case "acc16":
		fallthrough
	case "uint16":
		if p.Value() == uint16(math.MaxUint16) {
			f = math.NaN()
		}
	case "acc32":
		fallthrough
	case "uint32":
		if p.Value() == uint32(math.MaxUint32) {
			f = math.NaN()
		}
	case "acc64":
		fallthrough
	case "uint64":
		maxUint64 := uint64(math.MaxUint64)
		if p.Value() == maxUint64 {
			f = math.NaN()
		}
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

func scanCustom(client modbus.Client) {
	loop := uint16(base)
	loop += 2

	for {
		b, err := client.ReadHoldingRegisters(loop, 2)
		if err != nil {
			log.Fatal(err)
		}
		pf("loop: %d bytes: % x", loop, b)

		id := binary.BigEndian.Uint16(b)
		length := binary.BigEndian.Uint16(b[2:])
		pf("id/len: %d %d", id, length)

		if model, ok := custom.SunspecModels[int(id)]; ok {
			pf("model: %s", model)
		}

		if id == 0xffff {
			goto DONE
		}

		model := smdx.GetModel(id)
		if model != nil {
			pf("fixed length: %d blocks: %d", model.Length, len(model.Blocks))
			pf("%v", model)
		}

		b, err = client.ReadHoldingRegisters(loop+2, length)
		if err != nil {
			log.Fatal(err)
		}
		pf("data: % x", b)

		if id == 1 {
			core := customImpl.SunSpecCore{}
			suns := []byte{0x53, 0x75, 0x6e, 0x53, 0x00, 0x00, 0x00, 0x00}

			cb := append(suns, b...)
			d, err := core.DecodeSunSpecCommonBlock(cb)
			if err != nil {
				log.Fatal(err)
			}
			pf("%+v", d)
		}
		loop += length + 2
	}
DONE:
}

func scanSunspec(client modbus.Client) {
	in, err := bus.Open(client)
	if err != nil {
		log.Fatal(err)
	}

	in.Do(func(d sunspec.Device) {
		d.Do(func(m sunspec.Model) {
			pf("--------- Model %d %s ---------", m.Id(), modelName(m))

			if m.Id() == 11 {
				return
			}

			m.Do(func(b sunspec.Block) {
				err = b.Read()
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println("read points")
				b.Do(func(p sunspec.Point) {
					t := p.Type()[0:3]
					v := ""
					if t == "int" || t == "uin" || t == "acc" {
						// v = fmt.Sprintf("%.2f", p.ScaledValue())
						v = fmt.Sprintf("%.2f", scaledValue(p))
					}
					pf("%10s %-18s %8v %10s", p.Type(), p.Id(), p.Value(), v)
				})
			})

			printModel(smdx.GetModel(uint16(m.Id())))
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
	pf("----")
	pf("Model:  %d - %s", m.Id, m.Name)
	pf("Length: %d (0x%02x words, 0x%02x bytes)", m.Length, m.Length, 2*m.Length)
	pf("Blocks: %d", len(m.Blocks))

	for i, b := range m.Blocks {
		pf("--")
		pf("#%d - %s", i, b.Name)
		pf("Length: %d", b.Length)

		for _, p := range b.Points {
			u := ""
			if p.Units != "" {
				u = p.Units
			}
			pf("%4d %4d %8s %-4s %s", p.Offset, p.Length, p.Id, u, p.Type)
		}
	}
}

func main() {
	// model := smdx.GetModel(1)
	// printModel(model)
	// model = smdx.GetModel(101)
	// printModel(model)
	// model = smdx.GetModel(11)
	// printModel(model)
	// os.Exit(0)

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

	scanSunspec(client)
	// scanCustom(client)
}
