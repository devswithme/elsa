package generate

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type FuncInfo struct {
	FuncName       string
	PkgName        string
	PkgPath        string
	Params         []ParamInfo
	Results        []ResultInfo
	GenerateParams []string
}

type ParamInfo struct {
	Name string
	Type string
}

// parseElsaSets menerima path file dan return map[string][]string
// key = nama variabel (RepositorySet, ServicesSet, dll)
// value = list fungsi di dalam elsa.Set(...)
func (g *Generator) ParseElsaSets(path string) (map[string][]FuncInfo, error) {
	fset := token.NewFileSet()

	// parse file jadi AST
	node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	// --- Kumpulin imports ---
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

	results := make(map[string][]FuncInfo)

	// loop semua deklarasi di file
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

				// pastikan itu elsa.Set(...)
				funIdent, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					continue
				}
				if fmt.Sprint(funIdent.X) != "elsa" || funIdent.Sel.Name != "Set" {
					continue
				}

				// ambil semua argumen di dalam elsa.Set
				var funcs []FuncInfo
				for _, arg := range call.Args {
					switch v := arg.(type) {
					case *ast.SelectorExpr:
						// contoh: healthRepo.NewHealthRepositories
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

	return results, nil
}
