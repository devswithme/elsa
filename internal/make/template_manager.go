package make

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

// LoadProjectConfig loads the project configuration from .elsa-config.yaml
func (tm *TemplateManager) LoadProjectConfig(projectPath string) (*ProjectConfig, error) {
	configPath := filepath.Join(projectPath, ".elsa-config.yaml")

	if _, err := os.Stat(configPath); err != nil {
		return nil, fmt.Errorf("config file not found: %s", configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config ProjectConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("invalid config format: %v", err)
	}

	return &config, nil
}

// GenerateFile generates a file from template
func (tm *TemplateManager) GenerateFile(templateType, name string) error {
	// Load project config
	config, err := tm.LoadProjectConfig(".")
	if err != nil {
		return err
	}

	// Check if template type exists in config
	templateConfig, exists := config.Make[templateType]
	if !exists {
		// Show available template types
		availableTypes := make([]string, 0, len(config.Make))
		for t := range config.Make {
			availableTypes = append(availableTypes, t)
		}

		return fmt.Errorf("template type '%s' not found. Available types: %s",
			templateType, strings.Join(availableTypes, ", "))
	}

	// Check if output directory is configured, if not ask user
	if templateConfig.Output == "" {
		outputDir, err := tm.promptOutputDirectory(templateType)
		if err != nil {
			return err
		}
		templateConfig.Output = outputDir
	}

	// Parse template data
	data := tm.ParseTemplateData(name, templateConfig)

	// Create output directory if needed
	outputDir := filepath.Dir(data.OutputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory '%s': %v", outputDir, err)
	}

	// Resolve template path using template config
	templatePath := tm.resolveTemplatePath(templateConfig, config.Source)
	fmt.Printf("🔍 Using template path: %s\n", templatePath)

	// Load and execute template
	tmpl, err := tm.loadTemplate(templatePath)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	// Write to file
	if err := tm.writeFile(data.OutputPath, buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	fmt.Printf("✅ Generated: %s\n", data.OutputPath)
	return nil
}

// resolveTemplatePath resolves the template path based on priority
func (tm *TemplateManager) resolveTemplatePath(templateConfig MakeConfig, sourceInfo SourceInfo) string {
	// Get template type from template config
	templateType := tm.extractTemplateType(templateConfig.Template)
	templateFile := tm.extractTemplateFile(templateConfig.Template)

	// Priority 1: Local .stub directory (for development/testing)
	localStubPath := filepath.Join(".", ".stub", templateType, templateFile)
	if tm.templateExists(localStubPath) {
		return localStubPath
	}

	// Priority 2: Cache template
	cachePath := tm.getCachePath(sourceInfo.Name, sourceInfo.Version)
	cacheTemplatePath := filepath.Join(cachePath, ".stub", templateType, templateFile)
	if tm.templateExists(cacheTemplatePath) {
		return cacheTemplatePath
	}

	// Priority 3: Fallback to local .stub (for xarch template)
	return localStubPath
}

// templateExists checks if template exists
func (tm *TemplateManager) templateExists(templatePath string) bool {
	_, err := os.Stat(templatePath)
	return err == nil
}

// getCachePath returns the cache path for a template
func (tm *TemplateManager) getCachePath(templateName, version string) string {
	return filepath.Join(tm.cacheDir, "templates", templateName, version)
}

// loadTemplate loads and parses a template file
func (tm *TemplateManager) loadTemplate(templatePath string) (*template.Template, error) {
	// Check if template path exists
	if _, err := os.Stat(templatePath); err != nil {
		return nil, fmt.Errorf("template not found at %s: %v", templatePath, err)
	}

	content, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file %s: %v", templatePath, err)
	}

	return template.New("").Funcs(tm.createTemplateFunctions()).Parse(string(content))
}

// createTemplateFunctions creates custom template functions
func (tm *TemplateManager) createTemplateFunctions() template.FuncMap {
	return template.FuncMap{
		"title":    tm.toTitleCase,
		"lower":    strings.ToLower,
		"upper":    strings.ToUpper,
		"camel":    tm.toCamelCase,
		"snake":    tm.toSnakeCase,
		"pascal":   tm.toPascalCase,
		"plural":   tm.toPlural,
		"singular": tm.toSingular,
	}
}

// writeFile writes content to a file
func (tm *TemplateManager) writeFile(filePath string, content []byte) error {
	// Check if file already exists and ask for confirmation
	if _, err := os.Stat(filePath); err == nil {
		// File exists, ask for confirmation
		fmt.Printf("⚠️  File %s already exists. Do you want to replace it? (y/N): ", filePath)

		var response string
		fmt.Scanln(&response)

		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			fmt.Println("❌ Operation cancelled.")
			return nil
		}
	}

	return os.WriteFile(filePath, content, 0644)
}

// String manipulation functions
func (tm *TemplateManager) toTitleCase(s string) string {
	// Convert string to title case
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}

func (tm *TemplateManager) toCamelCase(s string) string {
	// Convert snake_case to camelCase
	parts := strings.Split(s, "_")
	if len(parts) == 0 {
		return s
	}

	result := strings.ToLower(parts[0])
	for _, part := range parts[1:] {
		result += tm.toTitleCase(part)
	}
	return result
}

func (tm *TemplateManager) toPascalCase(s string) string {
	// Convert snake_case to PascalCase
	parts := strings.Split(s, "_")
	var result string
	for _, part := range parts {
		result += tm.toTitleCase(part)
	}
	return result
}

func (tm *TemplateManager) toPlural(s string) string {
	s = strings.ToLower(s)

	// Simple pluralization
	if strings.HasSuffix(s, "y") {
		return strings.TrimSuffix(s, "y") + "ies"
	}
	if strings.HasSuffix(s, "s") || strings.HasSuffix(s, "sh") || strings.HasSuffix(s, "ch") {
		return s + "es"
	}
	return s + "s"
}

func (tm *TemplateManager) toSingular(s string) string {
	s = strings.ToLower(s)

	// Simple singularization
	if strings.HasSuffix(s, "ies") {
		return strings.TrimSuffix(s, "ies") + "y"
	}
	if strings.HasSuffix(s, "es") {
		return strings.TrimSuffix(s, "es")
	}
	if strings.HasSuffix(s, "s") {
		return strings.TrimSuffix(s, "s")
	}
	return s
}

// extractTemplateType extracts template type from template config
func (tm *TemplateManager) extractTemplateType(template string) string {
	// If template contains "/", extract the directory part
	if strings.Contains(template, "/") {
		parts := strings.Split(template, "/")
		return parts[0]
	}
	// If no "/", assume it's just the template type
	return template
}

// extractTemplateFile extracts template file name from template config
func (tm *TemplateManager) extractTemplateFile(template string) string {
	// If template contains "/", extract the file part
	if strings.Contains(template, "/") {
		parts := strings.Split(template, "/")
		return parts[1]
	}
	// If no "/", use default template file name
	return "template.go.tmpl"
}

// promptOutputDirectory prompts user for output directory
func (tm *TemplateManager) promptOutputDirectory(templateType string) (string, error) {
	fmt.Printf("📁 Output directory for '%s' is not configured.\n", templateType)
	fmt.Printf("📝 Please enter the output directory (e.g., domain/repositories): ")

	var outputDir string
	fmt.Scanln(&outputDir)

	// Validate input
	outputDir = strings.TrimSpace(outputDir)
	if outputDir == "" {
		return "", fmt.Errorf("output directory cannot be empty")
	}

	// Remove trailing slash if present
	outputDir = strings.TrimSuffix(outputDir, "/")

	fmt.Printf("📂 Using output directory: %s\n", outputDir)
	return outputDir, nil
}
