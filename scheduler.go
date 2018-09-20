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

	go qe.Transform(
		scheduler2queryengine, // input
		queryengine2scheduler, // error
		queryengine2tee,       // output
	)

	return scheduler, queryengine2tee
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
				log.Printf("Attempting to ping unavailable meter %d", meter.DeviceId)
				// inject probe snip - the re-enabling logic is in Run()
				operation := meter.Producer.Probe()
				snip := NewQuerySnip(meter.DeviceId, operation)
				q.out <- snip
			}
		}
		time.Sleep(15 * time.Minute)
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
			switch controlSnip.Type {
			case CONTROLSNIP_ERROR:
				log.Printf("Failure - deactivating meter %d: %s",
					controlSnip.DeviceId, controlSnip.Message)
				// search meter and deactivate it...
				if meter, ok := q.meters[controlSnip.DeviceId]; ok {
					meter.UpdateState(UNAVAILABLE)
				} else {
					log.Fatal("Internal device id mismatch - this should not happen!")
				}
			case CONTROLSNIP_OK:
				// search meter and reactivate it...
				if meter, ok := q.meters[controlSnip.DeviceId]; ok {
					if meter.GetState() != AVAILABLE {
						log.Printf("Reactivating meter %d", controlSnip.DeviceId)
						meter.UpdateState(AVAILABLE)
					}
				} else {
					log.Fatal("Internal device id mismatch - this should not happen!")
				}
			default:
				log.Fatal("Received unknown control snip - something weird happened.")
			}
		}
	}
}
