package elsafile

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
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
		if err := listElsafileCommands(showConflicts); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	ListCmd.Flags().BoolP("conflicts", "c", false, "Show only commands that conflict with built-in commands")
}

func listElsafileCommands(showConflictsOnly bool) error {
	// Create and use ConflictHandler
	handler := NewConflictHandler()

	commands, err := handler.ListCommands()
	if err != nil {
		return err
	}

	if len(commands) == 0 {
		fmt.Println("📝 No commands found in Elsafile")
		fmt.Println("💡 Run 'elsa init' to create an Elsafile with default commands")
		return nil
	}

	if showConflictsOnly {
		conflicts, err := handler.GetConflictingCommands()
		if err != nil {
			return err
		}

		if len(conflicts) == 0 {
			fmt.Println("✅ No command conflicts found")
			return nil
		}

		fmt.Println("⚠️  Commands that conflict with built-in Elsa commands:")
		fmt.Println("   (Use 'elsa run command_name' to execute these)")
		fmt.Println()
		for _, name := range conflicts {
			cmd := commands[name]
			fmt.Printf("  %s\n", name)
			if len(cmd.Commands) > 0 {
				fmt.Printf("    %s\n", strings.Join(cmd.Commands, " && "))
			}
			fmt.Println()
		}
		return nil
	}

	fmt.Println("📋 Available commands in Elsafile:")
	fmt.Println()

	// Show all commands
	for name, cmd := range commands {
		conflict := ""
		if handler.HasConflict(name) {
			conflict = " ⚠️  (conflicts with built-in, use 'run:' prefix)"
		}

		fmt.Printf("  %s%s\n", name, conflict)
		if len(cmd.Commands) > 0 {
			fmt.Printf("    %s\n", strings.Join(cmd.Commands, " && "))
		}
		fmt.Println()
	}

	// Show usage instructions
	fmt.Println("💡 Usage:")
	fmt.Println("  elsa run command_name    # Execute a command from Elsafile")
	fmt.Println("  elsa list --conflicts    # Show conflicting commands")
	fmt.Println("  elsa init                # Create a new Elsafile")

	return nil
}
