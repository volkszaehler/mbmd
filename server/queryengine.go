package server

import (
	"context"
	"sync"
	"time"

	"github.com/volkszaehler/mbmd/meters/connection"
)

// QueryEngine executes queries on connections and attached devices
type QueryEngine struct {
	sync.Mutex
	handlers map[string]*Handler
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
func NewQueryEngine(managers map[string]connection.Manager) *QueryEngine {
	handlers := make(map[string]*Handler)

	for conn, m := range managers {
		handlers[conn] = NewHandler(m)
	}

	qe := &QueryEngine{
		handlers: handlers,
	}
	return qe
}

// Run queries all connections and attached devices
func (q *QueryEngine) Run(
	ctx context.Context,
	control ControlSnipChannel,
	results QuerySnipChannel,
) {
	defer close(results)
	defer close(control)

	// run each connection manager inside separate goroutine
	var wg sync.WaitGroup
	for i, h := range q.handlers {
		wg.Add(1)
		go func(h *Handler, i string) {
			for {
				if sleepIsCancelled(ctx, 1) {
					wg.Done()
					return
				}
				h.Run(ctx, control, results)
			}
		}(h, i)
	}
	wg.Wait()
}
