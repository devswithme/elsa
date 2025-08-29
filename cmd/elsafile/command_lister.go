package elsafile

import (
	"github.com/risoftinc/elsa/internal/elsafile"
)

// CommandLister handles listing and displaying Elsafile commands
type CommandLister = elsafile.CommandLister

// NewCommandLister creates a new CommandLister instance
func NewCommandLister() *CommandLister {
	return elsafile.NewCommandLister()
}

// NewCommandListerWithHandler creates a new CommandLister instance with a specific conflict handler
func NewCommandListerWithHandler(handler *ConflictHandler) *CommandLister {
	return elsafile.NewCommandListerWithHandler(handler)
}
