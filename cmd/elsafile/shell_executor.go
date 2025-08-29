package elsafile

import (
	"github.com/risoftinc/elsa/internal/elsafile"
)

// ShellExecutor handles shell command execution with cross-platform support
type ShellExecutor = elsafile.ShellExecutor

// NewShellExecutor creates a new ShellExecutor instance
func NewShellExecutor() *ShellExecutor {
	return elsafile.NewShellExecutor()
}
