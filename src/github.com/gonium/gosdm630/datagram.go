package sdm630

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type ReadingChannel chan Readings

type Readings struct {
	Timestamp time.Time
	L1Voltage float32
	L2Voltage float32
	L3Voltage float32
	L1Current float32
	L2Current float32
	L3Current float32
	L1Power   float32
	L2Power   float32
	L3Power   float32
	L1CosPhi  float32
	L2CosPhi  float32
	L3CosPhi  float32
}

func (r *Readings) String() string {
	fmtString := "T: %s - L1: %.2fV %.2fA %.2fW %.2fcos | " +
		"L2: %.2fV %.2fA %.2fW %.2fcos | " +
		"L3: %.2fV %.2fA %.2fW %.2fcos"
	return fmt.Sprintf(fmtString,
		r.Timestamp.Format(time.RFC3339),
		r.L1Voltage,
		r.L1Current,
		r.L1Power,
		r.L1CosPhi,
		r.L2Voltage,
		r.L2Current,
		r.L2Power,
		r.L2CosPhi,
		r.L3Voltage,
		r.L3Current,
		r.L3Power,
		r.L3CosPhi,
	)
}

func (r *Readings) JSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(r)
}

/*
* Adds two readings. The individual values are added except for
* the time: the latter of the two times is copied over to the result
 */
func (lhs *Readings) add(rhs *Readings) (retval Readings) {
	retval = Readings{
		L1Voltage: lhs.L1Voltage + rhs.L1Voltage,
		L2Voltage: lhs.L2Voltage + rhs.L2Voltage,
		L3Voltage: lhs.L3Voltage + rhs.L3Voltage,
		L1Current: lhs.L1Current + rhs.L1Current,
		L2Current: lhs.L2Current + rhs.L2Current,
		L3Current: lhs.L3Current + rhs.L3Current,
		L1Power:   lhs.L1Power + rhs.L1Power,
		L2Power:   lhs.L2Power + rhs.L2Power,
		L3Power:   lhs.L3Power + rhs.L3Power,
		L1CosPhi:  lhs.L1CosPhi + rhs.L1CosPhi,
		L2CosPhi:  lhs.L2CosPhi + rhs.L2CosPhi,
		L3CosPhi:  lhs.L3CosPhi + rhs.L3CosPhi,
	}
	if lhs.Timestamp.After(rhs.Timestamp) {
		retval.Timestamp = lhs.Timestamp
	} else {
		retval.Timestamp = rhs.Timestamp
	}
	return retval
}

/*
* Dive a reading by an integer. The individual values are divided except
* for the time: it is simply copied over to the result
 */
func (lhs *Readings) divide(scalar float32) (retval Readings) {
	retval = Readings{
		L1Voltage: lhs.L1Voltage / scalar,
		L2Voltage: lhs.L2Voltage / scalar,
		L3Voltage: lhs.L3Voltage / scalar,
		L1Current: lhs.L1Current / scalar,
		L2Current: lhs.L2Current / scalar,
		L3Current: lhs.L3Current / scalar,
		L1Power:   lhs.L1Power / scalar,
		L2Power:   lhs.L2Power / scalar,
		L3Power:   lhs.L3Power / scalar,
		L1CosPhi:  lhs.L1CosPhi / scalar,
		L2CosPhi:  lhs.L2CosPhi / scalar,
		L3CosPhi:  lhs.L3CosPhi / scalar,
	}
	retval.Timestamp = lhs.Timestamp
	return retval
}
