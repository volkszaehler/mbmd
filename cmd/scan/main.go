package main

import (
	"errors"
	"log"
	"net/url"
	"os"
	"strconv"

	sunspec "github.com/crabmusket/gosunspec"
	_ "github.com/crabmusket/gosunspec/impl" // import models
	modbuswrapper "github.com/crabmusket/gosunspec/modbus"
	_ "github.com/crabmusket/gosunspec/models" // import models

	"github.com/grid-x/modbus"
	_ "github.com/volkszaehler/mbmd/meters/impl"
)

const (
	base = 40000
)

func uriDefaultSchemeAndPort(addr string, scheme string, port string) (string, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return addr, err
	}

	// set default scheme
	if u.Host == "" || u.Scheme == "" {
		if scheme == "" {
			return addr, errors.New("missing scheme")
		}
		addr = scheme + "://" + addr

		u, err = url.Parse(addr)
		if err != nil {
			return addr, err
		}
	}

	// set default port
	if port != "" && u.Port() == "" {
		u.Host = u.Host + ":" + port
		addr = u.String()
	}

	return addr, nil
}

func test(addr string) {
	s, err := uriDefaultSchemeAndPort(addr, "tcp", "")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s -> %s", addr, s)
}

func main() {
	addr := os.Args[1]
	_, err := url.Parse(addr)
	if err != nil {
		log.Fatal(err)
	}

	deviceID, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	handler := modbus.NewTCPClientHandler(addr)
	client := modbus.NewClient(handler)
	if err := handler.Connect(); err != nil {
		log.Fatal("Failed to connect: ", err)
	}

	handler.SetSlave(byte(deviceID))

	in, err := modbuswrapper.Open(client)
	if err != nil {
		log.Fatal(err)
	}

	in.Do(func(d sunspec.Device) {
		d.Do(func(m sunspec.Model) {
			log.Printf("Model: %d ----", m.Id())

			m.Do(func(b sunspec.Block) {
				err = b.Read()
				if err != nil {
					log.Fatal(err)
				}

				b.Do(func(p sunspec.Point) {
					log.Printf("%s %s %v", p.Type(), p.Id(), p.Value())
				})
			})
		})
	})

	/*
			loop := uint16(base)
			loop += 2

			for {
				b, err := client.ReadHoldingRegisters(loop, 2)
				if err != nil {
					log.Fatal(err)
				}
				log.Printf("loop: %d bytes: % x", loop, b)

				id := binary.BigEndian.Uint16(b)
				length := binary.BigEndian.Uint16(b[2:])
				log.Printf("id/len: %d %d", id, length)

				if model, ok := meters.SunspecModels[int(id)]; ok {
					log.Printf("model: %s", model)
				}

				if id == 0xffff {
					goto DONE
				}

				model := smdx.GetModel(id)
				if model != nil {
					log.Printf("fixed length: %d blocks: %d", model.Length, len(model.Blocks))
					log.Printf("%v", model)
				}

				b, err = client.ReadHoldingRegisters(loop+2, length)
				if err != nil {
					log.Fatal(err)
				}
				log.Printf("data: % x", b)

				if id == 1 {
					core := sunspec.SunSpecCore{}
					suns := []byte{0x53, 0x75, 0x6e, 0x53, 0x00, 0x00, 0x00, 0x00}

					cb := append(suns, b...)
					d, err := core.DecodeSunSpecCommonBlock(cb)
					if err != nil {
						log.Fatal(err)
					}
					log.Printf("%+v", d)
				}
				loop += length + 2
			}
		DONE:
	*/
}
