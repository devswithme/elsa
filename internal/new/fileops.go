package new

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// copyTemplate copies the template to the project directory
func (tm *TemplateManager) copyTemplate(cachedPath, projectPath string) error {
	// Create project directory
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return err
	}

	// Copy all files except .git
	return tm.copyDirectory(cachedPath, projectPath, []string{".git"})
}

// copyDirectory recursively copies directory contents, excluding specified directories
func (tm *TemplateManager) copyDirectory(src, dst string, excludeDirs []string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip excluded directories
		for _, excludeDir := range excludeDirs {
			if strings.Contains(path, excludeDir) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// Calculate relative path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Copy file
		return tm.copyFile(path, dstPath, info.Mode())
	})
}

// copyFile copies a single file
func (tm *TemplateManager) copyFile(src, dst string, mode os.FileMode) error {
	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	// Read source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create destination file
	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy content
	_, err = dstFile.ReadFrom(srcFile)
	return err
}

// updateGoMod updates the go.mod file with the new module name
// Returns the original module name for import replacement
func (tm *TemplateManager) updateGoMod(projectPath, moduleName string) (string, error) {
	goModPath := filepath.Join(projectPath, "go.mod")

	// Check if go.mod exists
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return "", nil // No go.mod file to update
	}

	// Read go.mod content
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return "", err
	}

	// Find and replace module line
	lines := strings.Split(string(content), "\n")
	var originalModule string

	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "module ") {
			// Extract original module name
			originalModule = strings.TrimSpace(strings.TrimPrefix(trimmedLine, "module "))
			// Replace with new module name
			lines[i] = "module " + moduleName
			break
		}
	}

	// Write updated content
	updatedContent := strings.Join(lines, "\n")
	err = os.WriteFile(goModPath, []byte(updatedContent), 0644)
	return originalModule, err
}

// updateImports updates all import statements in Go files
func (tm *TemplateManager) updateImports(projectPath, originalModule, newModule string) error {
	return filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-Go files
		if info.IsDir() || (!strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, ".proto")) {
			return nil
		}

		// Skip vendor directory
		if strings.Contains(path, "vendor") {
			return nil
		}

		// Update imports in this file
		return tm.updateFileImports(path, originalModule, newModule)
	})
}

// updateFileImports updates import statements in a single Go file
func (tm *TemplateManager) updateFileImports(filePath, originalModule, newModule string) error {
	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Convert to string and update imports
	fileContent := string(content)
	updatedContent := tm.replaceImports(fileContent, originalModule, newModule)

	// Only write if content changed
	if updatedContent != fileContent {
		return os.WriteFile(filePath, []byte(updatedContent), 0644)
	}

	return nil
}

// replaceImports replaces import statements in Go file content
func (tm *TemplateManager) replaceImports(content, originalModule, newModule string) string {
	lines := strings.Split(content, "\n")
	var updatedLines []string

	for _, line := range lines {
		updatedLine := line

		// Check if line contains the original module name
		if strings.Contains(line, originalModule) {
			updatedLine = tm.replaceImportLine(line, originalModule, newModule)
		}

		updatedLines = append(updatedLines, updatedLine)
	}

	return strings.Join(updatedLines, "\n")
}

// replaceImportLine replaces import paths in a single line
func (tm *TemplateManager) replaceImportLine(line, originalModule, newModule string) string {
	// Simple string replacement for the original module name
	updatedLine := strings.Replace(line, originalModule, newModule, -1)
	return updatedLine
}

// hasProtoFiles checks if the project contains .proto files
func (tm *TemplateManager) hasProtoFiles(projectPath string) bool {
	hasProto := false
	filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".proto") {
			hasProto = true
			return filepath.SkipDir // Stop walking once we find a proto file
		}
		return nil
	})
	return hasProto
}

// isProtocInstalled checks if protoc is installed and available
func (tm *TemplateManager) isProtocInstalled() bool {
	cmd := exec.Command("protoc", "--version")
	err := cmd.Run()
	return err == nil
}

// generateProtoFiles runs protoc to generate Go files from .proto files
func (tm *TemplateManager) generateProtoFiles(projectPath string) error {
	// Find all .proto files
	var protoFiles []string
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".proto") {
			protoFiles = append(protoFiles, path)
		}
		return nil
	})

	if err != nil {
		return err
	}

	if len(protoFiles) == 0 {
		return nil // No proto files found
	}

	// Run protoc for each .proto file
	for _, protoFile := range protoFiles {
		// Get the directory containing the proto file
		protoDir := filepath.Dir(protoFile)
		protoFileName := filepath.Base(protoFile)

		// Run protoc command
		cmd := exec.Command("protoc",
			"--go_out=.",
			"--go_opt=paths=source_relative",
			"--go-grpc_out=.",
			"--go-grpc_opt=paths=source_relative",
			protoFileName)
		cmd.Dir = protoDir

		// Capture output for better error reporting
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to generate proto files for %s: %v\nOutput: %s", protoFileName, err, string(output))
		}
	}

	return nil
}

// runGoModDownload runs 'go mod download' in the project directory
func (tm *TemplateManager) runGoModDownload(projectPath string) error {
	// Check if go.mod exists
	goModPath := filepath.Join(projectPath, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return nil // No go.mod file, skip download
	}

	// Run go mod download
	cmd := exec.Command("go", "mod", "download")
	cmd.Dir = projectPath

	// Capture output for better error reporting
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go mod download failed: %v\nOutput: %s", err, string(output))
	}

	return nil
}

// runGoModTidy runs 'go mod tidy' in the project directory
func (tm *TemplateManager) runGoModTidy(projectPath string) error {
	// Check if go.mod exists
	goModPath := filepath.Join(projectPath, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return nil // No go.mod file, skip tidy
	}

	// Run go mod tidy
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectPath

	// Capture output for better error reporting
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go mod tidy failed: %v\nOutput: %s", err, string(output))
	}

	return nil
}
