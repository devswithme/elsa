package cache

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetCacheDir returns the appropriate cache directory for the current platform
// This ensures consistency across all elsa commands (new, make, etc.)
func GetCacheDir() string {
	// Get user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home directory is not accessible
		homeDir = "."
	}

	// Platform-specific cache directory
	switch runtime.GOOS {
	case "windows":
		// Windows: %USERPROFILE%\.elsa-cache
		return filepath.Join(homeDir, ".elsa-cache")
	case "darwin":
		// macOS: ~/Library/Caches/elsa
		return filepath.Join(homeDir, "Library", "Caches", "elsa")
	case "linux":
		// Linux: ~/.cache/elsa
		return filepath.Join(homeDir, ".cache", "elsa")
	default:
		// Fallback for other platforms
		return filepath.Join(homeDir, ".elsa-cache")
	}
}

// GetTemplatesCacheDir returns the templates cache directory
func GetTemplatesCacheDir() string {
	return filepath.Join(GetCacheDir(), "templates")
}

// GetFilestubCacheDir returns the filestub cache directory
func GetFilestubCacheDir() string {
	return filepath.Join(GetCacheDir(), "filestub")
}
