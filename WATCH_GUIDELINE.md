# Elsa Watch - Complete Development Guide

## 🚀 Overview

**Elsa Watch** is a powerful file monitoring tool that automatically restarts your Go applications when source files change. It's designed to dramatically speed up the development process by eliminating the need to manually restart your application every time you make code changes.

## 🎯 Why Use Elsa Watch?

### The Development Problem
During Go development, developers typically need to:
1. Make code changes
2. Stop the running application (Ctrl+C)
3. Rebuild the application (`go build`)
4. Restart the application (`./app` or `go run main.go`)
5. Repeat this cycle for every change

This manual process is:
- **Time-consuming**: Each change requires 3-4 manual steps
- **Error-prone**: Easy to forget to restart after changes
- **Disruptive**: Breaks development flow and concentration
- **Inefficient**: Wastes valuable development time

### The Elsa Watch Solution
Elsa Watch solves these problems by:
- **Automatic Detection**: Monitors file changes in real-time
- **Smart Restarting**: Automatically restarts your application when Go files change
- **Zero Manual Intervention**: No need to stop/start manually
- **Development Flow**: Maintains continuous development momentum
- **Time Savings**: Saves 5-10 seconds per change, accumulating to hours over a development session
- **Auto Retry**: Automatically retries when Go program fails/errors and code is fixed
- **Seamless Development**: Focus on coding, not process management

## 🔧 How It Works

### Core Architecture
Elsa Watch consists of three main components:

1. **File Watcher** (`internal/watch/watcher.go`)
   - Uses efficient file system monitoring
   - Recursively watches directories while excluding unnecessary folders
   - Filters file changes by extension (default: `.go` files only)
   - Provides event channels for file change notifications

2. **Process Manager** (`internal/watch/process.go`)
   - Handles starting, stopping, and restarting processes
   - Cross-platform process management (Windows/Unix)
   - Graceful shutdown with fallback to force kill
   - Process monitoring and error handling

3. **Command Interface** (`cmd/watch/watch.go`)
   - CLI interface for user interaction
   - Configuration management (extensions, exclusions, delays)
   - Signal handling for clean shutdown
   - Event loop coordination

### File Change Detection Flow
```
File Change Detected
        ↓
Check File Extension (.go by default)
        ↓
Check if Directory is Excluded
        ↓
Trigger Restart with Debounce Delay
        ↓
Stop Current Process
        ↓
Wait for Port/Resource Release
        ↓
Start New Process
        ↓
Process Running Successfully?
        ↓
    Yes: Continue Monitoring
        ↓
    No: Wait for Code Fix
        ↓
    File Change Detected Again
        ↓
    Retry Process Start
        ↓
    Continue Monitoring
```

### Auto Retry & Code Fix Waiting
When your code has syntax errors or compilation issues, Elsa Watch intelligently handles the situation:

1. **Error Detection**: If the process fails to start due to code errors, Elsa Watch detects the failure
2. **Wait State**: Elsa Watch enters a waiting state, monitoring for your next file save
3. **Automatic Retry**: When the Go program fails/errors and you save the file again (after fixing the code), it automatically retries
4. **Success Continuation**: Once the code compiles and runs successfully, normal monitoring resumes

This means you can focus on fixing your code without worrying about manually restarting the application!

## 📋 Basic Usage

### Simple Command
```bash
# Watch and auto-restart a Go application
elsa watch "go run main.go"
```

### With Custom Extensions
```bash
# Watch both Go and Go module files
elsa watch "go run main.go" --ext ".go,.mod"
```

### With Excluded Directories
```bash
# Exclude test and vendor directories
elsa watch "go run main.go" --exclude "vendor,testdata,coverage"
```

### With Custom Restart Delay
```bash
# Wait 1 second before restarting
elsa watch "go run main.go" --delay 1s
```

## ⚙️ Configuration Options

### File Extensions (`--ext`, `-e`)
**Default**: `.go`

Controls which file types trigger restarts:
```bash
# Watch Go files only (default)
elsa watch "go run main.go"

# Watch Go and module files
elsa watch "go run main.go" --ext ".go,.mod"

# Watch Go, module, and template files
elsa watch "go run main.go" --ext ".go,.mod,.tmpl"
```

**Why Go files only by default?**
- Go files contain the actual application logic
- Other files (config, docs, etc.) rarely require application restart
- Prevents unnecessary restarts from temporary files

### Excluded Directories (`--exclude`, `-x`)
**Default**: `.git,vendor,tmp,temp,build,dist,bin,pkg,.vscode,.idea,coverage,testdata`

Prevents watching unnecessary directories:
```bash
# Use default exclusions
elsa watch "go run main.go"

# Add custom exclusions
elsa watch "go run main.go" --exclude "vendor,testdata,logs,cache"

# Watch everything except .git
elsa watch "go run main.go" --exclude ".git"
```

**Why exclude these directories?**
- **`.git`**: Version control files don't affect application behavior
- **`vendor`**: Third-party dependencies (managed by `go mod`)
- **`build`, `dist`, `bin`**: Build artifacts that change frequently
- **`tmp`, `temp`**: Temporary files
- **`.vscode`, `.idea`**: IDE configuration files
- **`coverage`**: Test coverage reports
- **`testdata`**: Test input files

### Restart Delay (`--delay`, `-d`)
**Default**: `500ms`

Controls how long to wait before restarting after file changes:
```bash
# Default 500ms delay
elsa watch "go run main.go"

# 1 second delay
elsa watch "go run main.go" --delay 1s

# 2 second delay
elsa watch "go run main.go" --delay 2s

# 100 milliseconds (very fast)
elsa watch "go run main.go" --delay 100ms
```

**Why use delays?**
- **Debouncing**: Prevents multiple rapid restarts from multiple file saves
- **Resource Release**: Gives time for ports and file handles to be released
- **Stability**: Ensures clean shutdown before restart

## 🎯 Common Use Cases

### 1. Web API Development
```bash
# REST API with hot reload
elsa watch "go run cmd/api/main.go"

# With custom port and environment
elsa watch "PORT=8080 go run cmd/api/main.go"
```

### 2. Microservice Development
```bash
# User service
elsa watch "go run services/user/main.go"

# Order service with custom config
elsa watch "CONFIG_PATH=./config go run services/order/main.go"
```

### 3. CLI Tool Development
```bash
# CLI application
elsa watch "go run cmd/cli/main.go"

# With arguments
elsa watch "go run cmd/cli/main.go --debug --verbose"
```

### 4. Testing and Development
```bash
# Run tests on file changes
elsa watch "go test ./..."

# Run specific test package
elsa watch "go test ./internal/auth"

# Run with coverage
elsa watch "go test -cover ./..."
```

### 5. Build and Run
```bash
# Build and run binary
elsa watch "go build -o app && ./app"

# Build with specific tags
elsa watch "go build -tags dev -o app && ./app"
```

## 🔧 Advanced Configuration

### Environment Variables
```bash
# Set environment variables for the watched process
elsa watch "go run main.go" --env "DEBUG=true,LOG_LEVEL=debug"

# Use .env file
elsa watch "go run main.go" --env-file ".env"
```

### Multiple Commands
```bash
# Chain multiple commands
elsa watch "go mod tidy && go run main.go"

# Run tests then application
elsa watch "go test ./... && go run main.go"
```

## 🚨 Troubleshooting

### Common Issues

#### 1. Port Already in Use
**Problem**: `bind: address already in use`
**Solution**: Increase restart delay to allow port release
```bash
elsa watch "go run main.go" --delay 2s
```

#### 2. Process Not Stopping
**Problem**: Old process still running after restart
**Solution**: Elsa Watch automatically handles this with graceful + force kill

#### 3. Too Many Restarts
**Problem**: Rapid restarts causing instability
**Solution**: Increase delay or check for file save loops
```bash
elsa watch "go run main.go" --delay 1s
```

#### 4. Not Detecting Changes
**Problem**: File changes not triggering restarts
**Solution**: Check file extensions and excluded directories
```bash
# Debug with verbose output
elsa watch "go run main.go" --ext ".go" --exclude ""
```

#### 5. Code Compilation Errors
**Problem**: Application fails to start due to syntax/compilation errors
**Solution**: Elsa Watch automatically handles this scenario
- **Automatic Detection**: Elsa Watch detects when your process fails to start
- **Wait for Fix**: It waits patiently for you to fix the code
- **Auto Retry**: Automatically retries when you save the file again
- **No Manual Restart**: You don't need to manually restart after fixing code

**Example Error Handling Flow:**
```bash
# 1. You save code with syntax error
# Elsa Watch tries to restart: "go run main.go"
# Process fails: "syntax error: unexpected token"

# 2. Elsa Watch shows error and waits
# Console shows: "Command exited with error: exit status 1"
# Elsa Watch continues monitoring for changes

# 3. You fix the code and save again
# Elsa Watch detects the change and retries: "go run main.go"
# Process starts successfully: "Server running on port 8080"

# 4. Normal monitoring resumes
```

### Performance Optimization

#### 1. Exclude Large Directories
```bash
# Exclude node_modules if present
elsa watch "go run main.go" --exclude "vendor,node_modules,dist"
```

#### 2. Use Appropriate Delays
```bash
# Fast development (if system can handle it)
elsa watch "go run main.go" --delay 200ms

# Stable development
elsa watch "go run main.go" --delay 1s
```

#### 3. Monitor Specific Extensions
```bash
# Only watch Go files (most efficient)
elsa watch "go run main.go" --ext ".go"
```

## 📊 Performance Impact

### Resource Usage
- **CPU**: Minimal overhead (~1-2% on modern systems)
- **Memory**: ~10-20MB additional memory usage
- **File Handles**: One handle per watched directory
- **Network**: No network usage (local file monitoring only)

### Optimization Tips
1. **Exclude unnecessary directories** to reduce file handle usage
2. **Use appropriate delays** to balance responsiveness and stability
3. **Monitor only essential file types** to reduce event processing
4. **Close unused applications** to free system resources

## 🎯 Best Practices

### 1. Development Environment
Create an `Elsafile` in your project root:

```bash
# Elsa - Engineer's Little Smart Assistant
# This file defines custom commands for your project

# Development with auto-restart
dev:
	elsa watch "go run main.go"

# Production run
prod:
	go run main.go

# Run tests with watch
test-watch:
	elsa watch "go test ./..."

# Build and run with watch
build-watch:
	elsa watch "go build -o app && ./app"
```

Then use:
```bash
# Development with auto-restart
elsa run dev

# Production run
elsa run prod

# Or run directly
elsa dev
elsa prod
```

### 2. Project Structure
```
project/
├── cmd/
│   └── api/
│       └── main.go          # Main application
├── internal/                # Internal packages
├── pkg/                     # Public packages
├── vendor/                  # Excluded by default
└── testdata/                # Excluded by default
```

### 3. Command Organization
Use `Elsafile` for common commands:

```bash
# Elsafile
dev:
	elsa watch "go run cmd/api/main.go"

test:
	elsa watch "go test ./..."

build:
	elsa watch "go build -o app && ./app"

# Run with: elsa dev, elsa test, elsa build
```


## 🤝 Support

### Getting Help
- **GitHub Issues**: [Report bugs and request features](https://github.com/risoftinc/elsa/issues)
- **Documentation**: Check this guide and related documentation
- **Community**: Join our developer community discussions

### Contributing
- **Bug Reports**: Help us improve by reporting issues
- **Feature Requests**: Suggest new functionality
- **Code Contributions**: Submit pull requests for improvements
- **Documentation**: Help improve our guides and examples

---

**Elsa Watch** - Accelerating Go development with intelligent file monitoring! 🚀
