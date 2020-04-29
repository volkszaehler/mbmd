package cmd

import (
	"fmt"
	golog "log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/grid-x/modbus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/volkszaehler/mbmd/meters"
	quirks "github.com/volkszaehler/mbmd/meters/sunspec"

	sunspec "github.com/andig/gosunspec"
	bus "github.com/andig/gosunspec/modbus"
	_ "github.com/andig/gosunspec/models" // import models
	"github.com/andig/gosunspec/smdx"
)

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect SunSpec device models and implemented values",
	Long: `Inspect will iterate across the SunSpec model definition for the specified device and print details about defined models and data points.
Devices are expected to be specified on command line- config file is being ignored.
Limited to SunSpec TCP devices (EXPERIMENTAL).`,
	Run: inspect,
}

func init() {
	rootCmd.AddCommand(inspectCmd)

	inspectCmd.PersistentFlags().StringSliceP(
		"devices", "d",
		[]string{},
		`MODBUS device type and ID to query, multiple devices separated by comma or by repeating the flag.
  Example: -d SDM:1,SDM:2 -d DZG:1.
Valid types are:`+meterHelp()+`
To use an adapter different from default, append RTU device or TCP address separated by @.
If the adapter is a TCP connection (identified by :port), the device type (SUNS) is ignored and
any type is considered valid.
  Example: -d SDM:1@/dev/USB11 -d SMA:126@localhost:502`,
	)
}

func pf(format string, v ...interface{}) {
	format = strings.TrimSuffix(format, "\n") + "\n"
	fmt.Printf(format, v...)
}

func scanSunspec(client modbus.Client) {
	in, err := bus.Open(client)
	if err != nil && in == nil {
		log.Fatal(err)
	} else if err != nil {
		log.Printf("warning: device opened with partial result: %v", err) // log error but continue
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	in.Do(func(d sunspec.Device) {
		d.Do(func(m sunspec.Model) {
			pf("--------- Model %d %s ---------", m.Id(), modelName(m))
			// if m.Id() != 113 && m.Id() != 103 {
			// 	return
			// }

			blocknum := 0
			m.Do(func(b sunspec.Block) {
				if blocknum > 0 {
					fmt.Fprintf(tw, "-- Block %d --\n", blocknum)
				}
				blocknum++

				err = b.Read()
				if err != nil {
					log.Printf("skipping due to read error: %v", err)
					return
				}

				b.Do(func(p sunspec.Point) {
					t := p.Type()[0:3]
					v := p.Value()
					if p.NotImplemented() {
						v = "n/a"
					} else if t == "int" || t == "uin" || t == "acc" {
						// for time being, always to this
						quirks.FixKostal(p)

						v = p.ScaledValue()
						v = fmt.Sprintf("%.2f", v)
					}

					vs := fmt.Sprintf("%17v", v)
					fmt.Fprintf(tw, "%s\t%s\t   %s\n", p.Id(), vs, p.Type())
				})
			})

			tw.Flush()
		})
	})
}

func modelName(m sunspec.Model) string {
	model := smdx.GetModel(uint16(m.Id()))
	if model == nil {
		return ""
	}
	return model.Name
}

func printModel(m *smdx.ModelElement) {
	pf("-- Definition --")
	// pf("----")
	// pf("Model:  %d - %s", m.Id, m.Name)
	pf("Length: %d (0x%02x words, 0x%02x bytes)", m.Length, m.Length, 2*m.Length)
	pf("Blocks: %d", len(m.Blocks))

	for i, b := range m.Blocks {
		pf("-- block #%d - %s", i, b.Name)
		pf("Length: %d", b.Length)

		for _, p := range b.Points {
			u := ""
			if p.Units != "" {
				u = p.Units
			}
			pf("%4d %4d %12s %-8s %s", p.Offset, p.Length, p.Id, u, p.Type)
		}
	}
}

func inspect(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		log.Fatalf("excess arguments, aborting: %v", args)
	}

	confHandler := NewDeviceConfigHandler()

	// create default adapter from configuration
	defaultDevice := viper.GetString("adapter")
	if defaultDevice != "" {
		confHandler.DefaultDevice = defaultDevice
		confHandler.ConnectionManager(defaultDevice, viper.GetBool("rtu"), viper.GetInt("baudrate"), viper.GetString("comset"))
	}

	// create devices from command line
	devices, _ := cmd.PersistentFlags().GetStringSlice("devices")
	if len(devices) == 0 {
		fmt.Fprint(os.Stderr, "config: no devices found - terminating")
		os.Exit(1)
	}
	for _, dev := range devices {
		if dev != "" {
			confHandler.CreateDeviceFromSpec(dev)
		}
	}

	// raw log
	if viper.GetBool("raw") {
		setLogger(confHandler.Managers, golog.New(os.Stderr, "", 0))
	}

	for _, m := range confHandler.Managers {
		m := m // pin!
		m.All(func(id uint8, dev meters.Device) {
			m.Conn.Slave(id)
			scanSunspec(m.Conn.ModbusClient())
		})
	}
}
