package server

import (
	"errors"
	"fmt"
	"math"
	"time"

	. "github.com/volkszaehler/mbmd/meters"
)

// UniqueIdFormat is a format string for unique ID generation.
// It expects one %d conversion specifier,
// which will be replaced with the device ID.
// The UniqueIdFormat can be changed on program startup,
// before any additional goroutines are started.
var UniqueIdFormat string = "Instrument%d"

// Readings combines readings of all measurements into one data structure
type Readings struct {
	UniqueId    string
	Timestamp   time.Time
	Unix        int64
	DeviceId    uint8 `json:"ModbusDeviceId"`
	Power       ThreePhaseReadings
	Voltage     ThreePhaseReadings
	Current     ThreePhaseReadings
	Cosphi      ThreePhaseReadings
	Import      ThreePhaseReadings
	TotalImport *float64
	Export      ThreePhaseReadings
	TotalExport *float64
	THD         THDInfo
	Frequency   *float64
}

type THDInfo struct {
	//	Current           ThreePhaseReadings
	//	AvgCurrent        float64
	VoltageNeutral    ThreePhaseReadings
	AvgVoltageNeutral *float64
}

type ThreePhaseReadings struct {
	L1 *float64
	L2 *float64
	L3 *float64
}

// F2fp helper converts float64 to *float64
func F2fp(x float64) *float64 {
	if math.IsNaN(x) {
		return nil
	}
	return &x
}

// Fp2f helper converts *float64 to float64, correctly handles uninitialized
// variables
func Fp2f(x *float64) float64 {
	if x == nil {
		// this is not initialized yet - return NaN
		return math.NaN()
	}
	return *x
}

func (r *Readings) String() string {
	fmtString := "%s " +
		"L1: %.1fV %.2fA %.0fW %.2fcos | " +
		"L2: %.1fV %.2fA %.0fW %.2fcos | " +
		"L3: %.1fV %.2fA %.0fW %.2fcos | " +
		"%.1fHz"
	return fmt.Sprintf(fmtString,
		r.UniqueId,
		Fp2f(r.Voltage.L1),
		Fp2f(r.Current.L1),
		Fp2f(r.Power.L1),
		Fp2f(r.Cosphi.L1),
		Fp2f(r.Voltage.L2),
		Fp2f(r.Current.L2),
		Fp2f(r.Power.L2),
		Fp2f(r.Cosphi.L2),
		Fp2f(r.Voltage.L3),
		Fp2f(r.Current.L3),
		Fp2f(r.Power.L3),
		Fp2f(r.Cosphi.L3),
		Fp2f(r.Frequency),
	)
}

// IsOlderThan returns true if the reading is older than the given timestamp.
func (r *Readings) IsOlderThan(ts time.Time) (retval bool) {
	return r.Timestamp.Before(ts)
}

func tpAdd(lhs ThreePhaseReadings, rhs ThreePhaseReadings) ThreePhaseReadings {
	res := ThreePhaseReadings{
		L1: F2fp(Fp2f(lhs.L1) + Fp2f(rhs.L1)),
		L2: F2fp(Fp2f(lhs.L2) + Fp2f(rhs.L2)),
		L3: F2fp(Fp2f(lhs.L3) + Fp2f(rhs.L3)),
	}
	return res
}

/*
* Adds two readings. The individual values are added except for
* the time: the latter of the two times is copied over to the result
 */
func (lhs *Readings) add(rhs *Readings) (*Readings, error) {
	if lhs.DeviceId != rhs.DeviceId {
		return &Readings{}, fmt.Errorf(
			"Cannot add readings of different devices - got IDs %d and %d",
			lhs.DeviceId, rhs.DeviceId)
	}

	res := &Readings{
		UniqueId: lhs.UniqueId,
		DeviceId: lhs.DeviceId,
		Voltage:  tpAdd(lhs.Voltage, rhs.Voltage),
		Current:  tpAdd(lhs.Current, rhs.Current),
		Power:    tpAdd(lhs.Power, rhs.Power),
		Cosphi:   tpAdd(lhs.Cosphi, rhs.Cosphi),
		Import:   tpAdd(lhs.Import, rhs.Import),
		TotalImport: F2fp(Fp2f(lhs.TotalImport) +
			Fp2f(rhs.TotalImport)),
		Export: tpAdd(lhs.Export, rhs.Export),
		TotalExport: F2fp(Fp2f(lhs.TotalExport) +
			Fp2f(rhs.TotalExport)),
		THD: THDInfo{
			VoltageNeutral: tpAdd(lhs.THD.VoltageNeutral, rhs.THD.VoltageNeutral),
			AvgVoltageNeutral: F2fp(Fp2f(lhs.THD.AvgVoltageNeutral) +
				Fp2f(rhs.THD.AvgVoltageNeutral)),
		},
		Frequency: F2fp(Fp2f(lhs.Frequency) +
			Fp2f(rhs.Frequency)),
	}

	if lhs.Timestamp.After(rhs.Timestamp) {
		res.Timestamp = lhs.Timestamp
		res.Unix = lhs.Unix
	} else {
		res.Timestamp = rhs.Timestamp
		res.Unix = rhs.Unix
	}

	return res, nil
}

func tpDivide(lhs ThreePhaseReadings, scaler float64) ThreePhaseReadings {
	res := ThreePhaseReadings{
		L1: F2fp(Fp2f(lhs.L1) / scaler),
		L2: F2fp(Fp2f(lhs.L2) / scaler),
		L3: F2fp(Fp2f(lhs.L3) / scaler),
	}
	return res
}

/*
 * Divide a reading by an integer. The individual values are divided except
 * for the time: it is simply copied over to the result
 */
func (lhs *Readings) divide(scaler float64) *Readings {
	res := &Readings{
		Timestamp: lhs.Timestamp,
		Unix:      lhs.Unix,
		DeviceId:  lhs.DeviceId,
		UniqueId:  lhs.UniqueId,

		Voltage:     tpDivide(lhs.Voltage, scaler),
		Current:     tpDivide(lhs.Current, scaler),
		Power:       tpDivide(lhs.Power, scaler),
		Cosphi:      tpDivide(lhs.Cosphi, scaler),
		Import:      tpDivide(lhs.Import, scaler),
		TotalImport: F2fp(Fp2f(lhs.TotalImport) / scaler),
		Export:      tpDivide(lhs.Export, scaler),
		TotalExport: F2fp(Fp2f(lhs.TotalExport) / scaler),
		THD: THDInfo{
			VoltageNeutral:    tpDivide(lhs.THD.VoltageNeutral, scaler),
			AvgVoltageNeutral: F2fp(Fp2f(lhs.THD.AvgVoltageNeutral) / scaler),
		},
		Frequency: F2fp(Fp2f(lhs.Frequency) / scaler),
	}
	return res
}

// MergeSnip adds the values represented by the QuerySnip to the
// Readings and updates the current time stamp
func (r *Readings) MergeSnip(q QuerySnip) {
	r.Timestamp = q.ReadTimestamp
	r.Unix = r.Timestamp.Unix()
	switch q.IEC61850 {
	case VoltageL1:
		r.Voltage.L1 = &q.Value
	case VoltageL2:
		r.Voltage.L2 = &q.Value
	case VoltageL3:
		r.Voltage.L3 = &q.Value
	case CurrentL1:
		r.Current.L1 = &q.Value
	case CurrentL2:
		r.Current.L2 = &q.Value
	case CurrentL3:
		r.Current.L3 = &q.Value
	case PowerL1:
		r.Power.L1 = &q.Value
	case PowerL2:
		r.Power.L2 = &q.Value
	case PowerL3:
		r.Power.L3 = &q.Value
	case CosphiL1:
		r.Cosphi.L1 = &q.Value
	case CosphiL2:
		r.Cosphi.L2 = &q.Value
	case CosphiL3:
		r.Cosphi.L3 = &q.Value
	case ImportL1:
		r.Import.L1 = &q.Value
	case ImportL2:
		r.Import.L2 = &q.Value
	case ImportL3:
		r.Import.L3 = &q.Value
	case Import:
		r.TotalImport = &q.Value
	case ExportL1:
		r.Export.L1 = &q.Value
	case ExportL2:
		r.Export.L2 = &q.Value
	case ExportL3:
		r.Export.L3 = &q.Value
	case Export:
		r.TotalExport = &q.Value
		//	case L1THDCurrent
		//		r.THD.Current.L1 = &q.Value
		//	case L2THDCurrent
		//		r.THD.Current.L2 = &q.Value
		//	case L3THDCurrent
		//		r.THD.Current.L3 = &q.Value
		//	case THDCurrent
		//		r.THD.AvgCurrent = &q.Value
	case THDL1:
		r.THD.VoltageNeutral.L1 = &q.Value
	case THDL2:
		r.THD.VoltageNeutral.L2 = &q.Value
	case THDL3:
		r.THD.VoltageNeutral.L3 = &q.Value
	case THD:
		r.THD.AvgVoltageNeutral = &q.Value
	case Frequency:
		r.Frequency = &q.Value
	default:
		// log.Fatalf("Cannot merge unknown IEC: %+v", q)
	}
}

// ReadingSlice is a type alias for a slice of readings.
type ReadingSlice []Readings

// NotOlderThan creates a new ReadingSlice of latest data
func (r ReadingSlice) NotOlderThan(ts time.Time) (res ReadingSlice) {
	res = ReadingSlice{}
	for _, reading := range r {
		if !reading.IsOlderThan(ts) {
			res = append(res, reading)
		}
	}
	return res
}

// Average calculates average across a ReadingSlice
func (r *ReadingSlice) Average() (avg *Readings, err error) {
	// check for panics
	defer func() {
		if r := recover(); r != nil {
			avg = nil
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
		}
	}()

	for idx, r := range *r {
		if idx == 0 {
			// This is the first element - initialize our accumulator
			avg = &r
		} else {
			avg, err = r.add(avg)
			if err != nil {
				return nil, err
			}
		}
	}

	if len(*r) == 0 {
		return nil, errors.New("readings empty")
	}

	return avg.divide(float64(len(*r))), nil
}
