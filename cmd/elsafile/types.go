package elsafile

import (
	"github.com/risoftinc/elsa/internal/elsafile"
)

// Re-export types from internal package
type (
	CommandResult    = elsafile.CommandResult
	ConflictInfo     = elsafile.ConflictInfo
	ElsafileInfo     = elsafile.ElsafileInfo
	ExecutionOptions = elsafile.ExecutionOptions
	ParseOptions     = elsafile.ParseOptions
	TemplateData     = elsafile.TemplateData
)
