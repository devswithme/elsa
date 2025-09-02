package generate

import (
	"fmt"
	"os"
	"path/filepath"
)

// FindElsabuildFiles searches for files with elsabuild build tags in the specified directory
// Recursively walks through the directory tree to find Go files containing the elsabuild tag
// Returns a slice of relative paths to files with the specified build tag
func (g *Generator) FindElsabuildFiles(targetDir string) ([]string, error) {
	searchDir, err := resolvePath(targetDir)
	if err != nil {
		return nil, err
	}

	foundFiles := []string{}

	err = safeWalk(searchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-Go files
		if info.IsDir() || !isGoFile(path) {
			return nil
		}

		// Read file content to check for elsabuild tag
		content, err := os.ReadFile(path)
		if err != nil {
			return nil // Skip files we can't read
		}

		// Check for elsabuild build tag
		if hasBuildTag(string(content), "elsabuild") {
			// Get relative path from search directory
			relPath, err := filepath.Rel(searchDir, path)
			if err != nil {
				relPath = path
			}
			foundFiles = append(foundFiles, relPath)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory: %v", err)
	}

	return foundFiles, nil
}
