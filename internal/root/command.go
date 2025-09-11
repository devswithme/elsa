package root

import (
	"github.com/spf13/cobra"
	"go.risoftinc.com/elsa/cmd/elsafile"
)

// CommandHandler handles root command logic and subcommand management
type CommandHandler struct {
	displayHelper *DisplayHelper
}

// NewCommandHandler creates a new command handler instance
func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		displayHelper: NewDisplayHelper(),
	}
}

// HandleRootCommand handles the root command execution logic
func (h *CommandHandler) HandleRootCommand(cmd *cobra.Command, args []string, version string) error {
	if len(args) > 0 {
		// Check if the first argument is a flag
		if h.displayHelper.IsFlag(args[0]) {
			// This is a flag, let cobra handle it normally
			return nil
		}

		// Try to handle as Elsafile command
		return h.handleElsafileCommand(cmd, args[0])
	}

	// Show custom root help if no arguments
	h.displayHelper.ShowRootHelp(cmd, version)
	return nil
}

// handleElsafileCommand handles unknown commands by delegating to Elsafile
func (h *CommandHandler) handleElsafileCommand(cmd *cobra.Command, commandName string) error {
	handler := elsafile.NewSimpleHandlerWithRoot(cmd)

	if err := handler.HandleUnknownCommand(commandName); err != nil {
		// If it's not an Elsafile command, show suggestions
		suggestions := handler.SuggestCommands(commandName)
		h.displayHelper.ShowSuggestions(suggestions)
		return err
	}

	return nil
}
