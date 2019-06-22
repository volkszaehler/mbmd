package sunspec

import (
	"errors"
	"math"
	"time"

	sunspec "github.com/crabmusket/gosunspec"
	sunspecbus "github.com/crabmusket/gosunspec/modbus"

	_ "github.com/crabmusket/gosunspec/models" // device tree parsing requires all models
	"github.com/crabmusket/gosunspec/models/model1"

	"github.com/grid-x/modbus"
	. "github.com/volkszaehler/mbmd/meters"
)

type sunSpec struct {
	models     []sunspec.Model
	descriptor DeviceDescriptor
}

func NewDevice() Device {
	return &sunSpec{}
}

func (d *sunSpec) Initialize(client modbus.Client) error {
	in, err := sunspecbus.Open(client)
	if err != nil {
		return err
	}

	devices := in.Collect(sunspec.AllDevices)
	if len(devices) == 0 {
		return errors.New("sunspec: device not found")
	}
	if len(devices) > 1 {
		return errors.New("sunspec: multiple devices found")
	}

	device := devices[0]

	// read common block
	if err := d.readCommonBlock(device); err != nil {
		return err
	}

	// collect relevant models
	if err := d.collectModels(device); err != nil {
		return err
	}

	return nil
}

func (d *sunSpec) readCommonBlock(device sunspec.Device) error {
	// TODO catch panic
	commonModel := device.MustModel(sunspec.ModelId(1))
	// TODO catch panic
	b := commonModel.MustBlock(0)
	if err := b.Read(); err != nil {
		return err
	}

	d.descriptor = DeviceDescriptor{
		Manufacturer: b.MustPoint(model1.Mn).StringValue(),
		Model:        b.MustPoint(model1.Md).StringValue(),
		Options:      b.MustPoint(model1.Opt).StringValue(),
		Version:      b.MustPoint(model1.Vr).StringValue(),
		Serial:       b.MustPoint(model1.SN).StringValue(),
	}

	return nil
}

// collect and sort supported models except for common
func (d *sunSpec) collectModels(device sunspec.Device) error {
	d.models = device.Collect(sunspec.OneOfSeveralModelIds(d.relevantModelIds()))
	if len(d.models) == 0 {
		return errors.New("sunspec: could not find supported model")
	}

	// sanitizeModels()
	return nil
}

func (d *sunSpec) relevantModelIds() []sunspec.ModelId {
	modelIds := make([]sunspec.ModelId, 0, len(modelPoints))
	for k := range modelPoints {
		modelIds = append(modelIds, sunspec.ModelId(k))
	}

	return modelIds
}

// remove model 101 if model 103 found
func (d *sunSpec) sanitizeModels() {
	m101 := -1
	for i, m := range d.models {
		if m.Id() == sunspec.ModelId(101) {
			m101 = i
		}
		if m101 >= 0 && m.Id() == sunspec.ModelId(103) {
			d.models = append(d.models[0:m101], d.models[m101+1:]...)
			break
		}
	}
}

func (d *sunSpec) Descriptor() DeviceDescriptor {
	return d.descriptor
}

func (d *sunSpec) Probe(client modbus.Client) (MeasurementResult, error) {
	return MeasurementResult{},nil
}

func (d *sunSpec) Query(client modbus.Client) ([]MeasurementResult, error) {
	res := make([]MeasurementResult, 0)

	for _, model := range d.models {
		// TODO catch panic
		b := model.MustBlock(0)
		if err := b.Read(); err != nil {
			return res, err
		}

		for _, pointID := range modelPoints[int(model.Id())] {
			// TODO catch panic
			p := b.MustPoint(pointID)

			if m, ok := opcodeMap[pointID]; !ok {
				panic("sunspec: no measurement for point id " + pointID)
			} else {
				v := ScaledValue(p)
				if math.IsNaN(v) {
					continue
				}

				mr := MeasurementResult{
					Measurement: m,
					Value:       v,
					Timestamp:   time.Now(),
				}

				res = append(res, mr)
			}
		}
	}

	return res, nil
}

// ScaledValue uses Point.ScaledValue and checks the result for NaN-encoded values
func ScaledValue(p sunspec.Point) float64 {
	f := p.ScaledValue()

	switch p.Type() {
	case "int16":
		if p.Value() == int16(math.MinInt16) {
			f = math.NaN()
		}
	case "int32":
		if p.Value() == int32(math.MinInt32) {
			f = math.NaN()
		}
	case "int64":
		if p.Value() == int64(math.MinInt64) {
			f = math.NaN()
		}
	case "uint16":
		if p.Value() == uint16(math.MaxUint16) {
			f = math.NaN()
		}
	case "uint32":
		if p.Value() == uint32(math.MaxUint32) {
			f = math.NaN()
		}
	case "uint64":
		if p.Value() == uint64(math.MaxUint64) {
			f = math.NaN()
		}
	}

	return f
}
