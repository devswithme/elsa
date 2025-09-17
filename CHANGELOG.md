# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- 

### Features
- 

## [1.0.1] - 2025-09-17

### Fixed
- **Migration SQL Parser**: Fixed PostgreSQL dollar-quoted string parsing in migration execution
  - Improved SQL statement splitting to properly handle `$$` and custom dollar-quoted strings
  - Resolved "unterminated dollar-quoted string" errors during migration execution
  - Enhanced parser to ignore semicolons within dollar-quoted function bodies

- **Project Template Copying**: Fixed hidden files not being copied during project creation
  - Ensures `.gitignore`, `.gitkeep`, and other hidden files are properly copied from templates
