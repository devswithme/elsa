package elsafile

import (
	"github.com/risoftinc/elsa/internal/elsafile"
)

// Parser handles parsing of Elsafile content
type Parser = elsafile.Parser

// NewParser creates a new Parser instance
func NewParser() *Parser {
	return elsafile.NewParser()
}
