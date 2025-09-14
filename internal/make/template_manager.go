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
func (tm *TemplateManager) GenerateFile(templateType, name string, refresh bool) error {
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
	templatePath := tm.resolveTemplatePath(templateConfig, config.Source, refresh)
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

	return nil
}

// resolveTemplatePath resolves the template path based on priority
func (tm *TemplateManager) resolveTemplatePath(templateConfig MakeConfig, sourceInfo SourceInfo, refresh bool) string {
	// Get template type from template config
	templateType := tm.extractTemplateType(templateConfig.Template)
	templateFile := tm.extractTemplateFile(templateConfig.Template)

	// Priority 1: Local .stub directory (for development/testing)
	localStubPath := filepath.Join(".", ".stub", templateType, templateFile)
	if tm.templateExists(localStubPath) {
		return localStubPath
	}

	// Priority 2: Filestub cache (new structure)
	filestubPath := tm.getFilestubCachePath(sourceInfo.GitURL, sourceInfo.GitCommit, templateType, templateFile)
	if !refresh && tm.templateExists(filestubPath) {
		return filestubPath
	}

	// Priority 3: Try to clone .stub if not found in filestub cache or if refresh is requested
	if refresh {
		fmt.Printf("🔄 Refresh requested, cloning .stub from remote repository\n")
	}

	if sourceInfo.GitCommit != "" {
		if tm.cloneStubToCache(sourceInfo) {
			// Try again after cloning
			if tm.templateExists(filestubPath) {
				return filestubPath
			}
		}
	} else if sourceInfo.GitURL != "" {
		// If git_commit is empty but git_url exists, get latest commit and clone
		fmt.Printf("🔍 Git commit is empty, getting latest commit from: %s\n", sourceInfo.GitURL)
		latestCommit := tm.getLatestCommitFromRemote(sourceInfo.GitURL)
		if latestCommit != "" {
			fmt.Printf("🔍 Found latest commit: %s\n", latestCommit)
			// Update config with latest commit
			if err := tm.updateConfigWithCommit(latestCommit); err != nil {
				fmt.Printf("⚠️  Failed to update config: %v\n", err)
			} else {
				fmt.Printf("✅ Config updated successfully\n")
				// Update sourceInfo with latest commit
				sourceInfo.GitCommit = latestCommit
				// Try to clone with latest commit
				if tm.cloneStubToCache(sourceInfo) {
					// Try again after cloning
					filestubPath = tm.getFilestubCachePath(sourceInfo.GitURL, sourceInfo.GitCommit, templateType, templateFile)
					if tm.templateExists(filestubPath) {
						return filestubPath
					}
				}
			}
		} else {
			fmt.Printf("❌ Could not get latest commit from remote\n")
		}
	}

	// Fallback to local .stub
	return localStubPath
}

// templateExists checks if template exists
func (tm *TemplateManager) templateExists(templatePath string) bool {
	_, err := os.Stat(templatePath)
	return err == nil
}

// getCachePath returns the cache path for a template
func (tm *TemplateManager) getCachePath(templateName, version string) string {
	return filepath.Join(tm.cacheDir, templateName, version)
}

// getFilestubCachePath returns the filestub cache path for a template
func (tm *TemplateManager) getFilestubCachePath(gitURL, commitHash, templateType, templateFile string) string {
	if commitHash == "" {
		return "" // Return empty if no commit hash
	}

	baseDir := getFilestubCacheDir(gitURL)
	return filepath.Join(baseDir, commitHash, ".stub", templateType, templateFile)
}

// cloneStubToCache clones only the .stub directory from the template repository
func (tm *TemplateManager) cloneStubToCache(sourceInfo SourceInfo) bool {
	// Create filestub cache directory using git URL
	filestubCacheDir := getFilestubCacheDir(sourceInfo.GitURL)
	stubDestPath := filepath.Join(filestubCacheDir, sourceInfo.GitCommit, ".stub")

	// Create temporary directory for cloning
	tempDir := filepath.Join(filestubCacheDir, "temp")
	defer os.RemoveAll(tempDir) // Clean up temp directory

	// Clone the repository to temp directory
	if err := tm.cloneRepository(sourceInfo.GitURL, tempDir, sourceInfo.GitCommit); err != nil {
		fmt.Printf("❌ Failed to clone repository: %v\n", err)
		return false
	}

	// Check if .stub directory exists in cloned repository
	stubSourcePath := filepath.Join(tempDir, ".stub")
	if _, err := os.Stat(stubSourcePath); os.IsNotExist(err) {
		fmt.Printf("⚠️  No .stub directory found in template\n")
		return false
	}

	// Create destination directory
	if err := os.MkdirAll(filestubCacheDir, 0755); err != nil {
		fmt.Printf("❌ Failed to create cache directory: %v\n", err)
		return false
	}

	// Remove existing .stub in cache if it exists
	if err := os.RemoveAll(stubDestPath); err != nil {
		fmt.Printf("❌ Failed to remove existing .stub cache: %v\n", err)
		return false
	}

	// Copy .stub directory to cache
	if err := tm.copyDirectory(stubSourcePath, stubDestPath, []string{}); err != nil {
		fmt.Printf("❌ Failed to copy .stub to cache: %v\n", err)
		return false
	}

	return true
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

// updateConfigWithCommit updates the .elsa-config.yaml with the latest commit hash
func (tm *TemplateManager) updateConfigWithCommit(commitHash string) error {
	configPath := ".elsa-config.yaml"

	// Read existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse YAML
	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config: %v", err)
	}

	// Update source section
	if source, ok := config["source"].(map[string]interface{}); ok {
		source["git_commit"] = commitHash
		config["source"] = source
	} else {
		// Create source section if it doesn't exist
		config["source"] = map[string]interface{}{
			"git_commit": commitHash,
		}
	}

	// Also update template section if it exists and has git_url
	if template, ok := config["template"].(map[string]interface{}); ok {
		if gitURL, exists := template["git_url"]; exists && gitURL != "" {
			template["git_commit"] = commitHash
			config["template"] = template
		}
	}

	// Write updated config
	updatedData, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal updated config: %v", err)
	}

	if err := os.WriteFile(configPath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write updated config: %v", err)
	}

	return nil
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
		"title":    toTitleCase,
		"lower":    strings.ToLower,
		"upper":    strings.ToUpper,
		"camel":    toCamelCase,
		"snake":    toSnakeCase,
		"pascal":   toPascalCase,
		"plural":   toPlural,
		"singular": toSingular,
	}
}

// writeFile writes content to a file
func (tm *TemplateManager) writeFile(filePath string, content []byte) error {
	// Check if file already exists and ask for confirmation
	if _, err := os.Stat(filePath); err == nil {
		// File exists, ask for confirmation
		fmt.Printf("⚠️  File %s already exists. Do you want to replace it? (Y/N): ", filePath)

		var response string
		fmt.Scanln(&response)

		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			fmt.Println("❌ Operation cancelled.")
			return nil
		}
	}

	fmt.Printf("✅ Generated: %s\n", filePath)
	return os.WriteFile(filePath, content, 0644)
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
