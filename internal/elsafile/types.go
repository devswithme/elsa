package elsafile

// Command represents a command defined in the Elsafile
type Command struct {
	Name        string
	Description string
	Commands    []string
}

// Manager handles parsing and execution of Elsafile commands
type Manager struct {
	commands    map[string]*Command
	filepath    string
	rootCommand interface{} // Will be *cobra.Command but we don't want to import cobra here
}

// NewManager creates a new Manager instance
func NewManager(filepath string) *Manager {
	return &Manager{
		commands: make(map[string]*Command),
		filepath: filepath,
	}
}

// NewManagerWithRoot creates a new Manager instance with root command for dynamic built-in detection
func NewManagerWithRoot(filepath string, rootCmd interface{}) *Manager {
	return &Manager{
		commands:    make(map[string]*Command),
		filepath:    filepath,
		rootCommand: rootCmd,
	}
}

// GetCommand returns a command by name
func (em *Manager) GetCommand(name string) (*Command, bool) {
	cmd, exists := em.commands[name]
	return cmd, exists
}

// ListCommands returns all available commands
func (em *Manager) ListCommands() map[string]*Command {
	return em.commands
}

// GetConflictingCommands returns a list of commands that conflict with built-ins
func (em *Manager) GetConflictingCommands() []string {
	var conflicts []string
	for name := range em.commands {
		if em.HasConflict(name) {
			conflicts = append(conflicts, name)
		}
	}
	return conflicts
}
