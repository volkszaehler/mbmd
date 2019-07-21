package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/volkszaehler/mbmd/meters"
	"github.com/volkszaehler/mbmd/meters/connection"
)

const (
	maxRetry   = 3
	retryDelay = 1 * time.Second
)

type Handler struct {
	manager connection.Manager
	status  map[uint8]*RuntimeInfo
}

func NewHandler(m connection.Manager) *Handler {
	handler := &Handler{
		manager: m,
		status:  make(map[uint8]*RuntimeInfo),
	}

	m.Devices(func(id uint8, dev meters.Device) {
		handler.status[id] = &RuntimeInfo{
			Online: true,
		}
	})

	return handler
}

// UniqueID creates a unique id per device
func (q *Handler) UniqueID(id uint8, dev meters.Device) string {
	uniqueID := fmt.Sprintf("%s:%d", dev.Descriptor().Manufacturer, id)

	// add a unique connection id
	// if len(q.handlers) > 1 {
	// if q.connections == nil {
	// 	q.connections = make([]string, len(q.handlers))
	// }

	// var connID int
	// for i, c := range q.connections {
	// 	if c == q.manager.String() {
	// 		connID = i + 1
	// 		break
	// 	}
	// }

	// if connID == 0 {
	// 	q.connections = append(q.connections, conn.String())
	// 	connID = len(q.connections)
	// }

	// uniqueID = fmt.Sprintf("%s@%d", uniqueID, connID)
	// }
	uniqueID = fmt.Sprintf("%s@%s", uniqueID, q.manager.Conn)

	return uniqueID
}

func (q *Handler) Run(
	ctx context.Context,
	control ControlSnipChannel,
	results QuerySnipChannel,
) {
	q.manager.Invoke(func(id uint8, dev meters.Device) {
		if sleepIsCancelled(ctx, 0) {
			return
		}

		status := q.status[id]
		if !status.initialized {
			if err := dev.Initialize(q.manager.Conn.ModbusClient()); err != nil {
				log.Printf("initializing device %d at %s failed: %v", id, q.manager.Conn, err)
				sleepIsCancelled(ctx, retryDelay)
				return
			}

			status.initialized = true
		}

		if queryable, wakeup := status.IsQueryable(); wakeup {
			log.Printf("device %d at %s is offline - reactivating", id, q.manager.Conn)
		} else if !queryable {
			return
		}

		q.queryDevice(ctx, control, results, id, dev)
	})
}

func (q *Handler) queryDevice(
	ctx context.Context,
	control ControlSnipChannel,
	results QuerySnipChannel,
	id uint8,
	dev meters.Device,
) {
	uniqueID := q.UniqueID(id, dev)

	status := q.status[id]

	for retry := 0; retry < maxRetry; retry++ {
		status.Requests++
		measurements, err := dev.Query(q.manager.Conn.ModbusClient())

		if err == nil {
			// send measurements
			for _, r := range measurements {
				snip := QuerySnip{
					Device:            uniqueID,
					MeasurementResult: r,
				}
				results <- snip
			}

			// send ok
			status.Available(true)

			control <- ControlSnip{
				Device: uniqueID,
				Status: *status,
			}

			return
		}

		status.Errors++
		log.Printf("device %d at %s did not respond (%d/%d)", id, q.manager.Conn, retry+1, maxRetry)
		if sleepIsCancelled(ctx, retryDelay) {
			return
		}
	}

	log.Printf("device %d at %s is offline", id, q.manager.Conn)

	// close connection to force modbus client to reopen
	q.manager.Conn.Close()
	status.Available(false)
	control <- ControlSnip{
		Device: uniqueID,
		Status: *status,
	}
}
