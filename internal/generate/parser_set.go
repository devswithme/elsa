package generate

import (
	"fmt"
	"go/ast"
	"go/token"
)

// FuncInfo represents information about a function found in elsa.Set calls
// Contains function name, package information, parameters, results, and generation parameters
type FuncInfo struct {
	FuncName       string       // Name of the function
	PkgName        string       // Package alias or name
	PkgPath        string       // Full package path
	Params         []ParamInfo  // Function parameters
	Results        []ResultInfo // Function return values
	GenerateParams []string     // Parameters passed to elsa.Generate
}

// ParamInfo represents information about a function parameter
// Contains the parameter name and type
type ParamInfo struct {
	Name string // Parameter name
	Type string // Parameter type
}

// ParseElsaSets parses a file and returns a map of variable names to function information
// Key = variable name (RepositorySet, ServicesSet, etc.)
// Value = list of functions inside elsa.Set(...) calls
// This function analyzes Go AST to find elsa.Set declarations and extract function information
func (g *Generator) ParseElsaSets(goModDir, path string) (map[string][]FuncInfo, error) {
	// Parse file into AST
	node, err := parseFile(path)
	if err != nil {
		return nil, err
	}

	// --- Collect imports ---
	g.imports = parseImports(node)

	results := make(map[string][]FuncInfo)

	// Loop through all declarations in the file
	for _, decl := range node.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.VAR {
			continue
		}

		for _, spec := range gen.Specs {
			vspec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for _, name := range vspec.Names {
				if len(vspec.Values) == 0 {
					continue
				}
				call, ok := vspec.Values[0].(*ast.CallExpr)
				if !ok {
					continue
				}

				// Ensure it's elsa.Set(...)
				funIdent, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					continue
				}
				if fmt.Sprint(funIdent.X) != "elsa" || funIdent.Sel.Name != "Set" {
					continue
				}

				// Extract all arguments inside elsa.Set
				var funcs []FuncInfo
				for _, arg := range call.Args {
					switch v := arg.(type) {
					case *ast.SelectorExpr:
						// Example: healthRepo.NewHealthRepositories
						pkgIdent, ok := v.X.(*ast.Ident)
						if !ok {
							continue
						}
						alias := pkgIdent.Name
						funcs = append(funcs, FuncInfo{
							FuncName: v.Sel.Name,
							PkgName:  alias,
							PkgPath:  g.imports[alias],
						})
					case *ast.Ident:
						funcs = append(funcs, FuncInfo{
							FuncName: v.Name,
						})
					}
				}
				results[name.Name] = funcs
			}
		}
	}

	for key, result := range results {
		for x, s := range result {
			constructors, err := g.LoadConstructors(goModDir, s.PkgPath, []string{s.FuncName})
			if err != nil {
				return nil, err
			}

			for _, c := range constructors.Params {
				results[key][x].Params = append(results[key][x].Params, ParamInfo{Type: c})
			}

			for _, c := range constructors.Results {
				results[key][x].Results = append(results[key][x].Results, ResultInfo{Type: c})
			}
		}
	}

	return results, nil
}
