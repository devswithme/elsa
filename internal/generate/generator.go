package generate

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

// Generator handles the generation process for finding elsabuild files
type Generator struct {
	imports map[string]string
}

// NewGenerator creates a new Generator instance
func NewGenerator() *Generator {
	return &Generator{}
}

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

func (g *Generator) processGenerateDependencies(target string) error {
	goModDir, err := g.FindGoModDir(target)
	if err != nil {
		return err
	}

	// Ekstrak fungsi yang mengandung elsa.Generate
	funcs, err := g.ExtractElsaGenerateFuncs(target)
	if err != nil {
		fmt.Printf("Warning: failed to extract generate functions: %v\n", err)
	} else {
		fmt.Println("=== FUNGSI YANG MENGANDUNG ELSA.GENERATE ===")
		for _, funcInfo := range funcs {
			fmt.Printf("Fungsi: %s\n", funcInfo.FuncName)
			fmt.Printf("Package: %s\n", funcInfo.PkgName)

			if len(funcInfo.Params) > 0 {
				fmt.Printf("Parameter:\n")
				for i, param := range funcInfo.Params {
					fmt.Printf("  [%d] %s: %s\n", i, param.Name, param.Type)
				}
			}

			if len(funcInfo.Results) > 0 {
				fmt.Printf("Return:\n")
				for _, result := range funcInfo.Results {
					fmt.Printf("  - %s", result.Type)
					if result.IsStruct {
						fmt.Printf(" (struct)")
						if len(result.StructFields) > 0 {
							fmt.Printf(":\n")
							for _, field := range result.StructFields {
								fmt.Printf("    * %s: %s", field.Name, field.Type)
								if field.Tag != "" {
									fmt.Printf(" `%s`", field.Tag)
								}
								fmt.Println()
							}
						} else {
							fmt.Printf(" (field details not available)")
						}
					}
					fmt.Println()
				}
			}

			if len(funcInfo.GenerateParams) > 0 {
				fmt.Printf("Elsa.Generate Parameters:\n")
				for i, param := range funcInfo.GenerateParams {
					fmt.Printf("  [%d] %s\n", i, param)
				}
			}
			fmt.Println()
		}
	}

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

// ResultInfo menyimpan informasi return value
type ResultInfo struct {
	Type         string
	IsStruct     bool
	StructFields []StructFieldInfo
}

// StructFieldInfo menyimpan informasi field struct
type StructFieldInfo struct {
	Name string
	Type string
	Tag  string
}

// ExtractElsaGenerateFuncs mengekstrak fungsi yang mengandung elsa.Generate
func (g *Generator) ExtractElsaGenerateFuncs(target string) ([]FuncInfo, error) {
	fset := token.NewFileSet()

	// Parse file menjadi AST
	node, err := parser.ParseFile(fset, target, nil, parser.AllErrors)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %v", err)
	}

	// Parse imports untuk mapping alias ke package path
	g.parseImports(node)

	var funcs []FuncInfo

	// Traverse AST untuk mencari fungsi
	ast.Inspect(node, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		// Cek apakah fungsi mengandung elsa.Generate
		if !containsElsaGenerate(fn) {
			return true
		}

		// Ekstrak informasi fungsi
		funcInfo := FuncInfo{
			FuncName: fn.Name.Name,
			PkgName:  node.Name.Name,
		}

		// Ekstrak parameter dengan package path yang lengkap
		funcInfo.Params = g.extractParamsWithImports(fn.Type.Params)

		// Ekstrak return values
		if fn.Type.Results != nil {
			funcInfo.Results = g.extractResults(fn.Type.Results, target)
		}

		// Ekstrak parameter dari elsa.Generate
		funcInfo.GenerateParams = g.extractGenerateParams(fn)

		funcs = append(funcs, funcInfo)
		return true
	})

	return funcs, nil
}

// parseImports mengumpulkan mapping alias ke package path
func (g *Generator) parseImports(node *ast.File) {
	g.imports = make(map[string]string)

	for _, imp := range node.Imports {
		pkgPath := strings.Trim(imp.Path.Value, `"`)
		if imp.Name != nil {
			// alias import, contoh: healthRepo "github.com/xxx/health"
			g.imports[imp.Name.Name] = pkgPath
		} else {
			// tanpa alias → ambil last segment
			parts := strings.Split(pkgPath, "/")
			g.imports[parts[len(parts)-1]] = pkgPath
		}
	}
}

// containsElsaGenerate mengecek apakah fungsi mengandung elsa.Generate
func containsElsaGenerate(fn *ast.FuncDecl) bool {
	if fn.Body == nil {
		return false
	}

	found := false
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			if x, ok := sel.X.(*ast.Ident); ok && x.Name == "elsa" && sel.Sel.Name == "Generate" {
				found = true
				return false
			}
		}
		return true
	})

	return found
}

// extractGenerateParams mengekstrak parameter dari elsa.Generate
func (g *Generator) extractGenerateParams(fn *ast.FuncDecl) []string {
	var params []string

	if fn.Body == nil {
		return params
	}

	ast.Inspect(fn.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			if x, ok := sel.X.(*ast.Ident); ok && x.Name == "elsa" && sel.Sel.Name == "Generate" {
				// Ekstrak semua argument dari elsa.Generate
				for _, arg := range call.Args {
					if ident, ok := arg.(*ast.Ident); ok {
						params = append(params, ident.Name)
					} else {
						// Fallback untuk expression yang bukan ident
						params = append(params, exprToString(arg))
					}
				}
				return false
			}
		}
		return true
	})

	return params
}

// extractParamsWithImports mengekstrak parameter dengan package path yang lengkap
func (g *Generator) extractParamsWithImports(fieldList *ast.FieldList) []ParamInfo {
	if fieldList == nil {
		return nil
	}

	var params []ParamInfo
	for _, field := range fieldList.List {
		typ := g.resolveTypeWithImports(field.Type)

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

// resolveTypeWithImports resolve tipe data dengan package path yang lengkap
func (g *Generator) resolveTypeWithImports(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		// Cek apakah ini tipe dari package yang di-import
		if pkgPath, exists := g.imports[t.Name]; exists {
			// Ini adalah package alias, gunakan package path lengkap
			return pkgPath + "." + t.Name
		}
		// Ini adalah tipe lokal atau built-in
		return t.Name
	case *ast.StarExpr:
		return "*" + g.resolveTypeWithImports(t.X)
	case *ast.SelectorExpr:
		// Handle selector expressions seperti package.Type
		if x, ok := t.X.(*ast.Ident); ok {
			if pkgPath, exists := g.imports[x.Name]; exists {
				// Ini adalah package alias, gunakan package path lengkap
				return pkgPath + "." + t.Sel.Name
			}
			// Ini adalah selector dari tipe lokal
			return x.Name + "." + t.Sel.Name
		}
		return exprToString(t)
	case *ast.ArrayType:
		return "[]" + g.resolveTypeWithImports(t.Elt)
	case *ast.MapType:
		return "map[" + g.resolveTypeWithImports(t.Key) + "]" + g.resolveTypeWithImports(t.Value)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.FuncType:
		return "func(...)"
	case *ast.ChanType:
		if t.Dir == ast.SEND {
			return "chan<- " + g.resolveTypeWithImports(t.Value)
		} else if t.Dir == ast.RECV {
			return "<-chan " + g.resolveTypeWithImports(t.Value)
		}
		return "chan " + g.resolveTypeWithImports(t.Value)
	default:
		return exprToString(t)
	}
}

// extractResults mengekstrak informasi return values
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

		// Cek apakah ini struct type
		if g.isStructType(field.Type, filePath) {
			resultInfo.IsStruct = true
			resultInfo.StructFields = g.extractStructFields(field.Type, filePath)
		}

		results = append(results, resultInfo)
	}
	return results
}

// isStructType mengecek apakah tipe adalah struct secara dinamis
func (g *Generator) isStructType(expr ast.Expr, filePath string) bool {
	switch t := expr.(type) {
	case *ast.StarExpr:
		return g.isStructType(t.X, filePath)
	case *ast.Ident:
		// Cek apakah ada definisi struct dengan nama ini di file
		return g.hasStructDefinition(t.Name, filePath)
	case *ast.SelectorExpr:
		// Ini bisa jadi struct dari package lain
		return true
	default:
		return false
	}
}

// hasStructDefinition mengecek apakah ada definisi struct dengan nama tertentu
func (g *Generator) hasStructDefinition(structName string, filePath string) bool {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
	if err != nil {
		return false
	}

	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		typeDecl, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		// Cek apakah nama struct cocok
		if typeDecl.Name.Name != structName {
			return true
		}

		// Cek apakah ini struct type
		_, ok = typeDecl.Type.(*ast.StructType)
		if ok {
			found = true
			return false // Stop traversal setelah ketemu
		}
		return true
	})

	return found
}

// extractStructFields mengekstrak field dari struct secara dinamis
func (g *Generator) extractStructFields(expr ast.Expr, filePath string) []StructFieldInfo {
	var fields []StructFieldInfo

	switch t := expr.(type) {
	case *ast.StarExpr:
		// Handle pointer types
		return g.extractStructFields(t.X, filePath)
	case *ast.Ident:
		// Cari definisi struct dengan nama ini
		return g.findStructDefinition(t.Name, filePath)
	case *ast.SelectorExpr:
		// Handle selector expressions seperti package.Type
		return g.findStructDefinition(exprToString(t), filePath)
	}

	return fields
}

// findStructDefinition mencari definisi struct berdasarkan nama
func (g *Generator) findStructDefinition(structName string, filePath string) []StructFieldInfo {
	var fields []StructFieldInfo

	// Parse file untuk mencari definisi struct
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
	if err != nil {
		return fields
	}

	// Traverse AST untuk mencari struct declaration
	ast.Inspect(node, func(n ast.Node) bool {
		typeDecl, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		// Cek apakah nama struct cocok
		if typeDecl.Name.Name != structName {
			return true
		}

		// Cek apakah ini struct type
		structType, ok := typeDecl.Type.(*ast.StructType)
		if !ok {
			return true
		}

		// Ekstrak field dari struct
		fields = g.extractStructFieldsFromAST(structType)
		return false // Stop traversal setelah ketemu
	})

	return fields
}

// extractStructFieldsFromAST mengekstrak field dari AST struct type
func (g *Generator) extractStructFieldsFromAST(structType *ast.StructType) []StructFieldInfo {
	var fields []StructFieldInfo

	if structType.Fields == nil {
		return fields
	}

	for _, field := range structType.Fields.List {
		// Gunakan resolveTypeWithImports untuk mendapatkan package path lengkap
		fieldType := g.resolveTypeWithImports(field.Type)

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

// extractFieldTag mengekstrak tag dari field
func extractFieldTag(tag *ast.BasicLit) string {
	if tag == nil {
		return ""
	}
	return tag.Value
}

// exprToString mengkonversi AST expression ke string
func exprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + exprToString(t.X)
	case *ast.SelectorExpr:
		return exprToString(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		return "[]" + exprToString(t.Elt)
	case *ast.MapType:
		return "map[" + exprToString(t.Key) + "]" + exprToString(t.Value)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.FuncType:
		return "func(...)"
	case *ast.ChanType:
		if t.Dir == ast.SEND {
			return "chan<- " + exprToString(t.Value)
		} else if t.Dir == ast.RECV {
			return "<-chan " + exprToString(t.Value)
		}
		return "chan " + exprToString(t.Value)
	default:
		return fmt.Sprintf("%T", t)
	}
}
