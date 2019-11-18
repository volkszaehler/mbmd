package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	latest "github.com/tcnksm/go-latest"

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
	Run: run,
}

// bindPflagsWithExceptions binds all pflags exception for exceptions
func bindPflagsWithExceptions(flags *pflag.FlagSet, exceptions ...string) {
	flags.VisitAll(func(flag *pflag.Flag) {
		for _, f := range exceptions {
			if flag.Name == f { // don't bind this key
				return
			}
		}
		_ = viper.BindPFlag(flag.Name, flag)
	})
}

func bindPFlagsWithPrefix(flags *pflag.FlagSet, prefix string, names ...string) {
	for _, f := range names {
		flag := flags.Lookup(prefix + "-" + f)
		if flag == nil {
			panic("pflag lookup failed for " + f)
		}
		_ = viper.BindPFlag(prefix+"."+f, flag)
	}
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.PersistentFlags().StringSliceP(
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
	runCmd.PersistentFlags().String(
		"api",
		"0.0.0.0:8080",
		"REST API url. Use 127.0.0.1:8080 to limit to localhost.",
	)
	runCmd.PersistentFlags().StringP(
		"mqtt-broker", "m",
		"",
		"MQTT broker URI. ex: tcp://10.10.1.1:1883",
	)
	runCmd.PersistentFlags().String(
		"mqtt-topic",
		"mbmd",
		"MQTT root topic. Set empty to disable publishing.",
	)
	runCmd.PersistentFlags().String(
		"mqtt-user",
		"",
		"MQTT user (optional)",
	)
	runCmd.PersistentFlags().String(
		"mqtt-password",
		"",
		"MQTT password (optional)",
	)
	runCmd.PersistentFlags().String(
		"mqtt-clientid",
		"mbmd",
		"MQTT client id",
	)
	runCmd.PersistentFlags().Bool(
		"mqtt-clean",
		false,
		"MQTT clean Session",
	)
	runCmd.PersistentFlags().Int(
		"mqtt-qos",
		0,
		"MQTT quality of service 0,1,2 (default 0)",
	)
	runCmd.PersistentFlags().String(
		"mqtt-homie",
		"homie",
		"MQTT Homie IoT discovery base topic (homieiot.github.io). Set empty to disable.",
	)
	runCmd.PersistentFlags().StringP(
		"influx-url", "i",
		"",
		"InfluxDB URL. ex: http://10.10.1.1:8086",
	)
	runCmd.PersistentFlags().String(
		"influx-database",
		"",
		"InfluxDB database",
	)
	runCmd.PersistentFlags().String(
		"influx-measurement",
		"data",
		"InfluxDB measurement",
	)
	runCmd.PersistentFlags().String(
		"influx-precision",
		"s",
		"InfluxDB precision",
	)
	runCmd.PersistentFlags().Duration(
		"influx-interval",
		30*time.Second,
		"InfluxDB write interval",
	)
	runCmd.PersistentFlags().String(
		"influx-user",
		"",
		"InfluxDB user (optional)",
	)
	runCmd.PersistentFlags().String(
		"influx-password",
		"",
		"InfluxDB password (optional)",
	)

	pflags := runCmd.PersistentFlags()

	// bind command line options to viper with exceptions
	bindPflagsWithExceptions(pflags, "devices")

	// mqtt
	bindPFlagsWithPrefix(pflags, "mqtt", "broker", "topic", "user", "password", "clientid", "clean", "qos", "homie")

	// influx
	bindPFlagsWithPrefix(pflags, "influx", "url", "database", "measurement", "precision", "interval", "user", "password")
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

func run(cmd *cobra.Command, args []string) {
	log.Printf("mbmd %s (%s)", server.Version, server.Commit)
	go checkVersion()

	confHandler := NewDeviceConfigHandler()

	// create default adapter from configuration
	defaultDevice := viper.GetString("adapter")
	if defaultDevice != "" {
		confHandler.DefaultDevice = defaultDevice
		confHandler.CreateAdapter(defaultDevice, viper.GetInt("baudrate"), viper.GetString("comset"))
	}

	// create devices from command line
	devices, _ := cmd.PersistentFlags().GetStringSlice("devices")
	for _, dev := range devices {
		if dev != "" {
			confHandler.CreateDeviceFromSpec(dev)
		}
	}

	if cfgFile != "" {
		// config file found
		log.Printf("config: using %s", viper.ConfigFileUsed())

		var conf Config
		if err := viper.GetViper().UnmarshalExact(&conf); err != nil {
			log.Fatalf("config: failed parsing config file %s: %v", cfgFile, err)
		}

		// create devices from config file only if not overridden on command line
		if len(devices) == 0 {
			// add adapters from configuration
			for _, a := range conf.Adapters {
				confHandler.CreateAdapter(a.Device, a.Baudrate, a.Comset)
			}

			// add devices from configuration
			for _, dev := range conf.Devices {
				confHandler.CreateDevice(dev)
			}
		}
	}

	if countDevices(confHandler.Managers) == 0 {
		log.Fatal("config: no devices found - terminiating")
	}

	// raw log
	if viper.GetBool("raw") {
		setLogger(confHandler.Managers, log.New(os.Stderr, "", log.LstdFlags))
	}

	// query engine
	qe := server.NewQueryEngine(confHandler.Managers)

	// results- and control channels
	rc := make(chan server.QuerySnip)
	cc := make(chan server.ControlSnip)

	// tee that broadcasts meter messages to multiple recipients
	tee := server.NewBroadcaster(server.FromSnipChannel(rc))
	go tee.Run()

	// tee that broadcasts control messages to multiple recipients
	teeC := server.NewBroadcaster(server.FromControlChannel(cc))
	go teeC.Run()

	// status cache (always needed to consume control messages)
	status := server.NewStatus(qe, server.ToControlChannel(teeC.Attach()))

	// web server
	if viper.GetString("api") != "" {
		// measurement cache for REST api
		cache := server.NewCache(cacheDuration, status, viper.GetBool("verbose"))
		tee.AttachRunner(server.NewSnipRunner(cache.Run))

		// websocket hub
		hub := server.NewSocketHub(status)
		tee.AttachRunner(server.NewSnipRunner(hub.Run))

		// http daemon
		httpd := server.NewHttpd(qe, cache)
		go httpd.Run(hub, status, viper.GetString("api"))
	}

	// MQTT client
	if viper.GetString("mqtt.broker") != "" {
		qos := byte(viper.GetInt("mqtt.qos"))
		verbose := viper.GetBool("verbose")

		// default mqtt runner
		if topic := viper.GetString("mqtt.topic"); topic != "" {
			options := server.NewMqttOptions(
				viper.GetString("mqtt.broker"),
				viper.GetString("mqtt.user"),
				viper.GetString("mqtt.password"),
				viper.GetString("mqtt.clientid"),
				viper.GetBool("mqtt.clean"),
			)
			mqttRunner := server.NewMqttRunner(options, qos, topic, verbose)
			tee.AttachRunner(server.NewSnipRunner(mqttRunner.Run))
		}

		// homie runner
		if topic := viper.GetString("mqtt.homie"); topic != "" {
			options := server.NewMqttOptions(
				viper.GetString("mqtt.broker"),
				viper.GetString("mqtt.user"),
				viper.GetString("mqtt.password"),
				viper.GetString("mqtt.clientid"),
				viper.GetBool("mqtt.clean"),
			)
			cc := server.ToControlChannel(teeC.Attach())
			homieRunner := server.NewHomieRunner(qe, cc, options, qos, topic, verbose)
			tee.AttachRunner(server.NewSnipRunner(homieRunner.Run))
		}
	}

	// InfluxDB client
	if viper.GetString("influx.url") != "" {
		influx := server.NewInfluxClient(
			viper.GetString("influx.url"),
			viper.GetString("influx.database"),
			viper.GetString("influx.measurement"),
			viper.GetString("influx.precision"),
			viper.GetDuration("influx.interval"),
			viper.GetString("influx.user"),
			viper.GetString("influx.password"),
			viper.GetBool("verbose"),
		)

		tee.AttachRunner(server.NewSnipRunner(influx.Run))
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
