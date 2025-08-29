package constants

// Shell execution constants
const (
	// WindowsShell is the shell command for Windows systems
	WindowsShell = "cmd"

	// WindowsShellArgs are the arguments for Windows shell execution
	WindowsShellArgs = "/C"

	// UnixShell is the shell command for Unix-like systems
	UnixShell = "/bin/sh"

	// UnixShellArgs are the arguments for Unix shell execution
	UnixShellArgs = "-c"
)

// Invalid character constants for command names
const (
	// InvalidCommandChars contains characters that are not allowed in command names
	InvalidCommandChars = ":/\\*?\"<>|"

	// WhitespaceChars contains whitespace characters that are not allowed in command names
	WhitespaceChars = " \t\n\r"
)
