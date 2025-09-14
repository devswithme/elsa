package make

import (
	"os"
	"path/filepath"
	"strings"
)

// gitURLToPath converts git URL to filesystem-safe path
func gitURLToPath(gitURL string) string {
	// Remove protocol prefix (https://, git@, etc.)
	url := strings.TrimPrefix(gitURL, "https://")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "git@")
	url = strings.TrimSuffix(url, ".git")

	// Replace : with / for SSH URLs (git@github.com:user/repo -> github.com/user/repo)
	url = strings.ReplaceAll(url, ":", "/")

	return url
}

// getFilestubCacheDir returns the filestub cache directory for a git URL
func getFilestubCacheDir(gitURL string) string {
	homeDir, _ := os.UserHomeDir()
	urlPath := gitURLToPath(gitURL)
	return filepath.Join(homeDir, ".elsa-cache", "filestub", urlPath)
}

// getCacheDir returns the cache directory path
func getCacheDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".elsa-cache", "templates")
}
