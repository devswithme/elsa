# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- 

### Features
- 

## [1.0.0] - 2025-09-16

### Added
- **Database Migration Management System**
  - DDL/DML separation for production safety
  - Multi-database support (MySQL, PostgreSQL, SQLite)
  - Migration status tracking and rollback control
  - Flexible naming with sequential or timestamp-based IDs
  - Interactive database connection setup
  - Environment-based configuration support

- **File Watching & Auto-Restart**
  - Smart file monitoring for Go files
  - Configurable file extensions and directory exclusions
  - Restart delays to prevent rapid restarts
  - Cross-platform file system watching

- **Elsafile Custom Commands System**
  - Custom command syntax definition
  - Command management and execution
  - Conflict detection with built-in commands
  - Project automation capabilities

- **Code Generation & Scaffolding**
  - Project template system with xarch template
  - Template caching with git-based cache paths
  - Cross-platform cache management
  - Module management with automatic go.mod creation
  - Dependency injection code generation

- **Make System for File Generation**
  - Dynamic template types support
  - Folder structure creation
  - Template versioning and override capabilities
  - Smart template resolution with priority-based loading
  - YAML configuration support

- **Project Creation Tools**
  - New project generation from templates
  - Module name auto-generation
  - Custom output directory support
  - Force overwrite and refresh options

### Features
- **Migration Commands**: `connect`, `create`, `up`, `down`, `status`, `refresh`
- **Watch Commands**: File monitoring with configurable options
- **Elsafile Commands**: `init`, `list`, `run` for custom command management
- **Generate Commands**: Dependency injection code generation
- **New Project Commands**: Template-based project creation
- **Make Commands**: File generation from templates
- **Root Commands**: Help and version information
