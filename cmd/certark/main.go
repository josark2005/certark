package main

import (
	"github.com/jokin1999/certark/certark"
	"github.com/jokin1999/certark/cmd"
)

var version = "dev"

func main() {
	initial()
	cmd.Execute(version)
}

// Initialization before running
func initial() {
	certark.LoadConfig(true)
}
