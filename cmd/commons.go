package cmd

import (
	"fmt"
	golog "log"
	"os"
	"sort"

	"github.com/volkszaehler/mbmd/meters"
	"github.com/volkszaehler/mbmd/meters/rs485"
)

// logger is the golang compatible logger interface used by all commands
type logger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatal(v ...interface{})
}

// quietlogger logs only fatal messages
type quietlogger struct {
	*golog.Logger
}

func (l *quietlogger) Printf(format string, v ...interface{}) {}
func (l *quietlogger) Println(v ...interface{})               {}

// log variable replaces golang log package functions
var log logger = golog.New(os.Stderr, "", golog.LstdFlags)

// configureLogger sets default logger.
// According to verbosity flag only fatal messages are shown
func configureLogger(verbose bool, flag int) {
	if verbose {
		log = golog.New(os.Stderr, "", flag)
	} else {
		log = &quietlogger{golog.New(os.Stderr, "", flag)}
	}
}

// meterHelp output list of supported devices
func meterHelp() string {
	s := fmt.Sprintf("\n  %s", "RTU")
	types := make([]string, 0)
	for t := range rs485.Producers {
		types = append(types, t)
	}

	sort.Strings(types)

	for _, t := range types {
		p := rs485.Producers[t]()
		s += fmt.Sprintf("\n    %-9s%s", t, p.Description())
	}

	s += fmt.Sprintf("\n  %s", "TCP")
	s += fmt.Sprintf("\n    %-9s%s", "SUNS", "Sunspec-compatible MODBUS TCP device (SMA, SolarEdge, KOSTAL, etc)")

	return s
}

// countDevices counts all devices for all meters
func countDevices(managers map[string]*meters.Manager) int {
	var count int
	for _, m := range managers {
		m.All(func(id uint8, dev meters.Device) {
			conf := dev.Descriptor()
			log.Printf("config: declared device %s:%d.%d", conf.Type, id, conf.SubDevice)
			count++
		})
	}
	return count
}

// setLogger enabled raw logging for all devices
func setLogger(managers map[string]*meters.Manager, logger meters.Logger) {
	for _, m := range managers {
		m.Conn.Logger(logger)
	}
}
