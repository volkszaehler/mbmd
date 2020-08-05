package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:               "mbmd",
	Short:             "ModBus Measurement Daemon",
	Long:              "Easily read and distribute data from ModBus meters and grid inverters",
	DisableAutoGenTag: true, // prevent changing timestamps
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile,
		"config", "c",
		"",
		"Config file (default is $HOME/mbmd.yaml)",
	)
	rootCmd.PersistentFlags().StringP(
		"adapter", "a",
		"",
		`Default MODBUS adapter. This option can be used if all devices are attached to a single adapter.
Can be either an RTU device (/dev/ttyUSB0) or TCP socket (localhost:502).
The default adapter can be overridden per device`,
	)
	rootCmd.PersistentFlags().IntP(
		"baudrate", "b",
		9600,
		`Serial interface baud rate`,
	)
	rootCmd.PersistentFlags().String(
		"comset",
		"8N1",
		`Communication parameters for default adapter, either 8N1 or 8E1.
Only applicable if the default adapter is an RTU device`,
	)
	rootCmd.PersistentFlags().Bool(
		"rtu",
		false,
		`Use RTU over TCP for default adapter.
Typically used with RS485 to Ethernet adapters that don't perform protocol conversion (e.g. USR-TCP232).
Only applicable if the default adapter is a TCP connection`,
	)
	rootCmd.PersistentFlags().BoolP(
		"help", "h",
		false,
		"Help for "+rootCmd.Name(),
	)
	rootCmd.PersistentFlags().BoolP(
		"verbose", "v",
		false,
		"Verbose mode",
	)
	rootCmd.PersistentFlags().Bool(
		"raw",
		false,
		"Log raw device data",
	)

	// bind command line options
	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		log.Fatal(err)
	}
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// search for config in home directory if available
		if home, err := os.UserHomeDir(); err == nil {
			viper.AddConfigPath(home)
		}

		viper.AddConfigPath(".")    // optionally look for config in the working directory
		viper.AddConfigPath("/etc") // path to look for the config file in

		viper.SetConfigName("mbmd")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		// using config file
		cfgFile = viper.ConfigFileUsed()
	} else {
		var configFileNotFound viper.ConfigFileNotFoundError
		var unsupportedConfig viper.UnsupportedConfigError
		if errors.As(err, &configFileNotFound) || errors.As(err, &unsupportedConfig) {
			// not using config file
			cfgFile = ""
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
