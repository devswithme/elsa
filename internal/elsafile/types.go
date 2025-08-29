package elsafile

// Command represents a command defined in the Elsafile
// This type is defined in manager.go to avoid duplication

// CommandResult represents the result of executing a command
type CommandResult struct {
	CommandName string
	Success     bool
	Output      string
	Error       error
	Duration    int64 // in milliseconds
}

// ConflictInfo represents information about command conflicts
type ConflictInfo struct {
	CommandName    string
	BuiltinCommand string
	Resolution     string
}

// ElsafileInfo represents information about the Elsafile
type ElsafileInfo struct {
	FilePath      string
	TotalCommands int
	Conflicts     int
	LastModified  int64
	IsValid       bool
}

// ExecutionOptions represents options for command execution
type ExecutionOptions struct {
	Silent        bool
	CaptureOutput bool
	Timeout       int64 // in seconds
	WorkingDir    string
	Environment   map[string]string
}

// ParseOptions represents options for parsing
type ParseOptions struct {
	SkipComments bool
	SkipEmpty    bool
	Validate     bool
}

// TemplateData represents data for template generation
type TemplateData struct {
	ProjectName     string
	Author          string
	Description     string
	DefaultCommands []string
	CustomCommands  map[string]string
}
