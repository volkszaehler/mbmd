package cmd

import (
	"fmt"
	"log"
	"sort"

	"github.com/volkszaehler/mbmd/meters"
	"github.com/volkszaehler/mbmd/meters/rs485"
)

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
func countDevices(managers map[string]meters.Manager) int {
	var count int
	for _, m := range managers {
		m.All(func(id uint8, dev meters.Device) {
			log.Printf("config: declared device %s:%d", dev.Descriptor().Manufacturer, id)
			count++
		})
	}
	return count
}

// setLogger enabled raw logging for all devices
func setLogger(managers map[string]meters.Manager, logger meters.Logger) {
	for _, m := range managers {
		m.Conn.Logger(logger)
	}
}
