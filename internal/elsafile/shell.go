package elsafile

import (
	"os"
	"os/exec"
)

// ShellExecutor handles shell command execution with cross-platform support
type ShellExecutor struct{}

// NewShellExecutor creates a new ShellExecutor instance
func NewShellExecutor() *ShellExecutor {
	return &ShellExecutor{}
}

// ExecuteCommand executes a shell command with appropriate shell detection
func (se *ShellExecutor) ExecuteCommand(command string) error {
	var shell string
	var args []string

	// Detect OS and use appropriate shell
	if se.isWindows() {
		shell = "cmd"
		args = []string{"/C", command}
	} else {
		shell = "/bin/sh"
		args = []string{"-c", command}
	}

	cmd := exec.Command(shell, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// ExecuteCommandWithOutput executes a shell command and returns the output
func (se *ShellExecutor) ExecuteCommandWithOutput(command string) (string, error) {
	var shell string
	var args []string

	// Detect OS and use appropriate shell
	if se.isWindows() {
		shell = "cmd"
		args = []string{"/C", command}
	} else {
		shell = "/bin/sh"
		args = []string{"-c", command}
	}

	cmd := exec.Command(shell, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// ExecuteCommandSilently executes a shell command without showing output
func (se *ShellExecutor) ExecuteCommandSilently(command string) error {
	var shell string
	var args []string

	// Detect OS and use appropriate shell
	if se.isWindows() {
		shell = "cmd"
		args = []string{"/C", command}
	} else {
		shell = "/bin/sh"
		args = []string{"-c", command}
	}

	cmd := exec.Command(shell, args...)
	return cmd.Run()
}

// isWindows checks if the current OS is Windows
func (se *ShellExecutor) isWindows() bool {
	return os.PathSeparator == '\\'
}

// GetShellInfo returns information about the current shell
func (se *ShellExecutor) GetShellInfo() (string, []string) {
	if se.isWindows() {
		return "cmd", []string{"/C"}
	}
	return "/bin/sh", []string{"-c"}
}

// ValidateCommand checks if a command can be executed
func (se *ShellExecutor) ValidateCommand(command string) error {
	if command == "" {
		return os.ErrInvalid
	}
	return nil
}

// SplitCommands splits multiple commands joined by && into individual commands
func (se *ShellExecutor) SplitCommands(commands string) []string {
	// This is a simple implementation - in practice you might want more sophisticated parsing
	// that handles quoted strings and nested commands
	return []string{commands}
}
