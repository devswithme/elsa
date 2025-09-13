package main

import (
	"fmt"
	"os"

	"go.risoftinc.com/elsa/cmd"
)

var (
	version = "0.8.2"
)

func main() {
	// Set version info
	cmd.SetVersionInfo(version)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
