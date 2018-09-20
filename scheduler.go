package sdm630

import (
	"log"
	"time"
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
			sniplist := meter.Producer.Produce(meter.DeviceId)
			for _, snip := range sniplist {
				// Check if meter is still valid
				if meter.GetState() != METERSTATE_UNAVAILABLE {
					q.out <- snip
				}
			}
		}
	}
}

func (q *MeterScheduler) supervisor() {
	for {
		for _, meter := range q.meters {
			if meter.GetState() == METERSTATE_UNAVAILABLE {
				log.Printf("Attempting to ping unavailable meter %d", meter.DeviceId)
				// inject probe snip - the re-enabling logic is in Run()
				q.out <- meter.Producer.Probe(meter.DeviceId)
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
				meter, ok := q.meters[controlSnip.DeviceId]
				if !ok {
					log.Fatal("Internal device id mismatch - this should not happen!")
				} else {
					meter.UpdateState(METERSTATE_UNAVAILABLE)
				}
			case CONTROLSNIP_OK:
				// search meter and reactivate it...
				meter, ok := q.meters[controlSnip.DeviceId]
				if !ok {
					log.Fatal("Internal device id mismatch - this should not happen!")
				} else {
					if meter.GetState() != METERSTATE_AVAILABLE {
						log.Printf("Re-activating meter %d", controlSnip.DeviceId)
						meter.UpdateState(METERSTATE_AVAILABLE)
					}
				}
			default:
				log.Fatal("Received unknown control snip - something weird happened.")
			}
		}
	}
}
