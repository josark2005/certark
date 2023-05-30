package main

import (
	"github.com/jokin1999/certark/cmd"
)

var version = "dev"

func main() {
	cmd.Execute(version)
}
