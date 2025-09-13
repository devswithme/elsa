package new

import (
	"fmt"
	"os"
	"path/filepath"

	"go.risoftinc.com/elsa/constants"
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

	// Clean git history
	if err := tm.cleanGitHistory(projectPath); err != nil {
		return fmt.Errorf(constants.NewErrorCleanupFailed, err)
	}

	return nil
}
