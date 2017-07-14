package sdm630

const (
	METER_JANITZA = "JANITZA"
	METER_SDM     = "SDM"
)

type Meter struct {
	Scheduler Scheduler
	DevId     uint8
	// TODO: Define state etc.
	State uint8
}

func NewMeter(
	devid uint8,
	scheduler Scheduler,
) *Meter {
	return &Meter{
		Scheduler: scheduler,
		DevId:     devid,
	}
}
