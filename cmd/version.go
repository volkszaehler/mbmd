package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/volkszaehler/mbmd/server"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show MBMD version",

	Run: func(cmd *cobra.Command, args []string) {
		displayVersion(rootCmd.Name())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func displayVersion(name string) {
	fmt.Printf(name+`
  version     : %s
  commit      : %s
  go version  : %s
  go compiler : %s
  platform    : %s/%s
`, server.Version, server.Commit, runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)
}
