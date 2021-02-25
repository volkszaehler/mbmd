package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	prometheusManager "github.com/volkszaehler/mbmd/prometheus_metrics"

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
	Manager *meters.Manager
	status  map[string]*RuntimeInfo
}

// NewHandler creates a connection handler. The handler is responsible
// for querying all devices attached to the connection.
func NewHandler(id int, m *meters.Manager) *Handler {
	handler := &Handler{
		ID:      id,
		Manager: m,
		status:  make(map[string]*RuntimeInfo),
	}

	// TODO prometheus: ConnectionHandlersCreated

	return handler
}

// deviceID creates a unique id per device
func (h *Handler) deviceID(id uint8, dev meters.Device) string {
	desc := dev.Descriptor()
	devID := fmt.Sprintf("%s%d.%d", desc.Type, h.ID, id)
	if desc.SubDevice > 0 {
		devID = fmt.Sprintf("%s.%d", devID, desc.SubDevice)
	}
	return devID
}

// Run initializes and queries every device attached to the handler's connection
func (h *Handler) Run(
	ctx context.Context,
	control chan<- ControlSnip,
	results chan<- QuerySnip,
) {
	h.Manager.All(func(id uint8, dev meters.Device) {
		// abort if context is cancelled
		select {
		case <-ctx.Done():
			return
		default:
		}

		// select device
		h.Manager.Conn.Slave(id)

		// initialize device
		deviceID := h.deviceID(id, dev)
		status, ok := h.status[deviceID]
		if !ok {
			var err error
			if status, err = h.initializeDevice(ctx, control, id, dev); err != nil {
				return
			}
			h.status[deviceID] = status
		}

		if queryable, wakeup := status.IsQueryable(); wakeup {
			log.Printf("device %s is offline - reactivating", deviceID)
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
	deviceID := h.deviceID(id, dev)
	// TODO prometheus: ConnectionHandlerDevicesInitializedRoutineStarted

	if err := dev.Initialize(h.Manager.Conn.ModbusClient()); err != nil {
		if !errors.Is(err, meters.ErrPartiallyOpened) {
			log.Printf("initializing device %s failed: %v", deviceID, err)

			// wait for error to settle
			ctx, cancel := context.WithTimeout(ctx, initDelay)
			defer cancel()
			<-ctx.Done()

			// TODO prometheus: ConnectionHandlerDevicesInitializationFailure
			return nil, err
		}
		log.Println(err) // log error but continue
	}

	log.Printf("initialized device %s: %v", deviceID, dev.Descriptor())
	// TODO prometheus: ConnectionHandlerDevicesInitializationSuccess

	// create status
	status := &RuntimeInfo{Online: true}

	// signal device online
	control <- ControlSnip{
		Device: deviceID,
		Status: *status,
	}

	// create Prometheus metrics
	prometheusManager.CreateMeasurementMetrics(dev)

	return status, nil
}

func (h *Handler) queryDevice(
	ctx context.Context,
	control chan<- ControlSnip,
	results chan<- QuerySnip,
	id uint8,
	dev meters.Device,
) {
	deviceID := h.deviceID(id, dev)
	deviceSerial := dev.Descriptor().Serial
	status := h.status[deviceID]

	for retry := 0; retry < maxRetry; retry++ {
		status.Requests++
		prometheusManager.DeviceQueriesTotal.WithLabelValues(deviceID, deviceSerial).Inc()

		measurements, err := dev.Query(h.Manager.Conn.ModbusClient())

		if err == nil {
			// send ok status
			status.Available(true)
			control <- ControlSnip{
				Device: deviceID,
				Status: *status,
			}
			prometheusManager.DeviceQueriesSuccessTotal.WithLabelValues(deviceID, deviceSerial).Inc()

			// send measurements
			for _, r := range measurements {
				if math.IsNaN(r.Value) {
					log.Printf("device %s skipping NaN for %s", deviceID, r.Measurement.String())
					prometheusManager.DeviceQueryMeasurementValueSkippedTotal.WithLabelValues(deviceID, deviceSerial).Inc()
					continue
				}

				snip := QuerySnip{
					Device:            deviceID,
					MeasurementResult: r,
				}
				results <- snip

				prometheusManager.UpdateMeasurementMetric(deviceID, deviceSerial, r)
			}

			return
		}

		status.Errors++
		prometheusManager.DeviceQueriesErrorTotal.WithLabelValues(deviceID, deviceSerial).Inc()
		log.Printf("device %s did not respond (%d/%d): %v", deviceID, retry+1, maxRetry, err)

		// wait for device to settle after error
		select {
		case <-ctx.Done():
			return
		case <-time.After(retryDelay):
		}
	}

	log.Printf("device %s is offline", deviceID)

	// close connection to force modbus client to reopen
	h.Manager.Conn.Close()
	// TODO prometheus: ModbusConnectionsClosed

	// send error status
	status.Available(false)
	control <- ControlSnip{
		Device: deviceID,
		Status: *status,
	}
}
