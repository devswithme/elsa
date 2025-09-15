# Elsa Init & Run Commands - Detailed Breakdown

## Overview

This document provides a comprehensive breakdown of Elsa's `init` and `run` commands, which form the core of Elsa's custom command system. These commands enable developers to define and execute project-specific automation tasks through the `Elsafile` configuration.

## Table of Contents

1. [Elsa Init Command](#elsa-init-command)
2. [Elsa Run Command](#elsa-run-command)
3. [Elsafile Structure](#elsafile-structure)
4. [Command Conflict Resolution](#command-conflict-resolution)
5. [Advanced Features](#advanced-features)
6. [Best Practices](#best-practices)
7. [Troubleshooting](#troubleshooting)

---

## Elsa Init Command

### Purpose
The `elsa init` command initializes a new `Elsafile` in the current directory, providing a foundation for project-specific automation commands.

### Syntax
```bash
elsa init
```

### What It Does

1. **Checks for Existing Elsafile**: Verifies that no `Elsafile` already exists in the current directory
2. **Creates Default Template**: Generates a new `Elsafile` with common Go project commands
3. **Sets Proper Permissions**: Creates the file with appropriate read/write permissions (0644)
4. **Provides Success Feedback**: Displays confirmation message with usage instructions

### Default Template Content

When you run `elsa init`, it creates an `Elsafile` with the following structure:

```bash
# Elsa - Engineer's Little Smart Assistant
# This file defines custom commands for your project
# Commands can be run using: elsa command_name or elsa run command_name

# Install elsa cli
install-tools:
	go install go.risoftinc.com/elsa/cmd/elsa@latest

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

### Error Handling

- **File Already Exists**: If `Elsafile` already exists, the command fails with an error message
- **Permission Issues**: If the directory is not writable, the command fails with a permission error
- **Template Generation**: If template generation fails, the command exits with an error

### Success Message

Upon successful creation, Elsa displays:

```
✅ Created Elsafile successfully!
📝 You can now define custom commands in the Elsafile
🚀 Run commands using: elsa run command_name

Example:
  elsa run build    # Run the build command
  elsa run test     # Run the test command
```

---

## Elsa Run Command

### Purpose
The `elsa run` command executes custom commands defined in the `Elsafile`, providing a flexible way to run project-specific automation tasks.

### Syntax
```bash
elsa run <command_name>
```

### What It Does

1. **Loads Elsafile**: Reads and parses the `Elsafile` from the current directory
2. **Validates Command**: Checks if the specified command exists in the Elsafile
3. **Executes Commands**: Runs the command(s) defined for the specified command name
4. **Handles Output**: Displays execution progress and results
5. **Error Handling**: Provides clear error messages for various failure scenarios

### Command Execution Process

1. **File Loading**: The system loads the `Elsafile` using the `Manager.Load()` method
2. **Command Lookup**: Searches for the specified command in the parsed command map
3. **Conflict Check**: Verifies if the command conflicts with built-in Elsa commands
4. **Execution**: Runs each command in the command definition sequentially
5. **Output Display**: Shows execution progress with emojis and formatted output

### Example Usage

```bash
# Run a build command
elsa run build

# Run a test command
elsa run test

# Run a custom command
elsa run deploy
```

### Command Execution Flow

```
🚀 Running Elsafile command: build
📝 Executing: go build -o bin/app .

[Command output follows...]
```

### Error Scenarios

1. **Command Not Found**: `Error: command 'xyz' not found in Elsafile`
2. **Elsafile Not Found**: `Elsafile not found at . Run 'elsa init' to create one`
3. **Execution Failure**: Command execution errors are displayed with context
4. **No Command Specified**: `Error: No command specified. Use 'elsa run command_name'`

---

## Elsafile Structure

### File Format
The `Elsafile` uses a simple, YAML-like syntax for defining commands:

```bash
# Comments start with #
command_name:
	command_to_execute

# Multiple commands in one definition
complex_command:
	command1 && command2 && command3

# Multi-line commands with backslash continuation
multi_line_command:
	echo "Starting process" && \
	cd /path/to/directory && \
	./script.sh
```

### Command Definition Rules

1. **Command Names**: Must end with a colon (`:`)
2. **Indentation**: Commands must be indented with tabs (not spaces)
3. **Comments**: Lines starting with `#` are treated as comments
4. **Empty Lines**: Empty lines are ignored
5. **Line Continuation**: Use backslash (`\`) at the end of a line to continue on the next line

### Variable Substitution

Elsafile supports variable substitution for dynamic command execution:

#### Environment Variables
```bash
# Use environment variables
test-env:
	echo "User: $USER"
	echo "Path: ${PATH}"
```

#### Interactive Input
```bash
# Prompt for user input if variable not in environment
migration-create:
	elsa migration create ddl ${?MIGRATION_NAME:Enter migration name}

# Simple prompt (uses variable name as prompt)
deploy:
	git push origin ${?BRANCH}
```

#### Variable Priority
1. **Environment Variables**: If the variable exists in the environment, use its value
2. **Interactive Input**: If not in environment, prompt the user for input

### Advanced Examples

```bash
# Complex setup with multiple variables
project-setup:
	echo "Setting up project: ${?PROJECT_NAME:Enter project name}" && \
	mkdir ${?PROJECT_NAME} && \
	cd ${?PROJECT_NAME} && \
	echo "Description: ${?DESCRIPTION:Enter description}"

# Mixed variable types
deploy-mixed:
	echo "User: $USER"  # Environment variable
	echo "Branch: ${?BRANCH:Enter branch name}"  # Interactive input
	echo "Environment: ${?ENV:Enter environment}"  # Interactive input
```

---

## Command Conflict Resolution

### Built-in Commands
Elsa has several built-in commands that take precedence over Elsafile commands:

- `init` - Initialize new Elsafile
- `run` - Run Elsafile commands
- `list` - List available commands
- `migrate` - Database migration commands
- `watch` - File watching commands
- `help` - Show help information
- `version` - Show version information

### Conflict Detection
When you try to run a command that conflicts with a built-in command, Elsa will:

1. **Detect the Conflict**: Identify that the command name matches a built-in command
2. **Show Warning**: Display a warning message about the conflict
3. **Provide Resolution**: Suggest using `elsa run <command>` to execute the Elsafile command
4. **Prevent Execution**: Block the execution to avoid confusion

### Conflict Resolution Examples

```bash
# This will show a conflict warning
elsa init  # Conflicts with built-in 'init' command

# Output:
# ⚠️  Command 'init' conflicts with a built-in Elsa command
# 💡 Use 'elsa run init' to execute the Elsafile command
#    Or use 'elsa init' to run the built-in command

# To run the Elsafile command:
elsa run init

# To run the built-in command:
elsa init
```

### Non-Conflicting Commands
Commands that don't conflict with built-in commands can be run directly:

```bash
# These work directly (no conflict)
elsa build
elsa test
elsa deploy
elsa clean
```

---

## Advanced Features

### Multi-Command Execution
You can define multiple commands to run sequentially:

```bash
# All commands run in sequence
full-setup:
	go mod download && \
	go mod tidy && \
	go build -o bin/app . && \
	go test ./...
```

### Conditional Execution
Use shell operators for conditional execution:

```bash
# Only run if previous command succeeds
build-and-test:
	go build -o bin/app . && go test ./...

# Run regardless of previous command result
build-or-test:
	go build -o bin/app . ; go test ./...
```

### Environment-Specific Commands
Create commands for different environments:

```bash
# Development
dev:
	go run . --env=development

# Production
prod:
	go build -o bin/app . && ./bin/app --env=production

# Testing
test-env:
	go test ./... --env=testing
```

### Nested Commands
You can call other Elsafile commands from within commands:

```bash
# Master command that calls others
all:
	elsa run clean && \
	elsa run build && \
	elsa run test

# Individual commands
clean:
	rm -rf bin/
	go clean

build:
	go build -o bin/app .

test:
	go test ./...
```

---

## Best Practices

### Command Naming
- Use descriptive, action-oriented names
- Use kebab-case for multi-word commands (`deploy-staging`, `run-tests`)
- Avoid names that conflict with built-in commands
- Use consistent naming conventions across your team

### Command Organization
- Group related commands together
- Use comments to explain complex commands
- Keep commands focused on single responsibilities
- Use meaningful variable names for interactive inputs

### Error Handling
- Include error checking in your commands
- Use appropriate exit codes
- Provide clear error messages
- Test commands thoroughly before committing

### Documentation
- Comment complex commands
- Document required environment variables
- Provide examples in comments
- Keep the Elsafile readable and maintainable

### Example Well-Structured Elsafile

```bash
# Elsa - Engineer's Little Smart Assistant
# This file defines custom commands for our Go API project
# Commands can be run using: elsa command_name or elsa run command_name

# =============================================================================
# Development Commands
# =============================================================================

# Start development server with hot reload
dev:
	elsa watch "go run cmd/api/main.go"

# Run all tests with coverage
test:
	go test -v -cover ./...

# Run tests for specific package
test-pkg:
	go test -v ./${?PACKAGE:Enter package name}

# =============================================================================
# Build Commands
# =============================================================================

# Build for development
build:
	go build -o bin/api cmd/api/main.go

# Build for production with optimizations
build-prod:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/api cmd/api/main.go

# =============================================================================
# Database Commands
# =============================================================================

# Create new migration
migrate-create:
	elsa migration create ddl ${?MIGRATION_NAME:Enter migration name}

# Apply all migrations
migrate-up:
	elsa migration up ddl

# =============================================================================
# Deployment Commands
# =============================================================================

# Deploy to staging
deploy-staging:
	git push origin staging && \
	ssh staging-server "cd /app && git pull && elsa run migrate-up"

# Deploy to production
deploy-prod:
	git push origin main && \
	ssh prod-server "cd /app && git pull && elsa run migrate-up && elsa run build-prod"
```

---

## Troubleshooting

### Common Issues

#### 1. "Elsafile not found" Error
**Problem**: Command fails with "Elsafile not found" error
**Solution**: 
- Ensure you're in the correct directory
- Run `elsa init` to create an Elsafile
- Check if the file is named exactly "Elsafile" (case-sensitive)

#### 2. "Command not found" Error
**Problem**: Command fails with "command 'xyz' not found in Elsafile"
**Solution**:
- Check the command name spelling
- Ensure the command is defined in the Elsafile
- Use `elsa list` to see available commands

#### 3. Permission Denied Error
**Problem**: Command execution fails with permission errors
**Solution**:
- Check file permissions on the Elsafile
- Ensure the commands have proper permissions
- Check if the target directories are writable

#### 4. Variable Substitution Not Working
**Problem**: Variables are not being substituted properly
**Solution**:
- Check variable syntax (`$VAR` or `${VAR}`)
- Ensure environment variables are set
- Use proper quoting for complex values

#### 5. Command Conflict Issues
**Problem**: Commands conflict with built-in commands
**Solution**:
- Use `elsa run <command>` for Elsafile commands
- Rename conflicting commands in Elsafile
- Use `elsa list --conflicts` to identify conflicts

### Debugging Tips

1. **Use Verbose Output**: Add `echo` statements to debug command execution
2. **Check Command Syntax**: Ensure proper indentation and formatting
3. **Test Commands Manually**: Run commands directly in the shell first
4. **Validate Elsafile**: Use `elsa list` to verify command parsing
5. **Check Environment**: Ensure required environment variables are set

### Getting Help

- Use `elsa --help` for general help
- Use `elsa run --help` for run command help
- Use `elsa list` to see available commands
- Use `elsa list --conflicts` to see conflicting commands
- Check the main README for additional documentation

---

## Conclusion

The `elsa init` and `elsa run` commands provide a powerful foundation for project automation. By understanding their capabilities, conflict resolution, and best practices, you can create efficient, maintainable automation workflows that enhance your development productivity.

The Elsafile system is designed to be simple yet powerful, allowing teams to standardize their development workflows while maintaining flexibility for project-specific needs.
