package sunspec

import (
	"fmt"
	"math"
	"time"

	sunspec "github.com/andig/gosunspec"
	sunspecbus "github.com/andig/gosunspec/modbus"

	_ "github.com/andig/gosunspec/models" // device tree parsing requires all models
	"github.com/andig/gosunspec/models/model1"
	"github.com/andig/gosunspec/models/model101"
	"github.com/pkg/errors"
	"github.com/volkszaehler/mbmd/meters"
)

type sunSpec struct {
	models     []sunspec.Model
	descriptor meters.DeviceDescriptor
}

// NewDevice creates a Sunspec device
func NewDevice() meters.Device {
	return &sunSpec{}
}

func (d *sunSpec) Initialize(client meters.ModbusClient) error {
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

	d.descriptor = meters.DeviceDescriptor{
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

func (d *sunSpec) Descriptor() meters.DeviceDescriptor {
	return d.descriptor
}

func (d *sunSpec) Probe(client meters.ModbusClient) (res meters.MeasurementResult, err error) {
	if d.notInitilized() {
		return res, errors.New("sunspec: not initialized")
	}

	for _, model := range d.models {
		if model.Id() != 101 && model.Id() != 103 {
			continue
		}

		b := model.MustBlock(0)
		if err = b.Read(); err != nil {
			return
		}

		pointID := model101.PhVphA
		p := b.MustPoint(pointID)

		if m, ok := opcodeMap[pointID]; !ok {
			panic("sunspec: no measurement for point id " + pointID)
		} else {
			v := p.ScaledValue()
			if math.IsNaN(v) {
				return res, errors.Wrapf(err, "sunspec: could not read probe snip")
			}

			mr := meters.MeasurementResult{
				Measurement: m,
				Value:       v,
				Timestamp:   time.Now(),
			}

			return mr, nil
		}
	}

	return res, fmt.Errorf("sunspec: could not find model for probe snip")
}

func (d *sunSpec) notInitilized() bool {
	return len(d.models) == 0
}

func (d *sunSpec) Query(client meters.ModbusClient) ([]meters.MeasurementResult, error) {
	res := make([]meters.MeasurementResult, 0)

	if d.notInitilized() {
		return res, errors.New("sunspec: not initialized")
	}

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
				v := p.ScaledValue()
				if math.IsNaN(v) {
					continue
				}

				mr := meters.MeasurementResult{
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
