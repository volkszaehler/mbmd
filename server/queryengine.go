package server

import (
	"fmt"

	"github.com/volkszaehler/mbmd/meters"
	"github.com/volkszaehler/mbmd/meters/connection"
)

const (
	maxRetry = 3
)

// QueryEngine executes queries on connections and attached devices
type QueryEngine struct {
	managers    map[string]connection.Manager
	status      map[string]MeterStatus
	connections []string
}

// NewQueryEngine creates new query engine
func NewQueryEngine(
	managers map[string]connection.Manager,
	status map[string]MeterStatus,
) *QueryEngine {
	return &QueryEngine{
		managers: managers,
		status:   status,
	}
}

func (q *QueryEngine) DeviceMap() map[string]meters.Device {
	res := make(map[string]meters.Device)

	for _, m := range q.managers {
		m.All(func(id uint8, dev meters.Device) {
			uniqueID := q.UniqueID(m.Conn, id, dev)
			res[uniqueID] = dev
		})
	}

	return res
}

// UniqueID creates a unique id per device
func (q *QueryEngine) UniqueID(conn connection.Connection, id uint8, dev meters.Device) string {
	uniqueID := fmt.Sprintf("%s%d", dev.Descriptor().Manufacturer, id)

	// add a unique connection id
	if len(q.managers) > 1 {
		if q.connections == nil {
			q.connections = make([]string, len(q.managers))
		}

		var connID int
		for i, c := range q.connections {
			if c == conn.String() {
				connID = i + 1
				break
			}
		}

		if connID == 0 {
			q.connections = append(q.connections, conn.String())
			connID = len(q.connections)
		}

		uniqueID = fmt.Sprintf("%s@%d", uniqueID, connID)
	}

	return uniqueID
}

// Run queries all connections and attached devices
func (q *QueryEngine) Run(
	controlChannel ControlSnipChannel,
	outputChannel QuerySnipChannel,
) {
	defer close(outputChannel)
	defer close(controlChannel)

	for _, m := range q.managers {
		fmt.Println(m.Conn)

		m.All(func(id uint8, dev meters.Device) {
			uniqueID := q.UniqueID(m.Conn, id, dev)
			fmt.Println(uniqueID)

			if _, ok := q.status[uniqueID]; !ok {
				q.status[uniqueID] = MeterStatus{}
			}

			if results, err := dev.Query(m.Conn.ModbusClient()); err != nil {
				// send results
				for _, r := range results {
					snip := QuerySnip{
						Device:            uniqueID,
						MeasurementResult: r,
					}
					outputChannel <- snip
				}

				// signal error
				controlChannel <- ControlSnip{
					Result:  failure,
					Device:  uniqueID,
					Message: fmt.Sprintf("device %d did not respond.", id),
				}
			} else {
				// 	for retry := 0; retry < maxRetry; retry++ {
				// 		bytes, err = q.Query(snip)
				// 		if err == nil {
				// 			break
				// 		}

				// 		q.status.IncreaseReconnectCounter()
				// 		log.Printf("Device %d failed to respond (%d/%d)",
				// 			id, retry+1, maxRetry)
				// 		time.Sleep(time.Duration(100) * time.Millisecond)
				// 	}

				// signal ok
				controlChannel <- ControlSnip{
					Result: ok,
					Device: uniqueID,
				}
			}
		})
	}
}

// func (q *QueryEngine) Scan() {
// 	type DeviceInfo struct {
// 		DeviceId  uint8
// 		MeterType string
// 	}

// 	var deviceId uint8
// 	deviceList := make([]DeviceInfo, 0)
// 	oldtimeout := q.setTimeout(50 * time.Millisecond)
// 	log.Printf("Starting bus scan")

// SCAN:
// 	// loop over all valid slave adresses
// 	for deviceId = 1; deviceId <= 247; deviceId++ {
// 		// give the bus some time to recover before querying the next device
// 		time.Sleep(time.Duration(40) * time.Millisecond)

// 		for _, factory := range Producers {
// 			producer := factory()
// 			operation := producer.Probe()
// 			snip := NewQuerySnip(deviceId, operation)

// 			value, err := q.Query(snip)
// 			if err == nil {
// 				log.Printf("Device %d: %s type device found, %s: %.2f\r\n",
// 					deviceId,
// 					producer.Type(),
// 					snip.IEC61850,
// 					snip.Transform(value))
// 				dev := DeviceInfo{
// 					DeviceId:  deviceId,
// 					MeterType: producer.Type(),
// 				}
// 				deviceList = append(deviceList, dev)
// 				continue SCAN
// 			}
// 		}

// 		log.Printf("Device %d: n/a\r\n", deviceId)
// 	}

// 	// restore timeout to old value
// 	q.setTimeout(oldtimeout)

// 	log.Printf("Found %d active devices:\r\n", len(deviceList))
// 	for _, device := range deviceList {
// 		log.Printf("* slave address %d: type %s\r\n", device.DeviceId,
// 			device.MeterType)
// 	}
// 	log.Println("WARNING: This lists only the devices that responded to " +
// 		"a known probe request. Devices with different " +
// 		"function code definitions might not be detected.")
// }
