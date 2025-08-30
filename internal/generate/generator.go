package generate

import (
	"fmt"
	"go/ast"
	"path/filepath"
)

// Generator handles the generation process for finding elsabuild files
// Manages the parsing and analysis of Go files with elsabuild build tags
type Generator struct {
	imports map[string]string
}

// NewGenerator creates a new Generator instance
// Initializes a new generator with an empty imports map
func NewGenerator() *Generator {
	return &Generator{}
}

// GenerateDependencies processes all elsabuild files in the target directory
// Finds files with elsabuild tags and processes their dependencies
// Returns an error if the process fails, but continues processing other files
func (g *Generator) GenerateDependencies(targetDir string) error {
	files, err := g.FindElsabuildFiles(targetDir)
	if err != nil {
		fmt.Printf("Warning: failed to find elsabuild files: %v\n", err)
		return nil
	}

	for _, file := range files {
		if err := g.processGenerateDependencies(filepath.Join(targetDir, file)); err != nil {
			fmt.Printf("Warning: failed to process generate dependencies: %v\n", err)
			return nil
		}
	}

	return nil
}

// processGenerateDependencies processes a single file for generation dependencies
// Extracts functions with elsa.Generate calls and analyzes elsa.Set declarations
// Loads constructor information for the found functions
func (g *Generator) processGenerateDependencies(target string) error {
	goModDir, err := g.FindGoModDir(target)
	if err != nil {
		return err
	}

	// Extract functions containing elsa.Generate
	funcs, err := g.ExtractElsaGenerateFuncs(target)
	if err != nil {
		return err
	}

	fmt.Println(funcs)

	sets, err := g.ParseElsaSets(target)
	if err != nil {
		return err
	}

	for key, set := range sets {
		fmt.Println(key, " => ", set)

		for _, s := range set {
			constructors, err := g.LoadConstructors(goModDir, s.PkgPath, []string{s.FuncName})
			if err != nil {
				return err
			}
			fmt.Println(constructors)
		}

	}

	return nil
}

// ResultInfo stores information about return values
// Contains type information and whether the return value is a struct
type ResultInfo struct {
	Type         string            // Return type
	IsStruct     bool              // Whether the return type is a struct
	StructFields []StructFieldInfo // Struct field information if applicable
}

// StructFieldInfo stores information about struct fields
// Contains field name, type, and tags
type StructFieldInfo struct {
	Name string // Field name
	Type string // Field type
	Tag  string // Field tags
}

// ExtractElsaGenerateFuncs extracts functions containing elsa.Generate calls
// Parses the target file and finds all functions that call elsa.Generate
// Returns detailed function information including parameters and return types
func (g *Generator) ExtractElsaGenerateFuncs(target string) ([]FuncInfo, error) {
	// Parse file into AST
	node, err := parseFile(target)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %v", err)
	}

	// Parse imports for mapping aliases to package paths
	g.imports = parseImports(node)

	var funcs []FuncInfo

	// Traverse AST to find functions
	ast.Inspect(node, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		// Check if function contains elsa.Generate
		if !containsElsaGenerate(fn) {
			return true
		}

		// Extract function information
		funcInfo := FuncInfo{
			FuncName: fn.Name.Name,
			PkgName:  node.Name.Name,
		}

		// Extract parameters with full package paths
		funcInfo.Params = g.extractParamsWithImports(fn.Type.Params)

		// Extract return values
		if fn.Type.Results != nil {
			funcInfo.Results = g.extractResults(fn.Type.Results, target)
		}

		// Extract parameters from elsa.Generate
		funcInfo.GenerateParams = g.extractGenerateParams(fn)

		funcs = append(funcs, funcInfo)
		return true
	})

	return funcs, nil
}

// containsElsaGenerate checks if a function contains elsa.Generate calls
// Searches through the function body for calls to elsa.Generate
// Returns true if found, false otherwise
func containsElsaGenerate(fn *ast.FuncDecl) bool {
	if fn.Body == nil {
		return false
	}

	calls := findCallExpressions(fn.Body)
	for _, call := range calls {
		if isElsaCall(call, "Generate") {
			return true
		}
	}

	return false
}

// extractGenerateParams extracts parameters from elsa.Generate calls
// Finds elsa.Generate calls in the function body and extracts their arguments
// Returns a slice of parameter names or string representations
func (g *Generator) extractGenerateParams(fn *ast.FuncDecl) []string {
	var params []string

	if fn.Body == nil {
		return params
	}

	calls := findCallExpressions(fn.Body)
	for _, call := range calls {
		if isElsaCall(call, "Generate") {
			// Extract all arguments from elsa.Generate
			for _, arg := range call.Args {
				if ident, ok := arg.(*ast.Ident); ok {
					params = append(params, ident.Name)
				} else {
					// Fallback for expressions that aren't identifiers
					params = append(params, exprToString(arg))
				}
			}
			break // Only process the first elsa.Generate call
		}
	}

	return params
}

// extractParamsWithImports extracts parameters with full package paths
// Uses TypeResolver to resolve types with import context
// Returns parameter information with fully qualified type names
func (g *Generator) extractParamsWithImports(fieldList *ast.FieldList) []ParamInfo {
	if fieldList == nil {
		return nil
	}

	resolver := NewTypeResolver(g.imports)
	var params []ParamInfo

	for _, field := range fieldList.List {
		typ := resolver.ResolveType(field.Type)

		if len(field.Names) > 0 {
			for _, name := range field.Names {
				params = append(params, ParamInfo{
					Name: name.Name,
					Type: typ,
				})
			}
		} else {
			params = append(params, ParamInfo{
				Name: "",
				Type: typ,
			})
		}
	}
	return params
}

// extractResults extracts information about return values
// Analyzes return types and determines if they are structs
// Extracts struct field information when applicable
func (g *Generator) extractResults(fieldList *ast.FieldList, filePath string) []ResultInfo {
	if fieldList == nil {
		return nil
	}

	var results []ResultInfo
	for _, field := range fieldList.List {
		typ := exprToString(field.Type)

		resultInfo := ResultInfo{
			Type: typ,
		}

		// Check if this is a struct type
		if g.isStructType(field.Type, filePath) {
			resultInfo.IsStruct = true
			resultInfo.StructFields = g.extractStructFields(field.Type, filePath)
		}

		results = append(results, resultInfo)
	}
	return results
}

// isStructType checks if a type is a struct dynamically
// Analyzes the type expression to determine if it represents a struct
// May require file parsing for complete type information
func (g *Generator) isStructType(expr ast.Expr, filePath string) bool {
	switch t := expr.(type) {
	case *ast.StarExpr:
		return g.isStructType(t.X, filePath)
	case *ast.Ident:
		// Check if there's a struct definition with this name in the file
		return g.hasStructDefinition(t.Name, filePath)
	case *ast.SelectorExpr:
		// This could be a struct from another package
		return true
	default:
		return false
	}
}

// hasStructDefinition checks if there's a struct definition with a specific name
// Parses the file to search for struct type declarations
// Returns true if a struct with the given name is found
func (g *Generator) hasStructDefinition(structName string, filePath string) bool {
	node, err := parseFile(filePath)
	if err != nil {
		return false
	}

	return findStructByName(node, structName) != nil
}

// extractStructFields extracts fields from structs dynamically
// Handles different type expressions and finds struct definitions
// Returns field information for the struct type
func (g *Generator) extractStructFields(expr ast.Expr, filePath string) []StructFieldInfo {
	var fields []StructFieldInfo

	switch t := expr.(type) {
	case *ast.StarExpr:
		// Handle pointer types
		return g.extractStructFields(t.X, filePath)
	case *ast.Ident:
		// Find struct definition with this name
		return g.findStructDefinition(t.Name, filePath)
	case *ast.SelectorExpr:
		// Handle selector expressions like package.Type
		return g.findStructDefinition(exprToString(t), filePath)
	}

	return fields
}

// findStructDefinition finds struct definition by name
// Parses the file to locate struct type declarations
// Returns field information if the struct is found
func (g *Generator) findStructDefinition(structName string, filePath string) []StructFieldInfo {
	var fields []StructFieldInfo

	// Parse file to find struct definition
	node, err := parseFile(filePath)
	if err != nil {
		return fields
	}

	// Find struct by name
	structType := findStructByName(node, structName)
	if structType != nil {
		fields = g.extractStructFieldsFromAST(structType)
	}

	return fields
}

// extractStructFieldsFromAST extracts fields from AST struct type
// Processes the AST representation of a struct to extract field information
// Uses TypeResolver to get full package paths for field types
func (g *Generator) extractStructFieldsFromAST(structType *ast.StructType) []StructFieldInfo {
	var fields []StructFieldInfo

	if structType.Fields == nil {
		return fields
	}

	resolver := NewTypeResolver(g.imports)
	for _, field := range structType.Fields.List {
		// Use TypeResolver to get full package paths
		fieldType := resolver.ResolveType(field.Type)

		if len(field.Names) > 0 {
			for _, name := range field.Names {
				fields = append(fields, StructFieldInfo{
					Name: name.Name,
					Type: fieldType,
					Tag:  extractFieldTag(field.Tag),
				})
			}
		} else {
			// Embedded field (anonymous field)
			fields = append(fields, StructFieldInfo{
				Name: "",
				Type: fieldType,
				Tag:  extractFieldTag(field.Tag),
			})
		}
	}

	return fields
}
