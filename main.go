package main

import (
	"embed"
	"io/fs"

	"github.com/volkszaehler/mbmd/cmd"
	"github.com/volkszaehler/mbmd/server"
)

//go:embed assets
var assets embed.FS

// init loads embedded assets unless live assets are already loaded
func init() {
	if server.Assets == nil {
		fsys, err := fs.Sub(assets, server.AssetsDir)
		if err != nil {
			panic(err)
		}
		server.Assets = fsys
	}
}

func main() {
	cmd.Execute()
}
