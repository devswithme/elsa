@echo off
setlocal enabledelayedexpansion

REM Simple Fixed Build script for Elsa CLI application
REM Supports: Linux, Windows, macOS
REM Architectures: amd64, arm64

REM Application info
set APP_NAME=elsa
set VERSION=0.2.0
set BUILD_DIR=build

echo Building Elsa CLI v%VERSION% for multiple platforms
echo =================================================

REM Create build directory
if not exist %BUILD_DIR% mkdir %BUILD_DIR%

REM Clean previous builds
echo Cleaning previous builds...
if exist %BUILD_DIR%\* del /q %BUILD_DIR%\*

REM Build for Linux
echo.
echo ========================================
echo Building for Linux...
echo ========================================
set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
go build -ldflags="-s -w -X main.version=%VERSION%" -o "%BUILD_DIR%\%APP_NAME%-%VERSION%-linux-amd64" .
if %ERRORLEVEL% equ 0 echo ✓ Built: %APP_NAME%-%VERSION%-linux-amd64

set CGO_ENABLED=0
set GOOS=linux
set GOARCH=arm64
go build -ldflags="-s -w -X main.version=%VERSION%" -o "%BUILD_DIR%\%APP_NAME%-%VERSION%-linux-arm64" .
if %ERRORLEVEL% equ 0 echo ✓ Built: %APP_NAME%-%VERSION%-linux-arm64

REM Build for Windows
echo.
echo ========================================
echo Building for Windows...
echo ========================================
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w -X main.version=%VERSION%" -o "%BUILD_DIR%\%APP_NAME%-%VERSION%-windows-amd64.exe" .
if %ERRORLEVEL% equ 0 echo ✓ Built: %APP_NAME%-%VERSION%-windows-amd64.exe

set GOOS=windows
set GOARCH=arm64
go build -ldflags="-s -w -X main.version=%VERSION%" -o "%BUILD_DIR%\%APP_NAME%-%VERSION%-windows-arm64.exe" .
if %ERRORLEVEL% equ 0 echo ✓ Built: %APP_NAME%-%VERSION%-windows-arm64.exe

REM Build for macOS
echo.
echo ========================================
echo Building for macOS...
echo ========================================
set GOOS=darwin
set GOARCH=amd64
go build -ldflags="-s -w -X main.version=%VERSION%" -o "%BUILD_DIR%\%APP_NAME%-%VERSION%-darwin-amd64" .
if %ERRORLEVEL% equ 0 echo ✓ Built: %APP_NAME%-%VERSION%-darwin-amd64

set GOOS=darwin
set GOARCH=arm64
go build -ldflags="-s -w -X main.version=%VERSION%" -o "%BUILD_DIR%\%APP_NAME%-%VERSION%-darwin-arm64" .
if %ERRORLEVEL% equ 0 echo ✓ Built: %APP_NAME%-%VERSION%-darwin-arm64

REM Show build summary
echo.
echo ========================================
echo Build Summary
echo ========================================
echo Build directory: %BUILD_DIR%
echo Version: %VERSION%

REM List all built files
echo.
echo Built executables:
dir %BUILD_DIR%

echo.
echo ========================================
echo Build completed!
echo You can find the executables in the '%BUILD_DIR%' directory
echo ========================================

pause
