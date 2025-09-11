package elsafile

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.risoftinc.com/elsa/internal/elsafile"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all commands defined in Elsafile",
	Long: `List displays all available commands defined in the Elsafile.
This is useful to see what custom commands are available for your project.

Examples:
  elsa list           # List all commands from Elsafile
  elsa list --conflicts  # Show only conflicting commands`,
	Run: func(cmd *cobra.Command, args []string) {
		showConflicts, _ := cmd.Flags().GetBool("conflicts")
		if err := listElsafileCommands(cmd.Root(), showConflicts); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	ListCmd.Flags().BoolP("conflicts", "c", false, "Show only commands that conflict with built-in commands")
}

func listElsafileCommands(rootCmd *cobra.Command, showConflictsOnly bool) error {
	commandLister := elsafile.NewCommandListerWithRoot(rootCmd)

	if showConflictsOnly {
		return commandLister.ListConflictingCommands()
	}

	return commandLister.ListAllCommands()
}
