package elsafile

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

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
	rootCommand *cobra.Command
}

// NewManager creates a new Manager instance
func NewManager(filepath string) *Manager {
	return &Manager{
		commands: make(map[string]*Command),
		filepath: filepath,
	}
}

// NewManagerWithRoot creates a new Manager instance with root command for dynamic built-in detection
func NewManagerWithRoot(filepath string, rootCmd *cobra.Command) *Manager {
	return &Manager{
		commands:    make(map[string]*Command),
		filepath:    filepath,
		rootCommand: rootCmd,
	}
}

// Load loads and parses the Elsafile
func (em *Manager) Load() error {
	if _, err := os.Stat(em.filepath); os.IsNotExist(err) {
		return fmt.Errorf("Elsafile not found at %s. Run 'elsa init' to create one", em.filepath)
	}

	file, err := os.Open(em.filepath)
	if err != nil {
		return fmt.Errorf("failed to open Elsafile: %v", err)
	}
	defer file.Close()

	return em.parseFile(file)
}

// parseFile parses the Elsafile content
func (em *Manager) parseFile(file *os.File) error {
	scanner := bufio.NewScanner(file)
	var currentCommand *Command

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check if this is a command definition (ends with :)
		if strings.HasSuffix(line, ":") {
			// Save previous command if exists
			if currentCommand != nil && len(currentCommand.Commands) > 0 {
				em.commands[currentCommand.Name] = currentCommand
			}

			// Start new command
			commandName := strings.TrimSuffix(line, ":")
			currentCommand = &Command{
				Name:     commandName,
				Commands: []string{},
			}
		} else if currentCommand != nil {
			// Add line to current command
			currentCommand.Commands = append(currentCommand.Commands, line)
		}
	}

	// Save last command
	if currentCommand != nil && len(currentCommand.Commands) > 0 {
		em.commands[currentCommand.Name] = currentCommand
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading Elsafile: %v", err)
	}

	return nil
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

// ExecuteCommand executes a command from the Elsafile
func (em *Manager) ExecuteCommand(name string) error {
	command, exists := em.GetCommand(name)
	if !exists {
		return fmt.Errorf("command '%s' not found in Elsafile", name)
	}

	fmt.Printf("🚀 Running Elsafile command: %s\n", name)
	fmt.Printf("📝 Executing: %s\n\n", strings.Join(command.Commands, " && "))

	// Join all commands with && to execute them sequentially
	fullCommand := strings.Join(command.Commands, " && ")
	return em.ExecuteShellCommand(fullCommand)
}

// ExecuteShellCommand executes a shell command
func (em *Manager) ExecuteShellCommand(command string) error {
	var shell string
	var args []string

	// Detect OS and use appropriate shell
	if os.PathSeparator == '\\' {
		// Windows
		shell = "cmd"
		args = []string{"/C", command}
	} else {
		// Unix-like systems
		shell = "/bin/sh"
		args = []string{"-c", command}
	}

	cmd := exec.Command(shell, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// HasConflict checks if a command name conflicts with built-in commands
func (em *Manager) HasConflict(name string) bool {
	// If we have root command, get built-in commands dynamically
	if em.rootCommand != nil {
		for _, cmd := range em.rootCommand.Commands() {
			if cmd.Name() == name && !cmd.Hidden {
				return true
			}
		}
		return false
	}

	// Fallback to static list if no root command available
	builtinCommands := []string{
		"init", "run", "list", "exec", "migrate", "watch", "help", "version",
	}

	for _, builtin := range builtinCommands {
		if name == builtin {
			return true
		}
	}

	return false
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
