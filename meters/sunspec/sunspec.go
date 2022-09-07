package sunspec

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	sunspec "github.com/andig/gosunspec"
	sunspecbus "github.com/andig/gosunspec/modbus"
	"github.com/grid-x/modbus"

	_ "github.com/andig/gosunspec/models" // device tree parsing requires all models
	"github.com/andig/gosunspec/models/model1"
	"github.com/andig/gosunspec/models/model101"
	"github.com/volkszaehler/mbmd/meters"
)

// SunSpec is the sunspec device implementation
type SunSpec struct {
	subdevice  int
	models     []sunspec.Model
	descriptor meters.DeviceDescriptor
}

// FixKostal implements workaround for negative KOSTAL values (https://github.com/volkszaehler/mbmd/pull/97)
func FixKostal(p sunspec.Point) {
	switch t := p.Value().(type) {
	case sunspec.Acc32:
		if t > sunspec.Acc32(math.MaxInt32) {
			p.SetAcc32(-p.Value().(sunspec.Acc32))
		}
	case sunspec.Acc64:
		if t > sunspec.Acc64(math.MaxInt64) {
			p.SetAcc64(-p.Value().(sunspec.Acc64))
		}
	}
}

// NewDevice creates a Sunspec device
func NewDevice(meterType string, subdevice ...int) *SunSpec {
	var dev int
	if len(subdevice) > 0 {
		dev = subdevice[0]
	}

	return &SunSpec{
		subdevice: dev,
		descriptor: meters.DeviceDescriptor{
			Type:         meterType,
			Manufacturer: meterType,
			SubDevice:    dev,
		},
	}
}

// Initialize implements the Device interface
func (d *SunSpec) Initialize(client modbus.Client) error {
	var partiallyOpen bool
	in, err := sunspecbus.Open(client)
	if err != nil {
		if in == nil {
			return err
		}

		partiallyOpen = true
	}

	devices := in.Collect(sunspec.AllDevices)
	if len(devices) == 0 {
		return errors.New("sunspec: device not found")
	}
	if len(devices) <= d.subdevice {
		return fmt.Errorf("sunspec: subdevice %d not found", d.subdevice)
	}

	device := devices[d.subdevice]

	// read common block
	if err := d.readCommonBlock(device); err != nil {
		return err
	}

	// collect relevant models
	if err := d.collectModels(device); err != nil {
		return err
	}

	// return partial open error if everything else went fine
	if partiallyOpen {
		err = fmt.Errorf("%w", meters.ErrPartiallyOpened)
	}

	return err
}

func stringVal(b sunspec.Block, point string) string {
	return strings.TrimSpace(b.MustPoint(point).StringValue())
}

func (d *SunSpec) readCommonBlock(device sunspec.Device) error {
	// TODO catch panic
	commonModel := device.MustModel(sunspec.ModelId(1))
	// TODO catch panic
	b := commonModel.MustBlock(0)
	if err := b.Read(); err != nil {
		return err
	}

	d.descriptor.Manufacturer = stringVal(b, model1.Mn)
	d.descriptor.Model = stringVal(b, model1.Md)
	d.descriptor.Options = stringVal(b, model1.Opt)
	d.descriptor.Version = stringVal(b, model1.Vr)
	d.descriptor.Serial = stringVal(b, model1.SN)

	return nil
}

// collect and sort supported models except for common
func (d *SunSpec) collectModels(device sunspec.Device) error {
	d.models = device.Collect(sunspec.AllModels)

	// don't error for sake of QueryPoint
	// if len(device.Collect(sunspec.OneOfSeveralModelIds(d.relevantModelIds()))) == 0 {
	// 	return errors.New("sunspec: could not find supported model")
	// }

	// sanitizeModels()
	return nil
}

func (d *SunSpec) relevantModelIds() []sunspec.ModelId {
	modelIds := make([]sunspec.ModelId, 0, len(modelMap))
	for k := range modelMap {
		modelIds = append(modelIds, sunspec.ModelId(k))
	}

	return modelIds
}

// remove model 101 if model 103 found
func (d *SunSpec) sanitizeModels() {
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

// Descriptor implements the Device interface
func (d *SunSpec) Descriptor() meters.DeviceDescriptor {
	return d.descriptor
}

// Probe implements the Device interface
func (d *SunSpec) Probe(client modbus.Client) (res meters.MeasurementResult, err error) {
	if d.notInitialized() {
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

		v := p.ScaledValue()
		if math.IsNaN(v) {
			return res, fmt.Errorf("%w", meters.ErrNaN)
		}

		mr := meters.MeasurementResult{
			Measurement: meters.Current,
			Value:       v,
			Timestamp:   time.Now(),
		}

		return mr, nil
	}

	return res, fmt.Errorf("sunspec: could not find model for probe snip")
}

func (d *SunSpec) notInitialized() bool {
	return len(d.models) == 0
}

func (d *SunSpec) convertPoint(b sunspec.Block, p sunspec.Point) (float64, error) {
	if d.descriptor.Manufacturer == "KOSTAL" {
		FixKostal(p)
	}

	v := p.ScaledValue()

	if math.IsNaN(v) {
		return 0, meters.ErrNaN
	}

	return v, nil
}

// QueryPoint executes a single query operation for model/block/point on the bus
func (d *SunSpec) QueryPointAny(client modbus.Client, modelID, blockID int, pointID string) (block sunspec.Block, point sunspec.Point, err error) {
	if d.notInitialized() {
		return block, point, errors.New("sunspec: not initialized")
	}

	for _, model := range d.models {
		if sunspec.ModelId(modelID) != model.Id() {
			continue
		}

		// read zero block to initialize scale factors
		if blockID > 0 {
			block, err := model.Block(0)
			if err == nil {
				err = block.Read()
			}
			if err != nil {
				return block, point, err
			}
		}

		block, err := model.Block(blockID)
		if err == nil {
			err = block.Read()
		}

		var point sunspec.Point
		if err == nil {
			point, err = block.Point(pointID)
		}

		return block, point, err
	}

	return block, point, fmt.Errorf("sunspec: %d:%d:%s not found", modelID, blockID, pointID)
}

// QueryPoint executes a single query operation for model/block/point on the bus.
// The result is returned as-is, i.e. not scaled.
func (d *SunSpec) QueryPoint(client modbus.Client, modelID, blockID int, pointID string) (float64, error) {
	block, point, err := d.QueryPointAny(client, modelID, blockID, pointID)
	if err != nil {
		return 0, err
	}

	return d.convertPoint(block, point)
}

func makeResult(v float64, m meters.Measurement) meters.MeasurementResult {
	// apply scale factor for energy
	if div, ok := dividerMap[m]; ok {
		v /= div
	}

	res := meters.MeasurementResult{
		Measurement: m,
		Value:       v,
		Timestamp:   time.Now(),
	}

	return res
}

// QueryOp queries all models and blocks until measurement is found.
// The result is scaled as defined in the divider map.
func (d *SunSpec) QueryOp(client modbus.Client, measurement meters.Measurement) (res meters.MeasurementResult, err error) {
	if d.notInitialized() {
		return res, errors.New("sunspec: not initialized")
	}

	for _, model := range d.models {
		for modelID, blockMap := range modelMap {
			if modelID != model.Id() {
				continue
			}

			for blockID, pointMap := range blockMap {
				if blockID >= model.Blocks() {
					continue
				}

				for pointID, m := range pointMap {
					if m == measurement {
						v, err := d.QueryPoint(client, int(modelID), blockID, pointID)

						var mr meters.MeasurementResult
						if err == nil {
							mr = makeResult(v, measurement)
						}

						return mr, err
					}
				}
			}
		}
	}

	return meters.MeasurementResult{}, fmt.Errorf("sunspec: %s not found", measurement)
}

// Query is called by the handler after preparing the bus by setting the device id and waiting for rate limit
// The results are scaled as defined in the divider map.
func (d *SunSpec) Query(client modbus.Client) (res []meters.MeasurementResult, err error) {
	if d.notInitialized() {
		return res, errors.New("sunspec: not initialized")
	}

	for _, model := range d.models {
		for modelID, blockMap := range modelMap {
			if modelID != model.Id() {
				continue
			}

			// sort blocks so block 0 is always read first
			sortedBlocks := make([]int, 0, len(blockMap))
			for k := range blockMap {
				sortedBlocks = append(sortedBlocks, k)
			}
			sort.Ints(sortedBlocks)

			// always add zero block
			if sortedBlocks[0] != 0 {
				sortedBlocks = append([]int{0}, sortedBlocks...)
			}

			for blockID := range sortedBlocks {
				if blockID >= model.Blocks() {
					continue
				}

				pointMap := blockMap[blockID]
				block := model.MustBlock(blockID)

				if err := block.Read(); err != nil {
					return res, err
				}

				for pointID, m := range pointMap {
					point := block.MustPoint(pointID)

					if v, err := d.convertPoint(block, point); err == nil {
						res = append(res, makeResult(v, m))
					}
				}
			}
		}
	}

	return res, nil
}
