# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- 

### Features
- 

## [1.0.2] - 2025-09-17

### Features
- **Improved Command Execution**: Enhanced `elsa run` command to handle single-line multi-command execution
  - Commands with `&&` operators are now executed as single shell commands instead of separate executions
  - Supports both single-line (`cd myapp && go run .`) and multi-line with backslash continuation
  - Maintains backward compatibility with existing separate command execution
  - Examples:
    - `cd myapp && go run .` - executes as single shell command
    - `echo "Setting up project" && \ mkdir project && cd project` - executes as single shell command
    - `cd myapp` followed by `go run .` - executes as separate commands (legacy behavior)
