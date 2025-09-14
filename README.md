# Elsa - Engineer's Little Smart Assistant

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![Version](https://img.shields.io/badge/Version-1.0.0-green.svg)](https://github.com/risoftinc/elsa)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

**Elsa** is a powerful developer productivity toolkit for Go that provides database migration management, file watching, custom command definitions, code generation, and project automation capabilities.

```
      ___           ___       ___           ___     
     /\  \         /\__\     /\  \         /\  \    
    /::\  \       /:/  /    /::\  \       /::\  \   
   /:/\:\  \     /:/  /    /:/\ \  \     /:/\:\  \  
  /::\~\:\  \   /:/  /    _\:\~\ \  \   /::\~\:\  \ 
 /:/\:\ \:\__\ /:/__/    /\ \:\ \ \__\ /:/\:\ \:\__\
 \:\~\:\ \/__/ \:\  \    \:\ \:\ \/__/ \/__\:\/:/  /
  \:\ \:\__\    \:\  \    \:\ \:\__\        \::/  / 
   \:\ \/__/     \:\  \    \:\/:/  /        /:/  /  
    \:\__\        \:\__\    \::/  /        /:/  /   
     \/__/         \/__/     \/__/         \/__/    V 1.0.0
(migration, scaffolding, project runner and task automation)
```

## 🚀 Key Features

### 📊 Database Migration Management
- **DDL Migrations**: Schema changes, table creation, modifications
- **DML Migrations**: Data seeding, updates, and transformations
- **Multi Database Support**: MySQL, PostgreSQL, SQLite
- **Migration Status Tracking**: View applied, pending, and rollback status
- **Sequential & Timestamp Formats**: Flexible migration naming formats

### 👀 File Watching & Auto-Restart
- **Smart File Monitoring**: Watch Go files and auto-restart on changes
- **Configurable Extensions**: Customize which file types to monitor
- **Directory Exclusion**: Exclude vendor, build, and other directories
- **Restart Delays**: Configurable delays to prevent rapid restarts

### 📝 Elsafile - Custom Commands
- **Custom Command Syntax**: Define custom commands for your project
- **Command Management**: List, run, and manage custom commands
- **Conflict Detection**: Identify conflicts with built-in commands
- **Project Automation**: Streamline your development workflow

### 🔧 Code Generation & Scaffolding
- **Project Templates**: Generate boilerplate code and project structures
- **Template Caching**: Smart caching system with 6-hour TTL
- **Cross-Platform Cache**: Platform-specific cache locations (Windows/macOS/Linux)
- **Version Support**: Support for specific tags, branches, and latest versions
- **Module Management**: Automatic go.mod module name creation

### 🏗️ Make System
- **Dynamic Template Types**: Generate files from configurable templates
- **Folder Structure Support**: Create files in custom directory structures
- **Template Versioning**: Support multiple template versions
- **Custom Template Override**: Override templates per project
- **Smart Template Resolution**: Priority-based template loading

## 📦 Installation

### Prerequisites
- Go 1.21 or higher
- Git

### Install from Repository
```bash
go install go.risoftinc.com/elsa/cmd/elsa@latest
```

### Verify Installation
```bash
elsa --version
```

## 🏃‍♂️ Quick Start

### 1. Initialize New Command
```bash
# Create new Elsafile for custom commands
elsa init

# This will create an Elsafile with common Go project commands
```

### 2. Database Migration
```bash
# Connect to your database
elsa migration connect

# Create new DDL migration
elsa migration create ddl create_users_table

# Apply migrations
elsa migration up ddl

# Check migration status
elsa migration status
```

### 3. File Watching
```bash
# Watch Go files and auto-restart your application
elsa watch "go run main.go"

# Watch with custom settings
elsa watch "go test ./..." --ext ".go,.mod" --exclude "vendor,testdata"
```

### 4. Create New Project
```bash
# Create new project from xarch template
elsa new xarch my-api --module "github.com/username/my-api"

# Create with specific version
elsa new xarch@v1.2.3 my-api --module "github.com/username/my-api"

# Create with custom output directory
elsa new xarch my-api --module "github.com/username/my-api" --output "./projects"
```

### 5. Generate Files with Make
```bash
# Generate repository
elsa make repository user_repository

# Generate service
elsa make service user_service

# Generate with folder structure
elsa make repository health/health_repository

# List available template types
elsa make list
```

**About xarch:**
[xarch](https://github.com/risoftinc/xarch) is a Go project template that provides a clean architecture structure with:
- Domain-driven design (DDD) pattern
- Repository and service layers
- HTTP handlers and middleware
- Database integration with GORM
- Dependency injection setup
- Health check endpoints
- Configuration management

### 5. Custom Commands
```bash
# List available commands from Elsafile
elsa list

# Run custom command
elsa run build
elsa run test
```

## 📚 Commands Reference

### Root Commands
| Command | Description |
|---------|-------------|
| `elsa --help` | Show help information |
| `elsa --version` | Show version information |

### Migration Commands
| Command | Description |
|---------|-------------|
| `elsa migration connect` | Connect to database interactively |
| `elsa migration create ddl <name>` | Create DDL migration |
| `elsa migration create dml <name>` | Create DML migration |
| `elsa migration up <type>` | Apply migrations (type: `ddl` or `dml`) |
| `elsa migration down <type>` | Rollback last migration (type: `ddl` or `dml`) |
| `elsa migration status` | Show migration status |
| `elsa migration refresh <type>` | Refresh all migrations (type: `ddl` or `dml`) |

**Migration Types:**
- `ddl`: Data Definition Language (schema changes, table creation, modifications)
- `dml`: Data Manipulation Language (data seeding, updates, transformations)

### Watch Commands
| Command | Description |
|---------|-------------|
| `elsa watch <command>` | Watch files and auto-restart command |
| `--ext <extensions>` | File extensions to watch (default: .go) |
| `--exclude <dirs>` | Directories to exclude |
| `--delay <duration>` | Restart delay (e.g., 500ms, 1s) |

### Elsafile Commands
| Command | Description |
|---------|-------------|
| `elsa init` | Create new Elsafile |
| `elsa list` | List all commands |
| `elsa list --conflicts` | Show conflicting commands |
| `elsa run <command>` | Execute custom command |

#### Advanced Elsafile Features

**Variable Substitution**
- Environment variables: `$VAR`, `${VAR}`
- Interactive input: `${?VAR:prompt}`, `${?VAR}`
- Quote handling: Supports both single and double quotes
- Line continuation: Use `\` for multi-line commands

**Example Elsafile with Advanced Features:**
```bash
# Environment variables (these variables exist in the environment)
test-env:
	echo "User: $USER"
	echo "Path: ${PATH}"

# Interactive input
# PRIORITY: Check OS environment first, if exists use environment value (NO prompt)
# If not in environment, then prompt user for manual input
migration-create:
	elsa migration create ddl ${?MIGRATION_NAME:Enter migration name}

# Multi-line commands with backslash continuation
complex-setup:
	echo "Setting up project: ${?PROJECT_NAME:Enter project name}" && \
	mkdir ${?PROJECT_NAME} && \
	cd ${?PROJECT_NAME} && \
	echo "Description: ${?DESCRIPTION:Enter description}"

# Mixed variables
test-mixed:
	echo "User: $USER"  # Environment variable
	echo "Migration: ${?MIGRATION_NAME:Enter migration name}"  
	echo "Custom: ${?CUSTOM_VAR:Enter custom value}"  
```

**Variable Substitution Priority:**
1. **Environment variables** (if already set)
2. **Interactive input** (if not in environment)

### Generate Commands
| Command | Description |
|---------|-------------|
| `elsa generate` | Generate dependency injection code |
| `elsa gen` | Short alias for generate |

#### Dependency Injection Example

Create a file with `//go:build elsabuild` tag (e.g., `dep_manager.go`):

```go
//go:build elsabuild
// +build elsabuild

package http

import (
    "go.risoftinc.com/elsa"
    "gorm.io/gorm"
)

type Dependencies struct {
    UserHandler UserHandler
}

func InitializeHandler(db *gorm.DB) *Dependencies {
    elsa.Generate(
        RepositorySet,
        ServicesSet,
        HandlerSet,
    )
    return nil
}

var RepositorySet = elsa.Set(
    NewUserRepository,
)

var ServicesSet = elsa.Set(
    NewUserService,
)

var HandlerSet = elsa.Set(
    NewUserHandler,
)
```

Run the generate command:
```bash
elsa generate
```

This will create `elsa_gen.go` with the generated dependency injection code:

```go
// Code generated by Elsa. DO NOT EDIT.

//go:generate go run -mod=mod go.risoftinc.com/elsa/cmd/elsa gen
//go:build !elsabuild
// +build !elsabuild

package http

func InitializeHandler(db *gorm.DB) *Dependencies {
    userRepo := NewUserRepository(db)
    userSvc := NewUserService(userRepo)
    userHandler := NewUserHandler(userSvc)

    elsa.Generate(userRepo, userSvc, userHandler)
    return &Dependencies{
        UserHandler: userHandler,
    }
}
```

### New Project Commands
| Command | Description |
|---------|-------------|
| `elsa new <template>[@version] <name>` | Create new project from template |
| `--module, -m` | Go module name (required) |
| `--output, -o` | Output directory (default: current) |
| `--force, -f` | Overwrite existing directory |
| `--refresh` | Force refresh template cache |

#### Project Creation Process
The `elsa new` command follows this optimized workflow:
1. **Clone/Update Template**: Download or refresh template from cache
2. **Copy Template**: Copy template files to project directory
3. **Update Module Name**: Replace module name in go.mod and all imports
4. **Generate Proto Files**: Generate Go files from .proto files (if present)
5. **Download Dependencies**: Run `go mod download` to fetch all dependencies
6. **Tidy Dependencies**: Run `go mod tidy` to clean and optimize dependencies
7. **Clean Git History**: Remove template git history for fresh start

This ensures your new project is ready to run immediately with all dependencies properly resolved.

### Make Commands
| Command | Description |
|---------|-------------|
| `elsa make <template-type> <name>` | Generate file from template |
| `elsa make list` | List available template types |
| `elsa make --help` | Show make command help |

#### Template Types
The available template types are defined in your project's `.elsa-config.yaml` file. Common examples:
- `repository` - Generate repository layer files
- `service` - Generate service layer files
- `handler` - Generate HTTP/GRPC handler files

#### Examples
```bash
# Generate repository
elsa make repository user_repository

# Generate service
elsa make service user_service

# Generate with folder structure
elsa make repository health/health_repository

# List available template types
elsa make list
```

#### Configuration
The make system uses `.elsa-config.yaml` in your project root:
```yaml
source:
  name: xarch
  version: v1.2.3
  git_url: https://github.com/risoftinc/xarch
  git_commit: abc123def456

make:
  repository:
    template: repository/template.go.tmpl
    output: domain/repositories
  service:
    template: service/template.go.tmpl
    output: domain/services
```

## 🔧 Configuration

### Elsafile Format
Create an `Elsafile` in your project root:

```bash
# Elsa - Engineer's Little Smart Assistant
# This file defines custom commands for your project
# Commands can be run using: elsa command_name or elsa run command_name

# Build the project
build:
	go build -o bin/app .

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run the application
run:
	go run .

# Format code
fmt:
	go fmt ./...
	go vet ./...
```

### Database Configuration
Elsa supports multiple database drivers:

#### MySQL
```bash
elsa migration connect -c "mysql://user:password@localhost:3306/database"
```

#### PostgreSQL
```bash
elsa migration connect -c "postgres://user:password@localhost:5432/database"
```

#### SQLite
```bash
elsa migration connect -c "sqlite://database.db"
```

## 📁 Project Structure

```
your-project/
├── Elsafile                    # Custom commands definition
├── database/
│   └── migration/
│       ├── ddl/               # DDL migrations
│       │   ├── 20240101120000_create_users.up.sql
│       │   └── 20240101120000_create_users.down.sql
│       └── dml/               # DML migrations
│           ├── 20240101120001_seed_users.up.sql
│           └── 20240101120001_seed_users.down.sql
├── .env                       # Environment variables (optional)
└── your-go-files...
```

## 🎯 Use Cases

### Development Workflow
1. **Project Setup**: Use `elsa init` to create project commands
2. **Database Management**: Create and manage migrations
3. **Development**: Use `elsa watch` for auto-restart during development
4. **Testing**: Run tests with custom commands
5. **Deployment**: Use custom build and deployment commands

### Team Collaboration
- **Standardized Commands**: Share Elsafile across team members
- **Database Consistency**: Use migrations for schema changes
- **Development Environment**: Consistent setup across different machines

## 🤝 Contributing

We welcome contributions! 
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Setup
```bash
# Clone repository
git clone https://github.com/risoftinc/elsa.git
cd elsa

# Install dependencies
go mod download

# Run tests
go test ./...

# Build project
go build -o elsa ./cmd/elsa
```

## 📚 Documentation

### Detailed Guides
- **[Make System Guideline](MAKE_GUIDELINE.md)** - Complete guide for `elsa make` system
  - Template development
  - YAML configuration
  - Advanced usage patterns
  - Troubleshooting guide

### Command Reference
- `elsa --help` - General help
- `elsa make --help` - Make system help
- `elsa migration --help` - Migration commands
- `elsa watch --help` - File watching options

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- **Get to Know**: [https://risoftinc.com](https://risoftinc.com)
- **GitHub**: [https://github.com/risoftinc/elsa](https://github.com/risoftinc/elsa)
- **Issues**: [https://github.com/risoftinc/elsa/issues](https://github.com/risoftinc/elsa/issues)

## 🙏 Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI framework
- Uses [GORM](https://gorm.io/) for database operations
- File watching powered by [fsnotify](https://github.com/fsnotify/fsnotify)

---

**Elsa** - Making Go development more productive, one command at a time! 🚀
