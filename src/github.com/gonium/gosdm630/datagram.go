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
	Power     PowerReadings
	Voltage   VoltageReadings
	Current   CurrentReadings
	Cosphi    CosphiReadings
}

type PowerReadings struct {
	L1 float32
	L2 float32
	L3 float32
}

type VoltageReadings struct {
	L1 float32
	L2 float32
	L3 float32
}

type CurrentReadings struct {
	L1 float32
	L2 float32
	L3 float32
}

type CosphiReadings struct {
	L1 float32
	L2 float32
	L3 float32
}

func (r *Readings) String() string {
	fmtString := "T: %s - L1: %.2fV %.2fA %.2fW %.2fcos | " +
		"L2: %.2fV %.2fA %.2fW %.2fcos | " +
		"L3: %.2fV %.2fA %.2fW %.2fcos"
	return fmt.Sprintf(fmtString,
		r.Timestamp.Format(time.RFC3339),
		r.Voltage.L1,
		r.Current.L1,
		r.Power.L1,
		r.Cosphi.L1,
		r.Voltage.L2,
		r.Current.L2,
		r.Power.L2,
		r.Cosphi.L2,
		r.Voltage.L3,
		r.Current.L3,
		r.Power.L3,
		r.Cosphi.L3,
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
		Voltage: VoltageReadings{
			L1: lhs.Voltage.L1 + rhs.Voltage.L1,
			L2: lhs.Voltage.L2 + rhs.Voltage.L2,
			L3: lhs.Voltage.L3 + rhs.Voltage.L3,
		},
		Current: CurrentReadings{
			L1: lhs.Current.L1 + rhs.Current.L1,
			L2: lhs.Current.L2 + rhs.Current.L2,
			L3: lhs.Current.L3 + rhs.Current.L3,
		},
		Power: PowerReadings{
			L1: lhs.Power.L1 + rhs.Power.L1,
			L2: lhs.Power.L2 + rhs.Power.L2,
			L3: lhs.Power.L3 + rhs.Power.L3,
		},
		Cosphi: CosphiReadings{
			L1: lhs.Cosphi.L1 + rhs.Cosphi.L1,
			L2: lhs.Cosphi.L2 + rhs.Cosphi.L2,
			L3: lhs.Cosphi.L3 + rhs.Cosphi.L3,
		},
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
		Voltage: VoltageReadings{
			L1: lhs.Voltage.L1 / scalar,
			L2: lhs.Voltage.L2 / scalar,
			L3: lhs.Voltage.L3 / scalar,
		},
		Current: CurrentReadings{
			L1: lhs.Current.L1 / scalar,
			L2: lhs.Current.L2 / scalar,
			L3: lhs.Current.L3 / scalar,
		},
		Power: PowerReadings{
			L1: lhs.Power.L1 / scalar,
			L2: lhs.Power.L2 / scalar,
			L3: lhs.Power.L3 / scalar,
		},
		Cosphi: CosphiReadings{
			L1: lhs.Cosphi.L1 / scalar,
			L2: lhs.Cosphi.L2 / scalar,
			L3: lhs.Cosphi.L3 / scalar,
		},
	}
	retval.Timestamp = lhs.Timestamp
	return retval
}
