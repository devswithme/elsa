package elsafile

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

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
