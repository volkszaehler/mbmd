package sdm630

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

/***
 * Opcodes as defined by Eastron.
 * See http://bg-etech.de/download/manual/SDM630Register.pdf
 */
const (
	OpCodeL1Voltage = 0x0000
	OpCodeL2Voltage = 0x0002
	OpCodeL3Voltage = 0x0004
	OpCodeL1Current = 0x0006
	OpCodeL2Current = 0x0008
	OpCodeL3Current = 0x000A
	OpCodeL1Power   = 0x000C
	OpCodeL2Power   = 0x000E
	OpCodeL3Power   = 0x0010
	OpCodeL1Import  = 0x015a
	OpCodeL2Import  = 0x015c
	OpCodeL3Import  = 0x015e
	OpCodeL1Export  = 0x0160
	OpCodeL2Export  = 0x0162
	OpCodeL3Export  = 0x0164
	OpCodeL1Cosphi  = 0x001e
	OpCodeL2Cosphi  = 0x0020
	OpCodeL3Cosphi  = 0x0022
)

/***
 * This is the definition of the Reading datatype. It combines readings
 * of all measurements into one data structure
 */

type ReadingChannel chan Readings

type Readings struct {
	Timestamp      time.Time
	Unix           int64
	ModbusDeviceId uint8
	Power          ThreePhaseReadings
	Voltage        ThreePhaseReadings
	Current        ThreePhaseReadings
	Cosphi         ThreePhaseReadings
	Import         ThreePhaseReadings
	Export         ThreePhaseReadings
}

type ThreePhaseReadings struct {
	L1 float32
	L2 float32
	L3 float32
}

func (r *Readings) String() string {
	fmtString := "ID: %d T: %s - L1: %.2fV %.2fA %.2fW %.2fcos | " +
		"L2: %.2fV %.2fA %.2fW %.2fcos | " +
		"L3: %.2fV %.2fA %.2fW %.2fcos"
	return fmt.Sprintf(fmtString,
		r.ModbusDeviceId,
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
 * Returns true if the reading is older than the given timestamp.
 */
func (r *Readings) IsOlderThan(ts time.Time) (retval bool) {
	return r.Timestamp.Before(ts)
}

/*
* Adds two readings. The individual values are added except for
* the time: the latter of the two times is copied over to the result
 */
func (lhs *Readings) add(rhs *Readings) (retval Readings, err error) {
	if lhs.ModbusDeviceId != rhs.ModbusDeviceId {
		return Readings{}, fmt.Errorf(
			"Cannot add readings of different devices - got IDs %d and %d",
			lhs.ModbusDeviceId, rhs.ModbusDeviceId)
	} else {
		retval = Readings{
			ModbusDeviceId: lhs.ModbusDeviceId,
			Voltage: ThreePhaseReadings{
				L1: lhs.Voltage.L1 + rhs.Voltage.L1,
				L2: lhs.Voltage.L2 + rhs.Voltage.L2,
				L3: lhs.Voltage.L3 + rhs.Voltage.L3,
			},
			Current: ThreePhaseReadings{
				L1: lhs.Current.L1 + rhs.Current.L1,
				L2: lhs.Current.L2 + rhs.Current.L2,
				L3: lhs.Current.L3 + rhs.Current.L3,
			},
			Power: ThreePhaseReadings{
				L1: lhs.Power.L1 + rhs.Power.L1,
				L2: lhs.Power.L2 + rhs.Power.L2,
				L3: lhs.Power.L3 + rhs.Power.L3,
			},
			Cosphi: ThreePhaseReadings{
				L1: lhs.Cosphi.L1 + rhs.Cosphi.L1,
				L2: lhs.Cosphi.L2 + rhs.Cosphi.L2,
				L3: lhs.Cosphi.L3 + rhs.Cosphi.L3,
			},
		}
		if lhs.Timestamp.After(rhs.Timestamp) {
			retval.Timestamp = lhs.Timestamp
			retval.Unix = lhs.Unix
		} else {
			retval.Timestamp = rhs.Timestamp
			retval.Unix = rhs.Unix
		}
		return retval, nil
	}
}

/*
* Dive a reading by an integer. The individual values are divided except
* for the time: it is simply copied over to the result
 */
func (lhs *Readings) divide(scalar float32) (retval Readings) {
	retval = Readings{
		Voltage: ThreePhaseReadings{
			L1: lhs.Voltage.L1 / scalar,
			L2: lhs.Voltage.L2 / scalar,
			L3: lhs.Voltage.L3 / scalar,
		},
		Current: ThreePhaseReadings{
			L1: lhs.Current.L1 / scalar,
			L2: lhs.Current.L2 / scalar,
			L3: lhs.Current.L3 / scalar,
		},
		Power: ThreePhaseReadings{
			L1: lhs.Power.L1 / scalar,
			L2: lhs.Power.L2 / scalar,
			L3: lhs.Power.L3 / scalar,
		},
		Cosphi: ThreePhaseReadings{
			L1: lhs.Cosphi.L1 / scalar,
			L2: lhs.Cosphi.L2 / scalar,
			L3: lhs.Cosphi.L3 / scalar,
		},
	}
	retval.Timestamp = lhs.Timestamp
	retval.Unix = lhs.Unix
	retval.ModbusDeviceId = lhs.ModbusDeviceId
	return retval
}

/* ReadingSlice is a type alias for a slice of readings.
 */
type ReadingSlice []Readings

func (r ReadingSlice) JSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(r)
}

func (r ReadingSlice) NotOlderThan(ts time.Time) (retval ReadingSlice) {
	retval = ReadingSlice{}
	for _, reading := range r {
		if !reading.IsOlderThan(ts) {
			retval = append(retval, reading)
		}
	}
	return retval
}

/***
 * A QuerySnip is just a little snippet of query information. It
 * encapsulates modbus query targets.
 */

type QuerySnip struct {
	DeviceId uint8
	OpCode   uint16
	Value    float32
}

func (r *Readings) MergeSnip(q QuerySnip) {
	switch q.OpCode {
	case OpCodeL1Voltage:
		r.Voltage.L1 = q.Value
	case OpCodeL2Voltage:
		r.Voltage.L2 = q.Value
	case OpCodeL3Voltage:
		r.Voltage.L3 = q.Value
	case OpCodeL1Current:
		r.Current.L1 = q.Value
	case OpCodeL2Current:
		r.Current.L2 = q.Value
	case OpCodeL3Current:
		r.Current.L3 = q.Value
	case OpCodeL1Cosphi:
		r.Cosphi.L1 = q.Value
	case OpCodeL2Cosphi:
		r.Cosphi.L2 = q.Value
	case OpCodeL3Cosphi:
		r.Cosphi.L3 = q.Value
	case OpCodeL1Import:
		r.Import.L1 = q.Value
	case OpCodeL2Import:
		r.Import.L2 = q.Value
	case OpCodeL3Import:
		r.Import.L3 = q.Value
	case OpCodeL1Export:
		r.Export.L1 = q.Value
	case OpCodeL2Export:
		r.Export.L2 = q.Value
	case OpCodeL3Export:
		r.Export.L3 = q.Value
	}
}
