package elsafile

import (
	"fmt"
	"strings"

	"github.com/risoftinc/elsa/constants"
)

// Formatter handles formatting of output messages
type Formatter struct{}

// NewFormatter creates a new Formatter instance
func NewFormatter() *Formatter {
	return &Formatter{}
}

// FormatCommand formats a single command for display
func (f *Formatter) FormatCommand(cmd *Command) string {
	if cmd == nil {
		return "  <nil command>"
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("  %s\n", cmd.Name))

	if len(cmd.Commands) > 0 {
		result.WriteString(fmt.Sprintf("    %s\n", strings.Join(cmd.Commands, constants.CommandSeparator)))
	}

	return result.String()
}

// FormatCommandList formats a list of commands for display
func (f *Formatter) FormatCommandList(commands map[string]*Command) string {
	if len(commands) == 0 {
		return fmt.Sprintf("%s %s", constants.PencilEmoji, constants.MsgNoCommandsFound)
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("%s Available commands in Elsafile:\n\n", constants.ClipboardEmoji))

	for _, cmd := range commands {
		result.WriteString(f.FormatCommand(cmd))
		result.WriteString("\n")
	}

	return result.String()
}

// FormatConflictList formats a list of conflicting commands for display
func (f *Formatter) FormatConflictList(conflicts []string) string {
	if len(conflicts) == 0 {
		return fmt.Sprintf("%s %s", constants.SuccessEmoji, constants.MsgNoConflictsFound)
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("%s Commands that conflict with built-in Elsa commands:\n", constants.WarningEmoji))
	result.WriteString("   (Use 'elsa run command_name' to execute these)\n\n")

	for _, name := range conflicts {
		result.WriteString(fmt.Sprintf("  %s\n", name))
	}
	result.WriteString("\n")

	return result.String()
}

// FormatSuccessMessage formats a success message
func (f *Formatter) FormatSuccessMessage(message string) string {
	return fmt.Sprintf("%s %s", constants.SuccessEmoji, message)
}

// FormatErrorMessage formats an error message
func (f *Formatter) FormatErrorMessage(message string) string {
	return fmt.Sprintf("%s %s", constants.ErrorEmoji, message)
}

// FormatWarningMessage formats a warning message
func (f *Formatter) FormatWarningMessage(message string) string {
	return fmt.Sprintf("%s %s", constants.WarningEmoji, message)
}

// FormatInfoMessage formats an info message
func (f *Formatter) FormatInfoMessage(message string) string {
	return fmt.Sprintf("%s %s", constants.InfoEmoji, message)
}

// FormatCommandExecution formats command execution information
func (f *Formatter) FormatCommandExecution(name string, commands []string) string {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("%s Running Elsafile command: %s\n", constants.RocketEmoji, name))
	result.WriteString(fmt.Sprintf("%s Executing: %s\n\n", constants.PencilEmoji, strings.Join(commands, constants.CommandSeparator)))
	return result.String()
}

// FormatUsageInstructions formats usage instructions
func (f *Formatter) FormatUsageInstructions() string {
	return fmt.Sprintf(`%s Usage:
  %s
  %s
  %s`, constants.InfoEmoji, constants.UsageRunCommand, constants.UsageListConflicts, constants.UsageInit)
}

// FormatCommandSummary formats a command summary
func (f *Formatter) FormatCommandSummary(commands map[string]*Command) string {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("%s Command Summary:\n", constants.ChartEmoji))
	result.WriteString(fmt.Sprintf("   Total commands: %d\n", len(commands)))

	// Count total command lines
	totalLines := 0
	for _, cmd := range commands {
		totalLines += len(cmd.Commands)
	}
	result.WriteString(fmt.Sprintf("   Total command lines: %d\n", totalLines))

	return result.String()
}

// FormatElsafileInfo formats Elsafile information
func (f *Formatter) FormatElsafileInfo(info *ElsafileInfo) string {
	if info == nil {
		return fmt.Sprintf("%s Elsafile: Not found", constants.DocumentEmoji)
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("%s Elsafile: %s\n", constants.DocumentEmoji, info.FilePath))
	result.WriteString(fmt.Sprintf("   Commands: %d\n", info.TotalCommands))
	result.WriteString(fmt.Sprintf("   Conflicts: %d\n", info.Conflicts))
	result.WriteString(fmt.Sprintf("   Valid: %t\n", info.IsValid))

	if info.LastModified > 0 {
		result.WriteString(fmt.Sprintf("   Last modified: %d\n", info.LastModified))
	}

	return result.String()
}

// TruncateString truncates a string to a maximum length
func (f *Formatter) TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength-3] + "..."
}

// PadString pads a string to a specific width
func (f *Formatter) PadString(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}
