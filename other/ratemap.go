package server

import (
	"time"
)

type RateMap map[string]int64

// Allowed checks if topic has been published longer than rate ago
func (r *RateMap) Allowed(rate int, topic string) bool {
	if rate == 0 {
		return true
	}

	t := (*r)[topic]
	now := time.Now().Unix()
	if now > t {
		(*r)[topic] = now + int64(rate)
		return true
	}

	return false
}

// WaitForCooldown waits until the rate limit has been honored
func (r *RateMap) WaitForCooldown(rate int, topic string) {
	if rate == 0 {
		return
	}

	t := (*r)[topic]
	waituntil := t + int64(rate)*1e9 // use ns
	now := time.Now().UnixNano()

	if waituntil > now {
		time.Sleep(time.Until(time.Unix(0, waituntil)))
		(*r)[topic] = waituntil
	} else {
		(*r)[topic] = now
	}
}

// CooldownDuration returns the time duration to wait for the cooldown period
// to expire. It updates the rate map assuming that the cooldown duration is honored.
func (r *RateMap) CooldownDuration(rate time.Duration, topic string) time.Duration {
	if rate == 0 {
		return time.Duration(0)
	}

	t := (*r)[topic]
	waituntil := time.Unix(0, t).Add(rate)
	remaining := time.Until(waituntil) // use ns

	if remaining <= 0 {
		(*r)[topic] = time.Now().UnixNano()
		return time.Duration(0)
	}

	(*r)[topic] = waituntil.UnixNano()
	return remaining
}
