package elsafile

import (
	"fmt"
	"strings"
)

// SimpleHandler handles unknown commands by checking Elsafile
type SimpleHandler struct {
	elsafileManager *Manager
}

// NewSimpleHandler creates a new SimpleHandler instance
func NewSimpleHandler() *SimpleHandler {
	return &SimpleHandler{
		elsafileManager: NewManager("Elsafile"),
	}
}

// HandleUnknownCommand handles unknown commands by checking Elsafile
func (h *SimpleHandler) HandleUnknownCommand(commandName string) error {
	// Try to load Elsafile
	if err := h.elsafileManager.Load(); err != nil {
		// If Elsafile doesn't exist, this is truly an unknown command
		return fmt.Errorf("unknown command '%s'", commandName)
	}

	// Check if command exists in Elsafile
	command, exists := h.elsafileManager.GetCommand(commandName)
	if !exists {
		return fmt.Errorf("unknown command '%s'", commandName)
	}

	// Check if there's a conflict with built-in commands
	if h.elsafileManager.HasConflict(commandName) {
		fmt.Printf("⚠️  Command '%s' conflicts with a built-in Elsa command\n", commandName)
		fmt.Printf("💡 Use 'elsa run %s' to execute the Elsafile command\n", commandName)
		fmt.Printf("   Or use 'elsa %s' to run the built-in command\n\n", commandName)
		return fmt.Errorf("command conflict detected. Use 'run:' prefix to execute Elsafile command")
	}

	// No conflict, execute the command
	fmt.Printf("🚀 Running Elsafile command: %s\n", commandName)
	fmt.Printf("📝 Executing: %s\n\n", strings.Join(command.Commands, " && "))

	// Join all commands with && to execute them sequentially
	fullCommand := strings.Join(command.Commands, " && ")
	return h.elsafileManager.ExecuteShellCommand(fullCommand)
}

// SuggestCommands suggests similar commands from Elsafile
func (h *SimpleHandler) SuggestCommands(commandName string) []string {
	if err := h.elsafileManager.Load(); err != nil {
		return nil
	}

	var suggestions []string
	commands := h.elsafileManager.ListCommands()

	for name := range commands {
		if strings.Contains(strings.ToLower(name), strings.ToLower(commandName)) ||
			strings.Contains(strings.ToLower(commandName), strings.ToLower(name)) {
			suggestions = append(suggestions, name)
		}
	}

	return suggestions
}
