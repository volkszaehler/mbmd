package sdm630

import (
	"log"
	"time"

	. "github.com/gonium/gosdm630/internal/meters"
)

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
		meters:  devices,
		control: control,
	}
}

// SetupScheduler creates a scheduler and its wiring
func SetupScheduler(meters map[uint8]*Meter, qe *ModbusEngine) (*MeterScheduler, QuerySnipChannel) {
	// Create Channels that link the goroutines
	var scheduler2queryengine = make(QuerySnipChannel)
	var queryengine2scheduler = make(ControlSnipChannel)
	var queryengine2tee = make(QuerySnipChannel)

	scheduler := NewMeterScheduler(
		scheduler2queryengine,
		queryengine2scheduler,
		meters,
	)

	go qe.Run(
		scheduler2queryengine, // input
		queryengine2scheduler, // error
		queryengine2tee,       // output
	)

	return scheduler, queryengine2tee
}

func (q *MeterScheduler) SetCache(mc *MeasurementCache) {
	q.mc = mc
}

func (q *MeterScheduler) produceSnips(out QuerySnipChannel) {
	for {
		for _, meter := range q.meters {
			operations := meter.Producer.Produce()
			for _, operation := range operations {
				// Check if meter is still valid
				if meter.GetState() != UNAVAILABLE {
					snip := NewQuerySnip(meter.DeviceId, operation)
					q.out <- snip
				}
			}
		}
	}
}

func (q *MeterScheduler) supervisor() {
	for {
		for _, meter := range q.meters {
			if meter.GetState() == UNAVAILABLE {
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

func (q *MeterScheduler) Run() {
	source := make(QuerySnipChannel)

	go q.supervisor()
	go q.produceSnips(source)

	for {
		select {
		case snip := <-source:
			q.out <- snip

		case controlSnip := <-q.control:
			meter, ok := q.meters[controlSnip.DeviceId]
			if !ok {
				log.Fatal("Internal device id mismatch")
			}

			switch controlSnip.Type {
			case CONTROLSNIP_ERROR:
				// search meter and deactivate it...
				log.Printf("Device %d failed terminally due to: %s",
					controlSnip.DeviceId, controlSnip.Message)
				state := meter.GetState()
				meter.UpdateState(UNAVAILABLE)
				if state == AVAILABLE && q.mc != nil {
					// purge cache if present
					q.mc.Purge(meter.DeviceId)
				}
			case CONTROLSNIP_OK:
				// search meter and reactivate it...
				if meter.GetState() != AVAILABLE {
					log.Printf("Reactivating device %d", controlSnip.DeviceId)
					meter.UpdateState(AVAILABLE)
				}
			default:
				log.Fatal("Unknown control snip")
			}
		}
	}
}
