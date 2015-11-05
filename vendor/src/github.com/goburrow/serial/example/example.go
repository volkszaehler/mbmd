package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/goburrow/serial"
)

var address = "/dev/ttyUSB0"

func main() {
	if len(os.Args) > 1 {
		address = os.Args[1]
	}
	port, err := serial.Open(&serial.Config{
		Address: address,
		Timeout: 30 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	if _, err = port.Write([]byte("serial")); err != nil {
		log.Fatal(err)
	}
	if _, err = io.Copy(os.Stdout, port); err != nil {
		log.Fatal(err)
	}
}
