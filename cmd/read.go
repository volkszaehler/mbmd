package cmd

import (
	"encoding/binary"
	"fmt"
	golog "log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/grid-x/modbus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/volkszaehler/mbmd/meters"
	"github.com/volkszaehler/mbmd/meters/rs485"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read [flags] " + strings.Join([]string{"register", "length"}, " "),
	Short: "Read register (EXPERIMENTAL)",
	Long: `Read reads a single register (input, holding, coil, discrete input)
and will return it according to defined format. Read will ignore the config file
and requires adapter configuration using command line.`,
	Args: cobra.ExactArgs(2),
	Run:  read,
}

func init() {
	rootCmd.AddCommand(readCmd)

	readCmd.PersistentFlags().StringP(
		"device", "d",
		"1",
		"MODBUS device ID to query. Only single device allowed.",
	)
	readCmd.PersistentFlags().StringP(
		"type", "t",
		"holding",
		"Register type to read: holding|input|coil|discrete",
	)
	readCmd.PersistentFlags().StringP(
		"encoding", "e",
		"int",
		"Data encoding: bit|int|uint|int32s|uint32s|hex|float|floats|string",
	)
}

func modbusClient() (meters.Connection, modbus.Client) {
	// create connection
	adapter := viper.GetString("adapter")
	if adapter == "" {
		log.Fatal("Missing adapter configuration")
	}

	// connection
	conn := createConnection(adapter, viper.GetBool("rtu"), viper.GetInt("baudrate"), viper.GetString("comset"))
	client := conn.ModbusClient()

	// raw log
	if viper.GetBool("raw") {
		conn.Logger(golog.New(os.Stderr, "", golog.LstdFlags))
	}

	return conn, client
}

// deviceIDFromSpec parses a device specification and extracts the slave id
func deviceIDFromSpec(meterDef string) uint8 {
	devID, err := strconv.Atoi(meterDef)
	if err != nil {
		meterSplit := strings.Split(meterDef, ":")
		if len(meterSplit) != 2 {
			log.Fatalf("Cannot parse device definition: %s. See -h for help.", meterDef)
		}

		devID, err = strconv.Atoi(meterSplit[1])
		if err != nil {
			log.Fatalf("Error parsing device id %s: %v. See -h for help.", meterSplit[1], err)
		}
	}

	if devID <= 0 {
		log.Fatal("Missing device id")
	}

	return uint8(devID)
}

func bytes2uint(b []byte, length int) uint64 {
	switch length {
	case 1:
		return uint64(binary.BigEndian.Uint16(b))
	case 2:
		return uint64(binary.BigEndian.Uint32(b))
	case 4:
		return binary.BigEndian.Uint64(b)
	}

	log.Fatal("Invalid length for (u)int encoding", "float")
	return 0
}

func bytes2float(b []byte, length int) float64 {
	switch length {
	case 2:
		return float64(math.Float32frombits(binary.BigEndian.Uint32(b)))
	case 4:
		return math.Float64frombits(binary.BigEndian.Uint64(b))
	}

	log.Fatal("Invalid length for float encoding")
	return 0
}

func decodeCoils(b []byte, length int) (s string) {
	idx := 0
BYTES:
	for _, byt := range b {
		for bit := 0; bit < 8; bit++ {
			if byt&(1<<bit) > 0 {
				s += "01 "
			} else {
				s += "00 "
			}

			idx++
			if idx >= length {
				break BYTES
			}
		}
	}
	return strings.TrimRight(s, " ")
}

func decode(b []byte, length int, encoding string) string {
	switch strings.ToLower(encoding) {
	case "bit":
		return decodeCoils(b, length)
	case "int":
		u := bytes2uint(b, length)
		return strconv.FormatInt(int64(u), 10)
	case "int32swapped", "int32s":
		if length != 2 {
			log.Fatal("Invalid length for int32(swapped) encoding")
		}
		u := rs485.BigEndianUint32Swapped(b)
		return strconv.FormatInt(int64(int32(u)), 10)
	case "uint":
		u := bytes2uint(b, length)
		return strconv.FormatUint(u, 10)
	case "uint32swapped", "uint32s":
		if length != 2 {
			log.Fatal("Invalid length for uint32(swapped) encoding")
		}
		u := rs485.BigEndianUint32Swapped(b)
		return strconv.FormatUint(uint64(uint32(u)), 10)
	case "hex":
		return fmt.Sprintf("%02x", b)
	case "string":
		return string(b)
	case "float":
		f := bytes2float(b, length)
		return fmt.Sprintf("%f", f)
	case "floatswapped", "floats":
		if length != 2 {
			log.Fatal("Invalid length for float(swapped) encoding")
		}
		f := rs485.RTUIeee754ToFloat64Swapped(b)
		return fmt.Sprintf("%f", f)
	}

	log.Fatal("Invalid encoding")
	return ""
}

func readFunction(client modbus.Client, typ string) (f func(address, quantity uint16) (results []byte, err error)) {
	switch strings.ToLower(typ) {
	case "holding":
		f = client.ReadHoldingRegisters
	case "input":
		f = client.ReadInputRegisters
	case "coil":
		f = client.ReadCoils
	case "discrete":
		f = client.ReadDiscreteInputs
	default:
		log.Fatalf("Invalid read type '%s'", typ)
	}

	return f
}

func parseArgs(args []string) (int, int) {
	register, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatal("Invalid register number")
	}

	length, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatal("Invalid length")
	}

	return register, length
}

func validateFlags(typ, encoding string) {
	if (typ == "discrete" || typ == "coil") && (encoding != "bit") {
		log.Fatal("Invalid encoding for register type")
	}
}

func read(cmd *cobra.Command, args []string) {
	// log only fatal messages
	configureLogger(viper.GetBool("verbose"), 0)

	// arguments
	register, length := parseArgs(args)

	// flags
	dev, _ := cmd.PersistentFlags().GetString("device")
	typ, _ := cmd.PersistentFlags().GetString("type")
	encoding, _ := cmd.PersistentFlags().GetString("encoding")
	validateFlags(typ, encoding)

	// parse modbus settings
	conn, client := modbusClient()
	conn.Slave(deviceIDFromSpec(dev))

	// raw log
	if viper.GetBool("raw") {
		conn.Logger(golog.New(os.Stderr, "", 0))
	}

	// execute read
	f := readFunction(client, typ)
	b, err := f(uint16(register), uint16(length))
	if err != nil {
		log.Fatal(err)
	}

	// result
	println(decode(b, length, encoding))
}
