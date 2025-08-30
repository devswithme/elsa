package generate

import (
	"fmt"
	"os"
	"path/filepath"
)

// Cari directory go.mod dimulai dari lokasi file tertentu (misal elsabuild.go)
func (g *Generator) FindGoModDir(start string) (string, error) {
	dir := filepath.Dir(start)

	for {
		// cek apakah ada go.mod di folder ini
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir, nil
		}

		// kalau sudah di root dan tidak ketemu
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found from %s upward", start)
		}
		dir = parent
	}
}
