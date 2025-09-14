package make

import (
	"fmt"
	"os"

	"go.risoftinc.com/elsa/internal/make"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: elsa make <template-type> <name>")
		fmt.Println("Example: elsa make repository user_repository")
		os.Exit(1)
	}

	command := make.NewMakeCommand()

	// Handle help flags
	if os.Args[1] == "--help" || os.Args[1] == "-h" || os.Args[1] == "help" {
		if err := command.Execute([]string{}); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Handle list command
	if os.Args[1] == "list" {
		if err := command.ListAvailableTypes(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Execute make command
	if err := command.Execute(os.Args[1:]); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
