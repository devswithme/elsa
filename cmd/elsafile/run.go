package elsafile

import (
	"fmt"
	"os"

	"github.com/risoftinc/elsa/constants"
	"github.com/risoftinc/elsa/internal/elsafile"
	"github.com/spf13/cobra"
)

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a command defined in Elsafile",
	Long: `Run executes a command defined in the Elsafile.
Commands are executed in the shell environment.

Usage:
  elsa run command_name    # Run a command from Elsafile
  elsa run command_name    # Alternative syntax (also supported)

Examples:
  elsa run build          # Run the build command from Elsafile
  elsa run test           # Run the test command from Elsafile
  elsa run clean          # Run the clean command from Elsafile

Note: If a command name conflicts with built-in Elsa commands,
      you must use 'run command_name' to execute the Elsafile command.
      Built-in commands take precedence when no prefix is used.`,
	DisableFlagParsing: true,
	Args:               cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Fprintf(os.Stderr, "Error: No command specified. Use 'elsa run command_name'\n")
			os.Exit(1)
		}

		// Check if the first argument starts with "run:"
		commandName := args[0]
		if len(commandName) > constants.ConflictResolutionPrefixLength && commandName[:constants.ConflictResolutionPrefixLength] == constants.ConflictResolutionPrefix {
			// Extract the actual command name after "run:"
			actualCommand := commandName[constants.ConflictResolutionPrefixLength:]
			if err := runElsafileCommand(actualCommand); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		} else {
			// If no "run:" prefix, treat the first argument as the command name
			if err := runElsafileCommand(commandName); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		}
	},
}

func runElsafileCommand(commandName string) error {
	// Create and use ConflictHandler from internal package
	handler := elsafile.NewConflictHandler()
	return handler.ExecuteElsafileCommand(commandName)
}
