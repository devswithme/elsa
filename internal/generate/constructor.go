package generate

import (
	"fmt"
	"go/types"
	"log"

	"golang.org/x/tools/go/packages"
)

// Constructor represents information about a constructor function
// Contains the function name, input parameters, and return types
type Constructor struct {
	Name    string   // Example: healthSvc.NewHealthService
	Params  []string // Input parameter types
	Results []string // Return output types
}

// LoadConstructors loads constructor information for specified functions from a package
// Uses the Go packages API to analyze function signatures and extract type information
// Returns a map of function names to their constructor information
func (g *Generator) LoadConstructors(goModDir string, pkgPath string, funcs []string) (map[string]Constructor, error) {
	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedDeps | packages.NeedTypesInfo,
		Dir:  goModDir,
	}
	pkgs, err := packages.Load(cfg, pkgPath)
	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 {
		return nil, fmt.Errorf("package not found: %s", pkgPath)
	}

	constructors := make(map[string]Constructor)

	// Loop through all target functions
	for _, fn := range funcs {
		// Function can be in format: alias.FuncName, example "healthSvc.NewHealthService"
		// For simplification, we extract the part after the dot
		fnName := extractFunctionName(fn, pkgPath)

		obj := pkgs[0].Types.Scope().Lookup(fnName)
		if obj == nil {
			log.Printf("function %s not found in %s", fnName, pkgPath)
			continue
		}

		sig, ok := obj.Type().(*types.Signature)
		if !ok {
			continue
		}

		// Extract parameters
		var params []string
		for i := 0; i < sig.Params().Len(); i++ {
			param := sig.Params().At(i)
			params = append(params, param.Type().String())
		}

		// Extract return values
		var results []string
		for i := 0; i < sig.Results().Len(); i++ {
			res := sig.Results().At(i)
			results = append(results, res.Type().String())
		}

		constructors[fn] = Constructor{
			Name:    fn,
			Params:  params,
			Results: results,
		}
	}

	return constructors, nil
}
