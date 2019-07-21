package server

import (
	"time"
)

const (
	retryTimeout = 1 * time.Second
)

// RuntimeInfo represents a single modbus device status
type RuntimeInfo struct {
	lastFailure time.Time
	initialized bool
	Online      bool
	Requests    uint64
	Errors      uint64
}

// Available sets the device online status
func (r *RuntimeInfo) Available(online bool) {
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
	retry := r.lastFailure.Add(retryTimeout).Before(time.Now())
	return r.Online || retry, !r.Online && retry
}
