package elsafile

import (
	"fmt"

	"github.com/risoftinc/elsa/constants"
)

// ConflictHandler handles command conflicts between built-in and Elsafile commands
type ConflictHandler struct {
	elsafileManager *Manager
}

// NewConflictHandler creates a new ConflictHandler instance
func NewConflictHandler() *ConflictHandler {
	return &ConflictHandler{
		elsafileManager: NewManager(constants.DefaultElsafileName),
	}
}

// NewConflictHandlerWithManager creates a new ConflictHandler instance with a specific manager
func NewConflictHandlerWithManager(manager *Manager) *ConflictHandler {
	return &ConflictHandler{
		elsafileManager: manager,
	}
}

// ExecuteElsafileCommand executes a command from Elsafile (used by run command)
func (h *ConflictHandler) ExecuteElsafileCommand(commandName string) error {
	if err := h.elsafileManager.Load(); err != nil {
		return err
	}

	return h.elsafileManager.ExecuteCommand(commandName)
}

// ListCommands lists all available commands from Elsafile
func (h *ConflictHandler) ListCommands() (map[string]*Command, error) {
	if err := h.elsafileManager.Load(); err != nil {
		return nil, err
	}

	return h.elsafileManager.ListCommands(), nil
}

// GetConflictingCommands returns commands that conflict with built-ins
func (h *ConflictHandler) GetConflictingCommands() ([]string, error) {
	if err := h.elsafileManager.Load(); err != nil {
		return nil, err
	}

	return h.elsafileManager.GetConflictingCommands(), nil
}

// HasConflict checks if a specific command conflicts with built-ins
func (h *ConflictHandler) HasConflict(commandName string) bool {
	if err := h.elsafileManager.Load(); err != nil {
		return false
	}

	return h.elsafileManager.HasConflict(commandName)
}

// GetConflictMessage returns a formatted message about command conflicts
func (h *ConflictHandler) GetConflictMessage(commandName string) string {
	return fmt.Sprintf(`⚠️  Command '%s' conflicts with a built-in Elsa command
💡 Use 'elsa run %s' to execute the Elsafile command
   Or use 'elsa %s' to run the built-in command`, commandName, commandName, commandName)
}

// GetConflictResolutionMessage returns a message explaining how to resolve conflicts
func (h *ConflictHandler) GetConflictResolutionMessage() string {
	return `💡 To resolve command conflicts:
   - Use 'elsa run command_name' to execute Elsafile commands
   - Use 'elsa command_name' to run built-in commands
   - Rename conflicting commands in your Elsafile to avoid conflicts`
}
