package new

// ProjectOptions contains all options for creating a new project
type ProjectOptions struct {
	TemplateName string
	Version      string
	ProjectName  string
	ModuleName   string
	OutputDir    string
	Force        bool
	Refresh      bool
}

// TemplateManager handles template operations
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

// NewProjectOptions creates a new ProjectOptions struct
func NewProjectOptions(templateName, projectName, moduleName, outputDir string, force, refresh bool) *ProjectOptions {
	return &ProjectOptions{
		TemplateName: templateName,
		ProjectName:  projectName,
		ModuleName:   moduleName,
		OutputDir:    outputDir,
		Force:        force,
		Refresh:      refresh,
	}
}
