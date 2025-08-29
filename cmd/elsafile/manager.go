package elsafile

import (
	"github.com/risoftinc/elsa/internal/elsafile"
	"github.com/spf13/cobra"
)

// Manager handles parsing and execution of Elsafile commands
type Manager = elsafile.Manager

// Command represents a command defined in the Elsafile
type Command = elsafile.Command

// NewManager creates a new Manager instance
func NewManager(filepath string) *Manager {
	return elsafile.NewManager(filepath)
}

// NewManagerWithRoot creates a new Manager instance with root command for dynamic built-in detection
func NewManagerWithRoot(filepath string, rootCmd interface{}) *Manager {
	if cmd, ok := rootCmd.(*cobra.Command); ok {
		return elsafile.NewManagerWithRoot(filepath, cmd)
	}
	return elsafile.NewManager(filepath)
}
