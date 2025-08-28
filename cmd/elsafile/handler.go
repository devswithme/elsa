package elsafile

// ConflictHandler handles command conflicts between built-in and Elsafile commands
type ConflictHandler struct {
	elsafileManager *Manager
}

// NewConflictHandler creates a new ConflictHandler instance
func NewConflictHandler() *ConflictHandler {
	return &ConflictHandler{
		elsafileManager: NewManager("Elsafile"),
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
