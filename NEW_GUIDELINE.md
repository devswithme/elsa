# Complete Elsa New Guide - Creating New Projects from Templates

## 📋 Table of Contents
1. [Introduction](#introduction)
2. [Command Syntax](#command-syntax)
3. [Features and Options](#features-and-options)
4. [Project Creation Process](#project-creation-process)
5. [Available Templates](#available-templates)
6. [Usage Examples](#usage-examples)
7. [Cache Configuration](#cache-configuration)
8. [Troubleshooting](#troubleshooting)
9. [Advanced Features](#advanced-features)

## 🚀 Introduction

The `elsa new` command is a core feature of Elsa that allows you to quickly create new Go projects using pre-configured templates. This feature is extremely useful for:

- **Starting new projects** with organized structure
- **Saving time** in boilerplate code setup
- **Ensuring consistency** in team project structure
- **Integrating best practices** built into templates

### 💡 Why Use Templates?

**📖 Real Experience Story:**
There was once a job testing opportunity with a 3-day deadline ⏰. Without a ready-to-use Go base code, most of the testing time was consumed by setting up project structure 🏗️, writing boilerplate code 📝, configuring database connections 🗄️, and various other technical setup tasks. As a result, the testing performance was suboptimal 😔 because energy and time were focused on building foundation code rather than completing the core tasks that should have been prioritized 🎯. From this experience, I realized the importance of having ready-to-use base code, and if you have any job opportunities for Go developers, please visit [my LinkedIn for hire me](https://www.linkedin.com/in/riskykurniawan15/) ✌️.

**🚀 With Elsa New:**
- ⏰ **Save Time**: Project ready in seconds
- 🎯 **Focus on Task**: Directly focus on business logic
- 🏗️ **Professional Structure**: Using proven best practices
- 🔧 **Production Ready**: Template already includes necessary configuration

### Elsa New Advantages
- ✅ **Template Caching**: Smart caching system for optimal performance
- ✅ **Auto Module Generation**: Automatically generate Go module names
- ✅ **Git History Clean**: Clean template git history
- ✅ **Dependency Management**: Automatic dependency download and tidy
- ✅ **Proto Support**: Generate Go files from .proto files
- ✅ **Cross Platform**: Works on Windows, macOS, and Linux

## 📝 Command Syntax

```bash
elsa new <template-name>[@version] <project-name> [flags]
```

### Required Parameters
- **`template-name`**: Name of the template to use
- **`project-name`**: Name of the project/directory to create

### Optional Parameters
- **`@version`**: Template version (tag, branch, or latest)
- **Flags**: Additional customization options

## ⚙️ Features and Options

### Available Flags

| Flag | Alias | Description | Default |
|------|-------|-------------|---------|
| `--module` | `-m` | Go module name for new project | Auto-generate from project name |
| `--output` | `-o` | Output directory | Current directory |
| `--force` | `-f` | Overwrite existing directory | `false` |
| `--refresh` | - | Force refresh template cache | `false` |

### Template Version Format

| Format | Description | Example |
|--------|-------------|---------|
| `@latest` or no `@` | Use latest version (main/master branch) | `xarch@latest` |
| `@v1.2.3` | Use specific tag/version | `xarch@v1.2.3` |
| `@main` | Use specific branch | `xarch@main` |
| `@develop` | Use develop branch | `xarch@develop` |

## 🔄 Project Creation Process

Elsa new follows an optimized workflow to ensure new projects are ready to use:

### 1. **Input Validation** 🔍
- Validate template name format
- Validate module name (if provided)
- Check output directory existence

### 2. **Template Management** 📦
- **Cache Check**: Check if template exists in cache
- **Clone/Update**: Download template if not exists or expired
- **Version Resolution**: Resolve requested template version

### 3. **Project Creation** 🏗️
- **Directory Setup**: Create project directory
- **File Copying**: Copy all template files (except .git and .stub)
- **Stub Caching**: Copy .stub directory to filestub cache

### 4. **Module Configuration** ⚙️
- **go.mod Update**: Update module name in go.mod
- **Dependency Resolution**: Download and tidy dependencies

### 5. **Code Generation** 🛠️
- **Proto Generation**: Generate Go files from .proto files (if present)
- **Config Generation**: Generate .elsa-config.yaml

### 6. **Cleanup** 🧹
- **Git History Clean**: Remove template git history
- **Final Validation**: Validate project is ready to use

## 📚 Available Templates

### Xarch Template
**Repository**: [https://github.com/risoftinc/xarch](https://github.com/risoftinc/xarch)

Xarch is a Go template that provides clean architecture structure with:

#### 🏗️ Project Structure
```
xarch/
├── cmd/                    # Application entry points
├── config/                 # Application configuration
├── constant/               # Application constants
├── database/               # Database migrations and seeders
├── domain/                 # Domain models and interfaces
│   ├── models/            # Entity models
│   ├── repositories/      # Repository interfaces
│   └── services/          # Service interfaces
├── driver/                 # External dependencies
├── infrastructure/         # Infrastructure layer
│   └── http/              # HTTP handlers and middleware
├── utils/                  # Utility functions
└── main.go                # Main application entry
```

#### ✨ Included Features
- **Domain-Driven Design (DDD)**: Clear layer separation
- **Repository Pattern**: Interface for data access
- **Service Layer**: Business logic abstraction
- **HTTP Handlers**: REST API endpoints
- **Middleware Support**: Request/response processing
- **Database Integration**: GORM integration
- **Dependency Injection**: Clean dependency management
- **Health Check**: Health monitoring endpoints
- **Configuration Management**: Environment-based config

#### 🛠️ Technologies Used
- **Go 1.21+**: Language runtime
- **GORM**: ORM for database operations
- **Gin**: HTTP web framework
- **Viper**: Configuration management
- **JWT**: Authentication support
- **Bcrypt**: Password hashing
- **Validator**: Input validation

## 💡 Usage Examples

### Basic Example
```bash
# Create new project with auto-generated module
elsa new xarch my-api

# Output:
# 🔧 Auto-generating module name: "my-api"
# 📁 Cache location: /Users/username/.elsa-cache
# ⚡ Using cached template "xarch"
# 🚀 Setting module name to "my-api"...
# 📥 Downloading Go modules...
# 🧹 Tidying Go modules...
# ✅ Project "my-api" created successfully!
```

### Example with Custom Module
```bash
# Create project with specific module name
elsa new xarch my-api --module "github.com/username/my-api"

# Output:
# 📁 Cache location: /Users/username/.elsa-cache
# ⚡ Using cached template "xarch"
# 🚀 Setting module name to "github.com/username/my-api"...
# 📥 Downloading Go modules...
# 🧹 Tidying Go modules...
# ✅ Project "my-api" created successfully!
```

### Example with Specific Version
```bash
# Use specific template version
elsa new xarch@v1.2.3 my-api --module "github.com/username/my-api"

# Use develop branch
elsa new xarch@develop my-api --module "github.com/username/my-api"
```

### Example with Output Directory
```bash
# Create project in specific directory
elsa new xarch my-api --module "github.com/username/my-api" --output "./projects"

# Result: ./projects/my-api/
```

### Example with Force Overwrite
```bash
# Overwrite existing directory
elsa new xarch my-api --module "github.com/username/my-api" --force
```

### Example with Refresh Cache
```bash
# Force refresh template from repository
elsa new xarch my-api --module "github.com/username/my-api" --refresh
```

## 💾 Cache Configuration

### Cache Location
Elsa uses an organized cache system based on platform:

| Platform | Cache Location |
|----------|----------------|
| **Windows** | `%USERPROFILE%\.elsa-cache` |
| **macOS** | `~/Library/Caches/elsa` |
| **Linux** | `~/.cache/elsa` |

### Cache Structure
```
# Windows: %USERPROFILE%\.elsa-cache
# macOS: ~/Library/Caches/elsa
# Linux: ~/.cache/elsa

elsa/
├── templates/                    # Template cache
│   └── xarch/                   # Template name
│       ├── main/                # Main branch version
│       ├── v1.2.3/              # Tag version
│       └── develop/             # Other branch versions
└── filestub/                    # Filestub cache
    └── github.com/
        └── risoftinc/
            └── xarch/
                └── abc123def/   # Commit hash
                    └── .stub/   # Stub files
```

### Cache Management
- **TTL**: Cache expires after 6 hours
- **Auto Refresh**: Automatically refresh if cache expired
- **Manual Refresh**: Use `--refresh` flag for force refresh
- **Git-based Paths**: Cache paths follow git URL structure

## 🔧 Troubleshooting

### Common Issues and Solutions

#### 1. Template Not Found
```bash
# Error: template "invalid-template" not found
# Solution: Ensure template name is correct
elsa new xarch my-api  # ✅ Correct
elsa new invalid my-api  # ❌ Wrong
```

#### 2. Directory Already Exists
```bash
# Error: directory "my-api" already exists. Use --force to overwrite
# Solution: Use --force flag or choose different name
elsa new xarch my-api --force  # ✅ Overwrite
elsa new xarch my-api-v2       # ✅ New name
```

#### 3. Invalid Module Name
```bash
# Error: invalid module name: module name must be a valid Go module path
# Solution: Use correct module format
elsa new xarch my-api --module "github.com/username/my-api"  # ✅ Correct
elsa new xarch my-api --module "invalid-module-name"        # ❌ Wrong
```

#### 4. Cache Corrupted
```bash
# Solution: Remove cache and refresh
rm -rf ~/.elsa-cache
elsa new xarch my-api --refresh
```

#### 5. Protoc Not Installed
```bash
# Warning: protoc not found, skipping proto generation
# Solution: Install protoc to generate .proto files
# macOS: brew install protobuf
# Ubuntu: sudo apt-get install protobuf-compiler
# Windows: Download from https://github.com/protocolbuffers/protobuf/releases
```

#### 6. Go Module Error
```bash
# Error: go mod download failed
# Solution: Ensure Go is installed and network is available
go version  # Check Go version
go env GOPROXY  # Check proxy setting
```

## 🚀 Advanced Features

### Elsa Ecosystem (In Development)

Currently Elsa only provides the `xarch` template, but an Elsa ecosystem is being developed that will enable:

#### 🔄 Template Sharing
- **Community Templates**: Templates created by the community
- **Template Registry**: Centralized registry for templates
- **Version Management**: Versioning system for templates
- **Quality Control**: Template validation according to best practices

#### 🏗️ Template Development and Distribution
- **Template Standards**: Standards for creating templates
- **Validation Tools**: Tools for template validation
- **Update Mechanism**: Template update mechanism
- **Package Management**: Package system for templates


### Best Practices for Templates

#### 1. **Template Structure**
```
template-name/
├── .stub/                    # Stub files for elsa make
├── .elsa-config.yaml        # Template configuration
├── README.md                # Template documentation
├── go.mod                   # Go module definition
├── main.go                  # Main entry point
└── ...                      # Template files
```


### Elsa New Roadmap

#### Phase 1: Current (v1.0.0)
- ✅ Basic template support (xarch)
- ✅ Cache management
- ✅ Module name generation
- ✅ Git history cleanup

#### Phase 2: Template Registry (v1.1.0)
- 🔄 Community template registry
- 🔄 Template marketplace
- 🔄 Custom template creation

## 📚 References

### Template Repository
- **[Xarch Template](https://github.com/risoftinc/xarch)** - Clean architecture template

### Community
- **[GitHub Issues](https://github.com/risoftinc/elsa/issues)** - Bug reports and feature requests
- **[Discussions](https://github.com/risoftinc/elsa/discussions)** - Community discussions

---

**Elsa New** - Creating Go projects quickly and efficiently! 🚀

*This documentation will be continuously updated as Elsa features evolve. If you find any errors or have suggestions, please create an issue in the GitHub repository.*
