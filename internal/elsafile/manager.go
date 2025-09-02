package elsafile

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/risoftinc/elsa/constants"
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
		return fmt.Errorf(constants.ErrElsafileNotFound, em.filepath)
	}

	file, err := os.Open(em.filepath)
	if err != nil {
		return fmt.Errorf(constants.ErrFailedToOpenFile, err)
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
		if line == "" || strings.HasPrefix(line, constants.CommentPrefix) {
			continue
		}

		// Check if this is a command definition (ends with :)
		if strings.HasSuffix(line, constants.CommandSuffix) {
			// Save previous command if exists
			if currentCommand != nil && len(currentCommand.Commands) > 0 {
				em.commands[currentCommand.Name] = currentCommand
			}

			// Start new command
			commandName := strings.TrimSuffix(line, constants.CommandSuffix)
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
		return fmt.Errorf(constants.ErrFailedToReadFile, err)
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
		return fmt.Errorf(constants.ErrCommandNotFound, name)
	}

	fmt.Printf("%s Running Elsafile command: %s\n", constants.RocketEmoji, name)
	fmt.Printf("%s Executing: %s\n\n", constants.PencilEmoji, strings.Join(command.Commands, constants.CommandSeparator))

	// Join all commands with && to execute them sequentially
	fullCommand := strings.Join(command.Commands, constants.CommandSeparator)
	return em.ExecuteShellCommand(fullCommand)
}

// ExecuteShellCommand executes a shell command
func (em *Manager) ExecuteShellCommand(command string) error {
	var shell string
	var args []string

	// Detect OS and use appropriate shell
	if os.PathSeparator == '\\' {
		// Windows
		shell = constants.WindowsShell
		args = []string{constants.WindowsShellArgs, command}
	} else {
		// Unix-like systems
		shell = constants.UnixShell
		args = []string{constants.UnixShellArgs, command}
	}

	cmd := exec.Command(shell, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Println(err.Error())
	}

	return nil
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
	builtinCommands := strings.Split(constants.BuiltinCommands, ",")

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
