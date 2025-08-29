package elsafile

import (
	"github.com/risoftinc/elsa/internal/elsafile"
)

// Formatter handles formatting of output messages
type Formatter = elsafile.Formatter

// NewFormatter creates a new Formatter instance
func NewFormatter() *Formatter {
	return elsafile.NewFormatter()
}
