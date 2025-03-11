package sunspec

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"

	"github.com/grid-x/modbus"
	"github.com/volkszaehler/mbmd/meters"
)

const (
	// Battery register addresses (in decimal, exactly as in the Perl script)
	batteryRegisterBase1 = 57600  // 0xE100 - Battery 1
	batteryRegisterBase2 = 57856  // 0xE200 - Battery 2
)

// SolarEdgeBattery implements the meters.Device interface for SolarEdge battery
type SolarEdgeBattery struct {
	SunSpec
}

// NewSolarEdgeBatteryDevice creates a SolarEdge battery device
func NewSolarEdgeBatteryDevice(subdevice int) *SolarEdgeBattery {
	dev := &SolarEdgeBattery{
		SunSpec: SunSpec{
			subdevice: subdevice,
			descriptor: meters.DeviceDescriptor{
				Type:         "SE-BAT",
				Manufacturer: "SolarEdge",
				Model:        "Home Battery",
				SubDevice:    subdevice,
			},
		},
	}

	return dev
}

// Initialize implements the meters.Device interface
func (d *SolarEdgeBattery) Initialize(client modbus.Client) error {
	// Nothing to initialize for SolarEdge battery
	return nil
}

// getBaseAddress returns the correct base register address for the battery
func (d *SolarEdgeBattery) getBaseAddress() uint16 {
	if d.subdevice == 0 {
		return batteryRegisterBase1
	}
	return batteryRegisterBase2
}

// readSEFloat reads a SolarEdge float value from the given data at the specified offset
func readSEFloat(data []byte, offset int) (float64, error) {
	if offset+3 >= len(data) {
		return 0, fmt.Errorf("not enough data")
	}

	// Extract as SEFLOAT (little endian)
	rawBytes := []byte{data[offset+1], data[offset], data[offset+3], data[offset+2]}
	bits := binary.LittleEndian.Uint32(rawBytes)
	value := float64(math.Float32frombits(bits))

	// Skip NaN and Infinity values
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, fmt.Errorf("invalid value")
	}

	return value, nil
}

// Probe implements the meters.Device interface
func (d *SolarEdgeBattery) Probe(client modbus.Client) (meters.MeasurementResult, error) {
	// Try to read the State of Energy register at offset 0x84 (132)
	baseAddr := d.getBaseAddress()

	// Read SOE register directly
	b, err := client.ReadHoldingRegisters(baseAddr + 0x84, 2)
	if err != nil {
		return meters.MeasurementResult{}, err
	}

	if len(b) < 4 {
		return meters.MeasurementResult{}, fmt.Errorf("insufficient data received")
	}

	// Extract the SOE value
	soc, err := readSEFloat(b, 0)
	if err != nil {
		return meters.MeasurementResult{}, err
	}

	return meters.MeasurementResult{
		Measurement: meters.BatterySOC,
		Value:       soc,
		Timestamp:   time.Now(),
	}, nil
}

// Query implements the meters.Device interface
func (d *SolarEdgeBattery) Query(client modbus.Client) ([]meters.MeasurementResult, error) {
	res := make([]meters.MeasurementResult, 0, 9)
	timestamp := time.Now()
	baseAddr := d.getBaseAddress()

	// Read registers in chunks to minimize Modbus traffic

	// First chunk: Read the rated energy and power capabilities
	ratedEnergyData, err := client.ReadHoldingRegisters(baseAddr + 0x42, 8)
	if err == nil && len(ratedEnergyData) >= 16 {
		// Rated Energy (0x42)
		if value, err := readSEFloat(ratedEnergyData, 0); err == nil {
			res = append(res, meters.MeasurementResult{
				Measurement: meters.BatteryRatedCapacity,
				Value:       value * 0.001, // Convert Wh to kWh
				Timestamp:   timestamp,
			})
		}

		// Max Charge Power (0x44)
		if value, err := readSEFloat(ratedEnergyData, 4); err == nil {
			res = append(res, meters.MeasurementResult{
				Measurement: meters.BatteryMaxPowerCharge,
				Value:       value,
				Timestamp:   timestamp,
			})
		}

		// Max Discharge Power (0x46)
		if value, err := readSEFloat(ratedEnergyData, 8); err == nil {
			res = append(res, meters.MeasurementResult{
				Measurement: meters.BatteryMaxPowerDischarge,
				Value:       value,
				Timestamp:   timestamp,
			})
		}
	}

	// Second chunk: Read the battery metrics (temperature through voltage, current, power)
	batteryData, err := client.ReadHoldingRegisters(baseAddr + 0x6C, 10)
	if err == nil && len(batteryData) >= 20 {
		// Temperature (0x6C)
		if value, err := readSEFloat(batteryData, 0); err == nil {
			res = append(res, meters.MeasurementResult{
				Measurement: meters.BatteryTemperature,
				Value:       value,
				Timestamp:   timestamp,
			})
		}

		// Voltage (0x70)
		if value, err := readSEFloat(batteryData, 8); err == nil {
			res = append(res, meters.MeasurementResult{
				Measurement: meters.BatteryVoltage,
				Value:       value,
				Timestamp:   timestamp,
			})
		}

		// Current (0x72)
		if value, err := readSEFloat(batteryData, 12); err == nil {
			res = append(res, meters.MeasurementResult{
				Measurement: meters.BatteryCurrent,
				Value:       value,
				Timestamp:   timestamp,
			})
		}

		// Power (0x74) - Negated to match the Perl script
		if value, err := readSEFloat(batteryData, 16); err == nil {
			res = append(res, meters.MeasurementResult{
				Measurement: meters.BatteryPower,
				Value:       -value, // Negate the value like the Perl script does
				Timestamp:   timestamp,
			})
		}
	}

	// Third chunk: Read the energy, SOH, SOC registers
	energyData, err := client.ReadHoldingRegisters(baseAddr + 0x80, 6)
	if err == nil && len(energyData) >= 12 {
		// Available Energy (0x80)
		if value, err := readSEFloat(energyData, 0); err == nil {
			res = append(res, meters.MeasurementResult{
				Measurement: meters.BatteryAvailableEnergy,
				Value:       value * 0.001, // Convert Wh to kWh
				Timestamp:   timestamp,
			})
		}

		// State of Health (0x82)
		if value, err := readSEFloat(energyData, 4); err == nil {
			res = append(res, meters.MeasurementResult{
				Measurement: meters.BatterySOH,
				Value:       value,
				Timestamp:   timestamp,
			})
		}

		// State of Energy (0x84)
		if value, err := readSEFloat(energyData, 8); err == nil {
			res = append(res, meters.MeasurementResult{
				Measurement: meters.BatterySOC,
				Value:       value,
				Timestamp:   timestamp,
			})
		}
	}

	return res, nil
}
