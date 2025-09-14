package make

import (
	"os"
	"path/filepath"
	"strings"
)

// ProjectConfig represents the configuration for a project
type ProjectConfig struct {
	Source SourceInfo            `yaml:"source"`
	Make   map[string]MakeConfig `yaml:"make"`
}

// SourceInfo contains source template information
type SourceInfo struct {
	Name      string `yaml:"name"`
	Version   string `yaml:"version"`
	GitURL    string `yaml:"git_url"`
	GitCommit string `yaml:"git_commit"`
}

// MakeConfig represents configuration for a specific make type
type MakeConfig struct {
	Template string `yaml:"template"`
	Output   string `yaml:"output"`
}

// TemplateData contains data for template generation
type TemplateData struct {
	PackageName string
	StructName  string
	FileName    string
	FolderPath  string
	OutputPath  string
	Imports     []string
	Fields      []Field
	Methods     []Method
}

// Field represents a struct field
type Field struct {
	Name string
	Type string
	Tag  string
}

// Method represents a method
type Method struct {
	Name    string
	Params  []Param
	Returns []Return
}

// Param represents a method parameter
type Param struct {
	Name string
	Type string
}

// Return represents a method return value
type Return struct {
	Name string
	Type string
}

// TemplateManager handles template operations for make commands
type TemplateManager struct {
	cacheDir string
}

// NewTemplateManager creates a new template manager
func NewTemplateManager() *TemplateManager {
	cacheDir := getCacheDir()
	return &TemplateManager{
		cacheDir: cacheDir,
	}
}

// ParseTemplateData parses the name and creates template data
func (tm *TemplateManager) ParseTemplateData(name string, config MakeConfig) *TemplateData {
	// Extract folder and file info
	folderPath, fileName := tm.extractFolderAndFile(name)
	structName := tm.extractStructName(fileName)

	// Build output path
	var outputPath string
	if folderPath != "" {
		outputPath = filepath.Join(config.Output, folderPath, fileName+".go")
	} else {
		outputPath = filepath.Join(config.Output, fileName+".go")
	}

	// Get package name - check existing files first, then fallback to extracted name
	outputDir := filepath.Dir(outputPath)
	packageName := tm.getExistingPackageName(outputDir)
	if packageName == "" {
		packageName = tm.extractPackageName(fileName)
	}

	return &TemplateData{
		PackageName: packageName,
		StructName:  structName,
		FileName:    fileName,
		FolderPath:  folderPath,
		OutputPath:  outputPath,
		Imports:     []string{},
		Fields:      []Field{},
		Methods:     []Method{},
	}
}

// extractFolderAndFile extracts folder path and file name from the input
func (tm *TemplateManager) extractFolderAndFile(name string) (string, string) {
	// health/health_repository -> (health, health_repository)
	// health_repository -> ("", health_repository)
	// UserService -> ("", user_service)

	if strings.Contains(name, "/") {
		parts := strings.Split(name, "/")
		if len(parts) == 2 {
			// Convert folder and file to snake_case
			folder := tm.toSnakeCase(parts[0])
			file := tm.toSnakeCase(parts[1])
			return folder, file
		}
	}

	// Convert to snake_case
	return "", tm.toSnakeCase(name)
}

// extractPackageName extracts package name from file name or existing Go files
func (tm *TemplateManager) extractPackageName(fileName string) string {
	// health_repository -> health
	// user_repository -> user
	parts := strings.Split(fileName, "_")
	return parts[0]
}

// getExistingPackageName reads package name from existing Go files in the directory
func (tm *TemplateManager) getExistingPackageName(outputDir string) string {
	// Read all .go files in the output directory
	files, err := filepath.Glob(filepath.Join(outputDir, "*.go"))
	if err != nil || len(files) == 0 {
		return ""
	}

	// Read the first .go file to get package name
	content, err := os.ReadFile(files[0])
	if err != nil {
		return ""
	}

	// Find package declaration
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "package ") {
			// Extract package name after "package "
			packageName := strings.TrimSpace(strings.TrimPrefix(line, "package "))
			// Remove comments if any
			if commentIndex := strings.Index(packageName, "//"); commentIndex != -1 {
				packageName = strings.TrimSpace(packageName[:commentIndex])
			}
			return packageName
		}
	}

	return ""
}

// extractStructName extracts struct name from file name
func (tm *TemplateManager) extractStructName(fileName string) string {
	// health_repository -> Health
	// user_repository -> User
	parts := strings.Split(fileName, "_")
	if len(parts) == 0 {
		return ""
	}
	// Convert first part to title case manually
	first := parts[0]
	if len(first) == 0 {
		return ""
	}
	return strings.ToUpper(string(first[0])) + strings.ToLower(first[1:])
}

// toSnakeCase converts PascalCase to snake_case
func (tm *TemplateManager) toSnakeCase(s string) string {
	// Convert PascalCase to snake_case
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}

// getCacheDir returns the cache directory path
func getCacheDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".elsa", "cache")
}
