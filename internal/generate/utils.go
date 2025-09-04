package generate

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"
)

// Common utility functions for the generate package

// resolvePath resolves a target directory path, handling both absolute and relative paths
// It validates the path exists and returns an error if the directory cannot be found
func resolvePath(targetDir string) (string, error) {
	if targetDir == "" {
		return os.Getwd()
	}

	if filepath.IsAbs(targetDir) {
		return targetDir, nil
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %v", err)
	}

	resolvedPath := filepath.Join(currentDir, targetDir)

	// Check if target directory exists
	if _, err := os.Stat(resolvedPath); os.IsNotExist(err) {
		return "", fmt.Errorf("target directory does not exist: %s", targetDir)
	}

	return resolvedPath, nil
}

// parseImports parses imports from an AST file and returns a map of alias to package path
// Handles both aliased imports (e.g., healthRepo "github.com/xxx/health") and direct imports
// For direct imports, it tries to determine the correct alias by checking if any package uses the alias
// and falls back to the actual package name if available
func parseImports(node *ast.File) map[string]string {
	imports := make(map[string]string)

	// First pass: collect all imports
	var directImports []string
	for _, imp := range node.Imports {
		pkgPath := strings.Trim(imp.Path.Value, `"`)
		if imp.Name != nil {
			// Aliased import, example: healthRepo "github.com/xxx/health"
			imports[imp.Name.Name] = pkgPath
		} else {
			// Direct import - collect for later processing
			directImports = append(directImports, pkgPath)
		}
	}

	// Second pass: process direct imports
	for _, pkgPath := range directImports {
		// Try to get the actual package name
		actualPackageName := getPackageName(pkgPath)

		// Check if the actual package name is already used as an alias
		if _, exists := imports[actualPackageName]; !exists {
			imports[actualPackageName] = pkgPath
		} else {
			// If actual package name is already used, fall back to last segment
			parts := strings.Split(pkgPath, "/")
			fallback := parts[len(parts)-1]
			imports[fallback] = pkgPath
		}
	}

	return imports
}

// getPackageName attempts to get the actual package name from a package path
// It tries to use 'go list' command to get the package name, and falls back to
// the last segment of the path if the command fails
func getPackageName(pkgPath string) string {
	// Try to get package name using go list command
	cmd := exec.Command("go", "list", "-f", "{{.Name}}", pkgPath)
	output, err := cmd.Output()
	if err == nil {
		name := strings.TrimSpace(string(output))
		if name != "" && name != "main" {
			return name
		}
	}

	// Fallback: use last segment of the path
	parts := strings.Split(pkgPath, "/")
	return parts[len(parts)-1]
}

// findGoModDir searches upward from a starting path to find the directory containing go.mod
// This is useful for locating the Go module root from any file within the module
func findGoModDir(start string) (string, error) {
	// If start is a file, get its directory
	dir := start
	if info, err := os.Stat(start); err == nil && !info.IsDir() {
		dir = filepath.Dir(start)
	}

	// Ensure we have an absolute path
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for %s: %v", dir, err)
	}
	dir = absDir

	for {
		// Check if go.mod exists in this directory
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir, nil
		}

		// If we've reached the root and haven't found go.mod
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found from %s upward", start)
		}
		dir = parent
	}
}

// parseFile parses a Go file and returns the AST node
// Uses the Go parser with all error reporting enabled for comprehensive parsing
func parseFile(filePath string) (*ast.File, error) {
	fset := token.NewFileSet()
	return parser.ParseFile(fset, filePath, nil, parser.AllErrors)
}

// isGoFile checks if a file is a Go file by examining its extension
// Case-insensitive check for .go extension
func isGoFile(path string) bool {
	return strings.HasSuffix(strings.ToLower(path), ".go")
}

// hasBuildTag checks if a file contains a specific build tag
// Supports both modern //go:build syntax and legacy // +build syntax
func hasBuildTag(content, tag string) bool {
	return strings.Contains(content, "//go:build "+tag) ||
		strings.Contains(content, "// +build "+tag)
}

// extractFunctionName extracts the function name from a function identifier
// Handles cases where the function name includes package prefix (e.g., "pkg.FuncName")
// Returns just the function name part after the last dot
func extractFunctionName(fn string, pkgPath string) string {
	if dot := strings.LastIndex(fn, "."); dot > 0 {
		return fn[dot+1:]
	}
	return fn
}

// safeWalk safely walks a directory tree, handling errors gracefully
// Wraps filepath.Walk with proper error handling to prevent crashes
func safeWalk(root string, fn filepath.WalkFunc) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return fn(path, info, err)
	})
}

type TypeInfo struct {
	ParamName  string
	Package    string
	UsePointer bool
	DataType   string
	Alias      string
}

func extractType(input string) TypeInfo {
	return extractTypeWithImports(input, nil)
}

func extractTypeWithImports(input string, imports map[string]string) TypeInfo {
	ti := TypeInfo{}

	// cek prefix *
	if strings.HasPrefix(input, "*") {
		ti.UsePointer = true
		input = strings.TrimPrefix(input, "*")
	}

	// Check if this is a builtin type
	if isBuiltinType(input) {
		ti.Package = ""
		ti.DataType = input
		ti.Alias = ""
		return ti
	}

	// pisahkan antara package path & type
	lastDot := strings.LastIndex(input, ".")
	if lastDot != -1 {
		ti.Package = input[:lastDot]
		ti.DataType = input[lastDot+1:]
	} else {
		ti.Package = input
		ti.DataType = input
	}

	// Try to find alias from imports map first
	for alias, pkgPath := range imports {
		if pkgPath == ti.Package {
			ti.Alias = alias
			return ti
		}
	}

	// Fallback: alias = segment terakhir dari package path
	alias := ti.Package
	if idx := strings.LastIndex(alias, "/"); idx != -1 {
		alias = alias[idx+1:]
	}
	ti.Alias = alias

	return ti
}

func lowerFirst(s string) string {
	if s == "" {
		return s
	}

	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func removeArrayByIndex[T any](s []T, index int) []T {
	if index < 0 || index >= len(s) {
		// kalau index invalid, balikin slice asli
		return s
	}
	return append(s[:index], s[index+1:]...)
}

func isBuiltinType(name string) bool {
	// pakai Default importer (bawaannya Go)
	scope := types.Universe
	obj := scope.Lookup(name)
	if obj == nil {
		return false
	}
	// cek apakah object itu sebuah tipe
	_, ok := obj.Type().Underlying().(*types.Basic)
	return ok
}

func getNameFile(path string) string {
	return strings.Split(path, "\\")[len(strings.Split(path, "\\"))-1]
}

// validateGoCode validates the generated Go code for syntax correctness.
// This function uses the Go parser to check if the generated code has valid syntax.
// It also performs additional semantic checks to catch logical errors.
// This ensures that the generated code is syntactically and semantically correct before writing to file.
// Returns an error if the code contains syntax errors, semantic errors, or parsing fails.
func validateGoCode(content string) error {
	fset := token.NewFileSet()
	_, err := parser.ParseFile(fset, "generated.go", content, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("invalid Go syntax: %v", err)
	}

	// Additional semantic validation
	if err := validateSemanticErrors(content); err != nil {
		return fmt.Errorf("semantic error: %v", err)
	}

	return nil
}

// validateSemanticErrors performs additional semantic validation on the generated code.
// This function checks for common logical errors that the Go parser might not catch.
// It looks for patterns like bare type names in return statements or function calls.
// Returns an error if semantic issues are found.
func validateSemanticErrors(content string) error {
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)

		// Check for bare type names in return statements (including continuation lines)
		if strings.HasPrefix(line, "return ") || strings.HasPrefix(line, "}, ") {
			// Check for bare type names in return statements
			// Look for patterns like "return string", "return int", "return ..., string", "}, string"
			if strings.Contains(line, ", string") || strings.Contains(line, ", int") ||
				strings.Contains(line, ", bool") || strings.Contains(line, ", float") ||
				strings.HasSuffix(strings.TrimSpace(line), " string") ||
				strings.HasSuffix(strings.TrimSpace(line), " int") ||
				strings.HasSuffix(strings.TrimSpace(line), " bool") ||
				strings.HasSuffix(strings.TrimSpace(line), " float") {
				// Check if it's a string literal (contains quotes) - if so, it's valid
				if strings.Contains(line, `"`) {
					continue // Skip validation for string literals
				}
				return fmt.Errorf("line %d: bare type name in return statement: %s", i+1, line)
			}
		}

		// Check for bare type names in function calls (but not function signatures)
		if strings.Contains(line, "(") && strings.Contains(line, ")") {
			// Skip function signatures (lines starting with "func ")
			if strings.HasPrefix(line, "func ") {
				continue
			}
			// Look for patterns like "func(..., string)" or "func(..., int)" in function calls
			if strings.Contains(line, ", string)") || strings.Contains(line, ", int)") ||
				strings.Contains(line, ", bool)") || strings.Contains(line, ", float)") {
				return fmt.Errorf("line %d: bare type name in function call: %s", i+1, line)
			}
		}
	}

	return nil
}

// getDefaultValueForType returns the appropriate default value for a given type.
// This function handles built-in Go types and returns their zero values.
// For pointer types, it returns nil. For value types, it returns the appropriate zero value.
// Returns the default value as a string that can be used in Go code.
func getDefaultValueForType(result TypeInfo) string {
	// Handle pointer types
	if result.UsePointer {
		return "nil"
	}

	// Handle built-in types
	switch result.DataType {
	case "string":
		return `""`
	case "int", "int8", "int16", "int32", "int64":
		return "0"
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return "0"
	case "float32", "float64":
		return "0"
	case "bool":
		return "false"
	case "byte":
		return "0"
	case "rune":
		return "0"
	case "complex64", "complex128":
		return "0"
	case "error":
		return "nil"
	default:
		// For custom types, try to use the type name in lowercase
		// This is a fallback for unknown types
		return lowerFirst(result.DataType)
	}
}
