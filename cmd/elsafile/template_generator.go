package elsafile

import (
	"github.com/risoftinc/elsa/internal/elsafile"
)

// TemplateGenerator handles creation of default Elsafile templates
type TemplateGenerator = elsafile.TemplateGenerator

// NewTemplateGenerator creates a new TemplateGenerator instance
func NewTemplateGenerator() *TemplateGenerator {
	return elsafile.NewTemplateGenerator()
}
