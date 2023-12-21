package main

import (
	"github.com/josark2005/certark/ark"
	"github.com/josark2005/certark/certark"
	"github.com/josark2005/certark/cmd"
)

var version = "dev"

func main() {
	initial()
	cmd.Execute(version)
}

// Initialization before running
func initial() {
	certark.LoadConfig(true)

	// log level
	if certark.CurrentConfig.Mode == certark.MODE_PROD {
		ark.SetLevel(ark.InfoLevel)
	} else {
		ark.SetLevel(ark.DebugLevel)
	}
}
