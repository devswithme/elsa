# Elsa Make Guideline

Complete guide for using and developing the `elsa make` system - a dynamic and flexible file generator from templates.

## Table of Contents

- [Overview](#overview)
- [Quick Start](#quick-start)
- [Template Structure](#template-structure)
- [YAML Configuration](#yaml-configuration)
- [Template Development](#template-development)
- [Advanced Usage](#advanced-usage)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

## Overview

`elsa make` is a file generator system that allows you to create Go files automatically based on configurable templates. The system supports:

- **Dynamic template types** - add new template types without code modification
- **Folder structure support** - generate files in desired directory structures
- **Template versioning** - support multiple template versions
- **Custom template override** - override templates per project
- **Cache management** - templates are cached for optimal performance
- **File replacement safety** - confirmation prompt before overwriting existing files
- **Smart package detection** - automatically detects existing package names in target directories
- **Interactive output dialog** - prompts for output directory when not configured

## Quick Start

### 1. Basic Usage

```bash
# Generate repository
elsa make repository user_repository

# Generate service
elsa make service user_service

# Generate with folder structure
elsa make repository health/health_repository
```

**File Replacement Safety:**
- If the target file already exists, Elsa will ask for confirmation before replacing it
- Type `y` or `yes` to confirm, anything else will cancel the operation
- This prevents accidental overwrites of existing code

### 2. List Available Templates

```bash
# List all available template types
elsa make list

# Or
elsa make --help
```

## Template Structure

### Directory Structure

```
template-repo/
├── .stub/                          # Template files (not copied to project)
│   ├── repository/
│   │   └── template.go.tmpl        # Repository template
│   ├── service/
│   │   └── template.go.tmpl        # Service template
│   └── handler-http/
│       └── template.go.tmpl        # HTTP handler template
├── .elsa-config.yaml              # Template configuration
├── domain/                         # Project files
├── infrastructure/
└── ...
```

### Template File Naming

- **Template file**: `template.go.tmpl`
- **Directory naming**: match template type (snake_case)
- **Extension**: `.tmpl` for Go templates

### Template File Structure

```go
// template.go.tmpl
package {{.PackageName}}

type (
	I{{.StructName | title}}Repositories interface {
		// Interface methods here
	}
	{{.StructName | title}}Repositories struct {
		// Struct fields here
	}
)

func New{{.StructName | title}}Repositories() I{{.StructName | title}}Repositories {
	return &{{.StructName | title}}Repositories{
		// Initialization here
	}
}
```

## YAML Configuration

### File: `.elsa-config.yaml`

```yaml
# Source template information (auto-generated)
source:
  name: xarch                    # Template name
  version: v1.2.3               # Template version
  git_url: https://github.com/risoftinc/xarch
  git_commit: abc123def456      # Git commit hash

# Make configuration (template-specific)
make:
  repository:                   # Template type name
    template: repository/template.go.tmpl  # Template file path
    output: domain/repositories # Output directory
  service:
    template: service/template.go.tmpl
    output: domain/services
  handler-http:
    template: handler-http/template.go.tmpl
    output: infrastructure/http/handlers
  handler-grpc:
    template: handler-grpc/template.go.tmpl
    output: infrastructure/grpc/handlers
```

### Configuration Rules

1. **Source section**: Auto-generated during `elsa new`, don't edit manually
2. **Make section**: Template-specific, can be customized per project
3. **Template paths**: Format `template_type/template_file.tmpl`
4. **Template type names**: Must match directory in `.stub/`
5. **Output paths**: Relative from project root
6. **Missing output**: If `output` is empty or missing, Elsa will prompt for input

### Template Resolution Process

The system resolves template files using the following priority order:

1. **Local .stub** (if exists in current directory)
   - Path: `./.stub/{template_type}/{template_file}`
   - Example: `./.stub/repository/template.go.tmpl`

2. **Cache templates** (downloaded from template repository)
   - Path: `~/.elsa/cache/templates/{template_name}/{version}/.stub/{template_type}/{template_file}`
   - Example: `~/.elsa/cache/templates/xarch/v1.2.3/.stub/repository/template.go.tmpl`

3. **Fallback** (local .stub for development)

**Template Path Parsing:**
- `template: repository/template.go.tmpl` → Type: `repository`, File: `template.go.tmpl`
- `template: service` → Type: `service`, File: `template.go.tmpl` (default)
- `template: handler-http/custom.tmpl` → Type: `handler-http`, File: `custom.tmpl`

## Template Development

### 1. Template Variables

The system provides the following variables for templates:

| Variable | Description | Example |
|----------|-------------|---------|
| `{{.PackageName}}` | Package name (extracted from file name) | `user` |
| `{{.StructName}}` | Struct name (extracted from file name) | `User` |
| `{{.FileName}}` | Full file name | `user_repository` |
| `{{.FolderPath}}` | Folder path (if any) | `health` |
| `{{.OutputPath}}` | Full output path | `domain/repositories/user_repository.go` |

### 2. Template Functions

The system provides helper functions:

| Function | Description | Example |
|----------|-------------|---------|
| `{{.StructName \| title}}` | Title case | `UserRepository` |
| `{{.StructName \| lower}}` | Lowercase | `userrepository` |
| `{{.StructName \| upper}}` | Uppercase | `USERREPOSITORY` |
| `{{.StructName \| camel}}` | Camel case | `userRepository` |
| `{{.StructName \| snake}}` | Snake case | `user_repository` |
| `{{.StructName \| pascal}}` | Pascal case | `UserRepository` |
| `{{.StructName \| plural}}` | Plural form | `users` |
| `{{.StructName \| singular}}` | Singular form | `user` |

### 3. Name Parsing Rules

The system parses file names with the following rules:

```bash
# Input: user_repository
PackageName: user
StructName: User
FileName: user_repository

# Input: health/health_repository  
PackageName: health
StructName: Health
FolderPath: health
FileName: health_repository

# Input: user_profile_repository
PackageName: user
StructName: User
FileName: user_profile_repository

# Input: UserService (PascalCase)
PackageName: user
StructName: User
FileName: user_service  # Converted to snake_case

# Input: Health/UserService (PascalCase with folder)
PackageName: user
StructName: User
FolderPath: health
FileName: user_service  # Converted to snake_case
```

**Name Conversion Rules:**
1. **PascalCase to snake_case** - `UserService` becomes `user_service`
2. **Folder names** - also converted to snake_case: `Health/UserService` becomes `health/user_service`
3. **Consistent naming** - all file names use snake_case convention

**Package Name Resolution:**
1. **Existing files** - if Go files exist in output directory, use their package name
2. **File name** - if no existing files, extract from file name

**Example:**
```bash
# If domain/services/user_service.go exists with "package services"
# Then elsa make repository health_repository will generate:
# package services  # (matches existing package)
# type HealthRepositories struct { ... }
```

### 4. Template Examples

#### Repository Template

```go
package {{.PackageName}}

import (
	"context"
)

type (
	I{{.StructName | title}}Repositories interface {
		Create(ctx context.Context, req *Create{{.StructName | title}}Request) error
		GetByID(ctx context.Context, id string) (*{{.StructName | title}}, error)
		Update(ctx context.Context, req *Update{{.StructName | title}}Request) error
		Delete(ctx context.Context, id string) error
		List(ctx context.Context, req *List{{.StructName | title}}Request) ([]*{{.StructName | title}}, error)
	}
	
	{{.StructName | title}}Repositories struct {
		db Database
	}
)

func New{{.StructName | title}}Repositories(db Database) I{{.StructName | title}}Repositories {
	return &{{.StructName | title}}Repositories{
		db: db,
	}
}

func (r *{{.StructName | title}}Repositories) Create(ctx context.Context, req *Create{{.StructName | title}}Request) error {
	// Implementation here
	return nil
}

func (r *{{.StructName | title}}Repositories) GetByID(ctx context.Context, id string) (*{{.StructName | title}}, error) {
	// Implementation here
	return nil, nil
}

func (r *{{.StructName | title}}Repositories) Update(ctx context.Context, req *Update{{.StructName | title}}Request) error {
	// Implementation here
	return nil
}

func (r *{{.StructName | title}}Repositories) Delete(ctx context.Context, id string) error {
	// Implementation here
	return nil
}

func (r *{{.StructName | title}}Repositories) List(ctx context.Context, req *List{{.StructName | title}}Request) ([]*{{.StructName | title}}, error) {
	// Implementation here
	return nil, nil
}
```

#### Service Template

```go
package {{.PackageName}}

import (
	"context"
)

type (
	I{{.StructName | title}}Service interface {
		Create{{.StructName | title}}(ctx context.Context, req *Create{{.StructName | title}}Request) (*{{.StructName | title}}Response, error)
		Get{{.StructName | title}}(ctx context.Context, id string) (*{{.StructName | title}}Response, error)
		Update{{.StructName | title}}(ctx context.Context, req *Update{{.StructName | title}}Request) (*{{.StructName | title}}Response, error)
		Delete{{.StructName | title}}(ctx context.Context, id string) error
		List{{.StructName | title | plural}}(ctx context.Context, req *List{{.StructName | title}}Request) ([]*{{.StructName | title}}Response, error)
	}
	
	{{.StructName | title}}Service struct {
		repo I{{.StructName | title}}Repositories
	}
)

func New{{.StructName | title}}Service(repo I{{.StructName | title}}Repositories) I{{.StructName | title}}Service {
	return &{{.StructName | title}}Service{
		repo: repo,
	}
}

func (s *{{.StructName | title}}Service) Create{{.StructName | title}}(ctx context.Context, req *Create{{.StructName | title}}Request) (*{{.StructName | title}}Response, error) {
	// Business logic here
	return nil, nil
}

func (s *{{.StructName | title}}Service) Get{{.StructName | title}}(ctx context.Context, id string) (*{{.StructName | title}}Response, error) {
	// Business logic here
	return nil, nil
}

func (s *{{.StructName | title}}Service) Update{{.StructName | title}}(ctx context.Context, req *Update{{.StructName | title}}Request) (*{{.StructName | title}}Response, error) {
	// Business logic here
	return nil, nil
}

func (s *{{.StructName | title}}Service) Delete{{.StructName | title}}(ctx context.Context, id string) error {
	// Business logic here
	return nil
}

func (s *{{.StructName | title}}Service) List{{.StructName | title | plural}}(ctx context.Context, req *List{{.StructName | title}}Request) ([]*{{.StructName | title}}Response, error) {
	// Business logic here
	return nil, nil
}
```

## Advanced Usage

### 1. Template Resolution

All templates are centralized in `.stub` directories:

```
template-repo/
├── .stub/                       # Centralized templates
│   ├── repository/
│   │   └── template.go.tmpl
│   └── service/
│       └── template.go.tmpl
└── domain/
```

**Priority order:**
1. Local `.stub` (if exists in project)
2. Cache templates
3. Fallback

**How Template Resolution Works:**

When you run `elsa make repository user_repository`, the system:

1. **Reads YAML config** from `.elsa-config.yaml`
2. **Finds template config** for `repository` type
3. **Parses template path** from `template: repository/template.go.tmpl`
   - Extracts type: `repository`
   - Extracts file: `template.go.tmpl`
4. **Searches for template** in priority order:
   - `./.stub/repository/template.go.tmpl` (local)
   - `~/.elsa/cache/templates/xarch/v1.2.3/.stub/repository/template.go.tmpl` (cache)
5. **Loads template** and generates file at `domain/repositories/user_repository.go`

### 2. Folder Structure Support

```bash
# Generate in folder
elsa make repository health/health_repository
# Output: domain/repositories/health/health_repository.go

# Nested folders
elsa make repository user/profile/user_profile_repository
# Output: domain/repositories/user/profile/user_profile_repository.go
```

### 3. Output Directory Dialog

When the `output` field in YAML config is empty or missing, Elsa will prompt for input:

```bash
$ elsa make repository user_repository
📁 Output directory for 'repository' is not configured.
📝 Please enter the output directory (e.g., domain/repositories): domain/repositories
📂 Using output directory: domain/repositories
🔍 Using template path: .stub/repository/template.go.tmpl
✅ Generated: domain/repositories/user_repository.go
```

**Dialog Features:**
- **Interactive prompt** - asks for output directory when missing
- **Input validation** - ensures directory is not empty
- **Path normalization** - removes trailing slashes automatically
- **Visual feedback** - shows confirmation with emoji indicators

**Use Cases:**
- **Quick setup** - generate files without pre-configuring YAML
- **One-off generation** - temporary output directories
- **Testing** - try different directory structures
- **Migration** - gradually move files to new structure

### 4. Template Versioning

Template versioning is supported through git tags/branches:

```bash
# Use specific version
elsa new myproject xarch@v1.2.3

# Use latest
elsa new myproject xarch@latest
```

## Best Practices

### 1. Template Naming

- **Template types**: use snake_case (`user_repository`, `auth_service`)
- **Directory names**: same as template type
- **File names**: `template.go.tmpl`

### 2. YAML Configuration

- **Source section**: don't edit manually, let it auto-generate
- **Make section**: customize according to project needs
- **Output paths**: use relative paths from project root

### 3. Template Development

- **Consistent naming**: use consistent naming conventions
- **Error handling**: add error handling in templates
- **Documentation**: add comments in generated code
- **Testing**: test templates with various inputs

### 4. Project Structure

```
project/
├── .elsa-config.yaml           # Auto-generated config
├── .stub/                      # Local templates (optional)
│   ├── repository/
│   └── service/
├── domain/
│   ├── repositories/           # Generated repositories
│   ├── services/              # Generated services
│   └── models/                # Generated models
└── infrastructure/
    ├── http/handlers/         # Generated HTTP handlers
    └── grpc/handlers/         # Generated gRPC handlers
```

## Troubleshooting

### Common Issues

#### 1. Template Not Found

```
Error: template not found at .stub/repository: The system cannot find the path specified.
```

**Solution:**
- Ensure template exists at `.stub/repository/template.go.tmpl`
- Check if template type is registered in YAML config

#### 2. Invalid Name Format

```
Error: invalid name: cannot start or end with '/'
```

**Solution:**
- Use format: `name` or `folder/name`
- Avoid special characters at start/end of name

#### 3. Output Directory Not Found

```
Error: failed to create output directory: The system cannot find the path specified.
```

**Solution:**
- Ensure output path is valid in YAML config
- Check permissions to create directory

#### 4. Template Parse Error

```
Error: failed to execute template: template: template:1: unexpected "}"
```

**Solution:**
- Check Go template syntax
- Ensure all `{{}}` tags are valid
- Test template with simple data

### Debug Mode

For debugging, add logging:

```bash
# Enable debug output
elsa make repository user_repository --verbose
```

### Cache Issues

If there are cache issues:

```bash
# Clear cache
rm -rf ~/.elsa/cache

# Force refresh template
elsa new myproject xarch --refresh
```

## Examples

### Complete Example

1. **Setup template repository:**

```bash
# Clone template
git clone https://github.com/risoftinc/xarch my-template
cd my-template

# Create .stub directory
mkdir -p .stub/repository .stub/service

# Create templates
cat > .stub/repository/template.go.tmpl << 'EOF'
package {{.PackageName}}

type I{{.StructName | title}}Repository interface {
	GetByID(id string) (*{{.StructName | title}}, error)
}

type {{.StructName | title}}Repository struct{}

func New{{.StructName | title}}Repository() I{{.StructName | title}}Repository {
	return &{{.StructName | title}}Repository{}
}
EOF

# Create YAML config
cat > .elsa-config.yaml << 'EOF'
source:
  name: my-template
  version: v1.0.0
  git_url: https://github.com/risoftinc/xarch
  git_commit: abc123

make:
  repository:
    template: repository/template.go.tmpl
    output: domain/repositories
  service:
    template: service/template.go.tmpl
    output: domain/services
EOF
```

2. **Use template:**

```bash
# Generate files
elsa make repository user_repository
elsa make service user_service
elsa make repository health/health_repository
```

3. **Generated files:**

```
domain/
├── repositories/
│   ├── user_repository.go
│   └── health/
│       └── health_repository.go
└── services/
    └── user_service.go
```

---

**Note:** This guideline is specific to `elsa make`. For general information about Elsa CLI, see [README.md](README.md).