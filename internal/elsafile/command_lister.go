package elsafile

import (
	"fmt"
	"strings"

	"github.com/risoftinc/elsa/constants"
	"github.com/spf13/cobra"
)

// CommandLister handles listing and displaying Elsafile commands
type CommandLister struct {
	conflictHandler *ConflictHandler
}

// NewCommandLister creates a new CommandLister instance
func NewCommandLister() *CommandLister {
	return &CommandLister{
		conflictHandler: NewConflictHandler(),
	}
}

// NewCommandListerWithHandler creates a new CommandLister instance with a specific conflict handler
func NewCommandListerWithHandler(handler *ConflictHandler) *CommandLister {
	return &CommandLister{
		conflictHandler: handler,
	}
}

// NewCommandListerWithRoot creates a new CommandLister instance with root command for dynamic built-in detection
func NewCommandListerWithRoot(rootCmd *cobra.Command) *CommandLister {
	conflictHandler := NewConflictHandlerWithRoot(rootCmd)
	return &CommandLister{
		conflictHandler: conflictHandler,
	}
}

// ListAllCommands lists all available commands from Elsafile
func (cl *CommandLister) ListAllCommands() error {
	commands, err := cl.conflictHandler.ListCommands()
	if err != nil {
		return err
	}

	if len(commands) == 0 {
		fmt.Printf("%s %s\n", constants.PencilEmoji, constants.MsgNoCommandsFound)
		fmt.Printf("%s Run 'elsa init' to create an Elsafile with default commands\n", constants.InfoEmoji)
		return nil
	}

	fmt.Printf("%s Available commands in Elsafile:\n", constants.ClipboardEmoji)
	fmt.Println()

	// Show all commands
	for name, cmd := range commands {
		conflict := ""
		if cl.conflictHandler.HasConflict(name) {
			conflict = " ⚠️  (conflicts with built-in, use 'run' prefix)"
		}

		fmt.Printf("  %s%s\n", name, conflict)
		if len(cmd.Commands) > 0 {
			fmt.Printf("    %s\n", strings.Join(cmd.Commands, constants.CommandSeparator))
		}
		fmt.Println()
	}

	// Show usage instructions
	fmt.Printf("%s Usage:\n", constants.InfoEmoji)
	fmt.Printf("  %s\n", constants.UsageRunCommand)
	fmt.Printf("  %s\n", constants.UsageListConflicts)
	fmt.Printf("  %s\n", constants.UsageInit)

	return nil
}

// ListConflictingCommands lists only commands that conflict with built-ins
func (cl *CommandLister) ListConflictingCommands() error {
	commands, err := cl.conflictHandler.ListCommands()
	if err != nil {
		return err
	}

	conflicts, err := cl.conflictHandler.GetConflictingCommands()
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
			fmt.Printf("    %s\n", strings.Join(cmd.Commands, constants.CommandSeparator))
		}
		fmt.Println()
	}
	return nil
}

// GetCommandSummary returns a summary of all commands
func (cl *CommandLister) GetCommandSummary() (map[string]string, error) {
	commands, err := cl.conflictHandler.ListCommands()
	if err != nil {
		return nil, err
	}

	summary := make(map[string]string)
	for name, cmd := range commands {
		summary[name] = strings.Join(cmd.Commands, constants.CommandSeparator)
	}

	return summary, nil
}

// GetConflictSummary returns a summary of conflicting commands
func (cl *CommandLister) GetConflictSummary() ([]string, error) {
	return cl.conflictHandler.GetConflictingCommands()
}

// FormatCommandDisplay formats a command for display
func (cl *CommandLister) FormatCommandDisplay(name string, cmd *Command, showConflict bool) string {
	var result strings.Builder

	conflict := ""
	if showConflict && cl.conflictHandler.HasConflict(name) {
		conflict = " ⚠️  (conflicts with built-in, use 'run' prefix)"
	}

	result.WriteString(fmt.Sprintf("  %s%s\n", name, conflict))
	if len(cmd.Commands) > 0 {
		result.WriteString(fmt.Sprintf("    %s\n", strings.Join(cmd.Commands, constants.CommandSeparator)))
	}
	result.WriteString("\n")

	return result.String()
}
