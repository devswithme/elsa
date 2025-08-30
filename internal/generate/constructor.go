package generate

import (
	"fmt"
	"go/types"
	"log"

	"golang.org/x/tools/go/packages"
)

// Constructor info
type Constructor struct {
	Name    string   // contoh: healthSvc.NewHealthService
	Params  []string // tipe param input
	Results []string // tipe return output
}

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

	// Loop semua fungsi target
	for _, fn := range funcs {
		// fn bisa dalam format: alias.FuncName, misal "healthSvc.NewHealthService"
		// untuk simplifikasi, kita ambil bagian setelah titik
		fnName := fn
		if dot := len(fn) - len(pkgPath); dot > 0 {
			fnName = fn[dot+1:]
		}

		obj := pkgs[0].Types.Scope().Lookup(fnName)
		if obj == nil {
			log.Printf("function %s not found in %s", fnName, pkgPath)
			continue
		}

		sig, ok := obj.Type().(*types.Signature)
		if !ok {
			continue
		}

		// Ambil parameter
		var params []string
		for i := 0; i < sig.Params().Len(); i++ {
			param := sig.Params().At(i)
			params = append(params, param.Type().String())
		}

		// Ambil return
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
