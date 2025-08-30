package main

import (
	"fmt"
	"os"

	cmd "github.com/risoftinc/elsa/cmd/root"
)

var (
	version = "0.5.5"
)

func main() {
	// Set version info
	cmd.SetVersionInfo(version)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
