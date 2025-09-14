package new

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"go.risoftinc.com/elsa/constants"
	"gopkg.in/yaml.v3"
)

// CreateProjectWithOutput creates a new project and handles all output messages
func (tm *TemplateManager) CreateProjectWithOutput(options *ProjectOptions) error {
	// Auto-generate module name if not provided
	if options.ModuleName == "" {
		options.ModuleName = generateModuleName(options.ProjectName)
		fmt.Printf(constants.NewInfoAutoModule+"\n", options.ModuleName)
	}

	// Show cache location
	fmt.Printf(constants.NewInfoCacheLocation+"\n", tm.GetCacheDir())

	// Create project
	if err := tm.CreateProjectFromOptions(options); err != nil {
		return err
	}

	// Show success message
	fmt.Printf(constants.NewSuccessProjectCreated+"\n", options.ProjectName)
	return nil
}

// getTemplateURL returns the repository URL for the template
func (tm *TemplateManager) getTemplateURL(templateName string) string {
	// For now, hardcode xarch template
	// In the future, this could be configurable or loaded from a registry
	templateURLs := map[string]string{
		"xarch": "https://github.com/risoftinc/xarch",
	}

	return templateURLs[templateName]
}

// CreateProjectFromOptions creates a new project from template using options
func (tm *TemplateManager) CreateProjectFromOptions(options *ProjectOptions) error {
	// Validate module name first
	if err := validateModuleName(options.ModuleName); err != nil {
		return fmt.Errorf(constants.NewErrorInvalidModuleName, err)
	}

	// Parse template name and version
	templateName, version, err := parseTemplateName(options.TemplateName)
	if err != nil {
		return fmt.Errorf(constants.NewErrorInvalidTemplate)
	}

	// Determine output directory
	outputDir := options.OutputDir
	if outputDir == "" {
		outputDir = "."
	}

	// Create full project path
	projectPath := filepath.Join(outputDir, options.ProjectName)

	// Check if directory exists and handle force flag
	if _, err := os.Stat(projectPath); err == nil {
		if !options.Force {
			return fmt.Errorf(constants.NewErrorDirExists, options.ProjectName)
		}
		// Remove existing directory if force is enabled
		if err := os.RemoveAll(projectPath); err != nil {
			return fmt.Errorf("failed to remove existing directory: %v", err)
		}
	}

	// Create project using internal logic
	return tm.createProject(templateName, version, projectPath, options.ModuleName, options.Refresh)
}

// createProject creates a new project from template (internal method)
func (tm *TemplateManager) createProject(templateName, version, projectPath, moduleName string, forceRefresh bool) error {
	// Get template repository URL
	templateURL := tm.getTemplateURL(templateName)
	if templateURL == "" {
		return fmt.Errorf(constants.NewErrorTemplateNotFound, templateName)
	}

	// Get cached template path
	cachedPath := tm.getCachedTemplatePath(templateName, version)

	// Check if we need to refresh cache
	needsRefresh := forceRefresh || tm.isCacheExpired(cachedPath)

	if needsRefresh {
		if forceRefresh {
			fmt.Printf(constants.NewInfoRefreshingCache + "\n")
		} else {
			fmt.Printf(constants.NewInfoCacheExpired + "\n")
		}

		// Clone/update template
		if err := tm.cloneTemplate(templateURL, cachedPath, version); err != nil {
			return fmt.Errorf(constants.NewErrorCloneFailed, err)
		}

		fmt.Printf(constants.NewSuccessTemplateCached+"\n", templateName)
	} else {
		fmt.Printf(constants.NewSuccessUsingCache+"\n", templateName)
	}

	// Copy template to project directory
	if err := tm.copyTemplate(cachedPath, projectPath); err != nil {
		return fmt.Errorf("failed to copy template: %v", err)
	}

	// Copy .stub to filestub cache
	if err := tm.copyStubToCache(templateName, version, cachedPath); err != nil {
		return fmt.Errorf("failed to copy .stub to cache: %v", err)
	}

	// Update go.mod module name and all import statements
	fmt.Printf(constants.NewInfoUpdatingModule+"\n", moduleName)
	originalModule, err := tm.updateGoMod(projectPath, moduleName)
	if err != nil {
		return fmt.Errorf(constants.NewErrorUpdateModuleFailed, err)
	}

	// Update all import statements in Go files
	if originalModule != "" {
		if err := tm.updateImports(projectPath, originalModule, moduleName); err != nil {
			return fmt.Errorf("failed to update imports: %v", err)
		}
	}

	// Generate proto files if .proto files are found
	if tm.hasProtoFiles(projectPath) {
		if tm.isProtocInstalled() {
			fmt.Printf(constants.NewInfoGeneratingProto + "\n")
			if err := tm.generateProtoFiles(projectPath); err != nil {
				return fmt.Errorf(constants.NewErrorProtoGeneration, err)
			}
		} else {
			fmt.Printf(constants.NewInfoProtocNotFound + "\n")
		}
	}

	// Download Go modules
	fmt.Printf(constants.NewInfoGoModDownload + "\n")
	if err := tm.runGoModDownload(projectPath); err != nil {
		return fmt.Errorf(constants.NewErrorGoModDownload, err)
	}

	// Tidy Go modules
	fmt.Printf(constants.NewInfoGoModTidy + "\n")
	if err := tm.runGoModTidy(projectPath); err != nil {
		return fmt.Errorf(constants.NewErrorGoModTidy, err)
	}

	// Generate .elsa-config.yaml if config file doesn't exist or has wrong format (before cleaning git history)
	if !tm.hasValidElsaConfig(projectPath) {
		if err := tm.generateElsaConfig(projectPath, templateName, version, templateURL, cachedPath); err != nil {
			return fmt.Errorf("failed to generate config: %v", err)
		}
	}

	// Clean git history
	if err := tm.cleanGitHistory(projectPath); err != nil {
		return fmt.Errorf(constants.NewErrorCleanupFailed, err)
	}

	return nil
}

// generateElsaConfig generates .elsa-config.yaml file for the project
func (tm *TemplateManager) generateElsaConfig(projectPath, templateName, version, templateURL, cachedPath string) error {
	// Get git commit hash from cached template (not from project)
	gitCommit := tm.getGitCommit(cachedPath)

	configPath := filepath.Join(projectPath, ".elsa-config.yaml")

	// Try to load existing config to preserve other sections
	var existingConfig map[string]interface{}
	if data, err := os.ReadFile(configPath); err == nil {
		// Parse existing config
		if err := yaml.Unmarshal(data, &existingConfig); err != nil {
			// If parsing fails, start with empty config
			existingConfig = make(map[string]interface{})
		}
	} else {
		// If file doesn't exist, start with empty config
		existingConfig = make(map[string]interface{})
	}

	// Create ordered config with source at the top
	config := make(map[string]interface{})

	// Add source section first
	config["source"] = map[string]string{
		"name":       templateName,
		"version":    version,
		"git_url":    templateURL,
		"git_commit": gitCommit,
	}

	// Add other sections from existing config (excluding source to avoid duplication)
	for key, value := range existingConfig {
		if key != "source" {
			config[key] = value
		}
	}

	// Convert to YAML
	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config to YAML: %v", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// getGitCommit gets the current git commit hash
func (tm *TemplateManager) getGitCommit(projectPath string) string {
	// Try to get git commit hash
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = projectPath
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	commit := strings.TrimSpace(string(output))
	if len(commit) > 7 {
		return commit[:7] // Return short commit hash
	}
	return commit
}

// hasValidElsaConfig checks if the project has a valid .elsa-config.yaml with correct format
func (tm *TemplateManager) hasValidElsaConfig(projectPath string) bool {
	configPath := filepath.Join(projectPath, ".elsa-config.yaml")

	// Check if file exists
	if _, err := os.Stat(configPath); err != nil {
		return false
	}

	// Read and parse the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return false
	}

	// Parse YAML to check format
	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return false
	}

	// Check if it has the correct "source" structure
	source, ok := config["source"].(map[string]interface{})
	if !ok {
		return false
	}

	// Check if source has required fields
	if _, hasName := source["name"]; !hasName {
		return false
	}
	if _, hasVersion := source["version"]; !hasVersion {
		return false
	}
	if _, hasGitURL := source["git_url"]; !hasGitURL {
		return false
	}
	if _, hasGitCommit := source["git_commit"]; !hasGitCommit {
		return false
	}

	return true
}
