package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
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
			fmt.Printf("Error: %s failed to process generate dependencies: %v\n", filepath.Join(targetDir, file), err)
			return nil
		}
	}

	return nil
}

type (
	ElsaGenFile struct {
		Target           string
		PackageName      string
		ImportedPackages map[string]ElsaImportedPackages
		StructsData      map[string][]StructFieldInfo
		Functions        map[string]ElsaGenFunction
		FuncGenerated    map[string][]FuncInfo
		Sets             map[string][]FuncInfo
	}

	ElsaImportedPackages struct {
		Alias    string
		UseAlias bool
	}

	ElsaGenFunction struct {
		SourcePackages map[string]ElsaSourceDetail
		Params         []TypeInfo
		Results        []TypeInfo
	}

	ElsaSourceDetail struct {
		VariableName string
		UsePointer   bool
	}
)

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

	if len(funcs) == 0 {
		return fmt.Errorf("no functions with elsa.Generate calls found")
	}

	sets, err := g.ParseElsaSets(goModDir, target)
	if err != nil {
		return err
	}

	var (
		elsaGenFile ElsaGenFile = ElsaGenFile{
			Target:           target,
			PackageName:      funcs[0].PkgName,
			ImportedPackages: make(map[string]ElsaImportedPackages),
			StructsData:      make(map[string][]StructFieldInfo),
			Functions:        make(map[string]ElsaGenFunction),
			FuncGenerated:    make(map[string][]FuncInfo),
			Sets:             sets,
		}
		counterAlias            = make(map[string]int)
		funcSetImportedPackages = func(typeInfo TypeInfo) {
			if _, ok := elsaGenFile.ImportedPackages[typeInfo.Package]; !ok {
				counterAlias[typeInfo.Alias]++
				useAlias := false
				if counterAlias[typeInfo.Alias] > 1 {
					typeInfo.Alias = fmt.Sprintf("%s%d", typeInfo.Alias, counterAlias[typeInfo.Alias])
					useAlias = true
				}
				elsaGenFile.ImportedPackages[typeInfo.Package] = ElsaImportedPackages{
					Alias:    typeInfo.Alias,
					UseAlias: useAlias,
				}
			}
		}
	)

	// Default import packages "elsa"
	counterAlias["elsa"]++

	for _, fn := range funcs {
		if _, ok := elsaGenFile.Functions[fn.FuncName]; !ok {
			elsaGenFile.Functions[fn.FuncName] = ElsaGenFunction{
				SourcePackages: make(map[string]ElsaSourceDetail),
				Params:         make([]TypeInfo, 0),
				Results:        make([]TypeInfo, 0),
			}
		}

		existingParams := make(map[string]bool)
		for _, param := range fn.Params {
			if _, ok := existingParams[param.Type]; ok {
				return fmt.Errorf("duplicate parameter on function %s: %s", fn.FuncName, param.Type)
			} else {
				existingParams[param.Type] = true
			}

			typeInfo := extractType(param.Type)
			typeInfo.ParamName = param.Name
			funcData := elsaGenFile.Functions[fn.FuncName]
			funcData.Params = append(funcData.Params, typeInfo)
			elsaGenFile.Functions[fn.FuncName] = funcData

			elsaGenFile.Functions[fn.FuncName].SourcePackages[fmt.Sprintf("%s.%s", typeInfo.Package, typeInfo.DataType)] = ElsaSourceDetail{
				VariableName: lowerFirst(param.Name),
				UsePointer:   typeInfo.UsePointer,
			}
			funcSetImportedPackages(typeInfo)
		}

		for _, result := range fn.Results {
			typeInfo := extractType(result.Type)
			funcData := elsaGenFile.Functions[fn.FuncName]
			funcData.Results = append(funcData.Results, typeInfo)
			elsaGenFile.Functions[fn.FuncName] = funcData

			if _, ok := elsaGenFile.StructsData[typeInfo.Package]; result.IsStruct && !ok {
				elsaGenFile.StructsData[typeInfo.Package] = result.StructFields
			}
		}

		var genFunctions []FuncInfo
		for _, gp := range fn.GenerateParams {
			genFunctions = append(genFunctions, sets[gp]...)
		}

		for {
			var isFounded bool
			for idx, gf := range genFunctions {
				isFounded = true
				for _, param := range gf.Params {
					typeInfo := extractType(param.Type)
					if elsaGenFile.Functions[fn.FuncName].SourcePackages[fmt.Sprintf("%s.%s", typeInfo.Package, typeInfo.DataType)].VariableName == "" {
						isFounded = false
						break
					}
				}

				if !isFounded {
					continue
				}

				for _, result := range gf.Results {
					typeInfo := extractType(result.Type)
					elsaGenFile.Functions[fn.FuncName].SourcePackages[fmt.Sprintf("%s.%s", typeInfo.Package, typeInfo.DataType)] = ElsaSourceDetail{
						VariableName: lowerFirst(typeInfo.DataType),
						UsePointer:   typeInfo.UsePointer,
					}
					funcSetImportedPackages(typeInfo)
				}

				isFounded = true
				elsaGenFile.FuncGenerated[fn.FuncName] = append(elsaGenFile.FuncGenerated[fn.FuncName], gf)
				genFunctions = removeArrayByIndex(genFunctions, idx)
				break
			}

			if len(genFunctions) == 0 {
				break
			}

			if !isFounded {
				return fmt.Errorf("failed get source from function %s", genFunctions[0].FuncName)
			}
		}

		variableName := make(map[string]int)
		for key, source := range elsaGenFile.Functions[fn.FuncName].SourcePackages {
			variableName[source.VariableName]++
			if variableName[source.VariableName] > 1 {
				source.VariableName = fmt.Sprintf("%s%d", source.VariableName, variableName[source.VariableName])
				elsaGenFile.Functions[fn.FuncName].SourcePackages[key] = source
			}
		}
	}

	if err := g.GenerateElsaGenFile(target, elsaGenFile); err != nil {
		return err
	}

	return nil
}

// GenerateElsaGenFile generates an elsa_gen.go file for the target
// Creates a minimal generated file with just the package declaration
// The package name is extracted from the target file
func (g *Generator) GenerateElsaGenFile(target string, elsaGenFile ElsaGenFile) error {
	// Generate the content for elsa_gen.go
	content := g.generateElsaGenContent(elsaGenFile)

	// Determine the output path (same directory as target)
	outputDir := filepath.Dir(target)
	outputPath := filepath.Join(outputDir, "elsa_gen.go")

	// Write the generated file
	err := os.WriteFile(outputPath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write elsa_gen.go: %v", err)
	}

	fmt.Printf("Generated: %s\n", outputPath)
	return nil
}

// generateElsaGenHeader generates the standard header for elsa_gen.go files
// Contains build tags and generation directives
func (g *Generator) generateElsaGenHeader() string {
	return `// Code generated by Elsa. DO NOT EDIT.

//go:generate go run -mod=mod github.com/risoftinc/elsa/cmd/elsa gen
//go:build !elsabuild
// +build !elsabuild

`
}

// generateElsaGenContent generates the content for elsa_gen.go
// Creates a file with header, package declaration, and optional additional content
func (g *Generator) generateElsaGenContent(elsaGenFile ElsaGenFile) string {
	var content string

	// Add header
	content += g.generateElsaGenHeader()

	// Add package declaration
	content += fmt.Sprintf("package %s\n\n", elsaGenFile.PackageName)

	// Add Import Selection
	content += g.generateElsaGenImportSelection(elsaGenFile)

	// Inject From
	content += fmt.Sprintf("// This file generated from %s at %s\n\n", getNameFile(elsaGenFile.Target), time.Now().Format("2006-01-02 15:04:05"))

	// Add Structs
	content += g.generateElsaGenStructs(elsaGenFile)

	// Add Functions
	content += g.generateElsaGenFunctions(elsaGenFile)

	// Add additional content (can be extended in the future)
	content += g.generateElsaGenAdditionalContent()

	return content
}

func (g *Generator) generateElsaGenImportSelection(elsaGenFile ElsaGenFile) string {
	content := "import (\n\t\"github.com/risoftinc/elsa\"\n\n"

	// Convert map to slice for consistent or[]string{}dering
	var packages []string
	for pkg := range elsaGenFile.ImportedPackages {
		if _, ok := elsaGenFile.StructsData[pkg]; ok || isBuiltinType(pkg) {
			continue
		}
		importedPackage := elsaGenFile.ImportedPackages[pkg]
		packages = append(packages, fmt.Sprintf("\t%s \"%s\"\n", importedPackage.Alias, pkg))
	}

	// Sort packages alphabetically for consistent order
	sort.Strings(packages)

	// Generate imports in sorted order
	for _, pkg := range packages {
		content += pkg
	}

	return content + ")\n\n"
}

func (g *Generator) generateElsaGenStructs(elsaGenFile ElsaGenFile) string {
	content := ""
	for pkg, structs := range elsaGenFile.StructsData {
		content += fmt.Sprintf("type %s struct {\n", pkg)

		// Find the maximum field name length for alignment
		maxNameLength := 0
		for _, structData := range structs {
			if len(structData.Name) > maxNameLength {
				maxNameLength = len(structData.Name)
			}
		}

		// Generate struct fields with proper alignment
		for _, structData := range structs {
			et := extractType(structData.Type)
			importedPackage := elsaGenFile.ImportedPackages[et.Package].Alias + "." + et.DataType

			// Create padding spaces for alignment
			padding := strings.Repeat(" ", maxNameLength-len(structData.Name))
			content += fmt.Sprintf("\t%s%s %s\n", structData.Name, padding, importedPackage)
		}
		content += "}\n\n"
	}
	return content
}

func (g *Generator) generateElsaGenFunctions(elsaGenFile ElsaGenFile) string {
	content := ""
	for name, function := range elsaGenFile.Functions {
		var params, results []string
		for _, param := range function.Params {
			importedPackage := elsaGenFile.ImportedPackages[param.Package].Alias
			if param.UsePointer {
				importedPackage = "*" + importedPackage
			}
			if _, ok := elsaGenFile.StructsData[param.Package]; !ok && !isBuiltinType(param.Package) {
				importedPackage += "." + param.DataType
			}

			params = append(params, fmt.Sprintf("%s %s",
				function.SourcePackages[param.Package+"."+param.DataType].VariableName, importedPackage,
			))
		}

		for _, result := range function.Results {
			if result.UsePointer {
				result.Package = "*" + result.Package
			}
			results = append(results, result.Package)
		}

		returnStr := strings.Join(results, ", ")
		if len(results) > 1 {
			returnStr = "(" + returnStr + ")"
		}

		content += fmt.Sprintf("func %s(%s) %s {\n", name, strings.Join(params, ", "), returnStr)

		var elsaGeneratedVariables []string
		for _, generated := range elsaGenFile.FuncGenerated[name] {
			importedPackage := elsaGenFile.ImportedPackages[generated.PkgPath].Alias
			var params, results []string
			for _, param := range generated.Params {
				et := extractType(param.Type)
				source := function.SourcePackages[et.Package+"."+et.DataType]
				if source.UsePointer == et.UsePointer {
					params = append(params, source.VariableName)
				} else if et.UsePointer && !source.UsePointer {
					params = append(params, "*"+source.VariableName)
				} else {
					params = append(params, "&"+source.VariableName)
				}
			}

			for _, result := range generated.Results {
				et := extractType(result.Type)
				results = append(results, function.SourcePackages[et.Package+"."+et.DataType].VariableName)
			}

			elsaGeneratedVariables = append(elsaGeneratedVariables, results...)

			content += fmt.Sprintf("\t%s := %s.%s(%s)\n", strings.Join(results, ", "), importedPackage, generated.FuncName, strings.Join(params, ", "))
		}

		content += fmt.Sprintf("\n\telsa.Generate(%s)\n", strings.Join(elsaGeneratedVariables, ", "))

		content += "\treturn "
		for _, result := range function.Results {
			stuctContent := ""
			for _, structData := range elsaGenFile.StructsData[result.Package] {
				stuctContent += fmt.Sprintf("\t\t%s: %s,\n", structData.Name, function.SourcePackages[structData.Type].VariableName)
			}
			if stuctContent != "" {
				stuctContent = "{\n" + stuctContent + "\t}"
			}
			if result.UsePointer {
				result.Package = "&" + result.Package
			}

			content += result.Package

			if stuctContent != "" {
				content += stuctContent
			}

			content += "\n"
		}

		content += "}\n\n"
	}
	return content
}

// generateElsaGenAdditionalContent generates additional content for elsa_gen.go
// This function can be easily extended to add more content like imports, types, etc.
// Currently returns empty string, but can be customized as needed
func (g *Generator) generateElsaGenAdditionalContent() string {
	// TODO: Add additional content here as needed
	// Examples:
	// - Import statements
	// - Type definitions
	// - Function declarations
	// - Constants
	// - Variables

	return ""
}
