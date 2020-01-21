package sunspec

import (
	"fmt"
	"math"
	"time"

	sunspec "github.com/andig/gosunspec"
	sunspecbus "github.com/andig/gosunspec/modbus"
	"github.com/grid-x/modbus"

	_ "github.com/andig/gosunspec/models" // device tree parsing requires all models
	"github.com/andig/gosunspec/models/model1"
	"github.com/pkg/errors"
	"github.com/volkszaehler/mbmd/meters"
)

type sunSpec struct {
	models     []sunspec.Model
	descriptor meters.DeviceDescriptor
}

// partialError can be behaviour-checked for SunSpecPartiallyInitialized()
// to indicate initialization error
type partialError struct {
	error
	cause error
}

// Cause implements errors.Causer()
func (e partialError) Cause() error {
	return e.cause
}

// PartiallyInitialized implements SunSpecPartiallyInitialized()
func (e partialError) PartiallyInitialized() {}

// NewDevice creates a Sunspec device
func NewDevice(meterType string) meters.Device {
	return &sunSpec{
		descriptor: meters.DeviceDescriptor{
			Manufacturer: meterType,
		},
	}
}

func (d *sunSpec) Initialize(client modbus.Client) error {
	in, err := sunspecbus.Open(client)
	if err != nil && in == nil {
		return err
	} else if err != nil {
		err = partialError{
			error: errors.New("sunspec: device opened partially"),
			cause: err,
		}
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

	return err
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
	modelIds := make([]sunspec.ModelId, 0, len(modelMap))
	for k := range modelMap {
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

func (d *sunSpec) Probe(client modbus.Client) (bool, error) {
	if d.notInitialized() {
		return false, errors.New("sunspec: not initialized")
	}

	for _, model := range d.models {
		for id := range modelMap {
			if id == model.Id() {
				return true, nil
			}
		}
	}

	return false, fmt.Errorf("sunspec: could not find suitable model")
}

func (d *sunSpec) notInitialized() bool {
	return len(d.models) == 0
}

func (d *sunSpec) convertPoint(b sunspec.Block, blockID int, pointID string, m meters.Measurement) (meters.MeasurementResult, error) {
	p := b.MustPoint(pointID)
	v := p.ScaledValue()

	if math.IsNaN(v) {
		return meters.MeasurementResult{}, errors.New("NaN value")
	}

	// apply scale factor for energy
	if div, ok := dividerMap[m]; ok {
		v /= div
	}

	mr := meters.MeasurementResult{
		Measurement: m,
		Value:       v,
		Timestamp:   time.Now(),
	}

	return mr, nil
}

func (d *sunSpec) Query(client modbus.Client) (res []meters.MeasurementResult, err error) {
	if d.notInitialized() {
		return res, errors.New("sunspec: not initialized")
	}

	for _, model := range d.models {
		blockID := 0

		model.Do(func(b sunspec.Block) {
			defer func() { blockID++ }()

			if err = b.Read(); err != nil {
				return
			}

			if bps, ok := modelMap[model.Id()][blockID]; ok {
				for pointID, m := range bps {
					if mr, err := d.convertPoint(b, blockID, pointID, m); err == nil {
						res = append(res, mr)
					}
				}
			}
		})
	}

	return res, nil
}
