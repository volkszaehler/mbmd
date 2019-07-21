package server

import (
	"log"
	"time"
)

func (q *QueryEngine) Scan() {
	type DeviceInfo struct {
		DeviceId  uint8
		MeterType string
	}

	var deviceId uint8
	deviceList := make([]DeviceInfo, 0)
	oldtimeout := q.setTimeout(50 * time.Millisecond)
	log.Printf("Starting bus scan")

SCAN:
	// loop over all valid slave adresses
	for deviceId = 1; deviceId <= 247; deviceId++ {
		// give the bus some time to recover before querying the next device
		time.Sleep(time.Duration(40) * time.Millisecond)

		for _, factory := range Producers {
			producer := factory()
			operation := producer.Probe()
			snip := NewQuerySnip(deviceId, operation)

			value, err := q.Query(snip)
			if err == nil {
				log.Printf("device %d: %s type device found, %s: %.2f\r\n",
					deviceId,
					producer.Type(),
					snip.IEC61850,
					snip.Transform(value))
				dev := DeviceInfo{
					DeviceId:  deviceId,
					MeterType: producer.Type(),
				}
				deviceList = append(deviceList, dev)
				continue SCAN
			}
		}

		log.Printf("device %d: n/a\r\n", deviceId)
	}

	// restore timeout to old value
	q.setTimeout(oldtimeout)

	log.Printf("Found %d active devices:\r\n", len(deviceList))
	for _, device := range deviceList {
		log.Printf("* slave address %d: type %s\r\n", device.DeviceId,
			device.MeterType)
	}
	log.Println("WARNING: This lists only the devices that responded to " +
		"a known probe request. Devices with different " +
		"function code definitions might not be detected.")
}
