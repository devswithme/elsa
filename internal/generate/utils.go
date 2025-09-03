package generate

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
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
// For direct imports, it extracts the last segment of the package path as the key
func parseImports(node *ast.File) map[string]string {
	imports := make(map[string]string)

	for _, imp := range node.Imports {
		pkgPath := strings.Trim(imp.Path.Value, `"`)
		if imp.Name != nil {
			// Aliased import, example: healthRepo "github.com/xxx/health"
			imports[imp.Name.Name] = pkgPath
		} else {
			// Direct import → extract last segment
			parts := strings.Split(pkgPath, "/")
			imports[parts[len(parts)-1]] = pkgPath
		}
	}

	return imports
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

	// alias = segment terakhir dari package path
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
