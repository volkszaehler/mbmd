package server

import (
	"context"
	"log"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/volkszaehler/mbmd/meters"
)

const (
	// deviceIDregex is the regex pattern that identifies valid device ids
	deviceIDregex = "\\w*(\\d+)\\.(\\d+)"
)

// DeviceInfo returns device descriptor by device id
type DeviceInfo interface {
	DeviceDescriptorByID(id string) meters.DeviceDescriptor
}

// QueryEngine executes queries on connections and attached devices
type QueryEngine struct {
	handlers map[string]*Handler
	re       *regexp.Regexp
}

// sleepIsCancelled waits for timeout to expire. If context is cancelled before
// timeout expires, it will return early and indicate so by returning true.
func sleepIsCancelled(ctx context.Context, timeout time.Duration) bool {
	timer := time.After(timeout)
	select {
	case <-ctx.Done():
		return true
	case <-timer:
		return false
	}
}

// NewQueryEngine creates new query engine
func NewQueryEngine(managers map[string]meters.Manager) *QueryEngine {
	handlers := make(map[string]*Handler)

	for conn, m := range managers {
		handlers[conn] = NewHandler(len(handlers)+1, m)
	}

	qe := &QueryEngine{
		handlers: handlers,
		re:       regexp.MustCompile(deviceIDregex),
	}
	return qe
}

// DeviceDescriptorByID implements DeviceInfo interface
func (q *QueryEngine) DeviceDescriptorByID(id string) (res meters.DeviceDescriptor) {
	match := q.re.FindStringSubmatch(id)
	if len(match) != 3 {
		log.Fatalf("unexpected device id %s", id)
	}

	handlerID, _ := strconv.Atoi(match[1])
	deviceID, _ := strconv.Atoi(match[2])

	for _, h := range q.handlers {
		if h.ID == handlerID {
			h.Manager.All(false, func(id uint8, dev meters.Device) {
				if id == uint8(deviceID) {
					res = dev.Descriptor()
				}
			})
		}
	}

	return res
}

// Run executes the query engine to produce measurement results
func (q *QueryEngine) Run(
	ctx context.Context,
	control chan<- ControlSnip,
	results chan<- QuerySnip,
) {
	defer close(control)
	defer close(results)

	// run each connection manager inside separate goroutine
	var wg sync.WaitGroup
	for i, h := range q.handlers {
		wg.Add(1)
		go func(h *Handler, i string) {
			for {
				h.Run(ctx, control, results)
				if sleepIsCancelled(ctx, time.Second) {
					wg.Done()
					return
				}
			}
		}(h, i)
	}
	wg.Wait()
}
