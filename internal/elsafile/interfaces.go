package elsafile

// ManagerInterface defines the interface for managing Elsafile operations
type ManagerInterface interface {
	Load() error
	GetCommand(name string) (*Command, bool)
	ListCommands() map[string]*Command
	ExecuteCommand(name string) error
	ExecuteShellCommand(command string) error
	HasConflict(name string) bool
	GetConflictingCommands() []string
}

// ParserInterface defines the interface for parsing Elsafile content
type ParserInterface interface {
	ParseFile(filepath string) (map[string]*Command, error)
	ParseContent(content string) (map[string]*Command, error)
	ValidateCommand(cmd *Command) error
	ValidateCommands(commands map[string]*Command) []error
}

// ShellExecutorInterface defines the interface for executing shell commands
type ShellExecutorInterface interface {
	ExecuteCommand(command string) error
	ExecuteCommandWithOutput(command string) (string, error)
	ExecuteCommandSilently(command string) error
	ValidateCommand(command string) error
}

// TemplateGeneratorInterface defines the interface for generating templates
type TemplateGeneratorInterface interface {
	CreateDefaultElsafile() error
	CreateCustomElsafile(content string) error
	GetDefaultTemplate() string
	GetSuccessMessage() string
}

// ConflictHandlerInterface defines the interface for handling command conflicts
type ConflictHandlerInterface interface {
	ExecuteElsafileCommand(commandName string) error
	ListCommands() (map[string]*Command, error)
	GetConflictingCommands() ([]string, error)
	HasConflict(commandName string) bool
	GetConflictMessage(commandName string) string
}

// SimpleHandlerInterface defines the interface for handling unknown commands
type SimpleHandlerInterface interface {
	HandleUnknownCommand(commandName string) error
	SuggestCommands(commandName string) []string
	GetSuggestionMessage(commandName string) string
}

// CommandListerInterface defines the interface for listing commands
type CommandListerInterface interface {
	ListAllCommands() error
	ListConflictingCommands() error
	GetCommandSummary() (map[string]string, error)
	GetConflictSummary() ([]string, error)
}

// UtilsInterface defines the interface for utility functions
type UtilsInterface interface {
	FileExists(filepath string) bool
	GetFileInfo(filepath string) (*ElsafileInfo, error)
	FindElsafile() (string, error)
	ValidateCommandName(name string) error
	SanitizeCommandName(name string) string
	IsValidCommand(cmd *Command) bool
}

// CommandValidatorInterface defines the interface for validating commands
type CommandValidatorInterface interface {
	ValidateCommand(cmd *Command) error
	ValidateCommandName(name string) error
	ValidateCommands(commands map[string]*Command) []error
}

// FileManagerInterface defines the interface for file operations
type FileManagerInterface interface {
	CreateFile(filepath string, content string) error
	ReadFile(filepath string) (string, error)
	WriteFile(filepath string, content string) error
	DeleteFile(filepath string) error
	FileExists(filepath string) bool
}

// OutputFormatterInterface defines the interface for formatting output
type OutputFormatterInterface interface {
	FormatCommand(cmd *Command) string
	FormatCommandList(commands map[string]*Command) string
	FormatConflictList(conflicts []string) string
	FormatSuccessMessage(message string) string
	FormatErrorMessage(message string) string
}
