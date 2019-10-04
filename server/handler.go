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
	initDelay  = 3 * time.Second
)

// Handler is responsible for querying a single connection
type Handler struct {
	ID      int
	Manager meters.Manager
	status  map[uint8]*RuntimeInfo
}

// NewHandler creates a connection handler. The handler is responsible
// for querying all devices attached to the connection.
func NewHandler(id int, m meters.Manager) *Handler {
	handler := &Handler{
		ID:      id,
		Manager: m,
		status:  make(map[uint8]*RuntimeInfo),
	}

	return handler
}

// uniqueID creates a unique id per device
func (h *Handler) uniqueID(id uint8, dev meters.Device) string {
	return fmt.Sprintf("%s%d.%d", dev.Descriptor().Manufacturer, h.ID, id)
}

// Run initializes and queries every device attached to the handler's connection
func (h *Handler) Run(
	ctx context.Context,
	control chan<- ControlSnip,
	results chan<- QuerySnip,
) {
	h.Manager.All(func(id uint8, dev meters.Device) {
		if sleepIsCancelled(ctx, 0) {
			return
		}

		// select device
		h.Manager.Conn.Slave(id)

		// initialize device
		status, ok := h.status[id]
		if !ok {
			var err error
			if status, err = h.initializeDevice(ctx, control, id, dev); err != nil {
				return
			}
			h.status[id] = status
		}

		uniqueID := h.uniqueID(id, dev)
		if queryable, wakeup := status.IsQueryable(); wakeup {
			log.Printf("device %s is offline - reactivating", uniqueID)
		} else if !queryable {
			return
		}

		// query device
		h.queryDevice(ctx, control, results, id, dev)
	})
}

func (h *Handler) initializeDevice(
	ctx context.Context,
	control chan<- ControlSnip,
	id uint8,
	dev meters.Device,
) (*RuntimeInfo, error) {
	uniqueID := h.uniqueID(id, dev)

	if err := dev.Initialize(h.Manager.Conn.ModbusClient()); err != nil {
		if _, partial := err.(meters.SunSpecPartiallyInitialized); !partial {
			log.Printf("initializing device %s failed: %v", uniqueID, err)
			sleepIsCancelled(ctx, initDelay)
			return nil, err
		}
		log.Println(err) // log error but continue
	}

	d := dev.Descriptor()
	uniqueID = h.uniqueID(id, dev) // update id
	log.Printf("initialized device %s: %v", uniqueID, d)

	// create status
	status := &RuntimeInfo{Online: true}

	// signal device online
	control <- ControlSnip{
		Device: uniqueID,
		Status: *status,
	}

	return status, nil
}

func (h *Handler) queryDevice(
	ctx context.Context,
	control chan<- ControlSnip,
	results chan<- QuerySnip,
	id uint8,
	dev meters.Device,
) {
	uniqueID := h.uniqueID(id, dev)
	status := h.status[id]

	for retry := 0; retry < maxRetry; retry++ {
		status.Requests++
		measurements, err := dev.Query(h.Manager.Conn.ModbusClient())

		if err == nil {
			// send ok status
			status.Available(true)
			control <- ControlSnip{
				Device: uniqueID,
				Status: *status,
			}

			// send measurements
			for _, r := range measurements {
				snip := QuerySnip{
					Device:            uniqueID,
					MeasurementResult: r,
				}
				results <- snip
			}

			return
		}

		status.Errors++
		log.Printf("device %s did not respond (%d/%d)", uniqueID, retry+1, maxRetry)
		if sleepIsCancelled(ctx, retryDelay) {
			return
		}
	}

	log.Printf("device %s is offline", uniqueID)

	// close connection to force modbus client to reopen
	h.Manager.Conn.Close()

	// send error status
	status.Available(false)
	control <- ControlSnip{
		Device: uniqueID,
		Status: *status,
	}
}
