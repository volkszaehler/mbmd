package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mbmd",
	Short: "ModBus Measurement Daemon",
	Long:  "Easily read and distribute data from ModBus meters and grid inverters",

	Run: func(cmd *cobra.Command, args []string) {
		run(cmd, args)
	},
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
		"/dev/ttyUSB0",
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
	rootCmd.PersistentFlags().BoolP(
		"verbose", "v",
		false,
		"Verbose mode",
	)

	// bind command line options
	_ = viper.BindPFlag("adapter", rootCmd.PersistentFlags().Lookup("adapter"))
	_ = viper.BindPFlag("baudrate", rootCmd.PersistentFlags().Lookup("baudrate"))
	_ = viper.BindPFlag("comset", rootCmd.PersistentFlags().Lookup("comset"))
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name "mbmd" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")    // optionally look for config in the working directory
		viper.AddConfigPath("/etc") // path to look for the config file in

		viper.SetConfigName("mbmd")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		// using config file
		cfgFile = viper.ConfigFileUsed()
	} else if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		// parsing failed - exit
		fmt.Println(err)
		os.Exit(1)
	} else {
		// not using config file
		cfgFile = ""
	}
}
