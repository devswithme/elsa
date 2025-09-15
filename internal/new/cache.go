package new

import (
	"os"
	"path/filepath"
	"time"

	"go.risoftinc.com/elsa/constants"
	"go.risoftinc.com/elsa/internal/cache"
)

// getCacheDir returns the appropriate cache directory for the current platform
func getCacheDir() string {
	return cache.GetTemplatesCacheDir()
}

// GetCacheDir returns the cache directory path
func (tm *TemplateManager) GetCacheDir() string {
	return tm.cacheDir
}

// getCachedTemplatePath returns the path to the cached template
func (tm *TemplateManager) getCachedTemplatePath(templateName, version string) string {
	if version == "" {
		version = "latest"
	}
	return filepath.Join(tm.cacheDir, templateName, version)
}

// isCacheExpired checks if the cached template is older than TTL
func (tm *TemplateManager) isCacheExpired(cachedPath string) bool {
	// Check if directory exists
	if _, err := os.Stat(cachedPath); os.IsNotExist(err) {
		return true
	}

	// Check if .git directory exists (indicates it's a valid git repository)
	gitPath := filepath.Join(cachedPath, ".git")
	if _, err := os.Stat(gitPath); os.IsNotExist(err) {
		return true
	}

	// Check modification time
	info, err := os.Stat(cachedPath)
	if err != nil {
		return true
	}

	// Check if cache is older than TTL
	expiryTime := info.ModTime().Add(time.Duration(constants.NewCacheTTLHours) * time.Hour)
	return time.Now().After(expiryTime)
}

// ClearCache clears the template cache
func (tm *TemplateManager) ClearCache() error {
	if _, err := os.Stat(tm.cacheDir); os.IsNotExist(err) {
		return nil // Cache directory doesn't exist
	}
	return os.RemoveAll(tm.cacheDir)
}

// GetCacheSize returns the size of the cache directory
func (tm *TemplateManager) GetCacheSize() (int64, error) {
	var size int64
	err := filepath.Walk(tm.cacheDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}
