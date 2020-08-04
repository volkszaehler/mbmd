package server

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/volkszaehler/mbmd/meters"
)

// DeviceInfo returns device descriptor by device id
type DeviceInfo interface {
	DeviceDescriptorByID(id string) meters.DeviceDescriptor
}

// QueryEngine executes queries on connections and attached devices
type QueryEngine struct {
	handlers    map[string]*Handler
	deviceCache map[string]meters.Device
}

// NewQueryEngine creates new query engine
func NewQueryEngine(managers map[string]*meters.Manager) *QueryEngine {
	handlers := make(map[string]*Handler)

	// sort handlers by name
	keys := make([]string, 0, len(managers))
	for conn := range managers {
		keys = append(keys, conn)
	}
	sort.Strings(keys)

	for _, conn := range keys {
		m := managers[conn]
		if m.Count() == 0 {
			// don't give ids to empty connections
			continue
		}

		handlers[conn] = NewHandler(len(handlers)+1, m)
	}

	qe := &QueryEngine{
		handlers:    handlers,
		deviceCache: make(map[string]meters.Device),
	}
	return qe
}

// DeviceDescriptorByID implements DeviceInfo interface
func (q *QueryEngine) DeviceDescriptorByID(id string) (res meters.DeviceDescriptor) {
	// already cached?
	if dev, ok := q.deviceCache[id]; ok {
		return dev.Descriptor()
	}

	for _, h := range q.handlers {
		h.Manager.Find(func(slaveID uint8, dev meters.Device) (found bool) {
			devID := h.deviceID(slaveID, dev)
			if id == devID {
				q.deviceCache[id] = dev
				res = dev.Descriptor()
				found = true
			}

			return
		})
	}

	return res
}

// Run executes the query engine to produce measurement results
func (q *QueryEngine) Run(
	ctx context.Context,
	rate time.Duration,
	control chan<- ControlSnip,
	results chan<- QuerySnip,
) {
	defer close(control)
	defer close(results)

	// run each connection manager inside separate goroutine
	var wg sync.WaitGroup
	for _, h := range q.handlers {
		wg.Add(1)

		go func(h *Handler) {
			ticker := time.NewTicker(rate)
			defer ticker.Stop()

			for {
				// run handlers
				h.Run(ctx, control, results)

				// wait for rate limit
				select {
				case <-ctx.Done():
					// abort if context is cancelled
					wg.Done()
					return
				case <-ticker.C:
				}
			}
		}(h)
	}

	wg.Wait()
}
