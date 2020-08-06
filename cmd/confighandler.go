package cmd

import (
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/volkszaehler/mbmd/meters"
	"github.com/volkszaehler/mbmd/meters/rs485"
	"github.com/volkszaehler/mbmd/meters/sunspec"
)

// Config describes the entire configuration
type Config struct {
	API      string
	Rate     time.Duration
	Mqtt     MqttConfig
	Influx   InfluxConfig
	Adapters []AdapterConfig
	Devices  []DeviceConfig
	Other    map[string]interface{} `mapstructure:",remain"`
}

// MqttConfig describes the mqtt broker configuration
type MqttConfig struct {
	Broker   string
	Topic    string
	User     string
	Password string
	ClientID string
	Qos      int
	Homie    string
}

// InfluxConfig describes the InfluxDB configuration
type InfluxConfig struct {
	URL          string
	Database     string
	Measurement  string
	Organization string
	Token        string
	User         string
	Password     string
}

// AdapterConfig describes device communication parameters
type AdapterConfig struct {
	Device   string
	RTU      bool
	Baudrate int
	Comset   string
}

// DeviceConfig describes a single device's configuration
type DeviceConfig struct {
	Type      string
	ID        uint8
	SubDevice int
	Name      string
	Adapter   string
}

// DeviceConfigHandler creates map of meter managers from given configuration
type DeviceConfigHandler struct {
	DefaultDevice string
	Managers      map[string]*meters.Manager
}

// NewDeviceConfigHandler creates a configuration handler
func NewDeviceConfigHandler() *DeviceConfigHandler {
	conf := &DeviceConfigHandler{
		Managers: make(map[string]*meters.Manager),
	}
	return conf
}

// createConnection parses adapter string to create TCP or RTU connection
func createConnection(device string, rtu bool, baudrate int, comset string) (res meters.Connection) {
	if device == "mock" {
		res = meters.NewMock(device) // mocked connection
	} else if tcp, _ := regexp.MatchString(":[0-9]+$", device); tcp {
		if rtu {
			// special case: RTU over TCP
			log.Printf("config: creating RTU over TCP connection for %s", device)
			res = meters.NewRTUOverTCP(device) // tcp connection
		} else {
			log.Printf("config: creating TCP connection for %s", device)
			res = meters.NewTCP(device) // tcp connection
		}
	} else {
		log.Printf("config: creating RTU connection for %s (%dbaud, %s)", device, baudrate, comset)
		if baudrate == 0 || comset == "" {
			log.Fatal("Missing comset configuration. See -h for help.")
		}
		if _, err := os.Stat(device); err != nil {
			log.Fatal(err)
		}
		res = meters.NewRTU(device, baudrate, comset) // serial connection
	}
	return res
}

// ConnectionManager returns connection manager from cache or creates new connection wrapped by manager
func (conf *DeviceConfigHandler) ConnectionManager(connSpec string, rtu bool, baudrate int, comset string) *meters.Manager {
	manager, ok := conf.Managers[connSpec]
	if !ok {
		conn := createConnection(connSpec, rtu, baudrate, comset)
		manager = meters.NewManager(conn)
		conf.Managers[connSpec] = manager
	}

	return manager
}

func (conf *DeviceConfigHandler) createDeviceForManager(
	manager *meters.Manager,
	meterType string,
	subdevice int,
) meters.Device {
	var meter meters.Device
	meterType = strings.ToUpper(meterType)

	var isSunspec bool
	sunspecTypes := []string{"FRONIUS", "KOSTAL", "KACO", "SE", "SMA", "SOLAREDGE", "STECA", "SUNS", "SUNSPEC"}
	for _, t := range sunspecTypes {
		if t == meterType {
			isSunspec = true
			break
		}
	}

	sort.SearchStrings(sunspecTypes, meterType)
	if isSunspec {
		meter = sunspec.NewDevice(meterType, subdevice)
	} else {
		if subdevice > 0 {
			log.Fatalf("Invalid subdevice number for device %s: %d", meterType, subdevice)
		}

		var err error
		meter, err = rs485.NewDevice(meterType)
		if err != nil {
			log.Fatalf("Error creating device %s: %v.", meterType, err)
		}
	}

	return meter
}

// CreateDevice creates new device and adds it to the connection manager
func (conf *DeviceConfigHandler) CreateDevice(devConf DeviceConfig) {
	if devConf.Adapter == "" {
		// find default adapter
		if len(conf.Managers) == 1 {
			for a := range conf.Managers {
				log.Printf("config: using default adapter %s for device %v", a, devConf)
				devConf.Adapter = a
			}
		} else {
			log.Fatalf("Missing adapter configuration for device %v", devConf)
		}
	}

	manager, ok := conf.Managers[devConf.Adapter]
	if !ok {
		log.Fatalf("Missing adapter configuration for device %v", devConf)
	}
	meter := conf.createDeviceForManager(manager, devConf.Type, devConf.SubDevice)

	if err := manager.Add(devConf.ID, meter); err != nil {
		log.Fatalf("Error adding device %v: %v.", devConf, err)
	}
}

// CreateDeviceFromSpec creates new device from specification string and adds
// it to the connection manager
func (conf *DeviceConfigHandler) CreateDeviceFromSpec(deviceDef string) {
	deviceSplit := strings.Split(deviceDef, "@")
	if len(deviceSplit) == 0 || len(deviceSplit) > 2 {
		log.Fatalf("Cannot parse connect string %s. See -h for help.", deviceDef)
	}

	meterDef := deviceSplit[0]
	connSpec := conf.DefaultDevice
	if len(deviceSplit) == 2 {
		connSpec = deviceSplit[1]
	}

	if connSpec == "" {
		log.Fatalf("Cannot parse connect string- missing physical device or connection for %s. See -h for help.", deviceDef)
	}

	meterSplit := strings.Split(meterDef, ":")
	if len(meterSplit) != 2 {
		log.Fatalf("Cannot parse device definition: %s. See -h for help.", meterDef)
	}

	meterType, devID := meterSplit[0], meterSplit[1]
	if len(strings.TrimSpace(meterType)) == 0 {
		log.Fatalf("Cannot parse device definition- meter type empty: %s. See -h for help.", meterDef)
	}

	var subdevice int
	devIDSplit := strings.SplitN(devID, ".", 2)
	if len(devIDSplit) == 2 {
		var err error
		subdevice, err = strconv.Atoi(devIDSplit[1])
		if err != nil {
			log.Fatalf("Error parsing device id %s: %v. See -h for help.", devID, err)
		}
	} else if len(devIDSplit) > 2 {
		log.Fatalf("Error parsing device id %s. See -h for help.", devID)
	}

	id, err := strconv.Atoi(devIDSplit[0])
	if err != nil {
		log.Fatalf("Error parsing device id %s: %v. See -h for help.", devID, err)
	}

	// If this is an RTU over TCP device, a default RTU over TCP should already
	// have been created of the --rtu flag was specified. We'll not re-check this here.
	manager := conf.ConnectionManager(connSpec, false, 0, "")

	meter := conf.createDeviceForManager(manager, meterType, subdevice)
	if err := manager.Add(uint8(id), meter); err != nil {
		log.Fatalf("Error adding device %s: %v. See -h for help.", meterDef, err)
	}
}
