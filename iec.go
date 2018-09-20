package sdm630

import "log"

var iec = map[string]string{
	"AmpLocPhsA":       "L1 Current (A)",
	"AmpLocPhsB":       "L2 Current (A)",
	"AmpLocPhsC":       "L3 Current (A)",
	"AngLocPhsA":       "L1 Cosphi",
	"AngLocPhsB":       "L2 Cosphi",
	"AngLocPhsC":       "L3 Cosphi",
	"Freq":             "Frequency of supply voltages",
	"ThdVol":           "Average voltage to neutral THD (%)",
	"ThdVolPhsA":       "L1 Voltage to neutral THD (%)",
	"ThdVolPhsB":       "L2 Voltage to neutral THD (%)",
	"ThdVolPhsC":       "L3 Voltage to neutral THD (%)",
	"TotkWhExport":     "Total Export (kWh)",
	"TotkWhExportPhsA": "L1 Export (kWh)",
	"TotkWhExportPhsB": "L2 Export (kWh)",
	"TotkWhExportPhsC": "L3 Export (kWh)",
	"TotkWhImport":     "Total Import (kWh)",
	"TotkWhImportPhsA": "L1 Import (kWh)",
	"TotkWhImportPhsB": "L2 Import (kWh)",
	"TotkWhImportPhsC": "L3 Import (kWh)",
	"VolLocPhsA":       "L1 Voltage (V)",
	"VolLocPhsB":       "L2 Voltage (V)",
	"VolLocPhsC":       "L3 Voltage (V)",
	"WLocPhsA":         "L1 Power (W)",
	"WLocPhsB":         "L2 Power (W)",
	"WLocPhsC":         "L3 Power (W)",
}

func GetIecDescription(key string) string {
	description, ok := iec[key]
	if !ok {
		log.Fatalf("Undefined IEC code %s", key)
	}
	return description
}
