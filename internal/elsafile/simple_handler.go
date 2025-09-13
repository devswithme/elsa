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

// NewSimpleHandlerWithRoot creates a new SimpleHandler instance with root command for dynamic built-in detection
func NewSimpleHandlerWithRoot(rootCmd interface{}) *SimpleHandler {
	// Type assertion to check if rootCmd has Commands() method
	if _, ok := rootCmd.(interface{ Commands() []interface{} }); ok {
		return &SimpleHandler{
			elsafileManager: NewManagerWithRoot("Elsafile", nil), // We'll need to handle this differently
		}
	}

	return &SimpleHandler{
		elsafileManager: NewManager("Elsafile"),
	}
}

// NewSimpleHandlerWithManager creates a new SimpleHandler instance with a specific manager
func NewSimpleHandlerWithManager(manager *Manager) *SimpleHandler {
	return &SimpleHandler{
		elsafileManager: manager,
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

	// Execute each command sequentially
	for _, cmd := range command.Commands {
		if err := h.elsafileManager.ExecuteShellCommand(cmd); err != nil {
			return err
		}
	}
	return nil
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

// GetSuggestionMessage returns a formatted message with command suggestions
func (h *SimpleHandler) GetSuggestionMessage(commandName string) string {
	suggestions := h.SuggestCommands(commandName)
	if len(suggestions) == 0 {
		return fmt.Sprintf("💡 No similar commands found for '%s'", commandName)
	}

	msg := "💡 Did you mean one of these commands?\n"
	for _, suggestion := range suggestions {
		msg += fmt.Sprintf("   %s\n", suggestion)
	}
	msg += "\n💡 Use 'elsa run command_name' to execute Elsafile commands"

	return msg
}
