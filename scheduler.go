package sdm630

import (
	"context"
	"log"
	"strconv"
	"time"

	. "github.com/gonium/gosdm630/internal/meters"
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

type MeterScheduler struct {
	out     QuerySnipChannel
	control ControlSnipChannel
	meters  map[uint8]*Meter
	mc      *MeasurementCache
}

func NewMeterScheduler(
	out QuerySnipChannel,
	control ControlSnipChannel,
	devices map[uint8]*Meter,
) *MeterScheduler {
	return &MeterScheduler{
		out:     out,
		control: control,
		meters:  devices,
	}
}

// SetupScheduler creates a scheduler and its wiring
func SetupScheduler(meters map[uint8]*Meter, qe *ModbusEngine) (*MeterScheduler, QuerySnipChannel) {
	// Create Channels that link the goroutines
	var out = make(QuerySnipChannel)
	var control = make(ControlSnipChannel)
	var tee = make(QuerySnipChannel)

	scheduler := NewMeterScheduler(
		out,
		control,
		meters,
	)

	go qe.Run(
		out,     // scheduler produceSnips output -> qe input
		control, // qe error -> scheduler handleControl input
		tee,     // qe output -> tee
	)

	return scheduler, tee
}

func (q *MeterScheduler) SetCache(mc *MeasurementCache) {
	q.mc = mc
}

// supervisor restarts failed meters by pinging them at regular interval
func (q *MeterScheduler) supervisor() {
	for {
		for _, meter := range q.meters {
			if meter.State() == UNAVAILABLE {
				log.Printf("Attempting to ping unavailable device %d", meter.DeviceId)
				// inject probe snip - the re-enabling logic is in Run()
				operation := meter.Producer.Probe()
				snip := NewQuerySnip(meter.DeviceId, operation)
				q.out <- snip
			}
		}
		time.Sleep(1 * time.Minute)
	}
}

// produceQuerySnips cycles all meters to create reading operations
func (q *MeterScheduler) produceQuerySnips(done <-chan bool, rate time.Duration) {
	defer close(q.out)
	rateMap := make(RateMap)

	for {
		select {
		case <-done:
			return // trigger closing out channel
		default:
			var meterAvailable bool
			for _, meter := range q.meters {
				// check if meter is still valid
				if meter.GetState() == UNAVAILABLE {
					continue
				}

				// rate limiting with early exit if signaled
				wait := rateMap.CooldownDuration(rate, strconv.Itoa(int(meter.DeviceId)))
				select {
				case <-time.After(wait):
				case <-done:
					return
				}

				meterAvailable = true
				operations := meter.Producer.Produce()
				for _, operation := range operations {
					snip := NewQuerySnip(meter.DeviceId, operation)
					q.out <- snip
				}
			}

			// wait before retry if no meter is available
			if !meterAvailable {
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

// handleControlSnips manages the meter status
func (q *MeterScheduler) handleControlSnips() {
	for controlSnip := range q.control {
		meter, ok := q.meters[controlSnip.DeviceId]
		if !ok {
			log.Fatal("Internal device id mismatch")
		}

		switch controlSnip.Type {
		case CONTROLSNIP_ERROR:
			// search meter and deactivate it...
			log.Printf("Device %d failed terminally due to: %s",
				controlSnip.DeviceId, controlSnip.Message)
			if meter.State() == AVAILABLE && q.mc != nil {
				// purge cache if present
				q.mc.Purge(meter.DeviceId)
			}
			meter.SetState(UNAVAILABLE)
		case CONTROLSNIP_OK:
			// search meter and reactivate it...
			if meter.State() != AVAILABLE {
				log.Printf("Reactivating device %d", controlSnip.DeviceId)
				meter.SetState(AVAILABLE)
			}
		default:
			log.Fatal("Unknown control snip")
		}
	}
}

// Run scheduler starts production of meter readings
func (q *MeterScheduler) Run(ctx context.Context, rate time.Duration) {
	done := make(chan bool)

	go q.supervisor()
	go q.produceQuerySnips(done, rate)
	go q.handleControlSnips()

	// wait for cancel
	select {
	case <-ctx.Done():
		done <- true
	}
}
