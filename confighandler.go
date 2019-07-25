package mbmd

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/volkszaehler/mbmd/meters"
	"github.com/volkszaehler/mbmd/meters/rs485"
	"github.com/volkszaehler/mbmd/meters/sunspec"
)

// DeviceConfigHandler creates map of meter managers from given configuration
type DeviceConfigHandler struct {
	defaultDevice string
	Managers      map[string]meters.Manager
}

// NewDeviceConfigHandler creates a configuration handler
func NewDeviceConfigHandler(defaultDevice string) *DeviceConfigHandler {
	conf := &DeviceConfigHandler{
		defaultDevice: defaultDevice,
		Managers:      make(map[string]meters.Manager),
	}
	return conf
}

// createConnection parses adapter string to create TCP or RTU connection
func createConnection(device string) (res meters.Connection) {
	if device == "mock" {
		res = meters.NewMock(device) // mocked connection
	} else if tcp, _ := regexp.MatchString(":[0-9]+$", device); tcp {
		res = meters.NewTCP(device) // tcp connection
	} else {
		res = meters.NewRTU(device, 1) // serial connection
	}
	return res
}

// CreateDevice creates new device and adds it to the
func (conf *DeviceConfigHandler) CreateDevice(deviceDef string) {
	deviceSplit := strings.Split(deviceDef, "@")
	if len(deviceSplit) == 0 || len(deviceSplit) > 2 {
		log.Fatalf("Cannot parse connect string %s. See -h for help.", deviceDef)
	}

	meterDef := deviceSplit[0]
	connSpec := conf.defaultDevice
	if len(deviceSplit) == 2 {
		connSpec = deviceSplit[1]
	}

	if connSpec == "" {
		log.Fatalf("Cannot parse connect string- missing physical device or connection for %s. See -h for help.", deviceDef)
	}

	manager, ok := conf.Managers[connSpec]
	if !ok {
		conn := createConnection(connSpec)
		manager = meters.NewManager(conn)
		conf.Managers[connSpec] = manager
	}

	meterSplit := strings.Split(meterDef, ":")
	if len(meterSplit) != 2 {
		log.Fatalf("Cannot parse device definition: %s. See -h for help.", meterDef)
	}

	meterType, devID := meterSplit[0], meterSplit[1]
	if len(strings.TrimSpace(meterType)) == 0 {
		log.Fatalf("Cannot parse device definition- meter type empty: %s. See -h for help.", meterDef)
	}

	id, err := strconv.Atoi(devID)
	if err != nil {
		log.Fatalf("Error parsing device id %s: %v. See -h for help.", devID, err)
	}

	var meter meters.Device
	if _, ok := manager.Conn.(*meters.TCP); ok {
		meter = sunspec.NewDevice()
	} else {
		meterType = strings.ToUpper(meterType)
		meter, err = rs485.NewDevice(meterType)
		if err != nil {
			log.Fatalf("Error creating device %s: %v. See -h for help.", meterDef, err)
		}
	}

	manager.Add(uint8(id), meter)
}
