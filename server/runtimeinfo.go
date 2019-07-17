package server

import (
	"sync"
	"sync/atomic"
	"time"
)

const (
	retryTimeout = 1 * time.Second
)

type RuntimeInfo struct {
	sync.Mutex
	lastFailure time.Time
	initialized bool
	Online      bool
	Requests    uint64
	Errors      uint64
}

func (r *RuntimeInfo) IncRequests() {
	r.Lock()
	atomic.AddUint64(&r.Requests, 1)
	r.Unlock()
}

func (r *RuntimeInfo) IncErrors() {
	r.Lock()
	defer r.Unlock()
	atomic.AddUint64(&r.Errors, 1)
}

func (r *RuntimeInfo) Status() bool {
	r.Lock()
	defer r.Unlock()
	return r.Online
}

func (r *RuntimeInfo) SetOnline(online bool) {
	r.Lock()
	defer r.Unlock()
	if !online {
		r.lastFailure = time.Now()
	}
	r.Online = online
}

// IsQueryable determines if a device can be queries.
// This is the case if either the device is online or
// the device is offline and the retryTimeout has elapsed.
// Returns queryable status and if the offline timeout has elapsed.
func (r *RuntimeInfo) IsQueryable() (queryable bool, elapsed bool) {
	r.Lock()
	defer r.Unlock()
	retry := r.lastFailure.Add(retryTimeout).Before(time.Now())
	return r.Online || retry, !r.Online && retry
}
