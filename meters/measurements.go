package meters

type Measurement int

//go:generate go run golang.org/x/tools/cmd/stringer -type=Measurement
const (
	// split signals block operations that need be split into individual results
	Split Measurement = iota

	// parameters
	Frequency

	Current
	CurrentL1
	CurrentL2
	CurrentL3

	// phases and sums
	Voltage
	VoltageL1
	VoltageL2
	VoltageL3

	Power
	PowerL1
	PowerL2
	PowerL3

	ImportPower
	ImportPowerL1
	ImportPowerL2
	ImportPowerL3

	ExportPower
	ExportPowerL1
	ExportPowerL2
	ExportPowerL3

	ActivePower
	ActivePowerL1
	ActivePowerL2
	ActivePowerL3

	ReactivePower
	ReactivePowerL1
	ReactivePowerL2
	ReactivePowerL3

	ApparentPower
	ApparentPowerL1
	ApparentPowerL2
	ApparentPowerL3

	Cosphi
	CosphiL1
	CosphiL2
	CosphiL3

	THD
	THDL1
	THDL2
	THDL3

	// energy
	Net
	NetL1
	NetL2
	NetL3

	Active
	ActiveNet
	ActiveNetL1
	ActiveNetL2
	ActiveNetL3

	Reactive
	ReactiveNet
	ReactiveNetL1
	ReactiveNetL2
	ReactiveNetL3

	Import
	ImportL1
	ImportL2
	ImportL3

	Export
	ExportL1
	ExportL2
	ExportL3

	ActiveImportT1
	ActiveImportT2

	ReactiveImportT1
	ReactiveImportT2

	ActiveExportT1
	ActiveExportT2

	ReactiveExportT1
	ReactiveExportT2

	// DC
	DCCurrent
	DCVoltage
	DCPower
	HeatSinkTemp
)

var iec = map[Measurement]string{
	CurrentL1:    "L1 Current (A)",
	CurrentL2:    "L2 Current (A)",
	CurrentL3:    "L3 Current (A)",
	CosphiL1:     "L1 Cosphi",
	CosphiL2:     "L2 Cosphi",
	CosphiL3:     "L3 Cosphi",
	Frequency:    "Frequency (Hz)",
	THD:          "Average voltage to neutral THD (%)",
	THDL1:        "L1 Voltage to neutral THD (%)",
	THDL2:        "L2 Voltage to neutral THD (%)",
	THDL3:        "L3 Voltage to neutral THD (%)",
	Export:       "Total Export (kWh)",
	ExportL1:     "L1 Export (kWh)",
	ExportL2:     "L2 Export (kWh)",
	ExportL3:     "L3 Export (kWh)",
	Import:       "Total Import (kWh)",
	ImportL1:     "L1 Import (kWh)",
	ImportL2:     "L2 Import (kWh)",
	ImportL3:     "L3 Import (kWh)",
	VoltageL1:    "L1 Voltage (V)",
	VoltageL2:    "L2 Voltage (V)",
	VoltageL3:    "L3 Voltage (V)",
	PowerL1:      "L1 Power (W)",
	PowerL2:      "L2 Power (W)",
	PowerL3:      "L3 Power (W)",
	DCCurrent:    "DC Current (A)",
	DCVoltage:    "DC Voltage (V)",
	DCPower:      "DC Power (W)",
	HeatSinkTemp: "Heat Sink Temperature (°C)",
}

func (m *Measurement) Description() string {
	if description, ok := iec[*m]; ok {
		return description
	}
	return m.String()
}
