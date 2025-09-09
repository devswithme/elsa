package generate

import (
	"fmt"
	"path/filepath"

	"github.com/risoftinc/elsa/constants"
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
			fmt.Printf("%s Error: %s failed to process generate dependencies: %v\n", constants.ErrorEmoji, filepath.Join(targetDir, file), err)
			return nil
		}
	}

	return nil
}
