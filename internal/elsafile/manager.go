package elsafile

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"go.risoftinc.com/elsa/constants"
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
	var currentLine strings.Builder

	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), " \t") // Trim only trailing spaces/tabs

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(strings.TrimSpace(line), constants.CommentPrefix) {
			continue
		}

		// Check if this is a command definition (ends with :)
		if strings.HasSuffix(line, constants.CommandSuffix) {
			// Save previous command if exists
			if currentCommand != nil {
				if currentLine.Len() > 0 {
					parsedCommands := parseCommandLine(currentLine.String())
					currentCommand.Commands = append(currentCommand.Commands, parsedCommands...)
				}
				em.commands[currentCommand.Name] = currentCommand
			}

			// Start new command
			commandName := strings.TrimSuffix(line, constants.CommandSuffix)
			currentCommand = &Command{
				Name:     commandName,
				Commands: []string{},
			}
			currentLine.Reset()
		} else if currentCommand != nil {
			// Check for line continuation with backslash
			if strings.HasSuffix(line, "\\") {
				// Remove the backslash and add to current line
				lineWithoutBackslash := strings.TrimSuffix(line, "\\")
				if currentLine.Len() > 0 {
					currentLine.WriteString(" ")
				}
				currentLine.WriteString(strings.TrimSpace(lineWithoutBackslash))
			} else {
				// Complete line, add to current line and process
				if currentLine.Len() > 0 {
					currentLine.WriteString(" ")
				}
				currentLine.WriteString(strings.TrimSpace(line))

				// Parse and add the complete command
				parsedCommands := parseCommandLine(currentLine.String())
				currentCommand.Commands = append(currentCommand.Commands, parsedCommands...)
				currentLine.Reset()
			}
		}
	}

	// Save last command if there's remaining content
	if currentCommand != nil {
		if currentLine.Len() > 0 {
			parsedCommands := parseCommandLine(currentLine.String())
			currentCommand.Commands = append(currentCommand.Commands, parsedCommands...)
		}
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

	// Check if we have a single command that contains && (should be executed as single shell command)
	if len(command.Commands) == 1 && strings.Contains(command.Commands[0], "&&") {
		// Execute as single shell command
		return em.ExecuteShellCommand(command.Commands[0])
	}

	// Execute each command sequentially (legacy behavior)
	for _, cmd := range command.Commands {
		if err := em.ExecuteShellCommand(cmd); err != nil {
			return err
		}
	}
	return nil
}

// ExecuteShellCommand executes a shell command
func (em *Manager) ExecuteShellCommand(command string) error {
	// Substitute variables in the command
	substitutedCommand := substituteVariables(command)
	// Parse command properly handling quotes
	parts := parseCommandArgs(substitutedCommand)

	var shell string
	var args []string

	// Detect OS and use appropriate shell
	if strings.HasPrefix(substitutedCommand, "elsa ") {
		// Use absolute path to elsa executable
		if sh, err := os.Executable(); err != nil {
			shell = "elsa" // fallback
		} else {
			shell = sh
		}
		args = parts[1:]
	} else if os.PathSeparator == '\\' {
		// Windows
		shell = constants.WindowsShell
		args = append([]string{constants.WindowsShellArgs}, parts...)
	} else {
		// Unix-like systems
		shell = constants.UnixShell
		args = append([]string{constants.UnixShellArgs}, parts...)
	}

	cmd := exec.Command(shell, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Inherit all environment variables including those set by os.Setenv()
	cmd.Env = os.Environ()

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

// parseCommandLine parses a command line, handling quoted strings properly
// It detects whether commands should be executed as a single shell command or separately
func parseCommandLine(line string) []string {
	// Check if the line contains && - if so, treat as single shell command
	if strings.Contains(line, "&&") {
		// Return the entire line as a single command to be executed by shell
		return []string{strings.TrimSpace(line)}
	}
	
	// For lines without &&, parse as separate commands (legacy behavior)
	var commands []string
	var current strings.Builder
	var inQuotes bool
	var quoteChar rune

	runes := []rune(line)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch r {
		case '"', '\'':
			if !inQuotes {
				// Start of quoted string
				inQuotes = true
				quoteChar = r
				current.WriteRune(r) // Keep the opening quote
			} else if r == quoteChar {
				// End of quoted string
				inQuotes = false
				quoteChar = 0
				current.WriteRune(r) // Keep the closing quote
			} else {
				// Different quote character inside quotes, treat as literal
				current.WriteRune(r)
			}
		case '&':
			if !inQuotes && i+1 < len(runes) && runes[i+1] == '&' {
				// Found && separator
				if current.Len() > 0 {
					commands = append(commands, strings.TrimSpace(current.String()))
					current.Reset()
				}
				i++ // Skip the next &
				continue
			} else {
				current.WriteRune(r)
			}
		default:
			current.WriteRune(r)
		}
	}

	// Add the last command if any
	if current.Len() > 0 {
		commands = append(commands, strings.TrimSpace(current.String()))
	}

	// If no commands were found, return the original line
	if len(commands) == 0 {
		return []string{line}
	}

	return commands
}

// substituteVariables replaces variables in command string with their values
func substituteVariables(command string) string {
	// Pattern for ${?VAR:prompt} syntax (interactive input with prompt)
	interactiveWithPromptPattern := regexp.MustCompile(`\$\{\?([^:}]+):([^}]+)\}`)
	// Pattern for ${?VAR} syntax (interactive input without prompt)
	interactivePattern := regexp.MustCompile(`\$\{\?([^}]+)\}`)
	// Pattern for ${VAR} syntax
	curlyPattern := regexp.MustCompile(`\$\{([^}]+)\}`)
	// Pattern for $VAR syntax (but not ${VAR})
	dollarPattern := regexp.MustCompile(`\$([A-Za-z_][A-Za-z0-9_]*)`)

	// First, handle interactive input with prompt ${?VAR:prompt}
	result := interactiveWithPromptPattern.ReplaceAllStringFunc(command, func(match string) string {
		// Extract variable name and prompt
		parts := strings.Split(match[3:len(match)-1], ":") // Remove ${? and }
		if len(parts) != 2 {
			return match
		}
		varName := parts[0]
		prompt := parts[1]

		// Check if already in environment
		if value := os.Getenv(varName); value != "" {
			return value
		}

		// Prompt user for input
		fmt.Printf("%s: ", prompt)
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Set environment variable for future use
		os.Setenv(varName, input)

		return input
	})

	// Then, handle interactive input without prompt ${?VAR}
	result = interactivePattern.ReplaceAllStringFunc(result, func(match string) string {
		varName := match[3 : len(match)-1] // Remove ${? and }

		// Check if already in environment
		if value := os.Getenv(varName); value != "" {
			return value
		}

		// Prompt user for input with default prompt
		fmt.Printf("Enter %s: ", varName)
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Set environment variable for future use
		os.Setenv(varName, input)

		return input
	})

	// Then, replace ${VAR} syntax
	result = curlyPattern.ReplaceAllStringFunc(result, func(match string) string {
		varName := match[2 : len(match)-1] // Remove ${ and }
		if value := os.Getenv(varName); value != "" {
			return value
		}
		// If not found in environment, return the original match
		return match
	})

	// Finally, replace $VAR syntax
	result = dollarPattern.ReplaceAllStringFunc(result, func(match string) string {
		varName := match[1:] // Remove $
		if value := os.Getenv(varName); value != "" {
			return value
		}
		// If not found in environment, return the original match
		return match
	})

	return result
}

// parseCommandArgs properly parses command arguments, handling quoted strings
func parseCommandArgs(command string) []string {
	var args []string
	var current strings.Builder
	var inQuotes bool
	var quoteChar rune

	runes := []rune(command)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch r {
		case '"', '\'':
			if !inQuotes {
				// Start of quoted string
				inQuotes = true
				quoteChar = r
				// Don't keep the opening quote - it's just for grouping
			} else if r == quoteChar {
				// End of quoted string
				inQuotes = false
				quoteChar = 0
				// Don't keep the closing quote - it's just for grouping
			} else {
				// Different quote character inside quotes, treat as literal
				current.WriteRune(r)
			}
		case ' ':
			if !inQuotes {
				// Space outside quotes - end of argument
				if current.Len() > 0 {
					args = append(args, current.String())
					current.Reset()
				}
			} else {
				// Space inside quotes - treat as literal
				current.WriteRune(r)
			}
		default:
			current.WriteRune(r)
		}
	}

	// Add the last argument if any
	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}
