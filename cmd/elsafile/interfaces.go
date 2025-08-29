package elsafile

import (
	"github.com/risoftinc/elsa/internal/elsafile"
)

// Re-export interfaces from internal package
type (
	ManagerInterface           = elsafile.ManagerInterface
	ParserInterface            = elsafile.ParserInterface
	ShellExecutorInterface     = elsafile.ShellExecutorInterface
	TemplateGeneratorInterface = elsafile.TemplateGeneratorInterface
	ConflictHandlerInterface   = elsafile.ConflictHandlerInterface
	SimpleHandlerInterface     = elsafile.SimpleHandlerInterface
	CommandListerInterface     = elsafile.CommandListerInterface
	UtilsInterface             = elsafile.UtilsInterface
	CommandValidatorInterface  = elsafile.CommandValidatorInterface
	FileManagerInterface       = elsafile.FileManagerInterface
	OutputFormatterInterface   = elsafile.OutputFormatterInterface
)
