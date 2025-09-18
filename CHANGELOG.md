# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- 

### Features
- 

## [1.0.3] - 2025-09-18

### Fixed
- **Fixed Migration Refresh Ordering**: Fixed migration refresh command to properly order rollback operations
  - Migration rollback now executes in reverse chronological order (newest to oldest)
  - Ensures proper dependency handling during migration refresh operations
  - Prevents potential errors when migrations have dependencies on each other
  - Consistent with migration down command behavior
