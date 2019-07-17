package server

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/volkszaehler/mbmd/meters"
	"github.com/volkszaehler/mbmd/meters/connection"
)

const (
	maxRetry     = 3
	errorTimeout = 100 * time.Millisecond
)

// QueryEngine executes queries on connections and attached devices
type QueryEngine struct {
	sync.Mutex
	managers    map[string]connection.Manager
	status      map[meters.Device]*RuntimeInfo
	connections []string
	control     ControlSnipChannel
	results     QuerySnipChannel
	ctx         context.Context
}

// NewQueryEngine creates new query engine
func NewQueryEngine(managers map[string]connection.Manager) *QueryEngine {
	runtimeStatus := make(map[meters.Device]*RuntimeInfo)
	return &QueryEngine{
		managers: managers,
		status:   runtimeStatus,
	}
}

func (q *QueryEngine) Connections(cb func(connection.Manager)) {
	for _, m := range q.managers {
		cb(m)
	}
}

// UniqueID creates a unique id per device
func (q *QueryEngine) UniqueID(conn connection.Connection, id uint8, dev meters.Device) string {
	q.Lock()
	defer q.Unlock()

	uniqueID := fmt.Sprintf("%s:%d", dev.Descriptor().Manufacturer, id)

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
	ctx context.Context,
	control ControlSnipChannel,
	results QuerySnipChannel,
) {
	defer close(results)
	defer close(control)

	q.control = control
	q.results = results
	q.ctx = ctx

	// create runtime status
	for _, m := range q.managers {
		m.All(func(id uint8, dev meters.Device) {
			ri := &RuntimeInfo{
				Online: true,
			}
			q.status[dev] = ri
		})
	}

	// run each connection manager inside separate gorouting
	var wg sync.WaitGroup
	for i, m := range q.managers {
		wg.Add(1)
		go func(m connection.Manager, i string) {
			for {
				if q.sleepOrCancelled(1) {
					wg.Done()
					return
				}
				q.runManager(m)
			}
		}(m, i)
	}
	wg.Wait()
}

// sleepOrCancelled waits for timeout to expire. If context is cancelled before
// timeout expires, it will return early and indicate so by returning true.
func (q *QueryEngine) sleepOrCancelled(timeout time.Duration) bool {
	timer := time.After(timeout)
	select {
	case <-q.ctx.Done():
		return true
	case <-timer:
		return false
	}
}

func (q *QueryEngine) runManager(m connection.Manager) {
	m.All(func(id uint8, dev meters.Device) {
		if q.sleepOrCancelled(0) {
			return
		}

		status := q.status[dev]

		if !status.initialized {
			if err := dev.Initialize(m.Conn.ModbusClient()); err != nil {
				log.Printf("initializing device %d at %s failed: %v", id, m.Conn, err)
				q.sleepOrCancelled(errorTimeout)
				return
			}

			status.initialized = true
		}

		if queryable, wakeup := status.IsQueryable(); wakeup {
			log.Printf("device %d at %s is offline - reactivating", id, m.Conn)
		} else if !queryable {
			return
		}

		q.queryDevice(m, id, dev, status)
	})
}

func (q *QueryEngine) queryDevice(
	m connection.Manager,
	id uint8, dev meters.Device,
	status *RuntimeInfo,
) {
	uniqueID := q.UniqueID(m.Conn, id, dev)
	fmt.Println(uniqueID)

	for retry := 0; retry < maxRetry; retry++ {
		status.IncRequests()
		measurements, err := dev.Query(m.Conn.ModbusClient())

		if err == nil {
			// send measurements
			for _, r := range measurements {
				snip := QuerySnip{
					Device:            uniqueID,
					MeasurementResult: r,
				}
				q.results <- snip
			}

			// send ok
			status.SetOnline(true)
			q.control <- ControlSnip{
				Result: ok,
				Device: uniqueID,
			}

			return
		}

		status.IncErrors()
		log.Printf("device %d at %s did not respond (%d/%d)", id, m.Conn, retry+1, maxRetry)
		if q.sleepOrCancelled(errorTimeout) {
			return
		}
	}

	log.Printf("device %d at %s is offline", id, m.Conn)

	// close connection to force modbus client to reopen
	m.Conn.Close()

	// signal error
	status.SetOnline(false)
	q.control <- ControlSnip{
		Result:  failure,
		Device:  uniqueID,
		Message: fmt.Sprintf("device %d at %s did not respond", id, m.Conn),
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
