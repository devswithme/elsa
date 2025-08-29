package elsafile

import (
	"github.com/risoftinc/elsa/internal/elsafile"
)

// SimpleHandler handles unknown commands by checking Elsafile
type SimpleHandler = elsafile.SimpleHandler

// NewSimpleHandler creates a new SimpleHandler instance
func NewSimpleHandler() *SimpleHandler {
	return elsafile.NewSimpleHandler()
}

// NewSimpleHandlerWithRoot creates a new SimpleHandler instance with root command for dynamic built-in detection
func NewSimpleHandlerWithRoot(rootCmd interface{}) *SimpleHandler {
	return elsafile.NewSimpleHandlerWithRoot(rootCmd)
}

// NewSimpleHandlerWithManager creates a new SimpleHandler instance with a specific manager
func NewSimpleHandlerWithManager(manager *elsafile.Manager) *SimpleHandler {
	return elsafile.NewSimpleHandlerWithManager(manager)
}
