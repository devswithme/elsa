package elsafile

import (
	"github.com/risoftinc/elsa/internal/elsafile"
)

// ConflictHandler handles command conflicts between built-in and Elsafile commands
type ConflictHandler = elsafile.ConflictHandler

// NewConflictHandler creates a new ConflictHandler instance
func NewConflictHandler() *ConflictHandler {
	return elsafile.NewConflictHandler()
}

// NewConflictHandlerWithManager creates a new ConflictHandler instance with a specific manager
func NewConflictHandlerWithManager(manager *Manager) *ConflictHandler {
	return elsafile.NewConflictHandlerWithManager(manager)
}
