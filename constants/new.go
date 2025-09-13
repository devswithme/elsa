package constants

// New command constants
const (
	// Command usage
	NewUse = "new <template-name>[@version] <project-name>"

	// Command descriptions
	NewShort = "Create a new project from template"
	NewLong  = `Create a new project from a template repository.

This command clones a template repository, updates the go.mod module name,
removes git history, and creates a fresh project ready for development.

Usage:
  elsa new <template-name>[@version] <project-name> [flags]

Examples:
  elsa new xarch my-api --module "github.com/username/my-api"
  elsa new xarch@v1.2.3 my-api --module "github.com/username/my-api"
  elsa new xarch@latest my-api --module "github.com/username/my-api"
  elsa new xarch@main my-api --module "github.com/username/my-api"

Template Format:
  template-name[@version]    Template name with optional version
  project-name               Name of the new project directory

Version Options:
  @latest or no @           Use latest version (main/master branch)
  @v1.2.3                   Use specific tag/version
  @main                     Use specific branch name
  @develop                  Use develop branch

Flags:
  -m, --module string       Go module name for the new project (auto-generated if not provided)
  -o, --output string       Output directory (default: current directory)
  -f, --force              Overwrite existing directory
  --refresh                Force refresh template cache
  -h, --help               Help for new`

	// Flag descriptions
	NewFlagModuleUsage  = "Go module name for the new project (auto-generated from project name if not provided)"
	NewFlagOutputUsage  = "Output directory (default: current directory)"
	NewFlagForceUsage   = "Overwrite existing directory"
	NewFlagRefreshUsage = "Force refresh template cache"

	// Error messages
	NewErrorModuleRequired     = "module name is required. Use --module flag"
	NewErrorInvalidTemplate    = "invalid template format. Use: template-name[@version]"
	NewErrorTemplateNotFound   = "template \"%s\" not found"
	NewErrorVersionNotFound    = "version \"%s\" not found for template \"%s\""
	NewErrorDirExists          = "directory \"%s\" already exists. Use --force to overwrite"
	NewErrorCloneFailed        = "failed to clone template: %v"
	NewErrorUpdateModuleFailed = "failed to update go.mod: %v"
	NewErrorCleanupFailed      = "failed to cleanup git: %v"
	NewErrorInvalidModuleName  = "invalid module name: %v"
	NewErrorProtoGeneration    = "failed to generate proto files: %v"
	NewErrorGoModDownload      = "failed to download Go modules: %v"
	NewErrorGoModTidy          = "failed to tidy Go modules: %v"

	// Success messages
	NewSuccessProjectCreated = "✅ Project \"%s\" created successfully!"
	NewSuccessTemplateCached = "📦 Template \"%s\" cached successfully"
	NewSuccessUsingCache     = "⚡ Using cached template \"%s\""

	// Info messages
	NewInfoCloningTemplate = "⬇️ Cloning template \"%s\"..."
	NewInfoUpdatingModule  = "🚀 Setting module name to \"%s\"..."
	NewInfoGeneratingProto = "🛠️ Generating Go files from .proto files..."
	NewInfoProtocNotFound  = "⚠️ protoc not found, skipping proto generation (install protoc to generate Go files from .proto)"
	NewInfoRefreshingCache = "🔄 Refreshing template cache..."
	NewInfoCacheExpired    = "⏰ Template cache expired, refreshing..."
	NewInfoAutoModule      = "🔧 Auto-generating module name: \"%s\""
	NewInfoGoModDownload   = "📥 Downloading Go modules..."
	NewInfoGoModTidy       = "🧹 Tidying Go modules..."

	// Cache settings
	NewCacheTTLHours = 6
	NewCacheDirName  = ".elsa-cache"
	NewTemplatesDir  = "templates"

	// Cache info messages
	NewInfoCacheLocation = "📁 Cache location: %s"
)
