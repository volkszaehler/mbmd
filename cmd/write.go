package cmd

import (
	"encoding/binary"
	"math"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// writeCmd represents the write command
var writeCmd = &cobra.Command{
	Use:   "write [flags] " + strings.Join([]string{"register", "length", "value"}, " "),
	Short: "Write register (EXPERIMENTAL)",
	Long: `Write writes a single register (holding, coil). Write will ignore the
config file and requires adapter configuration using command line.`,
	Args: cobra.ExactArgs(3),
	Run:  write,
}

func init() {
	rootCmd.AddCommand(writeCmd)

	writeCmd.PersistentFlags().StringP(
		"device", "d",
		"1",
		"MODBUS device ID to query. Only single device allowed.",
	)
	writeCmd.PersistentFlags().StringP(
		"type", "t",
		"holding",
		"Register type to write: holding|coil",
	)
	writeCmd.PersistentFlags().StringP(
		"encoding", "e",
		"int",
		"Data encoding: bit|int|uint|hex|float|string",
	)
}

func parseDecimalString(
	value string, length int, base int,
	parse func(string, int, int) (uint64, error),
) []byte {
	bytes := 2 * length
	b := make([]byte, bytes)

	u, err := parse(value, base, 8*bytes)
	if err != nil {
		log.Fatal(err)
	}

	switch length {
	case 1:
		binary.BigEndian.PutUint16(b, uint16(u))
	case 2:
		binary.BigEndian.PutUint32(b, uint32(u))
	case 4:
		binary.BigEndian.PutUint64(b, uint64(u))
	default:
		log.Fatal("Unsupported length")
	}

	return b
}

func parseFloatString(value string, length int) []byte {
	bytes := 2 * length
	b := make([]byte, bytes)

	f, err := strconv.ParseFloat(value, 8*bytes)
	if err != nil {
		log.Fatal(err)
	}

	switch length {
	case 2:
		binary.BigEndian.PutUint32(b, math.Float32bits(float32(f)))
	case 4:
		binary.BigEndian.PutUint64(b, math.Float64bits(f))
	default:
		log.Fatal("Unsupported length")
	}

	return b
}

func parseInt(s string, base int, bitSize int) (uint64, error) {
	i, err := strconv.ParseInt(s, base, bitSize)
	return uint64(i), err
}

// encodeCoil accepts 0 or 1 which is converted into 0x0000 or 0xFF00
func encodeCoil(value string, length int) (b []byte) {
	if length != 1 {
		log.Fatal("Invalid length")
	}

	u, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		log.Fatal(err)
	}

	switch u {
	case 0:
		b = []byte{0, 0}
	case 1, 0xFF00:
		b = []byte{0xFF, 0}
	default:
		log.Fatal("Invalid value")
	}

	return b
}

func encode(value string, length int, encoding string) []byte {
	switch strings.ToLower(encoding) {
	case "bit":
		return encodeCoil(value, length)
	case "int":
		return parseDecimalString(value, length, 10, parseInt)
	case "uint":
		return parseDecimalString(value, length, 10, strconv.ParseUint)
	case "hex":
		value = strings.TrimPrefix(value, "0x")
		return parseDecimalString(value, length, 16, strconv.ParseUint)
	case "string":
		if len(value) > 2*length {
			log.Fatal("Length exceeded")
		}
		// pad trailing zeros
		for len(value) < 2*length {
			value += "\000"
		}
		return []byte(value)
	case "float":
		return parseFloatString(value, length)
	}

	log.Fatal("Invalid encoding")
	return []byte{}
}

func write(cmd *cobra.Command, args []string) {
	// log only fatal messages
	configureLogger(viper.GetBool("verbose"), 0)

	// arguments
	register, length := parseArgs(args)
	value := args[2]

	// flags
	dev, _ := cmd.PersistentFlags().GetString("device")
	typ, _ := cmd.PersistentFlags().GetString("type")
	encoding, _ := cmd.PersistentFlags().GetString("encoding")
	validateFlags(typ, encoding)

	// parse modbus settings
	conn, client := modbusClient()
	conn.Slave(deviceIDFromSpec(dev))

	// encode argument to buffer
	b := encode(value, length, encoding)

	// execute write
	var err error
	switch strings.ToLower(typ) {
	case "holding":
		_, err = client.WriteMultipleRegisters(uint16(register), uint16(length), b)
	case "coil":
		_, err = client.WriteSingleCoil(uint16(register), binary.BigEndian.Uint16(b))
	default:
		log.Fatalf("Invalid write type '%s'", typ)
	}

	if err != nil {
		log.Fatal(err)
	}
}
