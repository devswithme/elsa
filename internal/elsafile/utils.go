package elsafile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/risoftinc/elsa/constants"
)

// Utils provides utility functions for Elsafile operations
type Utils struct{}

// NewUtils creates a new Utils instance
func NewUtils() *Utils {
	return &Utils{}
}

// FileExists checks if a file exists
func (u *Utils) FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil
}

// GetFileInfo returns information about a file
func (u *Utils) GetFileInfo(filepath string) (*ElsafileInfo, error) {
	info := &ElsafileInfo{
		FilePath: filepath,
	}

	if !u.FileExists(filepath) {
		return info, nil
	}

	// Get file stats
	stat, err := os.Stat(filepath)
	if err != nil {
		return nil, err
	}

	info.LastModified = stat.ModTime().Unix()
	info.IsValid = true

	// Parse file to get command count
	parser := NewParser()
	commands, err := parser.ParseFile(filepath)
	if err != nil {
		info.IsValid = false
		return info, nil
	}

	info.TotalCommands = len(commands)

	// Check for conflicts
	manager := NewManager(filepath)
	if err := manager.Load(); err != nil {
		info.IsValid = false
		return info, nil
	}

	conflicts := manager.GetConflictingCommands()
	info.Conflicts = len(conflicts)

	return info, nil
}

// FindElsafile searches for Elsafile in current and parent directories
func (u *Utils) FindElsafile() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Search in current directory and parent directories
	for {
		elsafilePath := filepath.Join(currentDir, constants.DefaultElsafileName)
		if u.FileExists(elsafilePath) {
			return elsafilePath, nil
		}

		// Move to parent directory
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// Reached root directory
			break
		}
		currentDir = parentDir
	}

	return "", fmt.Errorf(constants.ErrElsafileNotFoundInDirectories)
}

// ValidateCommandName validates a command name
func (u *Utils) ValidateCommandName(name string) error {
	if name == "" {
		return fmt.Errorf(constants.ErrCommandNameEmpty)
	}

	if strings.ContainsAny(name, constants.WhitespaceChars) {
		return fmt.Errorf(constants.ErrCommandNameWhitespace)
	}

	if strings.ContainsAny(name, constants.InvalidCommandChars) {
		return fmt.Errorf(constants.ErrCommandNameInvalidChars)
	}

	return nil
}

// SanitizeCommandName sanitizes a command name for safe use
func (u *Utils) SanitizeCommandName(name string) string {
	// Remove invalid characters
	invalidChars := strings.Split(constants.InvalidCommandChars+constants.WhitespaceChars, "")
	sanitized := name

	for _, char := range invalidChars {
		if char != "" {
			sanitized = strings.ReplaceAll(sanitized, char, "_")
		}
	}

	return sanitized
}

// FormatDuration formats duration in a human-readable format
func (u *Utils) FormatDuration(duration int64) string {
	if duration < constants.MillisecondThreshold {
		return fmt.Sprintf("%dms", duration)
	}

	seconds := duration / constants.MillisecondThreshold
	if seconds < constants.SecondThreshold {
		return fmt.Sprintf("%ds", seconds)
	}

	minutes := seconds / 60
	remainingSeconds := seconds % 60
	return fmt.Sprintf("%dm %ds", minutes, remainingSeconds)
}

// GetCurrentTimestamp returns current timestamp
func (u *Utils) GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

// IsValidCommand checks if a command is valid
func (u *Utils) IsValidCommand(cmd *Command) bool {
	if cmd == nil {
		return false
	}

	if err := u.ValidateCommandName(cmd.Name); err != nil {
		return false
	}

	if len(cmd.Commands) == 0 {
		return false
	}

	return true
}

// MergeCommands merges multiple command maps
func (u *Utils) MergeCommands(maps ...map[string]*Command) map[string]*Command {
	result := make(map[string]*Command)

	for _, m := range maps {
		for name, cmd := range m {
			result[name] = cmd
		}
	}

	return result
}
