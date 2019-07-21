package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/volkszaehler/mbmd/meters"
)

const (
	maxRetry   = 3
	retryDelay = 1 * time.Second
)

// Handler is responsible for querying a single connection
type Handler struct {
	manager meters.Manager
	status  map[uint8]*RuntimeInfo
}

// NewHandler creates a connection handler. The handler is responsible
// for querying all devices attached to the connection.
func NewHandler(m meters.Manager) *Handler {
	handler := &Handler{
		manager: m,
		status:  make(map[uint8]*RuntimeInfo),
	}

	m.All(false, func(id uint8, dev meters.Device) {
		handler.status[id] = &RuntimeInfo{
			Online: true,
		}
	})

	return handler
}

// UniqueID creates a unique id per device
func (h *Handler) UniqueID(id uint8, dev meters.Device) string {
	uniqueID := fmt.Sprintf("%s:%d@%s", dev.Descriptor().Manufacturer, id, h.manager.Conn)
	return uniqueID
}

// Run initializes and queries every device attached to the handler's connection
func (h *Handler) Run(
	ctx context.Context,
	control chan<- ControlSnip,
	results chan<- QuerySnip,
) {
	h.manager.All(true, func(id uint8, dev meters.Device) {
		if sleepIsCancelled(ctx, 0) {
			return
		}

		status := h.status[id]
		if !status.initialized {
			if err := dev.Initialize(h.manager.Conn.ModbusClient()); err != nil {
				log.Printf("initializing device %d at %s failed: %v", id, h.manager.Conn, err)
				sleepIsCancelled(ctx, retryDelay)
				return
			}

			d := dev.Descriptor()
			log.Printf("initialized device %d at %s %v", id, h.manager.Conn, d)

			status.initialized = true
		}

		if queryable, wakeup := status.IsQueryable(); wakeup {
			log.Printf("device %d at %s is offline - reactivating", id, h.manager.Conn)
		} else if !queryable {
			return
		}

		h.queryDevice(ctx, control, results, id, dev)
	})
}

func (h *Handler) queryDevice(
	ctx context.Context,
	control chan<- ControlSnip,
	results chan<- QuerySnip,
	id uint8,
	dev meters.Device,
) {
	uniqueID := h.UniqueID(id, dev)

	status := h.status[id]

	for retry := 0; retry < maxRetry; retry++ {
		status.Requests++
		measurements, err := dev.Query(h.manager.Conn.ModbusClient())

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
		log.Printf("device %d at %s did not respond (%d/%d)", id, h.manager.Conn, retry+1, maxRetry)
		if sleepIsCancelled(ctx, retryDelay) {
			return
		}
	}

	log.Printf("device %d at %s is offline", id, h.manager.Conn)

	// close connection to force modbus client to reopen
	h.manager.Conn.Close()
	status.Available(false)
	control <- ControlSnip{
		Device: uniqueID,
		Status: *status,
	}
}
