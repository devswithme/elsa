package elsafile

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"go.risoftinc.com/elsa/constants"
)

// Parser handles parsing of Elsafile content
type Parser struct{}

// NewParser creates a new Parser instance
func NewParser() *Parser {
	return &Parser{}
}

// ParseFile parses the Elsafile content and returns a map of commands
func (p *Parser) ParseFile(filepath string) (map[string]*Command, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	return p.parseFileContent(file)
}

// ParseContent parses Elsafile content from a string
func (p *Parser) ParseContent(content string) (map[string]*Command, error) {
	lines := strings.Split(content, "\n")
	return p.parseLines(lines)
}

// parseFileContent parses the file content line by line
func (p *Parser) parseFileContent(file *os.File) (map[string]*Command, error) {
	scanner := bufio.NewScanner(file)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return p.parseLines(lines)
}

// parseLines parses lines and builds command map
func (p *Parser) parseLines(lines []string) (map[string]*Command, error) {
	commands := make(map[string]*Command)
	var currentCommand *Command
	var currentLine strings.Builder

	for _, line := range lines {
		line = strings.TrimRight(line, " \t") // Trim only trailing spaces/tabs

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(strings.TrimSpace(line), constants.CommentPrefix) {
			continue
		}

		// Check if this is a command definition (ends with :)
		if strings.HasSuffix(line, constants.CommandSuffix) {
			// Save previous command if exists
			if currentCommand != nil && currentLine.Len() > 0 {
				parsedCommands := parseCommandLine(currentLine.String())
				currentCommand.Commands = append(currentCommand.Commands, parsedCommands...)
				commands[currentCommand.Name] = currentCommand
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
	if currentCommand != nil && currentLine.Len() > 0 {
		parsedCommands := parseCommandLine(currentLine.String())
		currentCommand.Commands = append(currentCommand.Commands, parsedCommands...)
		commands[currentCommand.Name] = currentCommand
	} else if currentCommand != nil && len(currentCommand.Commands) > 0 {
		commands[currentCommand.Name] = currentCommand
	}

	return commands, nil
}

// ValidateCommand validates a parsed command
func (p *Parser) ValidateCommand(cmd *Command) error {
	if cmd == nil {
		return fmt.Errorf("command is nil")
	}

	if cmd.Name == "" {
		return fmt.Errorf(constants.ErrCommandNameEmpty)
	}

	if len(cmd.Commands) == 0 {
		return fmt.Errorf("command '%s' has no commands to execute", cmd.Name)
	}

	return nil
}

// ValidateCommands validates all commands in a map
func (p *Parser) ValidateCommands(commands map[string]*Command) []error {
	var errors []error

	for name, cmd := range commands {
		if err := p.ValidateCommand(cmd); err != nil {
			errors = append(errors, fmt.Errorf("command '%s': %v", name, err))
		}
	}

	return errors
}

// GetCommandNames returns a slice of command names
func (p *Parser) GetCommandNames(commands map[string]*Command) []string {
	var names []string
	for name := range commands {
		names = append(names, name)
	}
	return names
}

// FilterCommandsByPrefix filters commands by a prefix
func (p *Parser) FilterCommandsByPrefix(commands map[string]*Command, prefix string) map[string]*Command {
	filtered := make(map[string]*Command)

	for name, cmd := range commands {
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
			filtered[name] = cmd
		}
	}

	return filtered
}
