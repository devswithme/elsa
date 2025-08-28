package elsafile

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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
