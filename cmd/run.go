package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	latest "github.com/tcnksm/go-latest"

	"github.com/volkszaehler/mbmd/meters/rs485"
	"github.com/volkszaehler/mbmd/server"
)

const (
	cacheDuration = 1 * time.Minute
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Read and publish measurements from all configured devices",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:`,
	Run: func(cmd *cobra.Command, args []string) {
		run(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.PersistentFlags().StringSliceP(
		"devices", "d",
		[]string{},
		`MODBUS device type and ID to query, separated by comma.
  Example: -d SDM:22,DZG:1.
Valid types are:`+meterHelp()+`
To use an adapter different from default, append RTU device or TCP address separated by @.
If the adapter is a TCP connection (identified by :port), the device type (SUNS) is ignored and
any type is considered valid.
  Example: -d SDM:22@/dev/USB1,SMA:126@localhost:502.`,
	)
	runCmd.PersistentFlags().StringP(
		"api", "",
		"0.0.0.0:8080",
		"REST API url. Use 127.0.0.1:8080 to limit to localhost.",
	)
	runCmd.PersistentFlags().StringP(
		"mqtt", "m",
		"",
		"MQTT broker URI. ex: tcp://10.10.1.1:1883",
	)
	runCmd.PersistentFlags().StringP(
		"topic", "t",
		"mbmd",
		"MQTT root topic. Set empty to disable publishing.",
	)
	runCmd.PersistentFlags().StringP(
		"user", "u",
		"",
		"MQTT user (optional)",
	)
	runCmd.PersistentFlags().StringP(
		"password", "p",
		"",
		"MQTT password (optional)",
	)
	runCmd.PersistentFlags().StringP(
		"clientid", "i",
		"mbmd",
		"MQTT client id",
	)
	runCmd.PersistentFlags().Bool(
		"clean",
		false,
		"MQTT clean Session",
	)
	runCmd.PersistentFlags().IntP(
		"qos", "q",
		0,
		"MQTT quality of service 0,1,2",
	)
	runCmd.PersistentFlags().String(
		"homie",
		"homie",
		"MQTT Homie IoT discovery base topic (homieiot.github.io). Set empty to disable.",
	)

	// bind command line options to viper wit exceptions
	runCmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		if flag.Name == "devices" { // don't bind this key
			return
		}
		_ = viper.BindPFlag(flag.Name, flag)
	})

	_ = viper.BindPFlag("mqtt.broker", runCmd.PersistentFlags().Lookup("mqtt"))
	mqtt := []string{"topic", "user", "password", "clientid", "clean", "qos", "homie"}
	for _, f := range mqtt {
		flag := runCmd.PersistentFlags().Lookup(f)
		if flag == nil {
			panic("pfllag lookup failed for " + f)
		}
		_ = viper.BindPFlag("mqtt."+f, flag)
	}
}

// checkVersion validates if updates are available
func checkVersion() {
	githubTag := &latest.GithubTag{
		Owner:      "volkszaehler",
		Repository: "mbmd",
	}

	if res, err := latest.Check(githubTag, server.Version); err == nil {
		if res.Outdated {
			log.Printf("updates available - please upgrade to ingress %s", res.Current)
		}
	}
}

// meterHelp output list of supported devices
func meterHelp() string {
	s := fmt.Sprintf("\n  %s", "RTU")
	types := make([]string, 0)
	for t := range rs485.Producers {
		types = append(types, t)
	}

	sort.Strings(types)

	for _, t := range types {
		p := rs485.Producers[t]()
		s += fmt.Sprintf("\n    %-9s%s", t, p.Description())
	}

	s += fmt.Sprintf("\n  %s", "TCP")
	s += fmt.Sprintf("\n    %-9s%s", "SUNS", "Sunspec-compatible MODBUS TCP device (SMA, SolarEdge, KOSTAL, etc)")

	return s
}

func run(cmd *cobra.Command, args []string) {
	log.Printf("mbmd %s (%s)", server.Version, server.Commit)
	go checkVersion()

	confHandler := NewDeviceConfigHandler()
	if cfgFile != "" {
		// config file found
		log.Printf("using %s", viper.ConfigFileUsed())

		var conf Config
		if err := viper.Unmarshal(&conf); err != nil {
			log.Fatalf("failed parsing config file: %v", err)
		}

		// add adapters from configuration
		for _, a := range conf.Adapters {
			confHandler.CreateAdapter(a.Device, a.Baudrate, a.Comset)
		}

		// add devices from configuration
		for _, dev := range conf.Devices {
			confHandler.CreateDevice(dev)
		}
	}

	// create default adapter from configuration
	defaultDevice := viper.GetString("adapter")
	if defaultDevice != "" {
		confHandler.DefaultDevice = defaultDevice
		confHandler.CreateAdapter(defaultDevice, viper.GetInt("baudrate"), viper.GetString("comset"))
	}

	// create remaining devices from command line
	devices, _ := cmd.PersistentFlags().GetStringSlice("devices")
	for _, dev := range devices {
		if dev != "" {
			confHandler.CreateDeviceFromSpec(dev)
		}
	}

	// query engine
	qe := server.NewQueryEngine(confHandler.Managers)

	// result channels
	rc := make(chan server.QuerySnip)
	cc := make(chan server.ControlSnip)

	// tee that broadcasts meter messages to multiple recipients
	tee := server.NewQuerySnipBroadcaster(rc)
	go tee.Run()

	// status cache
	status := server.NewStatus(qe, cc)

	// websocket hub
	hub := server.NewSocketHub(status)
	tee.AttachRunner(hub.Run)

	// measurement cache for REST api
	cache := server.NewCache(cacheDuration, status, viper.GetBool("verbose"))
	tee.AttachRunner(cache.Run)

	httpd := server.NewHttpd(qe, cache)
	go httpd.Run(hub, status, viper.GetString("api"))

	// MQTT client
	if viper.GetString("mqtt.broker") != "" {
		mqtt := server.NewMqttClient(
			viper.GetString("mqtt.broker"),
			viper.GetString("mqtt.topic"),
			viper.GetString("mqtt.user"),
			viper.GetString("mqtt.password"),
			viper.GetString("mqtt.clientid"),
			viper.GetInt("mqtt.qos"),
			viper.GetBool("mqtt.clean"),
			viper.GetBool("verbose"),
		)

		// homie needs to scan the bus, start it first
		if viper.GetString("mqtt.homie") != "" {
			homieRunner := server.NewHomieRunner(mqtt, qe, viper.GetString("mqtt.homie"))
			tee.AttachRunner(homieRunner.Run)
		}

		// start "normal" mqtt handler after homie setup
		if viper.GetString("mqtt.topic") != "" {
			mqttRunner := server.MqttRunner{MqttClient: mqtt}
			tee.AttachRunner(mqttRunner.Run)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	go qe.Run(ctx, cc, rc)

	// wait for signal on exit channel and cancel context
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit
	log.Println("received signal - stopping")
	cancel()

	// wait for Run methods attached to tee to finish
	<-tee.Done()
	log.Println("stopped")
}
