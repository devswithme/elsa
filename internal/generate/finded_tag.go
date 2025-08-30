package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FindElsabuildFiles searches for files with elsabuild build tags in the specified directory
func (g *Generator) FindElsabuildFiles(targetDir string) ([]string, error) {
	var searchDir string
	var err error

	if targetDir != "" {
		// If target directory is provided, resolve it relative to current directory
		if filepath.IsAbs(targetDir) {
			searchDir = targetDir
		} else {
			currentDir, err := os.Getwd()
			if err != nil {
				return []string{}, fmt.Errorf("failed to get current directory: %v", err)
			}
			searchDir = filepath.Join(currentDir, targetDir)
		}

		// Check if target directory exists
		if _, err := os.Stat(searchDir); os.IsNotExist(err) {
			return []string{}, fmt.Errorf("target directory does not exist: %s", targetDir)
		}
	} else {
		// Use current directory if no target specified
		searchDir, err = os.Getwd()
		if err != nil {
			return []string{}, fmt.Errorf("failed to get current directory: %v", err)
		}
	}

	foundFiles := []string{}
	err = filepath.Walk(searchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-Go files
		if info.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".go") {
			return nil
		}

		// Read file content to check for elsabuild tag
		content, err := os.ReadFile(path)
		if err != nil {
			return nil // Skip files we can't read
		}

		// Check for elsabuild build tag
		fileContent := string(content)
		if strings.Contains(fileContent, "//go:build elsabuild") ||
			strings.Contains(fileContent, "// +build elsabuild") {
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
		return []string{}, fmt.Errorf("error walking directory: %v", err)
	}

	if len(foundFiles) == 0 {
		return []string{}, nil
	}

	return foundFiles, nil
}
