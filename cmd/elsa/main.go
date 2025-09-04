package main

import (
	"fmt"
	"os"

	"github.com/risoftinc/elsa/cmd"
)

var (
	version = "0.6.7"
)

func main() {
	// Set version info
	cmd.SetVersionInfo(version)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
