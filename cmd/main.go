package main

import (
	"fmt"
	"os"

	"github.com/risoftinc/elsa/cmd/root"
)

var (
	version = "0.6.0"
)

func main() {
	// Set version info
	root.SetVersionInfo(version)

	if err := root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
